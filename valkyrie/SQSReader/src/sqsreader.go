package main

import (
	"os"
	"fmt"
	"cipherize"
	"bytes"
	"time"
	"strings"
	"strconv"

	"math/rand"
	"net/smtp"
	"net/url"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var sqsname string
var sqsregion string
var strkey string

var loggerurl string
var logkey string

var sleeptimeout int

var usesecurelogging bool

var key []byte

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
	RmPriority string
	UseSNS bool
	OnRemote string
	Image string
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

func Log(message string, level int) error {
	var Level string
	Service := "Valkyrie SQSReader"
	
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
	file, err := os.OpenFile("./sqsreader.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("failed to open log file: " + err.Error())
		return err
	}
	defer file.Close()

        current_time := time.Now().Local()
        t := current_time.Format("Jan 02 2006 03:04:05")
	_, err = file.WriteString(t + " - Valkyrie SQSReader: " + message + "\n")

	if err != nil {
		fmt.Println("failed to write to log file: " + err.Error())
		return err
	}

	return nil
}


func HandleWork(actnmsg *ActionMessage) error {

        enchost := url.QueryEscape(actnmsg.Host)
        encaction := url.QueryEscape(actnmsg.Action)
        encactionlevel := url.QueryEscape(actnmsg.ActionLevel)
        enccontact := url.QueryEscape(actnmsg.Contact)
        encmiscparams := url.QueryEscape(actnmsg.MiscParams)
        enconremote := url.QueryEscape(actnmsg.OnRemote)
        encimage := url.QueryEscape(actnmsg.Image)
	encgitrepo := url.QueryEscape(actnmsg.GitRepo)

        resp, err := http.Get(workers[current_worker] + "/trigger?host=" + enchost + "&action=" + encaction + "&actionlevel=" + encactionlevel + "&contact=" + enccontact + "&onremote=" + enconremote + "&image=" + encimage + "&gitrepo=" + encgitrepo + "&miscparams=" + encmiscparams)
        if err != nil {
                Log("Failed To Launch New Valkyrie: " + err.Error(), 0)
                return err
        }
        defer resp.Body.Close()

        if current_worker == (len(workers) -1) {
                current_worker = 0
        } else {
                current_worker++
        }

	Log("Launched New Valkyrie", 0)
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
	var err error

	loggerurl = os.Getenv("loggerurl")
	strusesecurelogging := os.Getenv("usesecurelogging")
	logkey = os.Getenv("logkey")

        sqsname = os.Getenv("sqsname")
        sqsregion = os.Getenv("sqsregion")
	strkey = os.Getenv("key")

	slptimeout := os.Getenv("sleeptimeout")
        whosts := os.Getenv("worker_hosts")

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


	if sqsname == "" {
		Log("sqsname not set: using default 'valkyrie'", 0)
		sqsname = "valkyrie"
	}

	if sqsregion == "" {
		Log("sqsregion not set: using default 'us-east-2'", 0)
		sqsregion = "us-east-2"
	}

	if slptimeout == "" {
		Log("sleep timeout not set: using default 10 seconds", 0)
		sleeptimeout = 10
	} else {
		sleeptimeout, err = strconv.Atoi(slptimeout)
		if err != nil {
			Log("Invalid value used as sleep time out for SQS Poll: " + err.Error(), 0)
			Log("---Using 10s Instead---", 0)
			sleeptimeout = 10
		}
	}

	if strkey == "" {
		Log("strkey not set: using insecure encrypt key 'LKHlhb899Y09olUi'", 0)
		strkey = "LKHlhb899Y09olUi"
	}

	key = []byte(strkey)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(sqsregion)},
	)
	if err != nil {
		Log("Failed to create new AWS Session | " + err.Error(), 0)
		return
	}

	svc := sqs.New(sess)
	queue, _ := GetSQSUrl(sqsname)
	
	for {

		q := string(queue)

		time.Sleep(time.Second * time.Duration(sleeptimeout))

		result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput {
			AttributeNames: aws.StringSlice([]string {
				"Host",
				"Action",
				"ActionLevel",
				"Contact",
				"MiscParams",
				"UseSNS",
				"OnRemote",
			}), 
		
			QueueUrl: &q,
			MaxNumberOfMessages: aws.Int64(1),
			VisibilityTimeout: aws.Int64(36000),
			WaitTimeSeconds: aws.Int64(0),
			MessageAttributeNames: aws.StringSlice([]string {
				"All",
			}),
		})

		if err != nil {
			Log("error getting message: " + err.Error(), 0)
			continue
		}

		if len(result.Messages) == 0 {
			continue
		}

		if (result.Messages[0] == nil) {
			Log("Length of result.Messages is Non-Zero, but result.Messages is nil???", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}

		if (result.Messages[0].MessageAttributes["Host"] == nil) {
			// Malformed Message
			Log("MessageAttribute 'Host' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}

		if (result.Messages[0].MessageAttributes["Action"] == nil) {
			// Malformed Message
			Log("MessageAttribute 'Action' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}

		if (result.Messages[0].MessageAttributes["ActionLevel"] == nil) {
			// Malformed Message
			Log("MessageAttribute 'ActionLevel' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}
		
		if (result.Messages[0].MessageAttributes["Contact"] == nil) {
			// Malformed Message
			Log("MessageAttribute 'Contact' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}

		if (result.Messages[0].MessageAttributes["MiscParams"] == nil) {
			// Malformed Message
			Log("MessageAttribute 'MiscParams' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}
		
		if (result.Messages[0].MessageAttributes["RmPriority"] == nil) {
			// Malformed Message
			Log("MessageAttribute 'RmPriority' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}

		if (result.Messages[0].MessageAttributes["UseSNS"] == nil) {
			// Malformed Message
			Log("MessageAttribute 'UseSNS' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}

		if (result.Messages[0].MessageAttributes["OnRemote"] == nil) {
			// Malformed Message
			Log("MessageAttribute 'OnRemote' is nil...Mailformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}

		if (result.Messages[0].MessageAttributes["Host"].StringValue == nil) {
			// Malformed Message
			Log("String Value of MessageAttribute 'Host' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}

		if (result.Messages[0].MessageAttributes["Action"].StringValue == nil) {
			// Malformed Message
			Log("String Value of MessageAttribute 'Action' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}

		if (result.Messages[0].MessageAttributes["ActionLevel"].StringValue == nil) {
			// Malformed Message
			Log("String Value of MessageAttribute 'ActionLevel' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}
		
		if (result.Messages[0].MessageAttributes["Contact"].StringValue == nil) {
			// Malformed Message
			Log("String Value of MessageAttribute 'Contact' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}

		if (result.Messages[0].MessageAttributes["MiscParams"].StringValue == nil) {
			// Malformed Message
			Log("String Value of MessageAttribute 'MiscParams' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}

		if (result.Messages[0].MessageAttributes["RmPriority"].StringValue == nil) {
			// Malformed Message
			Log("String Value of MessageAttribute 'RmPriority' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}

		if (result.Messages[0].MessageAttributes["UseSNS"].StringValue == nil) {
			// Malformed Message
			Log("String Value of MessageAttribute 'UseSNS' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}

		if (result.Messages[0].MessageAttributes["OnRemote"].StringValue == nil) {
			// Malformed Message
			Log("String Value of MessageAttribute 'OnRemote' is nil...Malformed Message", 0)
			Log("Deleting Message With No Action Taken", 0)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})

			continue
		}

		action := ActionMessage{}
		action.Host, _ = cipherize.Decrypt(key, *result.Messages[0].MessageAttributes["Host"].StringValue)
		action.Action, _ = cipherize.Decrypt(key, *result.Messages[0].MessageAttributes["Action"].StringValue)
		action.ActionLevel, _ = cipherize.Decrypt(key, *result.Messages[0].MessageAttributes["ActionLevel"].StringValue)
		action.Contact, _ = cipherize.Decrypt(key, *result.Messages[0].MessageAttributes["Contact"].StringValue)
		action.MiscParams, _ = cipherize.Decrypt(key, *result.Messages[0].MessageAttributes["MiscParams"].StringValue)
		action.RmPriority, _ = cipherize.Decrypt(key, *result.Messages[0].MessageAttributes["RmPriority"].StringValue)
		action.OnRemote, _ = cipherize.Decrypt(key, *result.Messages[0].MessageAttributes["OnRemote"].StringValue)


		if action.RmPriority == "immediate" {
			_, err := svc.DeleteMessage(&sqs.DeleteMessageInput {
				QueueUrl: &queue,
				ReceiptHandle: result.Messages[0].ReceiptHandle,
			})
			
			if err != nil {
				Log("failed to delete message from queue: " + err.Error(), 0)
				Log("NO ACTION HAS BEEN TAKEN", 0)
				continue
			}

			err = HandleWork(&action)
			if err != nil {
				Log("Couldn't Find, Or Run utility.  Doing Nothing", 0)
				Log("SQS Message ALREADY Deleted", 0)
			} else {
				jsn, _ := json.Marshal(action)
				Log("Action has been taken", 0)
				Log(string(jsn), 0)
			}

		} else {
			err = HandleWork(&action)
			if err != nil {
				Log("Couldn't Find, Or Run utility.  Doing Nothing", 0)
				Log("SQS Message Not Deleted", 0)
			} else {
				jsn, _ := json.Marshal(action)
				Log("Action has been taken", 0)
				Log(string(jsn), 0)

				_, err := svc.DeleteMessage(&sqs.DeleteMessageInput {
					QueueUrl: &queue,
					ReceiptHandle: result.Messages[0].ReceiptHandle,
				})

				if err != nil {
					Log("failed to delete message from queue: " + err.Error(), 0)
					Log("ACTION HAS ALREADY BEEN TAKEN!!!  Please remove from queue", 0)
					continue
				}

			}
		}
	}
}
