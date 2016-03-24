// Package api provide the server side Goed API
// via RPC over local socket.
// See client/ for the client implementation.
package api

import (
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

func (r *GoedRpc) Open(args []interface{}, _ *struct{}) error {
	vid := actions.Ar.EdOpen(args[1].(string), -1, args[0].(string), true)
	actions.Ar.EdActivateView(vid)
	actions.Ar.EdRender()
	return nil
}

func (r *GoedRpc) Edit(args []interface{}, _ *struct{}) error {
	curView := actions.Ar.EdCurView()
	vid := actions.Ar.EdOpen(args[1].(string), -1, args[0].(string), true)
	actions.Ar.EdActivateView(vid)
	actions.Ar.EdRender()
	// Wait til file closed
	for {
		v := core.Ed.ViewById(vid)
		if v.Terminated() {
			actions.Ar.EdActivateView(curView)
			actions.Ar.EdRender()
			return nil
		}
		time.Sleep(250 * time.Millisecond)
	}
}
