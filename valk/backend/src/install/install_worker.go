package install 

import (
	"vtypes"
	"errors"
	"strings"
	"sftpclient"
	"sshclient"
)

func Install_Worker(valkconfig vtypes.ValkConfig) error {
	keyrsa, _ := sshclient.SignerFromBytes([]byte(valkconfig.WorkerConfig.SSHPrivateKey))

	for _, hst := range valkconfig.LauncherConfig.Host {
		output, err := sshclient.RunOneCommand(hst, "[[ -e /var/run/docker.sock ]]; echo $?", 5, keyrsa)
		if err != nil {
			return err
		}

		output = strings.TrimSpace(output)
		if output != "0" {
			return errors.New("Use Docker Is True, But Docker Not Installed On: " + hst)
		}

		output, err = sftpclient.CopyFile("./deployment/worker.tar", hst, "/tmp/worker.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Couldn't Write Remote File: " + err.Error())
		}
		
		output, err = sshclient.RunOneCommand(hst, "docker load -i /tmp/worker.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Could Load Docker Image On: " + hst + ": " + err.Error())
		}

		output, err = sshclient.RunOneCommand(hst, "rm -f /tmp/worker.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Could Load Docker Image On: " + hst + ": " + err.Error())
		}

	}

	return nil
}

