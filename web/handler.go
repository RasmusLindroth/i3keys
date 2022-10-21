package web

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/RasmusLindroth/i3keys/keyboard"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// ...
type modeKeyboards struct {
	Name      string
	Keyboards []keyboard.Keyboard
}
type Layouts map[string][]modeKeyboards
type LayoutMap [][]string
type LayoutMaps map[string]LayoutMap

// Data is sent to render the page
type Data struct {
	CSS        template.CSS
	JS         template.JS
	Layouts    Layouts
	LayoutMaps LayoutMaps
	LayoutName string
}

// Handler holds data needed for the page
type Handler struct {
	Template *template.Template
	Data     Data
}

type KeyInfo struct {
	I, J, K   int
	Key_size  string
	Key_empty bool
	Key_usage int
	Key_enter int
	Key       keyboard.Key
}

var layoutMaps = LayoutMaps{
	"ISO": {
		{"single", "emptySingle", "single", "single", "single", "single", "emptySmall", "single", "single", "single", "single", "emptySmall", "single", "single", "single", "single", "emptySmall", "single", "single", "single"},
		{"single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "double", "emptySmall", "single", "single", "single", "emptySmall", "single", "single", "single", "single"},
		{"onehalf", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "enterUp", "emptySmall", "single", "single", "single", "emptySmall", "single", "single", "single", "doubleY"},
		{"semidouble", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "enterDown", "emptySmall", "emptySingle", "emptySingle", "emptySingle", "emptySmall", "single", "single", "single"},
		{"modifier", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "large", "emptySmall", "emptySingle", "single", "emptySingle", "emptySmall", "single", "single", "single", "doubleY"},
		{"modifier", "modifier", "modifier", "space", "modifier", "modifier", "modifier", "modifier", "emptySmall", "single", "single", "single", "emptySmall", "double", "single"},
	},
	"ANSI": {
		{"single", "emptySingle", "single", "single", "single", "single", "emptySmall", "single", "single", "single", "single", "emptySmall", "single", "single", "single", "single", "emptySmall", "single", "single", "single"},
		{"single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "double", "emptySmall", "single", "single", "single", "emptySmall", "single", "single", "single", "single"},
		{"onehalf", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "onehalf", "emptySmall", "single", "single", "single", "emptySmall", "single", "single", "single", "doubleY"},
		{"semidouble", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "semilarge", "emptySmall", "emptySingle", "emptySingle", "emptySingle", "emptySmall", "single", "single", "single"},
		{"semilarge", "single", "single", "single", "single", "single", "single", "single", "single", "single", "single", "large", "emptySmall", "emptySingle", "single", "emptySingle", "emptySmall", "single", "single", "single", "doubleY"},
		{"modifier", "modifier", "modifier", "space", "modifier", "modifier", "modifier", "modifier", "emptySmall", "single", "single", "single", "emptySmall", "double", "single"},
	},
}

func (h Handler) keyInfo(kbd keyboard.Keyboard) <-chan KeyInfo {
	//kbLayout := h.Data.Layouts[h.Data.LayoutName]
	kbLayoutMap := h.Data.LayoutMaps[h.Data.LayoutName]

	ki := make(chan KeyInfo)
	go func() {
		enterHit := 0
		for i, rowMap := range kbLayoutMap {
			k := 0
			for j, key_size := range rowMap {
				gHit := 0
				key_empty := (key_size == "emptySingle") || (key_size == "emptySmall") || (key_size == "enterDown")
				key := keyboard.Key{}
				if !key_empty {
					key = kbd.Keys[i][k]
					if key.Modifier {
						gHit = 1
					}
					if key.InUse {
						gHit = 2
					}
					if key_size == "enterUp" {
						enterHit = gHit
					}
					k++
				}
				if key_size == "enterDown" {
					gHit = enterHit
				}
				//println("keyinfo: ", i, j, k, key_size, key_empty, gHit, enterHit, key.Symbol)
				ki <- KeyInfo{i, j, k, key_size, key_empty, gHit, enterHit, key}
			}
		}
		close(ki)
	}()
	return ki
}
func mkSlice(args ...interface{}) []interface{} {
	return args
}

func readResource(filename string) (string, error) {
	i3keys_config := os.Getenv("HOME") + "/.config/i3keys/" // TODO: check XDG_CONFIG & Co., cache
	pathname := i3keys_config + filename
	pathname, _ = os.Readlink(pathname) // the blind readlink is *NOT* ok
	if buf, err := os.ReadFile(pathname); err == nil {
		return string(buf), nil
	} else {
		return pathname, err // horrible
	}
}

// New inits the handler for the web service
func New(layouts Layouts) Handler {
	println("Handler.New")
	handler := Handler{}

	handler.Data.Layouts = layouts
	handler.Data.LayoutMaps = layoutMaps

	return handler
}

// Start fires up the server
func (handler Handler) Start(port string) error {
	println("Handler.Start")
	r := mux.NewRouter()
	r.HandleFunc("/", handler.HomeHandler)

	gzip := handlers.CompressHandler(r)

	err := http.ListenAndServe(":"+port, handlers.RecoveryHandler()(gzip))

	return err
}

// HomeHandler serves root requests
func (handler *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	//println("Handler.HomeHandler")

	handler.Data.LayoutName = strings.ToUpper(r.URL.Query().Get("layout"))
	if handler.Data.LayoutName == "" {
		handler.Data.LayoutName = "ISO"
	}
	//println("Handler.HomeHandler: layout: ", handler.Data.LayoutName)

	if res, err := readResource("index.gohtml"); err == nil {
		indexTmplHTML = res
		//println("read HTML template from file")
	} else {
		println("could not read HTML from '" + res + "', using built-in")
	}
	handler.Template = template.Must(template.New("index").Funcs(template.FuncMap{
		"keyinfo": handler.keyInfo,
		"mkslice": mkSlice,
	}).Parse(indexTmplHTML))

	if res, err := readResource("index.css"); err == nil {
		indexTmplCSS = res
		//println("read CSS template from file")
	} else {
		println("could not read CSS from '" + res + "', using built-in")
	}
	handler.Data.CSS = template.CSS(indexTmplCSS)

	if res, err := readResource("index.js"); err == nil {
		indexTmplJS = res
		//println("read JS template from file")
	} else {
		println("could not read JS from '" + res + "', using built-in")
	}
	handler.Data.JS = template.JS(indexTmplJS)

	err := handler.Template.Execute(w, handler.Data)
	if err != nil {
		println("Handler.HomeHandler: failed to execute template:\n", err.Error())
	}
}
