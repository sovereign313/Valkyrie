package install 

import (
	"fmt"
	"vtypes"
	"errors"
	"strings"
	"sftpclient"
	"sshclient"
)

func Install_Dispatcher(valkconfig vtypes.ValkConfig) error {
	host := valkconfig.DispatcherConfig.Host
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

		output, err = sftpclient.CopyFile("./deployment/dispatcher.tar", host, "/tmp/dispatcher.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Couldn't Write Remote File: " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "docker load -i /tmp/dispatcher.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Could Load Docker Image On: " + host + ": " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "rm -f /tmp/dispatcher.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Could Load Docker Image On: " + host + ": " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "docker run -d --restart=always --name dispatcher -p 8090:8090 dispatcher", 5, keyrsa)
		if err != nil {
			return errors.New("Failed To Run Dispatcher In Docker On " + host + ": " + err.Error())
		}

		fmt.Println(output)
		return nil
	}

	_, err = sshclient.RunOneCommand(host, "[[ -d /opt/valkyrie/dispatcher ]] || mkdir /opt/valkyrie/dispatcher", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Create /opt/valkyrie/dispatcher on " + host + ": " + err.Error())
	}

	_, err = sftpclient.CopyFile("./deployment/dispatcher", host, "/opt/valkyrie/dispatcher/dispatcher", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Copy ./deployment/dispatcher to " + host + ": " + err.Error())
	}

	_, err = sftpclient.CopyFile("./deployment/start_dispatcher.sh", host, "/opt/valkyrie/dispatcher/dispatcher/start.sh", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Copy ./deployment/start_dispatcher.sh to " + host + ": " + err.Error())
	}

	_, err = sshclient.RunOneCommand(host, "/opt/valkyrie/dispatcher/start.sh", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Launch Dispatcher On " + host + ": " + err.Error())
	}

	return nil
}

