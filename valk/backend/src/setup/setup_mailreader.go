package setup

import (
	"os"
	"vtypes"
	"errors"
	"strings"

	"os/exec"
	"io/ioutil"
)

func Setup_MailReader(valkconfig vtypes.ValkConfig) error {

	err := WriteMailReaderFiles(valkconfig)
	if err != nil {
		return errors.New("Failed to write MailReader files: " + err.Error())
	}

	os.Mkdir("./deployment", 0755)

	if valkconfig.UseDocker {
		os.Chdir("./valkyrie/EmailReader/")
		output, err := exec.Command("/usr/bin/docker", "build", "-t", "mailreader", ".").Output()
		if err != nil {
			return errors.New("Failed to build MailReader: " + err.Error() + " - " + string(output))
		}

		os.Chdir("../..")

		_, err = exec.Command("/usr/bin/docker", "save", "-o", "./deployment/mailreader.tar", "mailreader").Output()
		if err != nil {
			return errors.New("Failed to write MailReader to tar: " + err.Error())
		}

		_, err = exec.Command("/usr/bin/docker", "rmi", "mailreader").Output()
		if err != nil {
			return errors.New("Failed to remove MailReader docker image: " + err.Error())
		}
	} else {
		err := CopyFile("./valkyrie/EmailReader/mailreader", "./deployment/mailreader")
		if err != nil {
			return errors.New("Failed to copy MailReader binary: " + err.Error())
		}

		err = CopyFile("./valkyrie/EmailReader/startup.sh", "./deployment/start_mailreader.sh")
		if err != nil {
			return errors.New("Failed to copy MailReader startup script: " + err.Error())
		}
	}

	return nil
}

