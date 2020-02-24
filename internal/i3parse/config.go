package i3parse

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"go.i3wm.org/i3/v4"
)

type wmType uint

const (
	wmNil wmType = iota
	wmi3
	wmSway
)

func testWM(program string) bool {
	_, err := exec.Command(program, "--get-socketpath").Output()
	return err == nil
}

func getWM() wmType {
	if testWM("i3") {
		return wmi3
	}
	if testWM("sway") {
		return wmSway
	}
	return wmNil
}

func getAutoWM() (*strings.Reader, error) {
	wm := getWM()
	switch wm {
	case wmi3:
		return getConfigFromRunningi3()
	case wmSway:
		return getConfigFromRunningSway()
	default:
		err := errors.New("couldn't determine if you're running i3 or Sway. Test the -i or -s flag")
		return nil, err
	}

}

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
