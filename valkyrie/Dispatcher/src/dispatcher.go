package main

import (
	"os"
	"fmt"
	"time"
	"errors"
	"strconv"
	"strings"
	"cipherize"

	"math/rand"
        "io/ioutil"
	"net/http"
	"net/url"
	"encoding/hex"
	"encoding/json"
	"crypto/md5"

	"github.com/jeffail/gabs"
	"github.com/gorilla/mux"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var sqsname string
var sqsregion string
var strkey string

var echoresponse string

var useaws bool
var dupprotection bool

var key []byte
var prottimes map[string]int64

var dupprottime int64

var usesecurelogging bool
var loggerurl string
var logkey string

var license string
var business string

type ActionMessage struct {
	Host string
	Action string
	ActionLevel string
	Contact string
	MiscParams string
	RmPriority string
	OnRemote string
	Image string
	GitRepo string
}

func handleSetUseAWS(w http.ResponseWriter, r *http.Request) {
	tmpuseaws := r.URL.Query().Get("useaws")
	if len(tmpuseaws) == 0 {
		fmt.Fprintf(w, "failed to set useaws")
		Log("failed to set useaws: left blank", 0)
		return
	}

	flag, err := strconv.ParseBool(tmpuseaws)
	if err != nil {
		fmt.Fprintf(w, "failed to set useaws: value is not true or false")
		Log("failed to set useaws: " + err.Error(), 0)
		return
	}

	useaws = flag 

	if flag {
		fmt.Fprintf(w, "now using aws")
	} else {
		fmt.Fprintf(w, "aws not in use any more")
	}
}

func handleSetDupProtection(w http.ResponseWriter, r *http.Request) {
	tmpdupprotection := r.URL.Query().Get("dupprotection")
	if len(tmpdupprotection) == 0 {
		fmt.Fprintf(w, "failed to set duplication protection")
		Log("failed to set dupprotection: left blank", 0)
		return
	}

	flag, err := strconv.ParseBool(tmpdupprotection)
	if err != nil {
		fmt.Fprintf(w, "failed to set duplication protection: value is not true or false")
		Log("failed to set dupprotection: " + err.Error(), 0)
		return
	}

	dupprotection = flag

	if flag {
		fmt.Fprintf(w, "duplication protection is now ON")
	} else {
		fmt.Fprintf(w, "duplication protection is now OFF")
	}
}


func handleSetDupProtTime(w http.ResponseWriter, r *http.Request) {
	tmpdupprottime := r.URL.Query().Get("dupprottime")
	if len(tmpdupprottime) == 0 {
		fmt.Fprintf(w, "failed to set duplication protection time")
		Log("failed to set dupprottime: left blank", 0)
		return
	}

	itime, err := strconv.Atoi(tmpdupprottime)
	if err != nil {
		fmt.Fprintf(w, "failed to set duplication protection time: value is not an integer number")
		Log("failed to set dupprottime: " + err.Error(), 0)
		return
	}

	dupprottime = int64(itime)
}

func handleSetSQSName(w http.ResponseWriter, r *http.Request) {
        tmpsqsname := r.URL.Query().Get("sqsname")
        if len(tmpsqsname) == 0 {
		fmt.Fprintf(w, "failed to set queue name")
		Log("failed to set queue name: left blank", 0)
		return
	}

	sqsname = tmpsqsname
	fmt.Fprintf(w, "queue name set to: " + sqsname)
}

func handleSetSQSRegion(w http.ResponseWriter, r *http.Request) {
        tmpsqsregion := r.URL.Query().Get("sqsregion")
	if len(tmpsqsregion) == 0 {
		fmt.Fprintf(w, "failed to set region")
		Log("failed to set queue region: left blank", 0)
		return
	}

	sqsregion = tmpsqsregion
	fmt.Fprintf(w, "region set to: " + sqsregion)
}

func handleSetSTRKey(w http.ResponseWriter, r *http.Request) {
	tmpstrkey := r.URL.Query().Get("strkey")
	if len(tmpstrkey) == 0 {
		fmt.Fprintf(w, "failed to encryption key")
		Log("failed to set encryption key: left blank", 0)
		return
	}

	strkey = tmpstrkey
	key = []byte(strkey)
	fmt.Fprintf(w, "encrypt key set")
}


func handleSetAWSAccessKey(w http.ResponseWriter, r *http.Request) {
	tmpawsakey := r.URL.Query().Get("awsaccesskey")
	if len(tmpawsakey) == 0 {
		fmt.Fprintf(w, "failed to set aws access key id")
		Log("failed to set aws access key id: left blank", 0)
		return
	}

	os.Setenv("AWS_ACCESS_KEY_ID", tmpawsakey)
	fmt.Fprintf(w, "set aws access id key")
}

func handleSetAWSSecretKey(w http.ResponseWriter, r *http.Request) {
	tmpawsskey := r.URL.Query().Get("awssecretkey")
	if len(tmpawsskey) == 0 {
		fmt.Fprintf(w, "failed to set aws secret key")
		Log("failed to set aws secret key: left blank", 0)
		return
	}

	os.Setenv("AWS_SECRET_ACCESS_KEY", tmpawsskey)
	fmt.Fprintf(w, "set aws secret key")
}

func handleWhoAreYou(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "dispatcher")
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, echoresponse)
}

