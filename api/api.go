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

func (r *GoedRpc) ViewReload(viewId int64, _ *struct{}) error {
	actions.ViewReloadAction(viewId)
	actions.EdRenderAction()
	return nil
}

func (r *GoedRpc) ViewSave(viewId int64, _ *struct{}) error {
	actions.ViewSaveAction(viewId)
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
