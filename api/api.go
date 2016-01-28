// Package api provide the server side Goed API
// via RPC over local socket.
// See client/ for the client implementation.
package api

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

type Api struct {
}

func (a *Api) Start() {
	r := new(GoedRpc)
	rpc.Register(r)
	rpc.HandleHTTP()
	l, err := net.Listen("unix", core.Socket)
	if err != nil {
		log.Fatalf("Socket listen error %s : \n", core.Socket, err.Error())
	}

	go func() {
		err = http.Serve(l, nil)
		if err != nil {
			panic(err)
		}
	}()
}

// Goed RPC functions holder
type GoedRpc struct{}

func (r *GoedRpc) ApiVersion(_ struct{}, version *string) error {
	*version = core.ApiVersion
	return nil
}

func (r *GoedRpc) Open(args []interface{}, _ *struct{}) error {
	actions.EdOpen(args[1].(string), -1, args[0].(string), true)
	return nil
}

func (r *GoedRpc) Edit(args []interface{}, _ *struct{}) error {
	vid := actions.EdOpen(args[1].(string), -1, args[0].(string), true)
	actions.EdRender()
	// Wait til file closed
	for {
		v := core.Ed.ViewById(vid)
		if v.Terminated() {
			return nil
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func (r *GoedRpc) ViewReload(viewId int64, _ *struct{}) error {
	actions.ViewReload(viewId)
	actions.EdRender()
	return nil
}

func (r *GoedRpc) ViewSave(viewId int64, _ *struct{}) error {
	actions.ViewSave(viewId)
	return nil
}

func (r *GoedRpc) ViewCwd(args []interface{}, _ *struct{}) error {
	actions.ViewSetWorkdir(args[0].(int64), args[1].(string))
	return nil
}

func (r *GoedRpc) ViewSrcLoc(viewId int64, srcLoc *string) error {
	v := core.Ed.ViewById(viewId)
	if v == nil {
		return fmt.Errorf("No such view : %d", viewId)
	}
	*srcLoc = v.Backend().SrcLoc()
	return nil
}
