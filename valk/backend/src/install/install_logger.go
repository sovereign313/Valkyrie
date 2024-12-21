package install 

import (
	"vtypes"
	"errors"
	"strings"
	"sftpclient"
	"sshclient"
)

func Install_Logger(valkconfig vtypes.ValkConfig) error {
	host := valkconfig.LoggerConfig.Host
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

		output, err = sftpclient.CopyFile("./deployment/logger.tar", host, "/tmp/logger.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Couldn't Write Remote File: " + err.Error())
		}
		
		output, err = sshclient.RunOneCommand(host, "docker load -i /tmp/logger.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Could Load Docker Image On: " + host + ": " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "rm -f /tmp/logger.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Could Load Docker Image On: " + host + ": " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "docker run -d --restart=always --name logger -p 8092:8092 logger", 5, keyrsa)
		if err != nil {
			return errors.New("Failed To Run Logger In Docker On " + host + ": " + err.Error())
		}

		return nil
	}

	_, err := sshclient.RunOneCommand(host, "[[ -d /opt/valkyrie/logger ]] || mkdir /opt/valkyrie/logger", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Create /opt/valkyrie/logger on " + host + ": " + err.Error())
	}

	_, err = sftpclient.CopyFile("./deployment/logger", host, "/opt/valkyrie/logger/logger", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Copy ./deployment/logger to " + host + ": " + err.Error())
	}

	_, err = sftpclient.CopyFile("./deployment/start_logger.sh", host, "/opt/valkyrie/logger/logger/start.sh", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Copy ./deployment/start_logger.sh to " + host + ": " + err.Error())
	}

	_, err = sshclient.RunOneCommand(host, "/opt/valkyrie/logger/start.sh", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Launch Logger On " + host + ": " + err.Error())
	}

	return nil
}

