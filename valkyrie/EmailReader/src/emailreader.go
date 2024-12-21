package main

import (
	"os"
	"fmt"
	"net"
	"time"
	"strconv"
	"strings"
//	"cipherize"

	"math/rand"
	"io/ioutil"
	"net/http"
	"net/url"
        "encoding/hex"
        "crypto/md5"

	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap"
)

var usesecurelogging bool
var loggerurl string
var logkey string

var dupprotection bool
var dupprottime int64
var prottimes map[string]int64

var workers []string
var current_worker int

var mailhost string
var mailuser string
var mailpass string
var trigger_subject string

var mailusetls bool 
var imapfolder string

var sleep_timeout int

var license string
var business string

type ActionMessage struct {
	Host string
	Action string
	ActionLevel string
	Contact string
	MiscParams string
	OnRemote string
	Image string
	GitRepo string
}

func Log(message string, level int) error {
	var Level string
	Service := "Valkyrie EmailReader"
	
	switch level {
		case 0:
			Level = "info"
		case 1:
			Level = "warning"
		case 2:
			Level = "error"
	}

	encservice := url.QueryEscape(Service)
	enclevel := url.QueryEscape(Level)
	encmessage := url.QueryEscape(message)
	enckey := url.QueryEscape(logkey)

	if usesecurelogging {
		_, err := http.Get(loggerurl + "log?service=" + encservice + "&loglevel=" + enclevel + "&message=" + encmessage + "&logkey=" + enckey)
		if err != nil {
			FileLog("Failed to log to logging server: " + err.Error())
			return err
		}
	} else {
		_, err := http.Get(loggerurl + "log?service=" + encservice + "&loglevel=" + enclevel + "&message=" + encmessage)
		if err != nil {
			FileLog("Failed to log to logging server: " + err.Error())
			return err
		}
	}

	return nil
}

func FileLog(message string) error {
	file, err := os.OpenFile("./emailreader.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("failed to open log file: " + err.Error())
		return err
	}
	defer file.Close()

        current_time := time.Now().Local()
        t := current_time.Format("Jan 02 2006 03:04:05")
	_, err = file.WriteString(t + " - Valkyrie EmailReader: " + message + "\n")

	if err != nil {
		fmt.Println("failed to write to log file: " + err.Error())
		return err
	}

	return nil
}

func HandleNewEmails(action ActionMessage) {

        // Duplication Protection
        if dupprotection == true {
                hasher := md5.New()
                hasher.Write([]byte(action.Host + action.Action))
                tmpmd5sum := hex.EncodeToString(hasher.Sum(nil))
                if val, ok := prottimes[tmpmd5sum]; ok {
                        now := time.Now()
                        epoch := now.Unix()

                        if epoch < val {
                                // Discard the message
                                Log("Duplicate Message Being Dropped | " + action.Host + " | " + action.Action, 0)
                                return
                        } else {
                                prottimes[tmpmd5sum] = (epoch + dupprottime)
                        }
                } else {
                        now := time.Now()
                        epoch := now.Unix()
                        prottimes[tmpmd5sum] = (epoch + dupprottime)
                }
        }

	enchost := url.QueryEscape(action.Host)
	encaction := url.QueryEscape(action.Action)
	encactionlevel := url.QueryEscape(action.ActionLevel)
	enccontact := url.QueryEscape(action.Contact)
	encmiscparams := url.QueryEscape(action.MiscParams)
	enconremote := url.QueryEscape(action.OnRemote)
	encimage := url.QueryEscape(action.Image)
	encgitrepo := url.QueryEscape(action.GitRepo)

	resp, err := http.Get(workers[current_worker] + "/trigger?host=" + enchost + "&action=" + encaction + "&actionlevel=" + encactionlevel + "&contact=" + enccontact + "&onremote=" + enconremote + "&image=" + encimage + "&gitrepo=" + encgitrepo + "&miscparams=" + encmiscparams)
	if err != nil {
		Log("Failed To Launch New Valkyrie: " + err.Error(), 0)
		return
	}
	defer resp.Body.Close()

	if current_worker == (len(workers) -1) {
		current_worker = 0
	} else {
		current_worker++
	}

	Log("Launched New Valkyrie", 0)
}

