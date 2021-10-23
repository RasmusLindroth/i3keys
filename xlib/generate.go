package xlib

import (
	"bufio"
	"bytes"
	"go/format"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/RasmusLindroth/i3keys/helpers"
)

var tmplstr = `package xlib

var KeySyms = map[string]string{#(range $k, $v := . )#
	"#( $k )#": "#( $v )#",#(end)#
}
`

//Generate generates keysyms.go, only used for making
func Generate() {
	keysyms := parseKeysymdef()

	t := template.New("")
	t = t.Delims("#(", ")#")
	t, err := t.Parse(tmplstr)
	if err != nil {
		panic(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("Can't get current dir")
	}

	path := filepath.Join(dir, "internal/xlib/keysyms.go")

	file, err := os.Create(path)

	if err != nil {
		log.Fatalln("Couldn't create keysyms.go")
	}

	var data bytes.Buffer
	t.Execute(&data, keysyms)
	content, err := format.Source(data.Bytes())

	if err != nil {
		log.Fatalln("Couldn't format the code in keysyms.go")
	}

	file.Write(content)
}

func parseKeysymdef() map[string]string {
	file, err := os.Open("/usr/include/X11/keysymdef.h")

	if err != nil {
		log.Fatalln("Can't open keysymdef.h")
	}

	reader := bufio.NewReader(file)
	var line string
	var readErr error
	var keysyms = make(map[string]string)

	for readErr != io.EOF {
		line, readErr = reader.ReadString('\n')

		parts := helpers.SplitBySpace(line)

		if len(parts) < 3 || parts[0] != "#define" {
			continue
		}

		if len(parts[1]) > 3 && parts[1][0:3] != "XK_" {
			continue
		}

		keysyms[parts[2]] = parts[1][3:]
	}

	return keysyms
}
