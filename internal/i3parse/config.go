package i3parse

import (
	"os"
	"strings"

	"go.i3wm.org/i3"
)

//getConfigFromRunning returns the i3 config file
func getConfigFromRunning() (*strings.Reader, error) {
	conf, err := i3.GetConfig()
	if err != nil {
		return nil, err
	}

	return strings.NewReader(conf.Config), nil
}

func getConfigFromFile(path string) (*os.File, error) {
	return os.Open(path)
}
