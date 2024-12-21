package install 

import (
	"vtypes"
	"errors"
	"strings"
	"sftpclient"
	"sshclient"
)

func Install_Alerter(valkconfig vtypes.ValkConfig) error {
	host := valkconfig.AlerterConfig.Host
	keyrsa, _ := sshclient.SignerFromBytes([]byte(valkconfig.WorkerConfig.SSHPrivateKey))

	if valkconfig.UseDocker {
		output, err := sshclient.RunOneCommand(host, "[[ -e /var/run/docker.sock ]]; echo $?", 5, keyrsa)
		if err != nil {
			return err
		}

		output = strings.TrimSpace(output)
		if output != "0" {
			return errors.New("Use Docker Is True, But Docker Not Installed On: " + host)
		}

		output, err = sftpclient.CopyFile("./deployment/alerter.tar", host, "/tmp/alerter.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Couldn't Write Remote File: " + err.Error())
		}
		
		output, err = sshclient.RunOneCommand(host, "docker load -i /tmp/alerter.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Could Load Docker Image On: " + host + ": " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "rm -f /tmp/alerter.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Could Load Docker Image On: " + host + ": " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "docker run -d --restart=always --name alerter -p 8093:8093 alerter", 5, keyrsa)
		if err != nil {
			return errors.New("Failed To Run Alerter In Docker On " + host + ": " + err.Error())
		}

		return nil
	}

	_, err := sshclient.RunOneCommand(host, "[[ -d /opt/valkyrie/alerter ]] || mkdir /opt/valkyrie/alerter", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Create /opt/valkyrie/alerter on " + host + ": " + err.Error())
	}

	_, err = sftpclient.CopyFile("./deployment/alerter", host, "/opt/valkyrie/alerter/alerter", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Copy ./deployment/alerter to " + host + ": " + err.Error())
	}

	_, err = sftpclient.CopyFile("./deployment/start_alerter.sh", host, "/opt/valkyrie/alerter/alerter/start.sh", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Copy ./deployment/start_alerter.sh to " + host + ": " + err.Error())
	}

	_, err = sshclient.RunOneCommand(host, "/opt/valkyrie/alerter/start.sh", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Launch Alerter On " + host + ": " + err.Error())
	}

	return nil
}

