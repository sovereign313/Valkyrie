package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"strings"
//	"cipherize"

	"math/rand"
	//"os/exec"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"golang.org/x/net/context"
)

var usesecurelogging bool
var loggerurl string
var logkey string

var license string
var business string
var metricurl string

type ActionMessage struct {
	Host string
	Action string
	ActionLevel string
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

	if len(loggerurl) == 0 {
		FileLog(message)
		return nil
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
	file, err := os.OpenFile("./launcher.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("failed to open log file: " + err.Error())
		return err
	}
	defer file.Close()

        current_time := time.Now().Local()
        t := current_time.Format("Jan 02 2006 03:04:05")
	_, err = file.WriteString(t + " - Valkyrie Launcher: " + message + "\n")

	if err != nil {
		fmt.Println("failed to write to log file: " + err.Error())
		return err
	}

	return nil
}

func SendMetric(message ActionMessage) error {
	enchost := url.QueryEscape(message.Host)
	encaction := url.QueryEscape(message.Action)
	encactionlevel := url.QueryEscape(message.ActionLevel)
	enccontact := url.QueryEscape(message.Contact)
	encmiscparams := url.QueryEscape(message.MiscParams)
	enconremote := url.QueryEscape(message.OnRemote)
	encimage := url.QueryEscape(message.Image)
	encgitrepo := url.QueryEscape(message.GitRepo)

	_, err := http.Get(metricurl + "record?host=" + enchost + "&action=" + encaction + "&actionlevel=" + encactionlevel + "&contact=" + enccontact + "&miscparams=" + encmiscparams + "&onremote=" + enconremote + "&image=" + encimage + "&gitrepo=" + encgitrepo)
	if err != nil {
		FileLog("Failed to log to logging server: " + err.Error())
		return err
	}
}

func CheckMetrics(message ActionMessage) error {

	_, err := http.Get(metricurl + "record?host=" + enchost + "&action=" + encaction + "&actionlevel=" + encactionlevel + "&contact=" + enccontact + "&miscparams=" + encmiscparams + "&onremote=" + enconremote + "&image=" + encimage + "&gitrepo=" + encgitrepo)
	if err != nil {
		FileLog("Failed to log to logging server: " + err.Error())
		return err
	}
}

func handleWhoAreYou(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Valkyrie Launcher")
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func handleDescription(w http.ResponseWriter, r *http.Request) {
	html := "Valkyrie Launcher - Tool for triggering workers directly (without a queue)\n"
	fmt.Fprintf(w, html)
}

func handleTrigger(w http.ResponseWriter, r *http.Request) {

	Host := r.URL.Query().Get("host")
	Action := r.URL.Query().Get("action")
	ActionLevel := r.URL.Query().Get("actionlevel")
	MiscParams := r.URL.Query().Get("miscparams")
	onremote := r.URL.Query().Get("onremote")
	GitRepo := r.URL.Query().Get("gitrepo")
	image := r.URL.Query().Get("image")

	if Host == "" {
		Log("failed to trigger: missing host data", 0)
		fmt.Fprintf(w, "failed to trigger: missing host data")
		return
	}

	if Action == "" {
		Log("failed to trigger: missing action data", 0)
		fmt.Fprintf(w, "failed to trigger: missing action data")
		return
	}

	if len(ActionLevel) == 0 {
		Log("failed to trigger: missing action level", 0)
		fmt.Fprintf(w, "failed to trigger: missing action level")
		return
	}

	if len(image) == 0 {
		image = "worker"
	}

	action := ActionMessage{}
	action.Host = Host
	action.Action = Action
	action.ActionLevel = ActionLevel
	action.MiscParams = MiscParams
	action.GitRepo = GitRepo
	action.Image = image
	action.OnRemote = onremote

	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, nil)
	if err != nil {
		Log("Failed To Make New Client | " + err.Error(), 0)
		fmt.Fprintf(w, "Failed To Make New Client: " + err.Error())
		return
	}

	ctx := context.Background()

	containerConfig := &container.Config{
		Image: image,
		Env: []string{"HOST=" + action.Host, "ACTION=" + action.Action, "ACTIONLEVEL=" + action.ActionLevel, "MISCPARAMS=" + action.MiscParams, "ONREMOTE=" + action.OnRemote, "GITREPO=" + action.GitRepo}, 
	}	

	resp, err := cli.ContainerCreate(ctx, containerConfig, nil, nil, "")
	if err != nil {
		Log("Failed To Create Container | " + err.Error(), 0)
		fmt.Fprintf(w, "Failed To Create Container: " + err.Error())
		return
	}

	/*
		Add Option Above To Remove Container When It's Complete
		Basically --rm To Docker Run
	*/

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		Log("Failed To Start Container | " + err.Error(), 0)
		fmt.Fprintf(w, "Failed to Start Container: " + err.Error())
		return
	}

	fmt.Fprintf(w, "Launched New Valkyrie")
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

func handleHelp(w http.ResponseWriter, r *http.Request) {
	html := "/            - This Help\n"
	html += "/ping        - Returns Pong (Ensures Service Is Working)\n"
	html += "/whoareyou   - Returns The Application (Valkyrie Launcher)\n"
	html += "/description - Returns A Description Of This Service\n"
	html += "/launch      - Executes A Worker Task \n"
	fmt.Fprintf(w, html)
}

func main() {
	var err error

        loggerurl = os.Getenv("loggerurl")
	metricurl = os.Getenv("metricurl")
        strusesecurelogging := os.Getenv("usesecurelogging")
        logkey = os.Getenv("logkey")
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

	if len(metricurl) == 0 {
		Log("Missing Metric Url", 0)
		fmt.Println("Missing Metric Url")
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

	router := mux.NewRouter()
        router.HandleFunc("/whoareyou", handleWhoAreYou)
        router.HandleFunc("/ping", handlePing)
	router.HandleFunc("/description", handleDescription)
	router.HandleFunc("/trigger", handleTrigger)
        router.HandleFunc("/", handleHelp)

        err = http.ListenAndServe(":8094", router)
        if err != nil {
                fmt.Println("ListenAndServe: ", err)
	}
}
