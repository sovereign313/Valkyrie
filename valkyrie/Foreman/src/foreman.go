package main

import (
	"os"
	"io"
	"fmt"
	"time"
	"strconv"
	"strings"
//	"cipherize"

	"io/ioutil"
	"net/http"
	"net/url"
        "encoding/hex"
        "crypto/md5"
	"math/rand"

	"github.com/jeffail/gabs"
	"github.com/gorilla/mux"
)

// I'm a monster, Quetzalcoatl

var dupprotection bool
var dupprottime int64
var prottimes map[string]int64

var usesecurelogging bool
var loggerurl string
var logkey string

var workers []string
var current_worker int

var license string
var business string

type ActionMessage struct {
	Host string
	Action string
	ActionLevel string
	Contact string
	MiscParams string
	OnRemote string
	GitRepo string
	Image string
}

func Log(message string, level int) error {
	var Level string
	Service := "Valkyrie Foreman"
	
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
	file, err := os.OpenFile("./foreman.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("failed to open log file: " + err.Error())
		return err
	}
	defer file.Close()

        current_time := time.Now().Local()
        t := current_time.Format("Jan 02 2006 03:04:05")
	_, err = file.WriteString(t + " - Valkyrie Foreman: " + message + "\n")

	if err != nil {
		fmt.Println("failed to write to log file: " + err.Error())
		return err
	}

	return nil
}

func handleWhoAreYou(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Valkyrie Foreman")
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func handleDescription(w http.ResponseWriter, r *http.Request) {
	html := "Valkyrie Foreman - Tool for triggering workers directly (without a queue)\n"
	fmt.Fprintf(w, html)
}

func handleTrigger(w http.ResponseWriter, r *http.Request) {

	Host := r.URL.Query().Get("host")
	Action := r.URL.Query().Get("action")
	ActionLevel := r.URL.Query().Get("actionlevel")
	Contact := r.URL.Query().Get("contact")
	MiscParams := r.URL.Query().Get("miscparams")
	onremote := r.URL.Query().Get("onremote")
	GitRepo := r.URL.Query().Get("gitrepo")
	Image := r.URL.Query().Get("image")

	if len(Host) == 0 {
		Log("failed to trigger: missing host data", 0)
		fmt.Fprintf(w, "failed to trigger: missing host data")
		return
	}

	if len(Action) == 0 {
		Log("failed to trigger: missing action data", 0)
		fmt.Fprintf(w, "failed to trigger: missing action data")
		return
	}

	if len(ActionLevel) == 0 {
		Log("failed to trigger: missing action level", 0)
		fmt.Fprintf(w, "failed to trigger: missing action level")
		return
	}

	if len(Contact) == 0 {
		Log("Missing contact. Not sending a response.", 0)
	}

	if len(Image) == 0 {
		Log("missing image parameter. Using default value of: worker", 0)
		Image = "worker"	
	}

	action := ActionMessage{}
	action.Host = Host
	action.Action = Action
	action.ActionLevel = ActionLevel
	action.Contact = Contact
	action.MiscParams = MiscParams
	action.OnRemote = onremote
	action.GitRepo = GitRepo
	action.Image = Image

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
                                fmt.Fprintf(w, "Duplicate Message Being Dropped: " + action.Host + " - " + action.Action)
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
	encgitrepo := url.QueryEscape(action.GitRepo)
	encimage := url.QueryEscape(action.Image)

	resp, err := http.Get(workers[current_worker] + "/trigger?host=" + enchost + "&action=" + encaction + "&actionlevel=" + encactionlevel + "&contact=" + enccontact + "&onremote=" + enconremote + "&gitrepo=" + encgitrepo + "&image=" + encimage + "&miscparams=" + encmiscparams)
	if err != nil {
		fmt.Fprintf(w, "failed to launch: " + err.Error())
		Log("Failed To Launch New Valkyrie: " + err.Error(), 0)
		return
	}
	defer resp.Body.Close()

	if current_worker == (len(workers) -1) {
		current_worker = 0
	} else {
		current_worker++
	}

	fmt.Fprintf(w, "Launched New Valkyrie")
}

func handleJSONBodyTrigger(w http.ResponseWriter, r *http.Request) {
	var dstimage string
	var dstcontact string

	bsplunkmsg, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		fmt.Fprintf(w, err.Error(), 500)
		Log("couldn't read JSON Post From Body: " + err.Error(), 0)
		return
	}

	jsonParsed, err := gabs.ParseJSON(bsplunkmsg)

	dsthost, ok := jsonParsed.Path("result.dsthost").Data().(string)
	if ! ok {
		fmt.Fprintf(w, "failed to parse result from POST json: missing dsthost")
		Log("failed to parse result from POST json: missing dsthost", 0)
		return
	}

	dstfn, ok := jsonParsed.Path("result.dstfn").Data().(string)
	if ! ok {
		fmt.Fprintf(w, "failed to parse result from POST json: missing dstfn")
		Log("failed to parse result from POST json: missing dstfn", 0)
		return
	}

	dstfnlevel, ok := jsonParsed.Path("result.dstfnlevel").Data().(string)
	if ! ok {
		fmt.Fprintf(w, "failed to parse result from POST json: missing dstfnlevel")
		Log("failed to parse result from POST json: missing dstfnlevel", 0)
		return
	}

	dstcontact, ok = jsonParsed.Path("result.dstcontact").Data().(string)
	if ! ok {
		Log("failed to parse result from POST json: missing dstcontact", 0)
		dstcontact = ""
	}

	dstimage, ok = jsonParsed.Path("result.dstimage").Data().(string)
	if ! ok {
		Log("dstimage not passed in. Using default value of: worker", 0)
		dstimage = "worker"
	}

	dstgitrepo, ok := jsonParsed.Path("result.dstgitrepo").Data().(string)
	dstonremote, ok := jsonParsed.Path("result.dstonremote").Data().(string)
	dstmiscparams, ok := jsonParsed.Path("result.dstmiscparams").Data().(string)

	_, err = strconv.ParseBool(dstonremote)
	if err != nil {
		Log("failed to convert onremote to boolean: " + err.Error(), 0)
		Log("using default value of true.", 0)
		dstonremote = "true"
	}

	action := ActionMessage{}
	action.Host = dsthost 
	action.Action = dstfn
	action.ActionLevel = dstfnlevel
	action.Contact = dstcontact 
	action.MiscParams = dstmiscparams
	action.OnRemote = dstonremote
	action.GitRepo = dstgitrepo
	action.Image = dstimage

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
				fmt.Fprintf(w, "Duplicate Message Being Dropped: " + action.Host + " - " + action.Action)
				Log("Duplicate Message Being Dropped: " + action.Host + " - " + action.Action, 0)
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
	encgitrepo := url.QueryEscape(action.GitRepo)
	encimage := url.QueryEscape(action.Image)

	resp, err := http.Get(workers[current_worker] + "/trigger?host=" + enchost + "&action=" + encaction + "&actionlevel=" + encactionlevel + "&contact=" + enccontact + "&onremote=" + enconremote + "&gitrepo=" + encgitrepo + "&image=" + encimage + "&miscparams=" + encmiscparams)
	if err != nil {
		fmt.Fprintf(w, "failed to launch: " + err.Error())
		Log("Failed To Launch New Valkyrie: " + err.Error(), 0)
		return
	}
	defer resp.Body.Close()

	if current_worker == (len(workers) -1) {
		current_worker = 0
	} else {
		current_worker++
	}

	fmt.Fprintf(w, "Launched New Valkyrie")
}