func handleTrigger(w http.ResponseWriter, r *http.Request) {
	var err error

	Host := r.URL.Query().Get("host")
	Action := r.URL.Query().Get("action")
	ActionLevel := r.URL.Query().Get("actionlevel")
	Contact := r.URL.Query().Get("contact")
	MiscParams := r.URL.Query().Get("miscparams")
	RmPriority := r.URL.Query().Get("rmpriority")
	usesns := r.URL.Query().Get("usesns")
	onremote := r.URL.Query().Get("onremote")
	Image := r.URL.Query().Get("image")
	GitRepo := r.URL.Query().Get("gitrepo")

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

	if Contact == "" {
		Log("failed to trigger: missing contact data", 0)
		fmt.Fprintf(w, "failed to trigger: missing contact data")
		return
	}

	if Image == "" {
		Image = "worker"
	}

	queue, _ := GetSQSUrl(sqsname)

	if queue == "" {
		queue, err = CreateSQS(sqsname)
		if err != nil {
			fmt.Fprintf(w, "failed to handle event")
			Log("failed to create Queue on AWS: " + err.Error(), 0)
			return
		}
	}

	_, err = strconv.ParseBool(usesns)
	if err != nil {
		Log("failed to convert usesns to boolean: " + err.Error(), 0)
		Log("using default value of true.", 0)
		usesns = "true"
	}

	action := ActionMessage{}
	action.Host = Host
	action.Action = Action
	action.ActionLevel = ActionLevel
	action.Contact = Contact
	action.MiscParams = MiscParams
	action.RmPriority = RmPriority
	action.OnRemote = onremote
	action.Image = Image
	action.GitRepo = GitRepo

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
	
	messageid, err := SendSQSMessage(queue, &action)
	if err != nil {
		fmt.Fprintf(w, "failed to insert message into sqs: " + err.Error())
		return
	}

	fmt.Fprintf(w, "successfully added message to sqs: " + messageid)

	jsn, _ := json.Marshal(action)
	Log("successfully added message to sqs: " + messageid, 0)
	Log(string(jsn), 0)

}

func handleJSONBodyTrigger(w http.ResponseWriter, r *http.Request) {
	bsplunkmsg, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		fmt.Fprintf(w, err.Error(), 500)
		Log("couldn't read JSON Post From Body: " + err.Error(), 0)
		return
	}

	echoresponse = string(bsplunkmsg)
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

	dstcontact, ok := jsonParsed.Path("result.dstcontact").Data().(string)
	if ! ok {
		Log("dstcontact missing", 0)
	}

	dstusesns, ok := jsonParsed.Path("result.dstusesns").Data().(string)
	dstonremote, ok := jsonParsed.Path("result.dstonremote").Data().(string)
	dstmiscparams, ok := jsonParsed.Path("result.dstmiscparams").Data().(string)
	dstrmpriority, ok := jsonParsed.Path("result.dstrmpriority").Data().(string)
	dstimage, ok := jsonParsed.Path("result.dstimage").Data().(string)
	dstgitrepo, ok := jsonParsed.Path("result.dstgitrepo").Data().(string)

	_, err = strconv.ParseBool(dstusesns)
	if err != nil {
		Log("failed to convert usesns to boolean: " + err.Error(), 0)
		Log("using default value of true.", 0)
		dstusesns = "true"
	}

	_, err = strconv.ParseBool(dstonremote)
	if err != nil {
		Log("failed to convert onremote to boolean: " + err.Error(), 0)
		Log("using default value of true.", 0)
		dstonremote = "true"
	}

	queue, _ := GetSQSUrl(sqsname)

	if queue == "" {
		queue, err = CreateSQS(sqsname)
		if err != nil {
			fmt.Fprintf(w, "failed to handle event")
			Log("failed to create Queue on AWS: " + err.Error(), 0)
			return
		}
	}

	action := ActionMessage{}
	action.Host = dsthost 
	action.Action = dstfn
	action.ActionLevel = dstfnlevel
	action.Contact = dstcontact 
	action.MiscParams = dstmiscparams
	action.RmPriority = dstrmpriority
	action.OnRemote = dstonremote
	action.Image = dstimage
	action.GitRepo = dstgitrepo

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

	messageid, err := SendSQSMessage(queue, &action)
	if err != nil {
		fmt.Fprintf(w, "failed to insert message into sqs: " + err.Error())
		return
	}

	Log("inserted message into SQS: " + messageid, 0)
	fmt.Fprintf(w, "successfully added message to sqs: " + messageid)
}

