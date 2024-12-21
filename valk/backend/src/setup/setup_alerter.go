package setup

import (
	"os"
	"vtypes"
	"errors"

	"os/exec"
	"io/ioutil"
)

func Setup_Alerter(valkconfig vtypes.ValkConfig) error {

	err := WriteAlerterFiles(valkconfig)
	if err != nil {
		return errors.New("Failed to write Alerter files: " + err.Error())
	}

	os.Mkdir("./deployment", 0755)

	if valkconfig.UseDocker {
		os.Chdir("./valkyrie/Alerter/")
		output, err := exec.Command("/usr/bin/docker", "build", "-t", "alerter", ".").Output()
		if err != nil {
			return errors.New("Failed to build Alerter: " + err.Error() + " - " + string(output))
		}

		os.Chdir("../..")

		_, err = exec.Command("/usr/bin/docker", "save", "-o", "./deployment/alerter.tar", "alerter").Output()
		if err != nil {
			return errors.New("Failed to write alerter to tar: " + err.Error())
		}

		_, err = exec.Command("/usr/bin/docker", "rmi", "alerter").Output()
		if err != nil {
			return errors.New("Failed to remove alerter docker image: " + err.Error())
		}
	} else {
		err := CopyFile("./valkyrie/Alerter/alerter", "./deployment/alerter")
		if err != nil {
			return errors.New("Failed to copy alerter binary: " + err.Error())
		}

		err = CopyFile("./valkyrie/Alerter/startup.sh", "./deployment/start_alerter.sh")
		if err != nil {
			return errors.New("Failed to copy alerter startup script: " + err.Error())
		}
	}

	return nil
	
}

func WriteAlerterFiles(valkconfig vtypes.ValkConfig) error {

	dockerfiledata := "# This Can Be Set FROM scratch, but if you plan to use git, it's better to use alpine.\n"
	dockerfiledata += "# FROM scratch\n"
	dockerfiledata += "FROM alpine\n"
	dockerfiledata += "MAINTAINER Ernest E. Teem III <eteem@valkyriesoftware.io>\n"
	dockerfiledata += "ADD alerter /alerter\n\n"

	dockerfiledata += "ENV logfilelocation /tmp/valkyrie.log\n"

	if valkconfig.AlerterConfig.FromAddress != "" {
		dockerfiledata += "ENV fromemail " + valkconfig.AlerterConfig.FromAddress + "\n"
	}

	if valkconfig.AlerterConfig.EmailServer != "" {
		dockerfiledata += "ENV mailserver " + valkconfig.AlerterConfig.EmailServer + "\n\n"
	}

	if valkconfig.AlerterConfig.TwilioAccount != "" {
		dockerfiledata += "ENV twilio_account " + valkconfig.AlerterConfig.TwilioAccount + "\n"
	}

	if valkconfig.AlerterConfig.TwilioToken != "" {
		dockerfiledata += "ENV twilio_token " + valkconfig.AlerterConfig.TwilioToken + "\n"
	}

	if valkconfig.AlerterConfig.TwilioPhoneNumber != "" {
		dockerfiledata += "ENV twilio_account " + valkconfig.AlerterConfig.TwilioPhoneNumber + "\n"
	}

	if valkconfig.AWSConfig.Region != "" {
		dockerfiledata += "ENV awsregion " + valkconfig.AWSConfig.Region + "\n"
	}

	if valkconfig.AWSConfig.AWSAccessKey != "" {
		dockerfiledata += "ENV aws_access_key_id " + valkconfig.AWSConfig.AWSAccessKey + "\n"
	}

	if valkconfig.AWSConfig.AWSSecretKey != "" {
		dockerfiledata += "ENV aws_secret_key_id " + valkconfig.AWSConfig.AWSSecretKey + "\n"
	}

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
	dockerfiledata += "CMD [\"/alerter\"]\n"
	
	err := ioutil.WriteFile("./valkyrie/Alerter/Dockerfile", []byte(dockerfiledata), 0755)
	if err != nil {
		return err
	}

	scriptfiledata := "#!/bin/bash\n\n"
	scriptfiledata += "export logfilelocation=\"/tmp/valkyrie.log\"\n"

	if valkconfig.AlerterConfig.FromAddress != "" {
		scriptfiledata += "export fromemail=\"" + valkconfig.AlerterConfig.FromAddress + "\"\n"
	}

	if valkconfig.AlerterConfig.EmailServer != "" {
		scriptfiledata += "export mailserver=\"" + valkconfig.AlerterConfig.EmailServer + "\"\n\n"
	}

	if valkconfig.AlerterConfig.TwilioAccount != "" {
		scriptfiledata += "export twilio_account=\"" + valkconfig.AlerterConfig.TwilioAccount + "\"\n"
	}

	if valkconfig.AlerterConfig.TwilioToken != "" {
		scriptfiledata += "export twilio_token=\"" + valkconfig.AlerterConfig.TwilioToken + "\"\n"
	}

	if valkconfig.AlerterConfig.TwilioPhoneNumber != "" {
		scriptfiledata += "export twilio_account=\"" + valkconfig.AlerterConfig.TwilioPhoneNumber + "\"\n"
	}

	if valkconfig.AWSConfig.Region != "" {
		scriptfiledata += "export awsregion=\"" + valkconfig.AWSConfig.Region + "\"\n"
	}

	if valkconfig.AWSConfig.AWSAccessKey != "" {
		scriptfiledata += "export aws_access_key_id=\"" + valkconfig.AWSConfig.AWSAccessKey + "\"\n"
	}

	if valkconfig.AWSConfig.AWSSecretKey != "" {
		scriptfiledata += "export aws_secret_key_id=\"" + valkconfig.AWSConfig.AWSSecretKey + "\"\n"
	}

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

	scriptfiledata += "nohup ./alerter &\n"
	err = ioutil.WriteFile("./valkyrie/Alerter/startup.sh", []byte(scriptfiledata), 0755)
	if err != nil {
		return err
	}

	return nil
}


