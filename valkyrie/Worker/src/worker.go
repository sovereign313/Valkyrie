package main

import (
	"os"
	"fmt"
	"sshclient"
	"sftpclient"
	"bytes"
	"time"
	"strconv"
	"strings"
//	"cipherize"

	"math/rand"
	"os/exec"
	"net/smtp"
	"net/url"
	"encoding/json"
	"net/http"

	"gopkg.in/src-d/go-git.v4"
)

var loggerurl string
var logkey string

var gitrepourl string
var externalpath string

var usegitrepo bool
var usesecurelogging bool

var key []byte

var license string
var business string

type ActionMessage struct {
	Host string
	Action string
	ActionLevel string
	Contact string
	MiscParams string
	OnRemote bool
	GitRepo string
}

func SendSMTPMessage(mailserver string, from string, to string, subject string, body string) error {
	connection, err := smtp.Dial(mailserver)
	if err != nil {
		return err
	}
	defer connection.Close()

	connection.Mail(from)
	connection.Rcpt(to)

	wc, err := connection.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	body = "To: " + to + "\r\nFrom: " + from + "\r\nSubject: " + subject + "\r\n\r\n" + body

	buf := bytes.NewBufferString(body)
	_, err = buf.WriteTo(wc)
	if err != nil {
		return err
	}

	return nil
}

func Log(message string, level int) error {
	var Level string
	Service := "Valkyrie Worker"
	
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
	file, err := os.OpenFile("./worker.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("failed to open log file: " + err.Error())
		return err
	}
	defer file.Close()

        current_time := time.Now().Local()
        t := current_time.Format("Jan 02 2006 03:04:05")
	_, err = file.WriteString(t + " - Valkyrie Worker: " + message + "\n")

	if err != nil {
		fmt.Println("failed to write to log file: " + err.Error())
		return err
	}

	return nil
}


func HandleWork(actnmsg *ActionMessage) error {

	if actnmsg.ActionLevel == "internal" {

	} else if actnmsg.ActionLevel == "plugin" {

	} else if actnmsg.ActionLevel == "instruction" {

	} else if actnmsg.ActionLevel == "external" {
		if actnmsg.OnRemote == true {
			if usegitrepo {
				_, err := os.Stat(externalpath + "/live/")
				if os.IsNotExist(err) {
					GetRepo(actnmsg)
				} else {
					UpdateRepo(actnmsg)
				}
			}

			pth := externalpath + "/live/" + actnmsg.Action
			_, err := os.Stat(pth)
			if os.IsNotExist(err) {
				Log("Can't Find External Utility: " + pth, 0)
				return err
			}

			keyrsa, _ := sftpclient.GetKeyFile("id_rsa")
			output, err := sftpclient.CopyFile(pth, actnmsg.Host, "/tmp/" + actnmsg.Action, 5, keyrsa)
			if err != nil {
				Log("Couldn't Write Remote File: " + pth + " to: " + actnmsg.Host, 0)
				return err
			}
			
			output, err = sshclient.RunOneCommand(actnmsg.Host, "/tmp/" + actnmsg.Action + " " + actnmsg.MiscParams, 5, keyrsa)
			_, err = sshclient.RunOneCommand(actnmsg.Host, "rm -f /tmp/" + actnmsg.Action, 5, keyrsa)
			Log("Command Output (" + actnmsg.Host + ":" + actnmsg.Action + "):" + output, 0)
		} else {
			if usegitrepo {
				_, err := os.Stat(externalpath + "/live/")
				if os.IsNotExist(err) {
					GetRepo(actnmsg)
				} else {
					UpdateRepo(actnmsg)
				}
			}

			pth := externalpath + "/live/" + actnmsg.Action
			_, err := os.Stat(pth)
			if os.IsNotExist(err) {
				Log("Can't Find External Utility: " + pth, 0)
				return err
			}

			out, err := exec.Command(pth).Output()
			if err != nil {
				Log("Error Running External Utility (" + pth + "): " + err.Error(), 0)
				return err
			}

			Log("Successfully Executed (" + pth + ")", 0)
			Log("Output: " + string(out), 0)
		}
	} else {
		if actnmsg.OnRemote == true {
			if usegitrepo {
				_, err := os.Stat(externalpath + "/live/")
				if os.IsNotExist(err) {
					GetRepo(actnmsg)
				} else {
					UpdateRepo(actnmsg)
				}
			}

			pth := externalpath + "/live/" + actnmsg.Action
			_, err := os.Stat(pth)
			if os.IsNotExist(err) {
				Log("Can't Find External Utility: " + pth, 0)
				return err
			}

			keyrsa, _ := sftpclient.GetKeyFile("id_rsa")
			output, err := sftpclient.CopyFile(pth, actnmsg.Host, "/tmp/" + actnmsg.Action, 5, keyrsa)
			if err != nil {
				Log("Couldn't Write Remote File: " + pth + " to: " + actnmsg.Host, 0)
				return err
			}
			
			output, err = sshclient.RunOneCommand(actnmsg.Host, "/tmp/" + actnmsg.Action + " " + actnmsg.MiscParams, 5, keyrsa)
			_, err = sshclient.RunOneCommand(actnmsg.Host, "rm -f /tmp/" + actnmsg.Action, 5, keyrsa)
			Log("Command Output (" + actnmsg.Host + ":" + actnmsg.Action + "):" + output, 0)

		} else {
			if usegitrepo {
				_, err := os.Stat(externalpath + "/live/")
				if os.IsNotExist(err) {
					GetRepo(actnmsg)
				} else {
					UpdateRepo(actnmsg)
				}
			}

			pth := externalpath + "/live/" + actnmsg.Action
			_, err := os.Stat(pth)
			if os.IsNotExist(err) {
				Log("Can't Find External Utility: " + pth, 0)
				return err
			}

			out, err := exec.Command(pth).Output()
			if err != nil {
				Log("Error Running External Utility (" + pth + "): " + err.Error(), 0)
				return err
			}

			Log("Successfully Executed (" + pth + ")", 0)
			Log("Output: " + string(out), 0)
		}
	}

	return nil
}