func SendSQSMessage(sqsname string, actionmsg *ActionMessage) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(sqsregion)},
	)

	// Create a SQS service client.
	svc := sqs.New(sess)

	if len(key) == 0 {
		return "", errors.New("Encryption key is not set")
	}


	enchst, _ := cipherize.Encrypt(key, actionmsg.Host)
	encaction, _ := cipherize.Encrypt(key, actionmsg.Action)
	encactionlevel, _ := cipherize.Encrypt(key, actionmsg.ActionLevel)
	enccontact, _ := cipherize.Encrypt(key, actionmsg.Contact)
	encdstparams, _ := cipherize.Encrypt(key, actionmsg.MiscParams)
	encrmpriority, _ := cipherize.Encrypt(key, actionmsg.RmPriority)
	enconremote, _ := cipherize.Encrypt(key, actionmsg.OnRemote)
	encimage, _ := cipherize.Encrypt(key, actionmsg.Image)
	encgitrepo, _ := cipherize.Encrypt(key, actionmsg.GitRepo)

	actionmsg.Host = enchst
	actionmsg.Action = encaction
	actionmsg.ActionLevel = encactionlevel
	actionmsg.Contact = enccontact
	actionmsg.MiscParams = encdstparams
	actionmsg.RmPriority = encrmpriority
	actionmsg.OnRemote = enconremote
	actionmsg.Image = encimage
	actionmsg.GitRepo = encgitrepo

	jsn, _ := json.Marshal(actionmsg)

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue {
			"Host": &sqs.MessageAttributeValue {
				DataType: aws.String("String"),
				StringValue: aws.String(actionmsg.Host),
			},
			"Action": &sqs.MessageAttributeValue {
				DataType: aws.String("String"),
				StringValue: aws.String(actionmsg.Action),
			},
			"ActionLevel": &sqs.MessageAttributeValue {
				DataType: aws.String("String"),
				StringValue: aws.String(actionmsg.ActionLevel),
			},
			"Contact": &sqs.MessageAttributeValue {
				DataType: aws.String("String"),
				StringValue: aws.String(actionmsg.Contact),
			},
			"MiscParams": &sqs.MessageAttributeValue {
				DataType: aws.String("String"),
				StringValue: aws.String(actionmsg.MiscParams),
			},
			"RmPriority": &sqs.MessageAttributeValue {
				DataType: aws.String("String"),
				StringValue: aws.String(actionmsg.RmPriority),
			},
			"OnRemote": &sqs.MessageAttributeValue {
				DataType: aws.String("String"),
				StringValue: aws.String(actionmsg.OnRemote),
			},
			"Image": &sqs.MessageAttributeValue {
				DataType: aws.String("String"),
				StringValue: aws.String(actionmsg.Image),
			},
			"GitRepo": &sqs.MessageAttributeValue {
				DataType: aws.String("String"),
				StringValue: aws.String(actionmsg.GitRepo),
			},
		},

		MessageBody: aws.String(string(jsn)),
		QueueUrl: &sqsname,
	})

	if err != nil {
		Log("failed to send sqs message: " + err.Error(), 0)
		return "", err
	}

	return *result.MessageId, nil

}

