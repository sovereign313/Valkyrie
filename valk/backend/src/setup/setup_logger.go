package setup

import (
	"os"
	"vtypes"
	"errors"

	"os/exec"
	"io/ioutil"
)

func Setup_Logger(valkconfig vtypes.ValkConfig) error {

	err := WriteLoggerFiles(valkconfig)
	if err != nil {
		return errors.New("Failed to write Logger files: " + err.Error())
	}

	os.Mkdir("./deployment", 0755)

	if valkconfig.UseDocker {
		os.Chdir("./valkyrie/Logger/")
		output, err := exec.Command("/usr/bin/docker", "build", "-t", "logger", ".").Output()
		if err != nil {
			return errors.New("Failed to build Logger: " + err.Error() + " - " + string(output))
		}

		os.Chdir("../..")

		_, err = exec.Command("/usr/bin/docker", "save", "-o", "./deployment/logger.tar", "logger").Output()
		if err != nil {
			return errors.New("Failed to write logger to tar: " + err.Error())
		}

		_, err = exec.Command("/usr/bin/docker", "rmi", "logger").Output()
		if err != nil {
			return errors.New("Failed to remove logger docker image: " + err.Error())
		}
	} else {
		err := CopyFile("./valkyrie/Logger/logger", "./deployment/logger")
		if err != nil {
			return errors.New("Failed to copy logger binary: " + err.Error())
		}

		err = CopyFile("./valkyrie/Logger/startup.sh", "./deployment/start_logger.sh")
		if err != nil {
			return errors.New("Failed to copy logger startup script: " + err.Error())
		}
	}

	return nil
	
}

func WriteLoggerFiles(valkconfig vtypes.ValkConfig) error {

	dockerfiledata := "# This Can Be Set FROM scratch, but if you plan to use git, it's better to use alpine.\n"
	dockerfiledata += "# FROM scratch\n"
	dockerfiledata += "FROM alpine\n"
	dockerfiledata += "MAINTAINER Ernest E. Teem III <eteem@valkyriesoftware.io>\n"
	dockerfiledata += "ADD logger /logger\n\n"

	dockerfiledata += "ENV useeventstreams " + valkconfig.LoggerConfig.UseEventStreams + "\n"

	if valkconfig.LoggerConfig.ESHostPort == "" {
		dockerfiledata += "ENV eshostport \"none\"\n"
	} else {
		dockerfiledata += "ENV eshostport " + valkconfig.LoggerConfig.ESHostPort + "\n"
	}

	dockerfiledata += "ENV usesecurelogging " + valkconfig.LoggerConfig.UseSecureLogging + "\n"
	dockerfiledata += "ENV logkey \"" + valkconfig.LoggerConfig.LogKey + "\"\n"
	dockerfiledata += "ENV logfilelocation \"" + valkconfig.LoggerConfig.LogFileLocation + "\"\n\n"

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
	dockerfiledata += "CMD [\"/logger\"]\n"
	
	err := ioutil.WriteFile("./valkyrie/Logger/Dockerfile", []byte(dockerfiledata), 0755)
	if err != nil {
		return err
	}

	scriptfiledata := "#!/bin/bash\n\n"

	scriptfiledata += "export useeventstreams=\"" + valkconfig.LoggerConfig.UseEventStreams + "\"\n"
	scriptfiledata += "export eshostport=\"" + valkconfig.LoggerConfig.ESHostPort + "\"\n"
	scriptfiledata += "export usesecurelogging=\"" + valkconfig.LoggerConfig.UseSecureLogging + "\"\n"
	scriptfiledata += "export logkey=\"" + valkconfig.LoggerConfig.LogKey + "\"\n"
	scriptfiledata += "export logfilelocation=\"http://" + valkconfig.LoggerConfig.LogFileLocation + ":8092/\"\n\n"

        scriptfiledata += "export license=" + valkconfig.LicenseKey + "\n"
        scriptfiledata += "export business=" + valkconfig.BusinessName + "\n\n"

	scriptfiledata += "nohup ./logger &\n"
	err = ioutil.WriteFile("./valkyrie/Logger/startup.sh", []byte(scriptfiledata), 0755)
	if err != nil {
		return err
	}

	return nil
}


