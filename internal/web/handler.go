package web

import (
	"html/template"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//Handler holds data needed for the page
type Handler struct {
	Template *template.Template
	CSS      template.CSS
	JS       template.JS
	JSData   template.JS
}

//New inits the handler for the web service
func New(js string) Handler {
	handler := Handler{}
	handler.Template = template.Must(template.New("index").Parse(indexTmplStr))
	handler.CSS = template.CSS(indexTmplCSS)
	handler.JS = template.JS(indexTmplJS)
	handler.JSData = template.JS(js)

	return handler
}

//Start fires up the server
func (handler Handler) Start(port string) error {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.HomeHandler)

	gzip := handlers.CompressHandler(r)

	err := http.ListenAndServe(":"+port, handlers.RecoveryHandler()(gzip))

	return err
}

//Data is sent to render the page
type Data struct {
	CSS    template.CSS
	JS     template.JS
	JSData template.JS
}

//HomeHandler serves root requests
func (handler *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := Data{
		CSS:    handler.CSS,
		JS:     handler.JS,
		JSData: handler.JSData,
	}
	handler.Template.Execute(w, data)
}