func GetSQSUrl(sqsname string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(sqsregion)},
	)

	// Create a SQS service client.
	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(sqsname),
	})

	if err != nil {
		Log("failed to get queue on aws, probably doesn't exist: " + err.Error(), 1)
		return "", err
	}

	return *result.QueueUrl, nil
}


func CreateSQS(sqsname string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(sqsregion)},
	)

	// Create a SQS service client.
	svc := sqs.New(sess)

	result, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(sqsname),
		Attributes: map[string]*string{
			"DelaySeconds":           aws.String("10"),
			"MessageRetentionPeriod": aws.String("86400"),
		},
	})

	if err != nil {
		Log("failed to create queue on aws... check creds file? : " + err.Error(), 0)
		return "", err
	}	

	return *result.QueueUrl, nil
}

func Log(message string, level int) error {
	var Level string
	Service := "Valkyrie Dispatcher"
	
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
		_, err := http.Get(loggerurl + "/log?service=" + encservice + "&loglevel=" + enclevel + "&message=" + encmessage + "&logkey=" + enckey)
		if err != nil {
			FileLog("Failed to log to logging server: " + err.Error())
			return err
		}
	} else {
		_, err := http.Get(loggerurl + "/log?service=" + encservice + "&loglevel=" + enclevel + "&message=" + encmessage)
		if err != nil {
			FileLog("Failed to log to logging server: " + err.Error())
			return err
		}
	}

	return nil
}

func FileLog(message string) error {
	file, err := os.OpenFile("./dispatcher.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("failed to open log file: " + err.Error())
		return err
	}
	defer file.Close()

        current_time := time.Now().Local()
        t := current_time.Format("Jan 02 2006 03:04:05")
	_, err = file.WriteString(t + " - Valkyrie Dispatcher: " + message + "\n")

	if err != nil {
		fmt.Println("failed to write to log file: " + err.Error())
		return err
	}

	return nil
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
	prottimes = make(map[string]int64)

	echoresponse = ""
        sqsname = os.Getenv("sqsname")
        sqsregion = os.Getenv("sqsregion")
	strkey = os.Getenv("key")

	strdupprotection := os.Getenv("dupprotection")
	strdupprottime := os.Getenv("dupprottime")

        loggerurl = os.Getenv("loggerurl")
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

        if sqsname == "" {
		sqsname = "dispatcher"
        }

        if sqsregion == "" {
		sqsregion = "us-east-2"
        }

	if strdupprottime == "" {
		dupprottime = 30
	}

	if strkey == "" {
		// Change This Thing
		strkey = "LKHlhb899Y09olUi"
	}
	// To Here

        flag, err := strconv.ParseBool(strdupprotection)
        if err != nil {
                Log("failed to set dupprotection: " + err.Error(), 0)
                Log("using default value of: true", 0)
		dupprotection = true
        } else {
		dupprotection = flag
	}

	itime, err := strconv.Atoi(strdupprottime)
	if err != nil {
		Log("failed to set dupprottime: " + err.Error(), 0)
		Log("using default value of: 30 minutes", 0)
		dupprottime = 30
	} else {
		dupprottime = int64(itime)
	}

	key = []byte(strkey)

        router := mux.NewRouter()
        router.HandleFunc("/whoareyou", handleWhoAreYou)
        router.HandleFunc("/ping", handlePing)
	router.HandleFunc("/echo", handleEcho)
	router.HandleFunc("/setsqsname", handleSetSQSName)
	router.HandleFunc("/setsqsregion", handleSetSQSRegion)
	router.HandleFunc("/setstrkey", handleSetSTRKey)
	router.HandleFunc("/setawsaccesskey", handleSetAWSAccessKey)
	router.HandleFunc("/setawssecretkey", handleSetAWSSecretKey)
	router.HandleFunc("/setuseaws", handleSetUseAWS)
	router.HandleFunc("/setdupprotection", handleSetDupProtection)
	router.HandleFunc("/setdupprottime", handleSetDupProtTime)
        router.HandleFunc("/trigger", handleTrigger)
        router.HandleFunc("/jsonbodytrigger", handleJSONBodyTrigger)

	Log("Dispatcher Started", 0)
	fmt.Println("Dispatcher Started")

        err = http.ListenAndServe(":8090", router)
        if err != nil {
                fmt.Println("ListenAndServe: ", err)
        }
}
