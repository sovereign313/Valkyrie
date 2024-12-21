package main

import (
	"os"
	"fmt"
	"net"
	"errors"
	"time"
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
)

var useeventstreams bool
var usesecurelogging bool

var eshostport string
var logfilelocation string

var logkey string

func Log(service string, loglevel string, message string) error {
	if useeventstreams {
		err := EventStream(service, loglevel, message)
		return err
	} else {
		err := FileLog(service, loglevel, message)
		return err
	}

	return nil
}

func EventStream(service string, loglevel string, message string) error {
	conn, err := net.DialTimeout("tcp", eshostport, 5 * time.Second)
	if err != nil {
		return err
	}

	current_time := time.Now().Local()
	t := current_time.Format("Jan 02 2006 03:04:05")

	msg := loglevel + " | " + t + " | " + service + " | " + message + "\n"
	fmt.Fprintf(conn, msg)
	conn.Close()

	FileLog(service, loglevel, message)
	return nil
}	

func FileLog(service string, loglevel string, message string) error {
	file, err := os.OpenFile(logfilelocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("Failed to open log file for writing: " + err.Error())
	}
	defer file.Close()

	current_time := time.Now().Local()
	t := current_time.Format("Jan 02 2006 03:04:05")
	_, err = file.WriteString(loglevel + " | " + t + " | " + service + " | " + message + "\n")

	if err != nil {
		return errors.New("Failed to write to log file: " + err.Error())
	}

	return nil
}

func handleWhoAreYou(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Valkyrie Logger")
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func handleDescription(w http.ResponseWriter, r *http.Request) {
	html := "Valkyrie Logger - Tool for handling logging for valkyrie micro-services\n"
	fmt.Fprintf(w, html)
}

func handleLog(w http.ResponseWriter, r *http.Request) {
	Service := r.URL.Query().Get("service")
	LogLevel := r.URL.Query().Get("loglevel")
	Message := r.URL.Query().Get("message")
	Key := r.URL.Query().Get("logkey")	

	if Service == "" {
		fmt.Fprintf(w, "failed to log: missing service")
		return
	}

	if LogLevel == "" {
		fmt.Fprintf(w, "failed to log: missing loglevel")
		return
	}

	if Message == "" {
		fmt.Fprintf(w, "failed to log: missing message")
		return
	}

	if usesecurelogging {
		if Key == "" {
			fmt.Fprintf(w, "failed to log: missing key")
			return
		}

		if Key != logkey {
			fmt.Fprintf(w, "failed to log: invalid key")
			return
		}
	}

	err := Log(Service, LogLevel, Message)
	if err != nil {
		html := "failed to log: " + err.Error() + "\n"
		html += "attempted to write: " + Service + ", " + LogLevel + ", " + Message + "\n"
		fmt.Fprintf(w, html) 
		return
	}

	fmt.Fprintf(w, "Logged Message")
}

func handleHelp(w http.ResponseWriter, r *http.Request) {
	html := "/            - This Help\n"
	html += "/ping        - Returns Pong (Ensures Service Is Working)\n"
	html += "/whoareyou   - Returns The Application (Valkyrie Foreman)\n"
	html += "/description - Returns A Description Of This Service\n"
	html += "/log         - Executes A Worker Task \n"
	fmt.Fprintf(w, html)
}

func main() {
	var err error

	struseeventstreams := os.Getenv("useeventstreams")
	eshostport = os.Getenv("eshostport")
	logfilelocation = os.Getenv("logfilelocation")
	strusesecurelogging := os.Getenv("usesecurelogging")
	logkey = os.Getenv("keylog")

	nflag, err := strconv.ParseBool(struseeventstreams)
	if err != nil {
		Log("Valkyrie Logger", "warning", "failed to set useeventstreams " + err.Error() + " (using default value of: false)")
		fmt.Println("failed to set useeventstreams " + err.Error() + " (using default value of: false)")
		useeventstreams = false
	} else {
		useeventstreams = nflag
	}

	if useeventstreams == true && len(eshostport) == 0 {
                Log("Valkyrie Logger", "error", "can't use event streams without a host/port to send data to")
		fmt.Println("can't use event streams without a host/port to send data to")
		return
	}

	fflag, err := strconv.ParseBool(strusesecurelogging)
	if err != nil {
                Log("Valkyrie Logger", "warning", "failed to set usesecurelogging (using default value of: false)")
		fmt.Println("failed to set usesecurelogging (using default value of: false")
		usesecurelogging = false
	} else {
		usesecurelogging = fflag
	}

	if usesecurelogging == true && len(logkey) == 0 {
                Log("Valkyrie Logger", "warning", "secure logging true, but no key set (using default value of: mykey")
		logkey = "mykey"
	}

	if len(logfilelocation) == 0 {
		logfilelocation = "/tmp/valkryie.log"
	}



	router := mux.NewRouter()
        router.HandleFunc("/whoareyou", handleWhoAreYou)
        router.HandleFunc("/ping", handlePing)
	router.HandleFunc("/description", handleDescription)
	router.HandleFunc("/log", handleLog)
        router.HandleFunc("/", handleHelp)

        err = http.ListenAndServe(":8092", router)
        if err != nil {
                fmt.Println("ListenAndServe: ", err)
	}

}
