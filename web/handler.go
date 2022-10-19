package web

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// ...

type Layouts map[string][]modeKeyboards
type LayoutMap [][]string
type LayoutMaps map[string]LayoutMap

/* these should become unnecessary when everything is in place */
func (this Layouts) ISO() []modeKeyboards {
	return this["ISO"]
}
func (this Layouts) ANSI() []modeKeyboards {
	return this["ANSI"]
}
func (this LayoutMaps) ISO() LayoutMap {
	return this["ISO"]
}
func (this LayoutMaps) ANSI() LayoutMap {
	return this["ANSI"]
}

// Data is sent to render the page
type Data struct {
	CSS        template.CSS
	JS         template.JS
	Layouts    Layouts
	LayoutMaps LayoutMaps
}

// Handler holds data needed for the page
type Handler struct {
	Template *template.Template
	Data     Data
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

func (h Handler) skipEmptyKeys(i, j int) int {
	// ugly and slow
	k := 0
	for ik := 0; ik <= j; ik++ {
		kmap := h.Data.LayoutMaps.ISO()[i][ik]
		if kmap != "emptySingle" && kmap != "emptySmall" && kmap != "enterDown" {
			k++
		}
	}
	return k
}

// New inits the handler for the web service
func New(layouts Layouts) Handler {
	handler := Handler{}

	if res, err := readResource("index.gohtml"); err == nil {
		indexTmplHTML = res
	} else {
		println("could not read HTML from '" + res + "', using built-in")
	}
	handler.Template = template.Must(template.New("index").Funcs(template.FuncMap{
		"skip": handler.skipEmptyKeys,
	}).Parse(indexTmplHTML))

	if res, err := readResource("index.css"); err == nil {
		indexTmplCSS = res
	} else {
		println("could not read CSS from '" + res + "', using built-in")
	}
	handler.Data.CSS = template.CSS(indexTmplCSS)

	if res, err := readResource("index.js"); err == nil {
		indexTmplJS = res
	} else {
		println("could not read JS from '" + res + "', using built-in")
	}
	handler.Data.JS = template.JS(indexTmplJS)

	handler.Data.Layouts = layouts
	handler.Data.LayoutMaps = LayoutMaps{
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
	println("Handler.HomeHandler")
	err := handler.Template.Execute(w, handler.Data)
	if err != nil {
		println("Handler.HomeHandler: failed to execute template:\n", err.Error())
	}
}