func HandleIMAP() {
	var c *client.Client 
	var err error

	if mailusetls {
		c, err = client.DialTLS(mailhost, nil)
		if err != nil {
			Log("Failed Dialing Imap Server: " + err.Error(), 0)
			return
		}
	} else {
		c, err = client.Dial(mailhost)
		if err != nil {
			Log("Failed Dialing Imap Server: " + err.Error(), 0)
			return
		}
	}

	defer c.Logout()

	if err := c.Login(mailuser, mailpass); err != nil {
		Log("Failed To Login To Server: " + err.Error(), 0)
		return
	}

	mbox, err := c.Select(imapfolder, false)
	if err != nil {
		Log("Failed To Select Folder: " + err.Error(), 0)
		return
	}

	from := uint32(1)
	to := mbox.Messages

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 10)
	go func() {
		if err := c.Fetch(seqset, items, messages); err != nil {
			Log("failed to get message: " + err.Error(), 0)
		}
	}()

	subj := ""
	host := ""
	action := ""
	actionlevel := ""
	contact := ""
	miscparams := ""
	onremote := ""
	image := ""
	gitrepo := ""

	for msg := range messages {
		r := msg.GetBody(section)
		str, err := ioutil.ReadAll(r)
		if err != nil {
			Log("Failed To Read Body: " + err.Error(), 0)
			continue
		}

		parts := strings.Split(string(str), "\n")
		for _, bpart := range(parts) {
			if strings.HasPrefix(strings.ToLower(bpart), "subject:") {
				subj = strings.TrimSpace(bpart[8:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "host:") {
				host = strings.TrimSpace(bpart[5:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "action:") {
				action = strings.TrimSpace(bpart[7:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "actionlevel:") {
				actionlevel = strings.TrimSpace(bpart[12:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "contact:") {
				contact = strings.TrimSpace(bpart[8:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "miscparams:") {
				miscparams = strings.TrimSpace(bpart[11:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "onremote:") {
				onremote = strings.TrimSpace(bpart[9:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "image:") {
				image = strings.TrimSpace(bpart[6:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "gitrepo:") {
				gitrepo = strings.TrimSpace(bpart[8:])
			}
		}

		if len(subj) == 0 || subj != trigger_subject {
			continue
		}
		
		if len(host) == 0 {
			continue
		}

		if len(action) == 0 {
			continue
		}

		if len(actionlevel) == 0 {
			continue
		}

		if len(image) == 0 {
			continue
		}

		actionmsg := ActionMessage{}
		actionmsg.Host = host
		actionmsg.Action = action
		actionmsg.ActionLevel = actionlevel
		actionmsg.Contact = contact
		actionmsg.MiscParams = miscparams
		actionmsg.OnRemote = onremote
		actionmsg.Image = image
		actionmsg.GitRepo = gitrepo
		HandleNewEmails(actionmsg)
	}

}

func HandlePop() {
	var deleteids []string

	conn, err := net.Dial("tcp", mailhost)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	buff := make([]byte, 2048) 

	conn.Read(buff)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conn.Write([]byte("USER " + mailuser + "\n"))
	conn.Read(buff)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conn.Write([]byte("PASS " + mailpass + "\n"))
	conn.Read(buff)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conn.Write([]byte("LIST\n"))
	conn.Read(buff)

	list := string(buff)
	parts := strings.Split(list, "\n")
	nparts := strings.Split(parts[0], " ")
	strcount := nparts[1]

	if strcount == "0" {
		return
	}

	count, err := strconv.Atoi(strcount)
	if err != nil {
		return
	}

	for i := 1; i <= count; i++ {
		msgid := strconv.Itoa(i)
		conn.Write([]byte("RETR " + msgid + "\n"))
		conn.Read(buff)
		body := string(buff)
		bparts := strings.Split(body, "\n")
		
		subj := ""
		host := ""
		action := ""
		actionlevel := ""
		contact := ""
		miscparams := ""
		onremote := ""
		image := ""
		gitrepo := ""

		for _, bpart := range bparts {
			if strings.HasPrefix(strings.ToLower(bpart), "subject:") {
				subj = strings.TrimSpace(bpart[8:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "host:") {
				host = strings.TrimSpace(bpart[5:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "action:") {
				action = strings.TrimSpace(bpart[7:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "actionlevel:") {
				actionlevel = strings.TrimSpace(bpart[12:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "contact:") {
				contact = strings.TrimSpace(bpart[8:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "miscparams:") {
				miscparams = strings.TrimSpace(bpart[11:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "onremote:") {
				onremote = strings.TrimSpace(bpart[9:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "image:") {
				image = strings.TrimSpace(bpart[6:])
			}

			if strings.HasPrefix(strings.ToLower(bpart), "gitrepo:") {
				gitrepo = strings.TrimSpace(bpart[8:])
			}
		}

		if len(subj) == 0 || subj != trigger_subject {
			continue
		}
		
		if len(host) == 0 {
			continue
		}

		if len(action) == 0 {
			continue
		}

		if len(actionlevel) == 0 {
			continue
		}

		if len(image) == 0 {
			continue
		}

		actionmsg := ActionMessage{}
		actionmsg.Host = host
		actionmsg.Action = action
		actionmsg.ActionLevel = actionlevel
		actionmsg.Contact = contact
		actionmsg.MiscParams = miscparams
		actionmsg.OnRemote = onremote
		actionmsg.Image = image
		actionmsg.GitRepo = gitrepo
		HandleNewEmails(actionmsg)

		deleteids = append(deleteids, msgid)
	}

	for _, dids := range deleteids {
		conn.Write([]byte("DELE " + dids + "\n"))
		conn.Read(buff)
	}

	conn.Write([]byte("QUIT\n"))
	conn.Read(buff)
}

func VerifyLicense(key string, business string) bool {
	secretkey := business
	padding := "nullsoftllcvalkyrie"
//	statickey := "static_key"

	if len(secretkey) < 16 {
		mdiff := 16 - len(secretkey)
		newchars := padding[0:mdiff]
		secretkey += newchars
	}

	if len(secretkey) > 16 {
		secretkey = secretkey[0:16]
	}


	nkey := strings.Replace(key, "-", "", -1)
	if len(nkey) != 24 {
		return false
	}

	now := time.Now().Unix()

	group1 := nkey[0:8]
	group2 := nkey[8:16]
	group3 := nkey[16:24]

	ogeneratedon, err := strconv.ParseInt("0x" + group1, 0, 64)
	if err != nil {
		return false
	}

	omid, err := strconv.ParseInt("0x" + group2, 0, 64)
	if err != nil {
		return false
	}

	oexptime, err := strconv.ParseInt("0x" + group3, 0, 64)
	if err != nil {
		return false
	}

	generatedon := ogeneratedon ^ omid
	exptime := oexptime ^ ogeneratedon
	diff := exptime - generatedon

	rand.Seed(generatedon)
	mid := rand.Int63n(generatedon - 5878423) + 5878423

	if mid != omid {
		return false
	}

	if diff > 31556952 {
		return false
	}

	if now > exptime {
		return false
	}

//	hiddenkey, _ := cipherize.Decrypt([]byte(secretkey), statickey)

//	if strings.ToUpper(hiddenkey) != nkey {
//		return false
//	}

	return true
}

func main() {
	var err error
        prottimes = make(map[string]int64)

        loggerurl = os.Getenv("loggerurl")
        strusesecurelogging := os.Getenv("usesecurelogging")
        logkey = os.Getenv("logkey")

        strdupprotection := os.Getenv("dupprotection")
        strdupprottime := os.Getenv("dupprottime")

	whosts := os.Getenv("worker_hosts")
	mailhost = os.Getenv("mailhost")
	mailproto := os.Getenv("mailproto")
	mailuser = os.Getenv("mailuser")
	mailpass = os.Getenv("mailpass")

	strmailusetls := os.Getenv("mailusetls")
	imapfolder = os.Getenv("imapfolder")

	trigger_subject = os.Getenv("trigger_subject")
	strsleep_timeout := os.Getenv("sleep_timeout")

	license = os.Getenv("license")
	business = os.Getenv("business")


	if len(license) == 0 {
		Log("Missing License Key", 0)
		fmt.Println("Missing License Key")
		return
	}

	if len(business) == 0 {
		Log("Missing Business Name", 0)
		fmt.Println("Missing Business Name")
		return
	}

	if ! VerifyLicense(license, business) {
		Log("License Is Invalid", 0)
		fmt.Println("License Is Invalid")
		return
	}

	if len(mailhost) == 0 {
		fmt.Println("Missing 'mailhost'. Exiting.")
		return
	}

	if len(mailproto) == 0 {
		fmt.Println("Missing 'mailproto'. Exiting.")
		return
	}

	if len(mailuser) == 0 {
		fmt.Println("Missing 'mailuser'. Exiting.")
		return
	}

	if len(mailpass) == 0 {
		fmt.Println("Missing 'mailpass'. Exiting.")
		return
	}

	if len(trigger_subject) == 0 {
		fmt.Println("Missing 'trigger_subject'. Exiting.")
		return
	}

	if len(whosts) == 0 {
		fmt.Println("Missing 'worker_hosts'. Exiting.")
		return
	}

	if len(strmailusetls) == 0 {
		fmt.Println("Missing 'mailusetls.  Using default value of false")
		mailusetls = false
	} else {
		mailusetls, err = strconv.ParseBool(strmailusetls)
		if err != nil {
			fmt.Println("Missing 'mailusetls.  Using default value of false: " + err.Error())
			mailusetls = false
		}
	}

	if len(strsleep_timeout) == 0 {
		fmt.Println("Missing 'sleep_timeout'.  Using default value of 1 minute")
		sleep_timeout = 1
	} else {
		sleep_timeout, err = strconv.Atoi(strsleep_timeout)
		if err != nil {
			fmt.Println("Error Converting sleep_timeout to integer (using default value of 1 min): " + err.Error())
			sleep_timeout = 1
		}	
	}


	// Probably a better way than using env 
	// for worker hosts list
	if ! strings.Contains(whosts, "|") {
		workers = append(workers, whosts)
	} else {
		workers = strings.Split(whosts, "|")
	}

        bflag, err := strconv.ParseBool(strusesecurelogging)
        if err != nil {
                usesecurelogging = false
        } else {
                usesecurelogging = bflag
        }

        if strdupprottime == "" {
                dupprottime = 30
        }

        flag, err := strconv.ParseBool(strdupprotection)
        if err != nil {
                Log("failed to set dupprotection | " + err.Error(), 0)
                Log("using default value of: true", 0)
                dupprotection = true
        } else {
                dupprotection = flag
        }

        itime, err := strconv.Atoi(strdupprottime)
        if err != nil {
                Log("failed to set dupprottime | " + err.Error(), 0)
                Log("using default value of: 30 minutes", 0)
                dupprottime = 1800
        } else {
                dupprottime = int64(itime)
        }

	if strings.ToLower(mailproto) == "imap" {
		if len(imapfolder) == 0 {
			fmt.Println("imapfolder not set, using default value of: INBOX")
			imapfolder = "INBOX"
		}

		for {
			HandleIMAP()
			time.Sleep(time.Minute * time.Duration(sleep_timeout))
		}

	} else if strings.ToLower(mailproto) == "pop" || strings.ToLower(mailproto) == "pop3" {
		for {
			HandlePop()
			time.Sleep(time.Minute * time.Duration(sleep_timeout))
		}
	} else {
		fmt.Println("Unsupported Mail Protocol Specified.  Exiting.")
		return
	}
}