func WriteMailReaderFiles(valkconfig vtypes.ValkConfig) error {

	dockerfiledata := "# This Can Be Set FROM scratch, but if you plan to use git, it's better to use alpine.\n"
	dockerfiledata += "# FROM scratch\n"
	dockerfiledata += "FROM alpine\n"
	dockerfiledata += "MAINTAINER Ernest E. Teem III <eteem@valkyriesoftware.io>\n"
	dockerfiledata += "ADD emailreader /emailreader\n\n"

        nhst := ""
        for _, hst := range valkconfig.LauncherConfig.Host {
                nhst += "http://" + hst + ":8094|"
        }
        nhst = strings.TrimSuffix(nhst, "|")

        dockerfiledata += "ENV dupprotection " + valkconfig.MailReaderConfig.DupProtection  + "\n"

	if valkconfig.MailReaderConfig.ProtectTime == "" {
	        dockerfiledata += "ENV dupprottime " + valkconfig.MailReaderConfig.ProtectTime + "\n\n"
	} else {
	        dockerfiledata += "ENV dupprottime 30\n\n"
	}

        dockerfiledata += "ENV worker_hosts \"" + nhst + "\"\n\n"

	dockerfiledata += "ENV mailhost " + valkconfig.MailReaderConfig.MailServer + "\n"
	dockerfiledata += "ENV mailuser " + valkconfig.MailReaderConfig.MailUser + "\n"
	dockerfiledata += "ENV mailpass " + valkconfig.MailReaderConfig.MailPassword + "\n"
	dockerfiledata += "ENV mailproto " + valkconfig.MailReaderConfig.MailProtocol + "\n"
	dockerfiledata += "ENV mailusetls " + valkconfig.MailReaderConfig.UseTLS + "\n"
	dockerfiledata += "ENV sleeptimeout " + valkconfig.MailReaderConfig.SleepTimeout + "\n"
	dockerfiledata += "ENV imapfolder INBOX\n"
	dockerfiledata += "ENV trigger_subject " + valkconfig.MailReaderConfig.MailSubjectTrigger + "\n"

        if valkconfig.LoggerConfig.Host != "" {
                dockerfiledata += "ENV loggerurl http://" + valkconfig.LoggerConfig.Host + ":8092/\n"
                dockerfiledata += "ENV usesecurelogging " + valkconfig.LoggerConfig.UseSecureLogging + "\n"
                dockerfiledata += "ENV logkey \"" + valkconfig.LoggerConfig.LogKey + "\"\n\n"
        } else {
                dockerfiledata += "ENV loggerurl http://localhost:8092/\n"
                dockerfiledata += "ENV usesecurelogging false\n"
                dockerfiledata += "ENV logkey \"\"\n\n"
        }

        dockerfiledata += "ENV license " + valkconfig.LicenseKey + "\n"
        dockerfiledata += "ENV business " + valkconfig.BusinessName + "\n\n"

	dockerfiledata += "# If Going Through A CNTLM Proxy\n"
	dockerfiledata += "#ENV HTTP_PROXY=\"http://localhost:3128\"\n"
	dockerfiledata += "#ENV HTTPS_PROXY=\"http://localhost:3128\"\n"
	dockerfiledata += "#ENV http_proxy=\"http://localhost:3128\"\n"
	dockerfiledata += "#ENV https_proxy=\"http://localhost:3128\"\n\n"

	dockerfiledata += "RUN apk update && apk add tzdata\n"
	dockerfiledata += "RUN cp /usr/share/zoneinfo/US/Eastern /etc/localtime\n"
	dockerfiledata += "RUN echo \"US/Eastern\" > /etc/timezone\n"
	dockerfiledata += "RUN apk del tzdata && apk add git && apk add openrc && apk add ca-certificates\n\n"
	dockerfiledata += "CMD [\"/mailreader\"]\n"
	
	err := ioutil.WriteFile("./valkyrie/EmailReader/Dockerfile", []byte(dockerfiledata), 0755)
	if err != nil {
		return err
	}

	scriptfiledata := "#!/bin/bash\n\n"
        scriptfiledata += "export dupprotection=\"" + valkconfig.MailReaderConfig.DupProtection  + "\"\n"

	if valkconfig.MailReaderConfig.ProtectTime == "" {
	        scriptfiledata += "export dupprottime=\"" + valkconfig.MailReaderConfig.ProtectTime + "\"\n"
	} else {
	        scriptfiledata += "export dupprottime=\"30\"\n\n"
	}

        scriptfiledata += "export worker_hosts=\"" + nhst + "\"\n\n"

	scriptfiledata += "export mailhost=\"" + valkconfig.MailReaderConfig.MailServer + "\"\n"
	scriptfiledata += "export mailuser=\"" + valkconfig.MailReaderConfig.MailUser + "\"\n"
	scriptfiledata += "export mailpass=\"" + valkconfig.MailReaderConfig.MailPassword + "\"\n"
	scriptfiledata += "export mailproto=\"" + valkconfig.MailReaderConfig.MailProtocol + "\"\n"
	scriptfiledata += "export mailusetls=\"" + valkconfig.MailReaderConfig.UseTLS + "\"\n"
	scriptfiledata += "export sleeptimeout=\"" + valkconfig.MailReaderConfig.SleepTimeout + "\"\n"
	scriptfiledata += "export imapfolder=\"INBOX\"\n"
	scriptfiledata += "export trigger_subject=\"" + valkconfig.MailReaderConfig.MailSubjectTrigger + "\"\n"

        if valkconfig.LoggerConfig.Host != "" {
                scriptfiledata += "export loggerurl=\"http://" + valkconfig.LoggerConfig.Host + ":8092/\"\n"
                scriptfiledata += "export usesecurelogging=\"" + valkconfig.LoggerConfig.UseSecureLogging + "\"\n"
                scriptfiledata += "export logkey=\"" + valkconfig.LoggerConfig.LogKey + "\"\n\n"
        } else {
                scriptfiledata += "export loggerurl=\"http://localhost:8092/\"\n"
                scriptfiledata += "export usesecurelogging=\"false\"\n"
                scriptfiledata += "export logkey=\"\"\n\n"
        }

        scriptfiledata += "export license=" + valkconfig.LicenseKey + "\n"
        scriptfiledata += "export business=" + valkconfig.BusinessName + "\n\n"

	scriptfiledata += "nohup ./mailreader &\n"
	err = ioutil.WriteFile("./valkyrie/EmailReader/startup.sh", []byte(scriptfiledata), 0755)
	if err != nil {
		return err
	}

	return nil
}


