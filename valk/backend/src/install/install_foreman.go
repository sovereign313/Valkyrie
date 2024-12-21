package install 

import (
	"fmt"
	"vtypes"
	"errors"
	"strings"
	"sftpclient"
	"sshclient"
)

func Install_Foreman(valkconfig vtypes.ValkConfig) error {
	host := valkconfig.ForemanConfig.Host
	keyrsa, err := sshclient.SignerFromBytes([]byte(valkconfig.WorkerConfig.SSHPrivateKey))
	if err != nil {
		fmt.Println(err.Error())
	}

	if valkconfig.UseDocker {
		output, err := sshclient.RunOneCommand(host, "[[ -e /var/run/docker.sock ]]; echo $?", 5, keyrsa)
		if err != nil {
			return err
		}

		output = strings.TrimSpace(output)
		if output != "0" {
			return errors.New("Use Docker Is True, But Docker Not Installed On: " + host)
		}

		output, err = sftpclient.CopyFile("./deployment/foreman.tar", host, "/tmp/foreman.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Couldn't Write Remote File: " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "docker load -i /tmp/foreman.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Could Load Docker Image On: " + host + ": " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "rm -f /tmp/foreman.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Could Load Docker Image On: " + host + ": " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "docker run -d --restart=always --name foreman -p 8091:8091 foreman", 5, keyrsa)
		if err != nil {
			return errors.New("Failed To Run Foreman In Docker On " + host + ": " + err.Error())
		}

		fmt.Println(output)
		return nil
	}

	_, err = sshclient.RunOneCommand(host, "[[ -d /opt/valkyrie/foreman ]] || mkdir /opt/valkyrie/foreman", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Create /opt/valkyrie/foreman on " + host + ": " + err.Error())
	}

	_, err = sftpclient.CopyFile("./deployment/foreman", host, "/opt/valkyrie/foreman/foreman", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Copy ./deployment/foreman to " + host + ": " + err.Error())
	}

	_, err = sftpclient.CopyFile("./deployment/start_foreman.sh", host, "/opt/valkyrie/foreman/foreman/start.sh", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Copy ./deployment/start_foreman.sh to " + host + ": " + err.Error())
	}

	_, err = sshclient.RunOneCommand(host, "/opt/valkyrie/foreman/start.sh", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Launch Foreman On " + host + ": " + err.Error())
	}

	return nil
}

