package swag

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	http_swagger "github.com/swaggo/http-swagger"
)

var (
	default_swag = New()
)

type Swag struct {
	is_running bool
	callbacks  map[string](map[string]func(http.ResponseWriter, *http.Request)) // method: (path: callback)
}

func New() *Swag {
	return &Swag{
		callbacks: make(map[string](map[string]func(http.ResponseWriter, *http.Request))),
	}
}

func (this *Swag) Register(method, path string, callback func(http.ResponseWriter, *http.Request)) {
	if this.is_running {
		return
	}

	if _, ok := this.callbacks[method]; !ok {
		this.callbacks[method] = make(map[string]func(http.ResponseWriter, *http.Request))
	}

	this.callbacks[method][path] = callback
}

func (this *Swag) Listen(service_name, port string) {
	if this.is_running {
		return
	}

	r := chi.NewRouter()
	for method, path_callback := range this.callbacks {
		for path, callback := range path_callback {
			switch method {
			case "GET":
				r.Get(path, callback)
			case "POST":
				r.Post(path, callback)
			}
		}
	}
	r.Get("/"+service_name+"/swagger/*", http_swagger.Handler())

	go func() {
		http.ListenAndServe(":"+port, r)
	}()
	this.is_running = true
}

func Register(method, path string, callback func(http.ResponseWriter, *http.Request)) {
	default_swag.Register(method, path, callback)
}

func Listen(service_name, port string) {
	default_swag.Listen(service_name, port)
}
