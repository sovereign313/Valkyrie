package setup

import (
	"os"
	"vtypes"
	"errors"
	"strings"

	"os/exec"
	"io/ioutil"
)

func Setup_Foreman(valkconfig vtypes.ValkConfig) error {
	foremanconfig := valkconfig.ForemanConfig

	err := WriteForemanFiles(valkconfig)
	if err != nil {
		return errors.New("Failed to write Foreman files: " + err.Error())
	}

	os.Mkdir("./deployment", 0755)

	if valkconfig.UseDocker {
		os.Chdir("./valkyrie/Foreman/")
		_, err := exec.Command("/usr/bin/docker", "build", "-t", "foreman", ".").Output()
		if err != nil {
			return errors.New("Failed to build foreman: " + err.Error())
		}

		os.Chdir("../..")

		_, err = exec.Command("/usr/bin/docker", "save", "-o", "./deployment/foreman.tar", "foreman").Output()
		if err != nil {
			return errors.New("Failed to write foreman to tar: " + err.Error())
		}

		_, err = exec.Command("/usr/bin/docker", "rmi", "foreman").Output()
		if err != nil {
			return errors.New("Failed to remove foreman docker image: " + err.Error())
		}
	} else {
		err := CopyFile("./valkyrie/Foreman/foreman", "./deployment/foreman")
		if err != nil {
			return errors.New("Failed to copy foreman binary: " + err.Error())
		}

		err = CopyFile("./valkyrie/Foreman/startup.sh", "./deployment/start_foreman.sh")
		if err != nil {
			return errors.New("Failed to copy foreman startup script: " + err.Error())
		}
	}

	_ = foremanconfig
	return nil
	
}

func CopyFile(src string, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}

	return nil
}


func WriteForemanFiles(valkconfig vtypes.ValkConfig) error {
	foremanconfig := valkconfig.ForemanConfig
	dockerfiledata := "# This Can Be Set FROM scratch, but if you plan to use git, it's better to use alpine.\n"
	dockerfiledata += "# FROM scratch\n"
	dockerfiledata += "FROM alpine\n"
	dockerfiledata += "MAINTAINER Ernest E. Teem III <eteem@valkyriesoftware.io>\n"
	dockerfiledata += "ADD foreman /foreman\n\n"

	dockerfiledata += "ENV dupprotection " + foremanconfig.DupProtection  + "\n"
	dockerfiledata += "ENV dupprottime " + foremanconfig.ProtectTime + "\n\n"

	if valkconfig.LoggerConfig.Host != "" {
		dockerfiledata += "ENV loggerurl http://" + valkconfig.LoggerConfig.Host + ":8092/\n"
		dockerfiledata += "ENV usesecurelogging " + valkconfig.LoggerConfig.UseSecureLogging + "\n"
		dockerfiledata += "ENV logkey \"" + valkconfig.LoggerConfig.LogKey + "\"\n\n"
	} else {
		dockerfiledata += "ENV loggerurl http://localhost:8092/\n"
		dockerfiledata += "ENV usesecurelogging false\n"
		dockerfiledata += "ENV logkey \"\"\n\n"
	}

	nhst := ""
	for _, hst := range valkconfig.LauncherConfig.Host {
		nhst += "http://" + hst + ":8094|"
	}	
	nhst = strings.TrimSuffix(nhst, "|")

        dockerfiledata += "ENV license " + valkconfig.LicenseKey + "\n"
        dockerfiledata += "ENV business " + valkconfig.BusinessName + "\n\n"

	dockerfiledata += "ENV worker_hosts \"" + nhst + "\"\n\n"

	dockerfiledata += "# If Going Through A CNTLM Proxy\n"
	dockerfiledata += "#ENV HTTP_PROXY=\"http://localhost:3128\"\n"
	dockerfiledata += "#ENV HTTPS_PROXY=\"http://localhost:3128\"\n"
	dockerfiledata += "#ENV http_proxy=\"http://localhost:3128\"\n"
	dockerfiledata += "#ENV https_proxy=\"http://localhost:3128\"\n\n"

	dockerfiledata += "RUN apk update && apk add tzdata\n"
	dockerfiledata += "RUN cp /usr/share/zoneinfo/US/Eastern /etc/localtime\n"
	dockerfiledata += "RUN echo \"US/Eastern\" > /etc/timezone\n"
	dockerfiledata += "RUN apk del tzdata && apk add git && apk add openrc && apk add ca-certificates\n\n"
	dockerfiledata += "CMD [\"/foreman\"]\n"
	
	err := ioutil.WriteFile("./valkyrie/Foreman/Dockerfile", []byte(dockerfiledata), 0755)
	if err != nil {
		return err
	}

	scriptfiledata := "#!/bin/bash\n\n"
	scriptfiledata += "export dupprotection=\"" + foremanconfig.DupProtection + "\"\n"
	scriptfiledata += "export dupprottime=\"" + foremanconfig.ProtectTime + "\"\n\n"

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

	scriptfiledata += "export worker_hosts=\"" + nhst + "\"\n\n"
	scriptfiledata += "nohup ./foreman &\n"
	err = ioutil.WriteFile("./valkyrie/Foreman/startup.sh", []byte(scriptfiledata), 0755)
	if err != nil {
		return err
	}

	return nil
}