func handleHelp(w http.ResponseWriter, r *http.Request) {
	html := "/                - This Help\n"
	html += "/ping            - Returns Pong (Ensures Service Is Working)\n"
	html += "/whoareyou       - Returns The Application (Valkyrie Foreman)\n"
	html += "/description     - Returns A Description Of This Service\n"
	html += "/trigger         - Executes A Worker Task \n"
	html += "/jsonbodytrigger - Executes A Worker Task With Params As POST In Body\n"

	fmt.Fprintf(w, html)
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

func HashMD5File(filePath string) (string, error) {
	var returnMD5String string

	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

func main() {
	var err error
        prottimes = make(map[string]int64)

	license = os.Getenv("license")
	business = os.Getenv("business")
        strdupprotection := os.Getenv("dupprotection")
        strdupprottime := os.Getenv("dupprottime")

        loggerurl = os.Getenv("loggerurl")
        strusesecurelogging := os.Getenv("usesecurelogging")
        logkey = os.Getenv("logkey")
	whosts := os.Getenv("worker_hosts")


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

	router := mux.NewRouter()
        router.HandleFunc("/whoareyou", handleWhoAreYou)
        router.HandleFunc("/ping", handlePing)
	router.HandleFunc("/description", handleDescription)
	router.HandleFunc("/trigger", handleTrigger)
        router.HandleFunc("/jsonbodytrigger", handleJSONBodyTrigger)
        router.HandleFunc("/", handleHelp)

        err = http.ListenAndServe(":8091", router)
        if err != nil {
                fmt.Println("ListenAndServe: ", err)
	}
}
