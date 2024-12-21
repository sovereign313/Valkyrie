package setup

import (
	"os"
	"vtypes"
	"errors"

	"os/exec"
	"io/ioutil"
)

func Setup_Worker(valkconfig vtypes.ValkConfig) error {

	err := WriteWorkerFiles(valkconfig)
	if err != nil {
		return errors.New("Failed to write Worker files: " + err.Error())
	}

	os.Mkdir("./deployment", 0755)

	if valkconfig.UseDocker {
		os.Chdir("./valkyrie/Worker/")
		output, err := exec.Command("/usr/bin/docker", "build", "-t", "worker", ".").Output()
		if err != nil {
			return errors.New("Failed to build Worker: " + err.Error() + " - " + string(output))
		}

		os.Chdir("../..")

		_, err = exec.Command("/usr/bin/docker", "save", "-o", "./deployment/worker.tar", "worker").Output()
		if err != nil {
			return errors.New("Failed to write worker to tar: " + err.Error())
		}

		_, err = exec.Command("/usr/bin/docker", "rmi", "worker").Output()
		if err != nil {
			return errors.New("Failed to remove worker docker image: " + err.Error())
		}
	} else {
		err := CopyFile("./valkyrie/Worker/worker", "./deployment/worker")
		if err != nil {
			return errors.New("Failed to copy worker binary: " + err.Error())
		}

		err = CopyFile("./valkyrie/Worker/startup.sh", "./deployment/start_worker.sh")
		if err != nil {
			return errors.New("Failed to copy worker startup script: " + err.Error())
		}
	}

	return nil
	
}

func WriteWorkerFiles(valkconfig vtypes.ValkConfig) error {
        if _, err := os.Stat("./valkyrie/Worker/keys"); os.IsNotExist(err) {
                os.Mkdir("./valkyrie/Worker/keys", 0700)
        }

	valkconfig.WorkerConfig.SSHPrivateKey += "\n"
	valkconfig.WorkerConfig.SSHPublicKey += "\n"

	err := ioutil.WriteFile("./valkyrie/Worker/keys/id_rsa", []byte(valkconfig.WorkerConfig.SSHPrivateKey), 0600)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./valkyrie/Worker/keys/id_rsa.pub", []byte(valkconfig.WorkerConfig.SSHPublicKey), 0600)
	if err != nil {
		return err
	}
	

	dockerfiledata := "# This Can Be Set FROM scratch, but if you plan to use git, it's better to use alpine.\n"
	dockerfiledata += "# FROM scratch\n"
	dockerfiledata += "FROM alpine\n"
	dockerfiledata += "MAINTAINER Ernest E. Teem III <eteem@valkyriesoftware.io>\n"
	dockerfiledata += "ADD keys/ /root/.ssh/\n"
	dockerfiledata += "ADD worker /worker\n\n"

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

	dockerfiledata += "ENV externalpath \"" + valkconfig.WorkerConfig.ExternalPath + "\"\n\n"
	dockerfiledata += "ENV usegitrepo true\n"
	dockerfiledata += "ENV gitrepourl " + valkconfig.WorkerConfig.GitRepoURL + "\n"
	dockerfiledata += "ENV sshuser " + valkconfig.WorkerConfig.SSHUser + "\n\n"

	dockerfiledata += "# If Going Through A CNTLM Proxy\n"
	dockerfiledata += "#ENV HTTP_PROXY=\"http://localhost:3128\"\n"
	dockerfiledata += "#ENV HTTPS_PROXY=\"http://localhost:3128\"\n"
	dockerfiledata += "#ENV http_proxy=\"http://localhost:3128\"\n"
	dockerfiledata += "#ENV https_proxy=\"http://localhost:3128\"\n\n"

	dockerfiledata += "RUN apk update && apk add tzdata\n"
	dockerfiledata += "RUN cp /usr/share/zoneinfo/US/Eastern /etc/localtime\n"
	dockerfiledata += "RUN echo \"US/Eastern\" > /etc/timezone\n"
	dockerfiledata += "RUN apk del tzdata && apk add git && apk add openrc && apk add ca-certificates\n\n"
	dockerfiledata += "CMD [\"/worker\"]\n"
	
	err = ioutil.WriteFile("./valkyrie/Worker/Dockerfile", []byte(dockerfiledata), 0755)
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

	scriptfiledata += "export externalpath=\"" + valkconfig.WorkerConfig.ExternalPath + "\"\n"
	scriptfiledata += "export usegitrepo=\"true\"\n"
	scriptfiledata += "export gitrepourl=\"" + valkconfig.WorkerConfig.GitRepoURL + "\"\n"
	scriptfiledata += "export sshuser=\"" + valkconfig.WorkerConfig.SSHUser + "\"\n\n"

        scriptfiledata += "export license=" + valkconfig.LicenseKey + "\n"
        scriptfiledata += "export business=" + valkconfig.BusinessName + "\n\n"

	scriptfiledata += "nohup ./worker &\n"
	err = ioutil.WriteFile("./valkyrie/Worker/startup.sh", []byte(scriptfiledata), 0755)
	if err != nil {
		return err
	}

	return nil
}


