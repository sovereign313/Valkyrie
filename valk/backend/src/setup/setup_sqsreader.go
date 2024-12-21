package setup

import (
	"os"
	"vtypes"
	"errors"

	"os/exec"
	"io/ioutil"
)

func Setup_SQSReader(valkconfig vtypes.ValkConfig) error {

	err := WriteSQSReaderFiles(valkconfig)
	if err != nil {
		return errors.New("Failed to write SQSReader files: " + err.Error())
	}

	os.Mkdir("./deployment", 0755)

	if valkconfig.UseDocker {
		os.Chdir("./valkyrie/SQSReader/")
		output, err := exec.Command("/usr/bin/docker", "build", "-t", "sqsreader", ".").Output()
		if err != nil {
			return errors.New("Failed to build SQSReader: " + err.Error() + " - " + string(output))
		}

		os.Chdir("../..")

		_, err = exec.Command("/usr/bin/docker", "save", "-o", "./deployment/sqsreader.tar", "sqsreader").Output()
		if err != nil {
			return errors.New("Failed to write SQSReader to tar: " + err.Error())
		}

		_, err = exec.Command("/usr/bin/docker", "rmi", "sqsreader").Output()
		if err != nil {
			return errors.New("Failed to remove SQSReader docker image: " + err.Error())
		}
	} else {
		err := CopyFile("./valkyrie/SQSReader/sqsreader", "./deployment/sqsreader")
		if err != nil {
			return errors.New("Failed to copy SQSReader binary: " + err.Error())
		}

		err = CopyFile("./valkyrie/SQSReader/startup.sh", "./deployment/start_sqsreader.sh")
		if err != nil {
			return errors.New("Failed to copy SQSReader startup script: " + err.Error())
		}
	}

	return nil
}

func WriteSQSReaderFiles(valkconfig vtypes.ValkConfig) error {

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
	dockerfiledata += "ADD sqsreader /sqsreader\n\n"

	dockerfiledata += "ENV sqsname " + valkconfig.AWSConfig.SQSName + "\n"
	dockerfiledata += "ENV sqsregion " + valkconfig.AWSConfig.Region + "\n"
	dockerfiledata += "ENV AWS_ACCESS_KEY_ID " + valkconfig.AWSConfig.AWSAccessKey + "\n"
	dockerfiledata += "ENV AWS_SECRET_ACCESS_KEY " + valkconfig.AWSConfig.AWSSecretKey + "\n"
	dockerfiledata += "ENV strkey " + valkconfig.AWSConfig.EncryptionKey + "\n"

	dockerfiledata += "ENV sleeptimeout " + valkconfig.SQSReaderConfig.SleepTimeout + "\n"

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
	dockerfiledata += "CMD [\"/sqsreader\"]\n"
	
	err := ioutil.WriteFile("./valkyrie/SQSReader/Dockerfile", []byte(dockerfiledata), 0755)
	if err != nil {
		return err
	}

	scriptfiledata := "#!/bin/bash\n\n"
	scriptfiledata += "export sqsname=\"" + valkconfig.AWSConfig.SQSName + "\"\n"
	scriptfiledata += "export sqsregion=\"" + valkconfig.AWSConfig.Region + "\"\n"
	scriptfiledata += "export AWS_ACCESS_KEY_ID=\"" + valkconfig.AWSConfig.AWSAccessKey + "\"\n"
	scriptfiledata += "export AWS_SECRET_ACCESS_KEY=\"" + valkconfig.AWSConfig.AWSSecretKey + "\"\n"
	scriptfiledata += "export strkey=\"" + valkconfig.AWSConfig.EncryptionKey + "\"\n"

	scriptfiledata += "export sleeptimeout=\"" + valkconfig.SQSReaderConfig.SleepTimeout + "\"\n"

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

	scriptfiledata += "nohup ./sqsreader &\n"
	err = ioutil.WriteFile("./valkyrie/SQSReader/startup.sh", []byte(scriptfiledata), 0755)
	if err != nil {
		return err
	}

	return nil
}