func GetRepo(action *ActionMessage) {
	if action.GitRepo == "" {
		_, err := git.PlainClone(externalpath, false, &git.CloneOptions{
			URL: gitrepourl,
		})

		if err != nil {
			Log("failed to pull repo: " + gitrepourl, 0)
			Log("Error: " + err.Error(), 0)
		}
	} else {
		_, err := git.PlainClone(externalpath, false, &git.CloneOptions{
			URL: action.GitRepo,
		})

		if err != nil {
			Log("failed to pull repo: " + gitrepourl, 0)
			Log("Error: " + err.Error(), 0)
		}
	}
}

func UpdateRepo(action *ActionMessage) {
	if action.GitRepo == "" {
		r, err := git.PlainOpen(externalpath)
		if err != nil {
			Log("failed to update repo: " + gitrepourl, 0)
			Log("Error: " + err.Error(), 0)
		}

		w, err := r.Worktree()
		if err != nil {
			Log("failed to update repo: " + gitrepourl, 0)
			Log("Error: " + err.Error(), 0)
		}

		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			Log("failed to update repo: " + gitrepourl, 0)
			Log("Error: " + err.Error(), 0)
		}

		ref, _ := r.Head()
		r.CommitObject(ref.Hash())
	} else {
		r, err := git.PlainOpen(externalpath)
		if err != nil {
			Log("failed to update repo: " + action.GitRepo, 0)
			Log("Error: " + err.Error(), 0)
		}

		w, err := r.Worktree()
		if err != nil {
			Log("failed to update repo: " + action.GitRepo, 0)
			Log("Error: " + err.Error(), 0)
		}

		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			Log("failed to update repo: " + action.GitRepo, 0)
			Log("Error: " + err.Error(), 0)
		}

		ref, _ := r.Head()
		r.CommitObject(ref.Hash())
	}
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

	loggerurl = os.Getenv("loggerurl")
	strusesecurelogging := os.Getenv("usesecurelogging")
	logkey = os.Getenv("logkey")
	externalpath = os.Getenv("externalpath")
	strusegitrepo := os.Getenv("usegitrepo")
	gitrepourl = os.Getenv("gitrepourl")

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

	bflag, err := strconv.ParseBool(strusesecurelogging)
	if err != nil {
		usesecurelogging = false
	} else {
		usesecurelogging = bflag
	}


	if externalpath == "" {
		Log("external path not set: using default /valhalla", 0)
		externalpath = "/valhalla"
	}

	rflag, err := strconv.ParseBool(strusegitrepo)
	if err != nil {
		Log("failed to set usegitrepo: " + err.Error(), 0)
		Log("using default value of: false", 0)
		usegitrepo = false
	} else {
		usegitrepo = rflag
	}

	if usegitrepo == true && len(gitrepourl) == 0 {
		fmt.Println("can't use gitrepo, if gitrepourl is not defined")
		fmt.Println("Fatal.  Crash.  Hang.  Boom.")
		Log("can't use gitrepo, if gitrepourl is not defined", 0)
		Log("Fatal.  Crash.  Hang.  Boom.", 0)
		os.Exit(1)
	}

	action := ActionMessage{}
	action.Host = os.Getenv("HOST")
	action.Action = os.Getenv("ACTION")
	action.ActionLevel = os.Getenv("ACTIONLEVEL")
	action.MiscParams = os.Getenv("MISCPARAMS")
	action.GitRepo = os.Getenv("GITREPO")
	onremote := os.Getenv("ONREMOTE")

	if len(action.Host) == 0 {
		Log("failed to get HOST Environment Variable", 0)
		return
	}

	if len(action.Action) == 0 {
		Log("failed to get ACTION Environment Variable", 0)
		return
	}

	if len(action.ActionLevel) == 0 {
		Log("failed to get ACTIONLEVEL Environment Variable", 0)
	}

	if len(action.MiscParams) == 0 {
		Log("failed to get MISCPARAMS Environment Variable", 0)
	}

	if len(onremote) == 0 {
		Log("failed to get ONREMOTE Environment Variable", 0)
	}

	action.OnRemote, err = strconv.ParseBool(onremote)
	if err != nil {
		Log("failed to convert OnRemote to Bool | " + err.Error(), 0)
		Log("using default value of true", 0)
		action.OnRemote = true
	}

	err = HandleWork(&action)
	if err != nil {
		Log("Couldn't Find, Or Run utility. Doing Nothing", 0)
	} else {
		jsn, _ := json.Marshal(action)
		Log("Action has been taken", 0)
		Log(string(jsn), 0)
	}
}
