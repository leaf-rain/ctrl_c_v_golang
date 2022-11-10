package putFile

import (
	"os"
	"os/exec"
	"path"
)

func PutGoLang(name, content string) error {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		if os.IsNotExist(err) {
			if _, err = os.Stat(path.Dir(name)); err != nil {
				if os.IsNotExist(err) {
					_ = os.MkdirAll(path.Dir(name), 0755)
				}
			}
			f, err = os.Create(name)
		}

	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	cmd := exec.Command("gofmt", "-w", f.Name())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err = cmd.Run(); err != nil {
		return err
	}
	return nil
}

func PutFile(name, content string) error {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		if os.IsNotExist(err) {
			if _, err = os.Stat(path.Dir(name)); err != nil {
				if os.IsNotExist(err) {
					_ = os.MkdirAll(path.Dir(name), 0755)
				}
			}
			f, err = os.Create(name)
		}

	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
