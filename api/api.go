package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bmizerany/pat"
	"github.com/tcolar/goed/core"
)

// Goed Api
type Api struct {
}

func (a *Api) Start(port int) {
	m := pat.New()

	m.Get("/api_version", http.HandlerFunc(a.ApiVersion))

	a.handleV1(m)

	// only listen to local connections
	http.Handle("/", m)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
	if err != nil {
		panic(err)
	}
}

// GET /api_version returns the server API version
func (a *Api) ApiVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, core.ApiVersion)
}

func (a *Api) handleV1(m *pat.PatternServeMux) {

	m.Get("/v1/cur_view", http.HandlerFunc(a.CurView))
	m.Put("/v1/cur_view/:id", http.HandlerFunc(a.TODO))
	m.Get("/v1/find_view/:identifier", http.HandlerFunc(a.TODO))
	// open file ?? inject content ?? file based ? memory based ?
	m.Post("/v1/new_view", http.HandlerFunc(a.TODO))
	m.Get("/v1/version", http.HandlerFunc(a.Version))
	m.Get("/v1/views", http.HandlerFunc(a.TODO))

	m.Get("/v1/view/:view/buffer_loc", http.HandlerFunc(a.TODO))
	m.Put("/v1/view/:view/close", http.HandlerFunc(a.TODO))
	m.Get("/v1/view/:view/cursor", http.HandlerFunc(a.TODO))
	m.Put("/v1/view/:view/cursor/:row/:col", http.HandlerFunc(a.TODO))
	m.Get("/v1/view/:view/dirty", http.HandlerFunc(a.Dirty))
	m.Put("/v1/view/:view/dirty/:val", http.HandlerFunc(a.TODO))
	m.Get("/v1/view/:view/line_count", http.HandlerFunc(a.TODO))
	m.Put("/v1/view/:view/refresh", http.HandlerFunc(a.TODO)) // or reset ??
	m.Put("/v1/view/:view/save", http.HandlerFunc(a.TODO))    // save as ??
	m.Get("/v1/view/:view/selections", http.HandlerFunc(a.Selections))
	m.Put("/v1/view/:view/selections", http.HandlerFunc(a.TODO))
	m.Get("/v1/view/:view/src_loc", http.HandlerFunc(a.SrcLoc))
	m.Get("/v1/view/:view/title", http.HandlerFunc(a.Title))
	m.Put("/v1/view/:view/title", http.HandlerFunc(a.TODO))
	m.Get("/v1/view/:view/workdir", http.HandlerFunc(a.WorkDir))
	m.Put("/v1/view/:view/workdir", http.HandlerFunc(a.TODO))
	// wipe, insert, remove
	// move / resize views ??
	// execute cmdbar items ? (search, goto etc...) ?
	// wm : size etc ...
	// setstatus, setstatuserr ?
}

// GET /v1/Version returns the Goed version
func (a *Api) Version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, core.Version)
}

// GET /v1/cur_view returns the Id of the currently active view
func (a *Api) CurView(w http.ResponseWriter, r *http.Request) {
	if core.Ed.CurView() == nil {
		http.Error(w, "No active view !", 500)
		return
	}
	fmt.Fprintf(w, "%d", core.Ed.CurView().Id())
}

// GET /v1/dirty returns whether the given view is dirty(1) or not(0)
func (a *Api) Dirty(w http.ResponseWriter, r *http.Request) {
	if core.Ed.CurView() == nil {
		http.Error(w, "No active view !", 500)
		return
	}
	dirty := "0"
	if core.Ed.CurView().Dirty() {
		dirty = "1"
	}
	fmt.Fprintf(w, dirty)
}

// GET /v1/view/<viewid>/selection returns the given view selection
// in the form:
// 	Row1 Col1 Row2 Col2
// If no selection returns an empty string
func (a *Api) Selections(w http.ResponseWriter, r *http.Request) {
	v, err := a.view(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	for _, s := range *v.Selections() {
		fmt.Fprintf(w, "%s\n", s.String())
	}
}

// GET /v1/view/<viewid>/src_loc returns the source location of the given view
// Note: Might be unavailable for memory/unsaved buffers
func (a *Api) SrcLoc(w http.ResponseWriter, r *http.Request) {
	v, err := a.view(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, v.Backend().SrcLoc())
}

// GET /v1/title returns the title of the given view
func (a *Api) Title(w http.ResponseWriter, r *http.Request) {
	v, err := a.view(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, v.Title())
}

// GET /v1/view/<viewid>/workdir returns the working directory of the given view
func (a *Api) WorkDir(w http.ResponseWriter, r *http.Request) {
	v, err := a.view(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, v.WorkDir())
}

func (a *Api) view(r *http.Request) (core.Viewable, error) {
	id, err := strconv.Atoi(r.URL.Query().Get(":view"))
	if err != nil {
		return nil, err
	}
	v := core.Ed.ViewById(id)
	if v == nil {
		return nil, fmt.Errorf("No such view %d", id)
	}
	return v, nil
}

func (a *Api) TODO(w http.ResponseWriter, r *http.Request) {
}
