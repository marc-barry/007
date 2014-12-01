package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net"
	"net/http"
	"time"
)

const (
	// Route vars are gorilla/mux paths variables.
	FileRouteVar     = "file"
	TemplateRouteVar = "template"
	TypeRouteVar     = "type"

	WebAssetsDir = "./assets"
	TemplatesDir = "./templates"
)

func StartHTTPServer(port int) chan error {
	r := mux.NewRouter()

	// This is the asset sub-router. It routes the "/assets" path prefix.
	// Assets are found in sub-directories under /assets (i.e. css, js...)
	assetsRouter := r.PathPrefix("/assets").Methods("GET").Subrouter()
	assetsRouter.Handle("/{"+TypeRouteVar+"}/{"+FileRouteVar+"}", http.StripPrefix("/assets/", http.FileServer(http.Dir(WebAssetsDir))))

	thandler := NewTemplateHandler()
	r.Handle("/", thandler)
	r.Handle("/{"+TemplateRouteVar+"}", thandler)

	http.Handle("/", r)

	done := make(chan error)
	go withLogging(func() {
		Log.WithField("port", port).Info("HTTP listen and serve.")
		done <- http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	})
	return done
}

type TemplateArgs struct {
	Interfaces []net.Interface
}

type TemplateHandler struct {
	baseTemplate *template.Template
}

func getInterfaceIPAddressString(iface net.Interface) string {
	ip, err := getInterfaceIPAddress(iface)

	if err != nil {
		return ""
	}

	return ip.String()
}

func newTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"timeNow":                     time.Now,
		"localHostname":               GetLocalHostname,
		"getInterfaceIPAddressString": getInterfaceIPAddressString,
		"readCPUInfo":                 readCPUInfo,
	}
}

func NewTemplateHandler() *TemplateHandler {
	funcs := newTemplateFuncMap()
	return &TemplateHandler{baseTemplate: template.Must(template.New("base").Funcs(funcs).ParseGlob(TemplatesDir + "/*.html"))}
}

func (handler *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	template := vars[TemplateRouteVar]

	if template == "" {
		template = "interfaces.html"
	}

	args := &TemplateArgs{
		Interfaces: IfaceList.All(),
	}

	Log.WithField("RemoteAddr", r.RemoteAddr).WithField("url", r.URL).Info("Serving template.")

	if err := handler.baseTemplate.ExecuteTemplate(w, template, args); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
