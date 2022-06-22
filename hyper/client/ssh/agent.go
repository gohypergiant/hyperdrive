package ssh

import (
	"fmt"
	"os/exec"
)

func AddKeySshAgent(keyPath string) error {
	bin, err := exec.LookPath("ssh-add")
	if err != nil {
		return err
	}

	cmd := exec.Command(bin, keyPath)
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	return nil
}
func RemoveKeySshAgent(keyPath string) error {
	bin, err := exec.LookPath("ssh-add")
	if err != nil {
		return err
	}

	cmd := exec.Command(bin, fmt.Sprintf("-d %s", keyPath))
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	return nil
}
