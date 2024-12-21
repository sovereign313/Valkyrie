package install 

import (
	"vtypes"
	"errors"
	"strings"
	"sftpclient"
	"sshclient"
)

func Install_Launcher(valkconfig vtypes.ValkConfig) error {
	keyrsa, _ := sshclient.SignerFromBytes([]byte(valkconfig.WorkerConfig.SSHPrivateKey))

	if valkconfig.UseDocker {
		for _, hst := range valkconfig.LauncherConfig.Host {
			output, err := sshclient.RunOneCommand(hst, "[[ -e /var/run/docker.sock ]]; echo $?", 5, keyrsa)
			if err != nil {
				return err
			}

			output = strings.TrimSpace(output)
			if output != "0" {
				return errors.New("Use Docker Is True, But Docker Not Installed On: " + hst)
			}

			output, err = sftpclient.CopyFile("./deployment/launcher.tar", hst, "/tmp/launcher.tar", 5, keyrsa)
			if err != nil {
				return errors.New("Couldn't Write Remote File: " + err.Error())
			}
			
			output, err = sshclient.RunOneCommand(hst, "docker load -i /tmp/launcher.tar", 5, keyrsa)
			if err != nil {
				return errors.New("Could Load Docker Image On: " + hst + ": " + err.Error())
			}

			output, err = sshclient.RunOneCommand(hst, "rm -f /tmp/launcher.tar", 5, keyrsa)
			if err != nil {
				return errors.New("Could Load Docker Image On: " + hst + ": " + err.Error())
			}

			output, err = sshclient.RunOneCommand(hst, "docker run -d --restart=always --name launcher -v /var/run/docker.sock:/var/run/docker.sock -p 8094:8094 launcher", 5, keyrsa)
			if err != nil {
				return errors.New("Failed To Run Launcher In Docker On " + hst + ": " + err.Error())
			}
		}

		return nil
	}

	for _, hst := range valkconfig.LauncherConfig.Host {
		_, err := sshclient.RunOneCommand(hst, "[[ -d /opt/valkyrie/launcher ]] || mkdir /opt/valkyrie/launcher", 5, keyrsa)
		if err != nil {
			return errors.New("Failed To Create /opt/valkyrie/launcher on " + hst + ": " + err.Error())
		}

		_, err = sftpclient.CopyFile("./deployment/launcher", hst, "/opt/valkyrie/launcher/launcher", 5, keyrsa)
		if err != nil {
			return errors.New("Failed To Copy ./deployment/launcher to " + hst + ": " + err.Error())
		}

		_, err = sftpclient.CopyFile("./deployment/start_launcher.sh", hst, "/opt/valkyrie/launcher/launcher/start.sh", 5, keyrsa)
		if err != nil {
			return errors.New("Failed To Copy ./deployment/start_launcher.sh to " + hst + ": " + err.Error())
		}

		_, err = sshclient.RunOneCommand(hst, "/opt/valkyrie/launcher/start.sh", 5, keyrsa)
		if err != nil {
			return errors.New("Failed To Launch Launcher On " + hst + ": " + err.Error())
		}
	}

	return nil
}

