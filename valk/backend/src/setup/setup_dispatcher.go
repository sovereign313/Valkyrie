package setup

import (
	"os"
	"vtypes"
	"errors"

	"os/exec"
	"io/ioutil"
)

func Setup_Dispatcher(valkconfig vtypes.ValkConfig) error {

	err := WriteDispatcherFiles(valkconfig)
	if err != nil {
		return errors.New("Failed to write Dispatcher files: " + err.Error())
	}

	os.Mkdir("./deployment", 0755)

	if valkconfig.UseDocker {
		os.Chdir("./valkyrie/Dispatcher/")
		output, err := exec.Command("/usr/bin/docker", "build", "-t", "dispatcher", ".").Output()
		if err != nil {
			return errors.New("Failed to build Dispatcher: " + err.Error() + " - " + string(output))
		}

		os.Chdir("../..")

		_, err = exec.Command("/usr/bin/docker", "save", "-o", "./deployment/dispatcher.tar", "dispatcher").Output()
		if err != nil {
			return errors.New("Failed to write Dispatcher to tar: " + err.Error())
		}

		_, err = exec.Command("/usr/bin/docker", "rmi", "dispatcher").Output()
		if err != nil {
			return errors.New("Failed to remove Dispatcher docker image: " + err.Error())
		}
	} else {
		err := CopyFile("./valkyrie/Dispatcher/dispatcher", "./deployment/dispatcher")
		if err != nil {
			return errors.New("Failed to copy Dispatcher binary: " + err.Error())
		}

		err = CopyFile("./valkyrie/Dispatcher/startup.sh", "./deployment/start_dispatcher.sh")
		if err != nil {
			return errors.New("Failed to copy Dispatcher startup script: " + err.Error())
		}
	}

	return nil
}

func WriteDispatcherFiles(valkconfig vtypes.ValkConfig) error {

	if valkconfig.AWSConfig.Region == "" {
		return errors.New("AWS Not Configured (Region)")
	}

	if valkconfig.AWSConfig.SQSName == "" {
		return errors.New("AWS Not Configured (SQS Name)")
	}

	if valkconfig.AWSConfig.AWSAccessKey == "" {
		return errors.New("AWS Not Configured (AWS Access Key)")
	}

	if valkconfig.AWSConfig.AWSSecretKey == "" {
		return errors.New("AWS Not Configured (AWS Secret Key)")
	}

	if valkconfig.AWSConfig.EncryptionKey == "" {
		return errors.New("AWS Not Configured (Encryption Key)")
	}

	dockerfiledata := "# This Can Be Set FROM scratch, but if you plan to use git, it's better to use alpine.\n"
	dockerfiledata += "# FROM scratch\n"
	dockerfiledata += "FROM alpine\n"
	dockerfiledata += "MAINTAINER Ernest E. Teem III <eteem@valkyriesoftware.io>\n"
	dockerfiledata += "ADD dispatcher /dispatcher\n\n"

	dockerfiledata += "ENV sqsname " + valkconfig.AWSConfig.SQSName + "\n"
	dockerfiledata += "ENV sqsregion " + valkconfig.AWSConfig.Region + "\n"
	dockerfiledata += "ENV AWS_ACCESS_KEY_ID " + valkconfig.AWSConfig.AWSAccessKey + "\n"
	dockerfiledata += "ENV AWS_SECRET_ACCESS_KEY " + valkconfig.AWSConfig.AWSSecretKey + "\n"
	dockerfiledata += "ENV strkey " + valkconfig.AWSConfig.EncryptionKey + "\n"

	dockerfiledata += "ENV dupprotection " + valkconfig.DispatcherConfig.DupProtection + "\n"
	dockerfiledata += "ENV dupprottime " + valkconfig.DispatcherConfig.ProtectTime + "\n"

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
	dockerfiledata += "CMD [\"/dispatcher\"]\n"
	
	err := ioutil.WriteFile("./valkyrie/Dispatcher/Dockerfile", []byte(dockerfiledata), 0755)
	if err != nil {
		return err
	}

	scriptfiledata := "#!/bin/bash\n\n"
	scriptfiledata += "export sqsname=\"" + valkconfig.AWSConfig.SQSName + "\"\n"
	scriptfiledata += "export sqsregion=\"" + valkconfig.AWSConfig.Region + "\"\n"
	scriptfiledata += "export AWS_ACCESS_KEY_ID=\"" + valkconfig.AWSConfig.AWSAccessKey + "\"\n"
	scriptfiledata += "export AWS_SECRET_ACCESS_KEY=\"" + valkconfig.AWSConfig.AWSSecretKey + "\"\n"
	scriptfiledata += "export strkey=\"" + valkconfig.AWSConfig.EncryptionKey + "\"\n"

	scriptfiledata += "export dupprotection=\"" + valkconfig.DispatcherConfig.DupProtection + "\"\n"
	scriptfiledata += "export dupprottime=\"" + valkconfig.DispatcherConfig.ProtectTime + "\"\n"

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

	scriptfiledata += "nohup ./dispatcher &\n"
	err = ioutil.WriteFile("./valkyrie/Dispatcher/startup.sh", []byte(scriptfiledata), 0755)
	if err != nil {
		return err
	}

	return nil
}


