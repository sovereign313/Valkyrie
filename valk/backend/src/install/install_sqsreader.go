package install 

import (
	"fmt"
	"vtypes"
	"errors"
	"strings"
	"sftpclient"
	"sshclient"
)

func Install_SQSReader(valkconfig vtypes.ValkConfig) error {
	host := valkconfig.SQSReaderConfig.Host
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

		output, err = sftpclient.CopyFile("./deployment/sqsreader.tar", host, "/tmp/sqsreader.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Couldn't Write Remote File: " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "docker load -i /tmp/sqsreader.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Could Load Docker Image On: " + host + ": " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "rm -f /tmp/sqsreader.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Could Load Docker Image On: " + host + ": " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "docker run -d --restart=always --name sqsreader sqsreader", 5, keyrsa)
		if err != nil {
			return errors.New("Failed To Run SQSReader In Docker On " + host + ": " + err.Error())
		}

		fmt.Println(output)
		return nil
	}

	_, err = sshclient.RunOneCommand(host, "[[ -d /opt/valkyrie/sqsreader ]] || mkdir /opt/valkyrie/sqsreader", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Create /opt/valkyrie/sqsreader on " + host + ": " + err.Error())
	}

	_, err = sftpclient.CopyFile("./deployment/sqsreader", host, "/opt/valkyrie/sqsreader/sqsreader", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Copy ./deployment/sqsreader to " + host + ": " + err.Error())
	}

	_, err = sftpclient.CopyFile("./deployment/start_sqsreader.sh", host, "/opt/valkyrie/sqsreader/sqsreader/start.sh", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Copy ./deployment/start_sqsreader.sh to " + host + ": " + err.Error())
	}

	_, err = sshclient.RunOneCommand(host, "/opt/valkyrie/sqsreader/start.sh", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Launch SQSReader On " + host + ": " + err.Error())
	}

	return nil
}

