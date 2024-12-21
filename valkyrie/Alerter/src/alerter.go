package main

import (
	"os"
	"fmt"
	"net"
	"bytes"
	"errors"
	"strings"
	"time"
	"strconv"
//	"cipherize"

	"math/rand"
	"net/url"
	"net/http"
	"net/smtp"
	"encoding/json"

	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var fromemail string
var mailserver string

var twilio_account string
var twilio_token string
var twilio_from string

var awsregion string
var aws_access_key_id string
var aws_secret_key_id string

var logfilelocation string
var logkey string
var loggerurl string
var usesecurelogging bool

var license string
var business string

func Log(message string, level int) error {
        var Level string
        Service := "Valkyrie Alerter"

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
        file, err := os.OpenFile("./alerter.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
                fmt.Println("failed to open log file: " + err.Error())
                return err
        }
        defer file.Close()

        current_time := time.Now().Local()
        t := current_time.Format("Jan 02 2006 03:04:05")
        _, err = file.WriteString(t + " - Valkyrie Alerter: " + message + "\n")

        if err != nil {
                fmt.Println("failed to write to log file: " + err.Error())
                return err
        }

        return nil
}

func SendSMTPMessage(mailserver string, from string, to string, subject string, body string) error {
	conn, err := net.DialTimeout("tcp", mailserver, 10 * time.Second)
	if err != nil {
		return err
	}	

	host, _, _ := net.SplitHostPort(mailserver)
	connection, err := smtp.NewClient(conn, host)
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

func SendSMSTwilio(phonenumber string, fromnumber string, body string) error {
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + twilio_account + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("To","+" + phonenumber)
	msgData.Set("From","+" + fromnumber)
	msgData.Set("Body", "totally a test")
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(twilio_account, twilio_token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if (resp.StatusCode >= 200 && resp.StatusCode < 300) {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if (err != nil) {
			return err
		}

		return nil
	} else {
		Log("Failed To Alert.  Response From Twilio: " + resp.Status, 0)
		return errors.New("failed to alert")
	}

	return nil
}

func SendSMSSNS(region string, phonenumber string, body string) error {
	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_key_id, "")
	_, err := creds.Get()
	if err != nil {
		return err
	}

	cfg := aws.NewConfig().WithRegion(region).WithCredentials(creds)
	sess, err := session.NewSession(cfg)
	svc := sns.New(sess)

        params := &sns.PublishInput{
                Message: aws.String(body),
		PhoneNumber: aws.String("+" + phonenumber), 
        }

        _, err = svc.Publish(params)
        if err != nil {
                return err
        }

        return nil 
}


func handleWhoAreYou(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Valkyrie Alerter")
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func handleDescription(w http.ResponseWriter, r *http.Request) {
	html := "Valkyrie Alerter - Tool for alerting admins from valkyrie micro-services\n"
	fmt.Fprintf(w, html)
}

func handleAlert(w http.ResponseWriter, r *http.Request) {
	var Contact string
	var Subject string
	var Message string
	var Method string
	var AppID string

	Method = r.URL.Query().Get("method")

	if len(Method) == 0 {
		fmt.Fprintf(w, "failed to alert: missing method [email, text, app]")
		return
	}

	switch Method {
		case "email":
			Contact = r.URL.Query().Get("email")
			Subject = r.URL.Query().Get("subject")
			Message = r.URL.Query().Get("message")

			if len(Contact) == 0 {
				fmt.Fprintf(w, "failed to alert: missing email")
				return
			}

			if len(Subject) == 0 {
				fmt.Fprintf(w, "failed to alert: missing subject")
				return
			}

			if len(Message) == 0 {
				fmt.Fprintf(w, "failed to alert: missing message")
				return
			}

			err := SendSMTPMessage(mailserver, fromemail, Contact, Subject, Message)
			if err != nil {
				Log("failed to send email: " + err.Error(), 0)
				fmt.Fprintf(w, "failed to send email: " + err.Error())
				return
			}

			Log("Email Alert Sent To: " + Contact, 0)
			fmt.Fprintf(w, "success")
			return
		case "text":
			Contact = r.URL.Query().Get("phonenumber")
			Message = r.URL.Query().Get("message")
			Provider := r.URL.Query().Get("provider")

			if len(Contact) == 0 {
				fmt.Fprintf(w, "failed to alert: missing phonenumber")
				return
			}

			if len(Message) == 0 {
				fmt.Fprintf(w, "failed to alert: missing message")
				return
			}

			if len(Provider) == 0 {
				fmt.Fprintf(w, "failed to alert: missing provider")
				return
			}

			if Provider == "aws" {
				err := SendSMSSNS(awsregion, Contact, Message)
				if err != nil {
					Log("failed to send SNS Text: " + err.Error(), 0)
					fmt.Fprintf(w, "failed to send SNS Text: " + err.Error())
					return
				}

				Log("Text Alert Via SNS Sent To: " + Contact, 0)
				fmt.Fprintf(w, "success")
			} else if Provider == "twilio" {
				err := SendSMSTwilio(Contact, twilio_from, Message)
				if err != nil {
					Log("failed to send Twilio Text: " + err.Error(), 0)
					fmt.Fprintf(w, "failed to send Twilio Text: " + err.Error())
					return
				}

				Log("Text Alert Via Twilio Sent To: " + Contact, 0)
				fmt.Fprintf(w, "success")
			} else {
				fmt.Fprintf(w, "Failed: Provider Not Supported")
				
			}

			return
		case "app":
			AppID = r.URL.Query().Get("appid")
			if len(AppID) == 0 {
				fmt.Fprintf(w, "failed to alert: method is app, but no appid specified")
				return
			}
		default:
			fmt.Fprintf(w, "Invalid Contact Method")
			return
	}

	return
}

func handleHelp(w http.ResponseWriter, r *http.Request) {
	html := "/            - This Help\n"
	html += "/ping        - Returns Pong (Ensures Service Is Working)\n"
	html += "/whoareyou   - Returns The Application (Valkyrie Alerter)\n"
	html += "/description - Returns A Description Of This Service\n"
	html += "/alert       - Triggers An Alert \n"
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

func main() {
	var err error

	strusesecurelogging := os.Getenv("usesecurelogging")
	logfilelocation = os.Getenv("logfilelocation")
	loggerurl = os.Getenv("loggerurl")	
	logkey = os.Getenv("keylog")

	fromemail = os.Getenv("fromemail")
	mailserver = os.Getenv("mailserver")

	twilio_account = os.Getenv("twilio_account")
	twilio_token = os.Getenv("twilio_token")
	twilio_from = os.Getenv("twilio_from")

	awsregion = os.Getenv("awsregion")
	aws_access_key_id = os.Getenv("aws_access_key_id")
	aws_secret_key_id = os.Getenv("aws_secret_key_id")

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


	if len(mailserver) == 0 {
		Log("failed to set mailserver. Email is not functional.", 0)
		fmt.Println("failed to set mailserver. Email is not functional.")
	}

	if len(fromemail) == 0 {
		if len(mailserver) != 0 {
			Log("failed to set fromemail (using default value of: valkalerter@" + mailserver + ")", 0)
			fmt.Println("failed to set fromeamil (using default value of: valkalerter@" + mailserver)
		}
	}

	fflag, err := strconv.ParseBool(strusesecurelogging)
	if err != nil {
                Log("failed to set usesecurelogging (using default value of: false)", 0)
		fmt.Println("failed to set usesecurelogging (using default value of: false")
		usesecurelogging = false
	} else {
		usesecurelogging = fflag
	}

	if usesecurelogging == true && len(logkey) == 0 {
                Log("secure logging true, but no key set (using default value of: mykey)", 0)
		logkey = "mykey"
	}

	if len(logfilelocation) == 0 {
		logfilelocation = "/tmp/valkryie.log"
	}


	router := mux.NewRouter()
        router.HandleFunc("/whoareyou", handleWhoAreYou)
        router.HandleFunc("/ping", handlePing)
	router.HandleFunc("/description", handleDescription)
	router.HandleFunc("/alert", handleAlert)
        router.HandleFunc("/", handleHelp)

        err = http.ListenAndServe(":8093", router)
        if err != nil {
                fmt.Println("ListenAndServe: ", err)
	}

}
