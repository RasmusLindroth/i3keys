package web

import (
	"html/template"
	"net/http"
	"os"
	"path"
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

func realName(pathname string) string {
	if fileinfo, err := os.Lstat(pathname); err == nil {
		if fileinfo.Mode()&os.ModeSymlink == os.ModeSymlink {
			if realname, err := os.Readlink(pathname); err == nil {
				return realname
			}
		}
	}
	return pathname
}

func readResource(filename string) (string, bool) {
	pathname := path.Join(os.Getenv("HOME"), ".config")
	if xdg_config_home, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		pathname = xdg_config_home
	}
	pathname = realName(pathname)
	pathname = realName(path.Join(pathname, "i3keys"))
	pathname = realName(path.Join(pathname, filename))

	if buf, err := os.ReadFile(pathname); err == nil {
		//println("read resource from '" + pathname + "'")
		return string(buf), true
	} else {
		//println("could not read resource from '" + pathname + "', using built-in")
		return "", false
	}
}

// New inits the handler for the web service
func New(layouts Layouts) Handler {
	handler := Handler{}

	handler.Data.Layouts = layouts
	handler.Data.LayoutMaps = layoutMaps

	return handler
}

// Start fires up the server
func (handler Handler) Start(port string) error {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.HomeHandler)

	gzip := handlers.CompressHandler(r)

	err := http.ListenAndServe(":"+port, handlers.RecoveryHandler()(gzip))

	return err
}

// HomeHandler serves root requests
func (handler *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {

	// LayoutName must be set first or the html template fails subtly
	handler.Data.LayoutName = strings.ToUpper(r.URL.Query().Get("layout"))
	if handler.Data.LayoutName == "" {
		handler.Data.LayoutName = "ISO" // because i'm selfish
	}

	if res, ok := readResource("index.gohtml"); ok {
		indexTmplHTML = res
	}
	handler.Template = template.Must(template.New("index").Funcs(template.FuncMap{
		"keyinfo": handler.keyInfo,
		"mkslice": mkSlice,
	}).Parse(indexTmplHTML))

	if res, ok := readResource("index.css"); ok {
		indexTmplCSS = res
	}
	handler.Data.CSS = template.CSS(indexTmplCSS)

	if res, ok := readResource("index.js"); ok {
		indexTmplJS = res
	}
	handler.Data.JS = template.JS(indexTmplJS)

	err := handler.Template.Execute(w, handler.Data)
	if err != nil {
		println("Handler.HomeHandler: failed to execute template:\n", err.Error())
	}
}
