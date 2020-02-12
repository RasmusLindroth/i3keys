package i3parse

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"go.i3wm.org/i3/v4"
)

//getConfigFromRunningi3 returns the i3 config file
func getConfigFromRunningi3() (*strings.Reader, error) {
	conf, err := i3.GetConfig()
	if err != nil {
		return nil, err
	}

	return strings.NewReader(conf.Config), nil
}

//getConfigFromRunningSway returns the sway config file
func getConfigFromRunningSway() (*strings.Reader, error) {
	//Discard log saying sway is unsupported
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)
	i3.SocketPathHook = func() (string, error) {
		out, err := exec.Command("sway", "--get-socketpath").CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("getting sway socketpath: %v (output: %s)", err, out)
		}
		return string(out), nil
	}
	conf, err := i3.GetConfig()
	if err != nil {
		return nil, err
	}
	return strings.NewReader(conf.Config), nil
}

func getConfigFromFile(path string) (*os.File, error) {
	return os.Open(path)
}
