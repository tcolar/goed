// Package api provide the server side Goed API
// via RPC over local socket.
// See client/ for the client implementation.
package api

import (
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

type RpcStruct struct {
	Data []string
}

func (r *GoedRpc) Action(args RpcStruct, res *RpcStruct) error {
	results, err := actions.Exec(args.Data[0], args.Data[1:])
	for _, r := range results {
		res.Data = append(res.Data, r)
	}
	return err
}

/*
func (r *GoedRpc) ApiVersion(_ struct{}, version *string) error {
	*version = core.ApiVersion
	return nil
}

func (r *GoedRpc) Open(args []interface{}, _ *struct{}) error {
	actions.EdOpen(args[1].(string), -1, args[0].(string), true)
	return nil
}

func (r *GoedRpc) Edit(args []interface{}, _ *struct{}) error {
	curView := actions.EdCurView()
	vid := actions.EdOpen(args[1].(string), -1, args[0].(string), true)
	actions.EdActivateView(vid, 0, 0)
	actions.EdRender()
	// Wait til file closed
	for {
		v := core.Ed.ViewById(vid)
		if v.Terminated() {
			actions.EdActivateView(curView, 0, 0)
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

func (r *GoedRpc) ViewRows(viewId int64, rows *int) error {
	v := core.Ed.ViewById(viewId)
	if v == nil {
		return fmt.Errorf("No such view : %d", viewId)
	}
	*rows = v.LastViewLine()
	return nil
}

func (r *GoedRpc) ViewCols(viewId int64, cols *int) error {
	v := core.Ed.ViewById(viewId)
	if v == nil {
		return fmt.Errorf("No such view : %d", viewId)
	}
	*cols = v.LastViewCol()
	return nil
}

func (r *GoedRpc) ViewVtCols(args []interface{}, _ *struct{}) error {
	actions.ViewSetVtCols(args[0].(int64), args[1].(int))
	return nil
}*/
