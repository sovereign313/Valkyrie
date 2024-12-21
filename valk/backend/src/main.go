package main

import (
	"os" 
	"io"
	"fmt" 
	"net"
	"strings"
	"sshclient"
//	"cipherize"
	"time"
	"strconv"
	"vtypes"
	"setup"
	"install"
	"s3downloader"

	"os/exec"
	"math/rand"
	"io/ioutil"
	"net/http"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"

	"github.com/gorilla/mux"
	"gopkg.in/src-d/go-git.v4"
)

type Response struct {
	Message string
	Code string
}

var md5list map[string]string
var status string

func handleStatus(w http.ResponseWriter, r *http.Request) {
	if status != "" {
		fmt.Fprintf(w, status)
	} 
}

func handleVerifyKey(w http.ResponseWriter, r *http.Request) {
	if ! VerifyMD5Sums() {
		resp := Response{Message: "Corrupted Client Side Files", Code: "502"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return;
	}

	key := r.FormValue("licensekey")
	business := r.FormValue("business")

	if len(key) == 0 {
		resp := Response{Message: "Missing licensekey parameter", Code: "500"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return;
	}

	if len(business) == 0 {
		resp := Response{Message: "Missing business parameter", Code: "500"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return;
	}

	if VerifyLicense(key, business) {
		resp := Response{Message: "Key Is Valid", Code: "200"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return
	} else {
		resp := Response{Message: "Key Is Invalid", Code: "501"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return
	}
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	if ! VerifyMD5Sums() {
		resp := Response{Message: "Corrupted Client Side Files", Code: "502"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return;
	}

	fname := r.FormValue("fname")
	lname := r.FormValue("lname")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	company := r.FormValue("company")
	issue := r.FormValue("issue")
	license := r.FormValue("license")
	business := r.FormValue("business")

	if ! VerifyLicense(license, business) {
		resp := Response{Message: "Key Is Invalid!", Code: "501"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return
	}

	if len(fname) == 0 {
		fmt.Fprintf(w, "missing fname")
		return
	}

	if len(lname) == 0 {
		fmt.Fprintf(w, "missing lname")
		return
	}

	if len(phone) == 0 {
		fmt.Fprintf(w, "missing phone")
		return
	}

	if len(email) == 0 {
		fmt.Fprintf(w, "missing email")
		return
	}

	if len(company) == 0 {
		fmt.Fprintf(w, "missing company")
		return
	}

	if len(issue) == 0 {
		fmt.Fprintf(w, "missing issue")
		return
	}

	if len(license) == 0 {
		fmt.Fprintf(w, "missing license")
		return
	}

	body := strings.NewReader(`fname=` + fname + `&lname=` + lname + `&phone=` + phone + `&email=` + email + `&company=` + company + `&issue=` + issue + `&license=` + license)
	req, err := http.NewRequest("POST", "https://valkyriesoftware.io/contact", body)
	if err != nil {
		resp := Response{Message: err.Error(), Code: "500"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return;
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		resp := Response{Message: err.Error(), Code: "500"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return;
	}
	defer resp.Body.Close()

	retval, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resp := Response{Message: err.Error(), Code: "500"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return;
	}

	if string(retval) == "success" {
		resp := Response{Message: "success", Code: "200"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
	} else {
		resp := Response{Message: err.Error(), Code: "500"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return;
	}
}

func handleInstall(w http.ResponseWriter, r *http.Request) {
	if ! VerifyMD5Sums() {
		status = "Bailing out: JS files have been changed\n"
		resp := Response{Message: "Corrupted Client Side Files", Code: "502"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return;
	}

	license := r.FormValue("licensekey")
	business := r.FormValue("business")
	hostconfig := r.FormValue("hostconfig")
	foremanconfig := r.FormValue("foremanconfig")
	launcherconfig := r.FormValue("launcherconfig")
	workerconfig := r.FormValue("workerconfig")
	loggerconfig := r.FormValue("loggerconfig")
	alerterconfig := r.FormValue("alerterconfig")
	awsconfig := r.FormValue("awsconfig")
	dispatcherconfig := r.FormValue("dispatcherconfig")
	sqsreaderconfig := r.FormValue("sqsreaderconfig")
	mailreaderconfig := r.FormValue("mailreaderconfig")

	if ! VerifyLicense(license, business) {
		status = "Bailing out: invalid key\n"
		resp := Response{Message: "Key Is Invalid!", Code: "501"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return
	}

	if hostconfig == "none" {
		status = "Bailing Out: Missing Host Config!\n"
		resp := Response{Message: "Missing Host Config!", Code: "503"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return
	}

	if launcherconfig == "none" {
		status = "Bailing Out: Missing Launcher Config!\n"
		resp := Response{Message: "Missing Launcher Config!", Code: "504"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return
	}

	if workerconfig == "none" {
		status = "Bailing Out: Missing Worker Config!\n"
		resp := Response{Message: "Missing Worker Config!", Code: "505"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return
	}

	if foremanconfig == "none" {
		if awsconfig == "none" || dispatcherconfig == "none" || sqsreaderconfig == "none" {
			status = "Bailing Out: Foreman And AWS Are Not Configured!\n"
			resp := Response{Message: "Foreman OR AWS Needs To Be Configured!", Code: "506"}
			jsn, _ := json.Marshal(resp)
			fmt.Fprintf(w, string(jsn))
			return
		}
	}

	if awsconfig == "none" && dispatcherconfig == "none" && sqsreaderconfig == "none" {
		if foremanconfig == "none" {
			status = "Bailing Out: Foreman And AWS Are Not Configured!\n"
			resp := Response{Message: "Foreman OR AWS Needs To Be Configured!", Code: "506"}
			jsn, _ := json.Marshal(resp)
			fmt.Fprintf(w, string(jsn))
			return
		}
	}

	valkconfig := vtypes.ValkConfig{}
	valkconfig.LicenseKey = license
	valkconfig.BusinessName = business

	if hostconfig != "none" {
		err := json.Unmarshal([]byte(hostconfig), &valkconfig.HostConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal HostConfig: " + err.Error())
		}
	}

	if foremanconfig != "none" {
		err := json.Unmarshal([]byte(foremanconfig), &valkconfig.ForemanConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal ForemanConfig: " + err.Error())
		}
	}

	if launcherconfig != "none" {
		err := json.Unmarshal([]byte(launcherconfig), &valkconfig.LauncherConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal LauncherConfig: " + err.Error())
		}
	}

	if workerconfig != "none" {
		err := json.Unmarshal([]byte(workerconfig), &valkconfig.WorkerConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal WorkerConfig: " + err.Error())
		}
	}

	if loggerconfig != "none" {
		err := json.Unmarshal([]byte(loggerconfig), &valkconfig.LoggerConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal LoggerConfig: " + err.Error())
		}
	}

	if alerterconfig != "none" {
		err := json.Unmarshal([]byte(alerterconfig), &valkconfig.AlerterConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal AlerterConfig: " + err.Error())
		}
	}

	if awsconfig != "none" {
		err := json.Unmarshal([]byte(awsconfig), &valkconfig.AWSConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal AWSConfig: " + err.Error())
		}
	}

	if dispatcherconfig != "none" {
		err := json.Unmarshal([]byte(dispatcherconfig), &valkconfig.DispatcherConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal DispatcherConfig: " + err.Error())
		}
	}

	if sqsreaderconfig != "none" {
		err := json.Unmarshal([]byte(sqsreaderconfig), &valkconfig.SQSReaderConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal SQSReaderConfig: " + err.Error())
		}
	}

	if mailreaderconfig != "none" {
		err := json.Unmarshal([]byte(mailreaderconfig), &valkconfig.MailReaderConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal MailReaderConfig: " + err.Error())
		}
	}

	go func() {
		status = "Checking SSH Connections..."
		results := checkSSHConnection(valkconfig)
		flag := false
		hst := ""

		for host, value := range results {
			if value != "success" {
				flag = true
				hst = host
				break;
			}
		}

		if flag {
			status += "Failed: " + hst
			return
		}

		status += "All Good.\n"

		status += "Checking For Docker..."
		dockerresults := CheckForDocker(valkconfig)
		flag = true 
		for _, value := range dockerresults {
			if value != "success" {
				flag = false
			}
		}

		if flag {
			status += "found it\n"
		} else {
			status += "Nope.  Using scripts instead!\n"
		}

		if ! flag {
			status += "Checking For Docker On Launcher/Worker Hosts..."
			if ! CheckForDockerOnWorkers(valkconfig) {
				status += "Nope - :(\n"
				status += "Terminating Install. Docker Must Be Installed On Worker Hosts!\n"
				return
			}
			status += "Found! Proceeding.\n"
		}

		status += "Downloading Valkyrie..."
		ok, err := s3downloader.DownloadFile("valkyrie.tar.gz", "valkdls", "us-east-1")
		if err != nil {
			status += "Booo. Failed. " + err.Error() + "\n"
			return
		}

		if ! ok {
			status += "Booo. Failed.\n" 
			return
		}

		status += "Yay. Got it!\n"
		valkconfig.UseDocker = flag

		status += "Extracting Valkyrie..."
	        _, err = exec.Command("/usr/bin/tar", "xfz", "valkyrie.tar.gz").Output()
		status += "Done\n"

		if valkconfig.ForemanConfig.Host != "" {
			status += "Setting Up Foreman..."
			err = setup.Setup_Foreman(valkconfig)
			if err != nil {
				status += "Booo.  Failed: " + err.Error() + "\n"
				return
			}
			status += "Done!\n"
		} else {
			status += "Skipping Foreman (Not Configured)\n"
		}

		status += "Setting Up Launcher..."
		err = setup.Setup_Launcher(valkconfig)
		if err != nil {
			status += "Booo. Failed: " + err.Error() + "\n"
			return
		}
		status += "Done!\n"

		status += "Setting Up Worker..."
		err = setup.Setup_Worker(valkconfig)
		if err != nil {
			status += "Booo. Failed: " + err.Error() + "\n"
			return
		}
		status += "Done!\n"

		if valkconfig.LoggerConfig.Host != "" {
			status += "Setting Up Logger..."
			err = setup.Setup_Logger(valkconfig)
			if err != nil {
				status += "Booo. Failed: " + err.Error() + "\n"
				return
			}
			status += "Done!\n"
		} else {
			status += "Skipping Logger (Not Configured)\n"
		}

		if valkconfig.AlerterConfig.Host != "" {
			status += "Setting Up Alerter..."
			err = setup.Setup_Alerter(valkconfig)
			if err != nil {
				status += "Booo. Failed: " + err.Error() + "\n"
				return
			}
			status += "Done!\n"
		} else {
			status += "Skipping Alerter (Not Configured)\n"
		}

		if valkconfig.AWSConfig.Region != "" && valkconfig.AWSConfig.SQSName != "" && valkconfig.AWSConfig.AWSAccessKey != "" && valkconfig.AWSConfig.AWSSecretKey != "" && valkconfig.AWSConfig.EncryptionKey != "" {
			if valkconfig.DispatcherConfig.Host != "" {
				status += "Setting Up Dispatcher..."
				err = setup.Setup_Dispatcher(valkconfig)
				if err != nil {
					status += "Booo. Failed: " + err.Error() + "\n"
					return
				}
				status += "Done!\n"
			} else {
				status += "Skipping Dispatcher (Not Configured)\n"
			}
		} else {
			status += "Skipping Dispatcher (AWS Not Configured)\n"
		}

		if valkconfig.AWSConfig.Region != "" && valkconfig.AWSConfig.SQSName != "" && valkconfig.AWSConfig.AWSAccessKey != "" && valkconfig.AWSConfig.AWSSecretKey != "" && valkconfig.AWSConfig.EncryptionKey != "" {
			if valkconfig.SQSReaderConfig.Host != "" {
				status += "Setting Up SQSReader..."
				err = setup.Setup_SQSReader(valkconfig)
				if err != nil {
					status += "Booo. Failed: " + err.Error() + "\n"
					return
				}
				status += "Done!\n"
			} else {
				status += "Skipping SQSReader (Not Configured)\n"
			}
		} else {
			status += "Skipping SQSReader (AWS Not Configured)\n"
		}

		if valkconfig.MailReaderConfig.Host != "" {
			status += "Setting Up MailReader..."
			err = setup.Setup_MailReader(valkconfig)
			if err != nil {
				status += "Booo. Failed: " + err.Error() + "\n"
				return
			}
			status += "Done!\n"
		} else {
			status += "Skipping MailReader (Not Configured)\n"
		}

		status += "Cleaning Up My Setup Mess..."
		os.RemoveAll("./valkyrie")
		os.Remove("./valkyrie.tar.gz")
		status += "Done!\n"

		status += "Starting Installation.\n"

/*	************************************************************* */

		if valkconfig.ForemanConfig.Host != "" {
			status += "Installing Foreman..."
			err = install.Install_Foreman(valkconfig)
			if err != nil {
				status += "Booo.  Failed: " + err.Error() + "\n"
				return
			}
			status += "Done!\n"
		} else {
			status += "Skipping Foreman (Not Configured)\n"
		}

		status += "Installing Launcher..."
		err = install.Install_Launcher(valkconfig)
		if err != nil {
			status += "Booo. Failed: " + err.Error() + "\n"
			return
		}
		status += "Done!\n"

		status += "Installing Worker..."
		err = install.Install_Worker(valkconfig)
		if err != nil {
			status += "Booo. Failed: " + err.Error() + "\n"
			return
		}
		status += "Done!\n"

		if valkconfig.LoggerConfig.Host != "" {
			status += "Installing Logger..."
			err = install.Install_Logger(valkconfig)
			if err != nil {
				status += "Booo. Failed: " + err.Error() + "\n"
				return
			}
			status += "Done!\n"
		} else {
			status += "Skipping Logger (Not Configured)\n"
		}

		if valkconfig.AlerterConfig.Host != "" {
			status += "Installing Alerter..."
			err = install.Install_Alerter(valkconfig)
			if err != nil {
				status += "Booo. Failed: " + err.Error() + "\n"
				return
			}
			status += "Done!\n"
		} else {
			status += "Skipping Alerter (Not Configured)\n"
		}

		if valkconfig.AWSConfig.Region != "" && valkconfig.AWSConfig.SQSName != "" && valkconfig.AWSConfig.AWSAccessKey != "" && valkconfig.AWSConfig.AWSSecretKey != "" && valkconfig.AWSConfig.EncryptionKey != "" {
			if valkconfig.DispatcherConfig.Host != "" {
				status += "Installing Dispatcher..."
				err = install.Install_Dispatcher(valkconfig)
				if err != nil {
					status += "Booo. Failed: " + err.Error() + "\n"
					return
				}
				status += "Done!\n"
			} else {
				status += "Skipping Dispatcher (Not Configured)\n"
			}
		} else {
			status += "Skipping Dispatcher (AWS Not Configured)\n"
		}

		if valkconfig.AWSConfig.Region != "" && valkconfig.AWSConfig.SQSName != "" && valkconfig.AWSConfig.AWSAccessKey != "" && valkconfig.AWSConfig.AWSSecretKey != "" && valkconfig.AWSConfig.EncryptionKey != "" {
			if valkconfig.SQSReaderConfig.Host != "" {
				status += "Installing SQSReader..."
				err = install.Install_SQSReader(valkconfig)
				if err != nil {
					status += "Booo. Failed: " + err.Error() + "\n"
					return
				}
				status += "Done!\n"
			} else {
				status += "Skipping SQSReader (Not Configured)\n"
			}
		} else {
			status += "Skipping SQSReader (AWS Not Configured)\n"
		}

		if valkconfig.MailReaderConfig.Host != "" {
			status += "Installing MailReader..."
			err = install.Install_MailReader(valkconfig)
			if err != nil {
				status += "Booo. Failed: " + err.Error() + "\n"
				return
			}
			status += "Done!\n"
		} else {
			status += "Skipping MailReader (Not Configured)\n"
		}

		status += "Cleaning Up My Installation Mess..."
		os.RemoveAll("./dist")
		status += "Done!\n"

		status += "All Done.\n"
		return
	}()

	resp := Response{Message: "Successfully Launched Install! It Can Take A Long Time! Please Wait For It To Complete.", Code: "200"}
	jsn, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(jsn))
}

func GetRepo(gitrepourl string) error {
	if _, err := os.Stat("./valkyrie"); ! os.IsNotExist(err) {
		os.RemoveAll("./valkyrie")
	}

        _, err := git.PlainClone("./valkyrie", false, &git.CloneOptions{
                URL: gitrepourl,
        })

        if err != nil {
		return err
        }

	return nil
}

func CheckForDockerOnWorkers(vc vtypes.ValkConfig) bool {
	keyrsa, _ := sshclient.SignerFromBytes([]byte(vc.WorkerConfig.SSHPrivateKey))

	for _, host := range vc.LauncherConfig.Host {
		output, err := sshclient.RunOneCommand(host, "[[ -e /var/run/docker.sock ]]; echo $?", 5, keyrsa)
		output = strings.TrimSpace(output)
		if err != nil {
			return false
		}

		if output != "0" {
			return false
		}
	}

	return true
}

func CheckForDocker(vc vtypes.ValkConfig) map[string]string {
	retval := make(map[string]string)

	keyrsa, _ := sshclient.SignerFromBytes([]byte(vc.WorkerConfig.SSHPrivateKey))
	for _, host := range vc.HostConfig {
		output, err := sshclient.RunOneCommand(host, "[[ -e /var/run/docker.sock ]]; echo $?", 5, keyrsa)
		output = strings.TrimSpace(output)
		if err != nil {
			retval[host] = err.Error()
			continue
		}

		if output != "0" {
			retval[host] = output
			continue
		}

		retval[host] = "success"
	}

	return retval

}

func checkSSHConnection(vc vtypes.ValkConfig) map[string]string {
	retval := make(map[string]string)

	keyrsa, _ := sshclient.SignerFromBytes([]byte(vc.WorkerConfig.SSHPrivateKey))
	for _, host := range vc.HostConfig {
		output, err := sshclient.RunOneCommand(host, "/usr/bin/uptime", 5, keyrsa)
		if err != nil {
			fmt.Println(err.Error())
			retval[host] = err.Error()
			continue
		}

		if ! strings.Contains(output, "load average") {
			fmt.Println(output)
			retval[host] = output
			continue
		}

		retval[host] = "success"
	}

	return retval
}

func handleGenerateManifest(w http.ResponseWriter, r *http.Request) {
	if ! VerifyMD5Sums() {
		resp := Response{Message: "Corrupted Client Side Files", Code: "502"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return;
	}

	license := r.FormValue("licensekey")
	business := r.FormValue("business")
	hostconfig := r.FormValue("hostconfig")
	foremanconfig := r.FormValue("foremanconfig")
	launcherconfig := r.FormValue("launcherconfig")
	workerconfig := r.FormValue("workerconfig")
	loggerconfig := r.FormValue("loggerconfig")
	alerterconfig := r.FormValue("alerterconfig")
	awsconfig := r.FormValue("awsconfig")
	dispatcherconfig := r.FormValue("dispatcherconfig")
	sqsreaderconfig := r.FormValue("sqsreaderconfig")
	mailreaderconfig := r.FormValue("mailreaderconfig")

	if ! VerifyLicense(license, business) {
		resp := Response{Message: "Key Is Invalid!", Code: "501"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return
	}

	if hostconfig == "none" {
		resp := Response{Message: "Missing Host Config!", Code: "503"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return
	}

	if launcherconfig == "none" {
		resp := Response{Message: "Missing Launcher Config!", Code: "504"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return
	}

	if workerconfig == "none" {
		resp := Response{Message: "Missing Worker Config!", Code: "505"}
		jsn, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsn))
		return
	}

	if foremanconfig == "none" {
		if awsconfig == "none" || dispatcherconfig == "none" || sqsreaderconfig == "none" {
			resp := Response{Message: "Foreman OR AWS Needs To Be Configured!", Code: "506"}
			jsn, _ := json.Marshal(resp)
			fmt.Fprintf(w, string(jsn))
			return
		}
	}

	if awsconfig == "none" && dispatcherconfig == "none" && sqsreaderconfig == "none" {
		if foremanconfig == "none" {
			resp := Response{Message: "Foreman OR AWS Needs To Be Configured!", Code: "506"}
			jsn, _ := json.Marshal(resp)
			fmt.Fprintf(w, string(jsn))
			return
		}
	}

	valkconfig := vtypes.ValkConfig{}
	valkconfig.LicenseKey = license
	valkconfig.BusinessName = business

	if hostconfig != "none" {
		err := json.Unmarshal([]byte(hostconfig), &valkconfig.HostConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal HostConfig: " + err.Error())
		}
	}

	if foremanconfig != "none" {
		err := json.Unmarshal([]byte(foremanconfig), &valkconfig.ForemanConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal ForemanConfig: " + err.Error())
		}
	}

	if launcherconfig != "none" {
		err := json.Unmarshal([]byte(launcherconfig), &valkconfig.LauncherConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal LauncherConfig: " + err.Error())
		}
	}

	if workerconfig != "none" {
		err := json.Unmarshal([]byte(workerconfig), &valkconfig.WorkerConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal WorkerConfig: " + err.Error())
		}
	}

	if loggerconfig != "none" {
		err := json.Unmarshal([]byte(loggerconfig), &valkconfig.LoggerConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal LoggerConfig: " + err.Error())
		}
	}

	if alerterconfig != "none" {
		err := json.Unmarshal([]byte(alerterconfig), &valkconfig.AlerterConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal AlerterConfig: " + err.Error())
		}
	}

	if awsconfig != "none" {
		err := json.Unmarshal([]byte(awsconfig), &valkconfig.AWSConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal AWSConfig: " + err.Error())
		}
	}

	if dispatcherconfig != "none" {
		err := json.Unmarshal([]byte(dispatcherconfig), &valkconfig.DispatcherConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal DispatcherConfig: " + err.Error())
		}
	}

	if sqsreaderconfig != "none" {
		err := json.Unmarshal([]byte(sqsreaderconfig), &valkconfig.SQSReaderConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal SQSReaderConfig: " + err.Error())
		}
	}

	if mailreaderconfig != "none" {
		err := json.Unmarshal([]byte(mailreaderconfig), &valkconfig.MailReaderConfig);
		if err != nil {
			fmt.Println("Couldn't Marshal MailReaderConfig: " + err.Error())
		}
	}

	jsn, err := json.Marshal(valkconfig)
	if err != nil {
		fmt.Println("failed to marshal valkconfig: " + err.Error())
		return;
	}

	resp := Response{Message: string(jsn), Code: "200"}
	rjsn, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(rjsn))
}

func VerifyMD5Sums() bool {
	md5list = make(map[string]string)

	// Get This From The Server
	md5list["valk.alerter.js"] = "8a22a66ae6f39650a8e2cd679e78f554"
	md5list["valk.aws.js"] = "bb97939dad283b3855112ae3ca3719f8"
	md5list["valk.dispatcher.js"] = "effd9da9369680a09884456229aa922e"
	md5list["valk.foreman.js"] = "2d1ba83ae81218cab6a0c8a7f2a23654"
	md5list["valk.hosts.js"] = "d01b31cd9c1989a046017f81d0add906"
	md5list["valk.js"] = "d8813b86a8a575610524e4ca34badb23"
	md5list["valk.launcher.js"] = "f2485661ba5740470ade72608d33a7d1"
	md5list["valk.logger.js"] = "cb552e9a9830dc043f1b21547b74a77e"
	md5list["valk.mailreader.js"] = "bb2088408e88684c36fe3cf9fc474e3e"
	md5list["valk.sqsreader.js"] = "e4608e2b4159729677a45a4bc80ab102"
	md5list["valk.support.js"] = "f43ec5d24b10cba8e0ebd9d374f4cbc8"
	md5list["valk.worker.js"] = "abc9b81397035ff50b266bab370ea393"

	for file, hash := range md5list {
		if _, err := os.Stat("./static/js/" + file); os.IsNotExist(err) {
			return false
		}

		md5sum, err := HashMD5File("./static/js/" + file)
		if err != nil {
			return false
		}

		if md5sum != hash {
			return false
		}
	}
	
	return true	
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
	var ips []string

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}

	router := mux.NewRouter()
	router.HandleFunc("/verifykey", handleVerifyKey)
	router.HandleFunc("/contact", handleContact)
	router.HandleFunc("/generatemanifest", handleGenerateManifest)
	router.HandleFunc("/status", handleStatus)
	router.HandleFunc("/install", handleInstall)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	hname, _ := os.Hostname()
	fmt.Println("If you are running this on a remote system (via SSH) then please")
	fmt.Println("go to the IP or hostname of this system in your web browser")
	fmt.Println("** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** ** **")
	fmt.Println("Maybe one of these:")
	fmt.Println("http://" + hname + ":5080/")
	for _, ip := range ips {
		fmt.Println("http://" + ip + ":5080/")
	}

	err = http.ListenAndServe(":5080", router)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
