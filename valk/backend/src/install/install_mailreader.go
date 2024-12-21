package install 

import (
	"fmt"
	"vtypes"
	"errors"
	"strings"
	"sftpclient"
	"sshclient"
)

func Install_MailReader(valkconfig vtypes.ValkConfig) error {
	host := valkconfig.MailReaderConfig.Host
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

		output, err = sftpclient.CopyFile("./deployment/mailreader.tar", host, "/tmp/mailreader.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Couldn't Write Remote File: " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "docker load -i /tmp/mailreader.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Could Load Docker Image On: " + host + ": " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "rm -f /tmp/mailreader.tar", 5, keyrsa)
		if err != nil {
			return errors.New("Could Load Docker Image On: " + host + ": " + err.Error())
		}

		output, err = sshclient.RunOneCommand(host, "docker run -d --restart=always --name mailreader mailreader", 5, keyrsa)
		if err != nil {
			return errors.New("Failed To Run Mailreader In Docker On " + host + ": " + err.Error())
		}

		fmt.Println(output)
		return nil
	}

	_, err = sshclient.RunOneCommand(host, "[[ -d /opt/valkyrie/mailreader ]] || mkdir /opt/valkyrie/mailreader", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Create /opt/valkyrie/mailreader on " + host + ": " + err.Error())
	}

	_, err = sftpclient.CopyFile("./deployment/mailreader", host, "/opt/valkyrie/mailreader/mailreader", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Copy ./deployment/mailreader to " + host + ": " + err.Error())
	}

	_, err = sftpclient.CopyFile("./deployment/start_mailreader.sh", host, "/opt/valkyrie/mailreader/mailreader/start.sh", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Copy ./deployment/start_mailreader.sh to " + host + ": " + err.Error())
	}

	_, err = sshclient.RunOneCommand(host, "/opt/valkyrie/mailreader/start.sh", 5, keyrsa)
	if err != nil {
		return errors.New("Failed To Launch Mailreader On " + host + ": " + err.Error())
	}

	return nil
}

