package web

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Data is sent to render the page
type Data struct {
	CSS    template.CSS
	JS     template.JS
	JSData template.JS
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

// New inits the handler for the web service
func New(js string) Handler {
	handler := Handler{}
	handler.Template = template.Must(template.New("index").Parse(indexTmplHTML))

	if res, err := readResource("index.css"); err == nil {
		indexTmplCSS = res
	} else {
		println("could not read CSS from '" + res + "', using built-in")
	}
	handler.Data.CSS = template.CSS(indexTmplCSS)

	if res, err := readResource("index.js"); err == nil {
		indexTmplCSS = res
	} else {
		println("could not read JS from '" + res + "', using built-in")
	}
	handler.Data.JS = template.JS(indexTmplCSS)

	handler.Data.JSData = template.JS(js)

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
	data := handler.Data
	handler.Template.Execute(w, data)
}
