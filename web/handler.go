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

// New inits the handler for the web service
func New(js string) Handler {
	handler := Handler{}
	handler.Template = template.Must(template.New("index").Parse(indexTmplStr))

	i3keys_config := os.Getenv("HOME") + "/.config/i3keys/" // TODO: check XDG_CONFIG & Co.
	index_css := i3keys_config + "index.css"
	index_css, _ = os.Readlink(index_css) // is the blind readlink ok? i think not...
	if buf, err := os.ReadFile(index_css); err == nil {
		indexTmplCSS = string(buf) // assigning to indexTmplCSS is also not good, will do for now...
	} else {
		println("could not read ", index_css, err)
	}

	handler.Data.CSS = template.CSS(indexTmplCSS)
	handler.Data.JS = template.JS(indexTmplJS)
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
