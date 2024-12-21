package setup

import (
	"os"
	"fmt"
	"vtypes"
	"errors"

	"os/exec"
	"io/ioutil"
)

func Setup_Launcher(valkconfig vtypes.ValkConfig) error {

	err := WriteLauncherFiles(valkconfig)
	if err != nil {
		return errors.New("Failed to write Launcher files: " + err.Error())
	}

	os.Mkdir("./deployment", 0755)

	if valkconfig.UseDocker {
		os.Chdir("./valkyrie/Launcher/")
		_, err := exec.Command("/usr/bin/docker", "build", "-t", "launcher", ".").Output()
		if err != nil {
			fmt.Println(os.Getwd())
			return errors.New("Failed to build Launcher: " + err.Error())
		}

		os.Chdir("../..")

		_, err = exec.Command("/usr/bin/docker", "save", "-o", "./deployment/launcher.tar", "launcher").Output()
		if err != nil {
			return errors.New("Failed to write launcher to tar: " + err.Error())
		}

		_, err = exec.Command("/usr/bin/docker", "rmi", "launcher").Output()
		if err != nil {
			return errors.New("Failed to remove launcher docker image: " + err.Error())
		}
	} else {
		err := CopyFile("./valkyrie/Launcher/launcher", "./deployment/launcher")
		if err != nil {
			return errors.New("Failed to copy launcher binary: " + err.Error())
		}

		err = CopyFile("./valkyrie/Launcher/startup.sh", "./deployment/start_launcher.sh")
		if err != nil {
			return errors.New("Failed to copy launcher startup script: " + err.Error())
		}
	}

	return nil
	
}

func WriteLauncherFiles(valkconfig vtypes.ValkConfig) error {
	dockerfiledata := "# This Can Be Set FROM scratch, but if you plan to use git, it's better to use alpine.\n"
	dockerfiledata += "# FROM scratch\n"
	dockerfiledata += "FROM alpine\n"
	dockerfiledata += "MAINTAINER Ernest E. Teem III <eteem@valkyriesoftware.io>\n"
	dockerfiledata += "ADD launcher /launcher\n\n"

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

	dockerfiledata += "ENV defaultimage \"" + valkconfig.LauncherConfig.DefaultImage + "\"\n\n"

	dockerfiledata += "# If Going Through A CNTLM Proxy\n"
	dockerfiledata += "#ENV HTTP_PROXY=\"http://localhost:3128\"\n"
	dockerfiledata += "#ENV HTTPS_PROXY=\"http://localhost:3128\"\n"
	dockerfiledata += "#ENV http_proxy=\"http://localhost:3128\"\n"
	dockerfiledata += "#ENV https_proxy=\"http://localhost:3128\"\n\n"

	dockerfiledata += "RUN apk update && apk add tzdata\n"
	dockerfiledata += "RUN cp /usr/share/zoneinfo/US/Eastern /etc/localtime\n"
	dockerfiledata += "RUN echo \"US/Eastern\" > /etc/timezone\n"
	dockerfiledata += "RUN apk del tzdata && apk add git && apk add openrc && apk add ca-certificates\n\n"
	dockerfiledata += "CMD [\"/launcher\"]\n"
	
	err := ioutil.WriteFile("./valkyrie/Launcher/Dockerfile", []byte(dockerfiledata), 0755)
	if err != nil {
		return err
	}

	scriptfiledata := "#!/bin/bash\n\n"

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

	scriptfiledata += "export defaultimage=\"" + valkconfig.LauncherConfig.DefaultImage + "\"\n\n"
	scriptfiledata += "nohup ./launcher &\n"
	err = ioutil.WriteFile("./valkyrie/Launcher/startup.sh", []byte(scriptfiledata), 0755)
	if err != nil {
		return err
	}

	return nil
}


