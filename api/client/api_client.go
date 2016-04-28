// Package client provides client functions to the Goed API server
// via RPC over a socket.
package client

import (
	"fmt"
	"net/rpc"
	"os"
	"path"

	"github.com/tcolar/goed/api"
	"github.com/tcolar/goed/core"
)

func Action(instanceId int64, strs []string) ([]string, error) {
	c := getClient(instanceId)
	defer c.Close()
	args := api.RpcStruct{Data: strs}
	results := api.RpcStruct{}
	err := c.Call("GoedRpc.Action", args, &results)
	return results.Data, err
}

func Open(instance int64, cwd, loc string) error {
	c := getClient(instance)
	defer c.Close()
	err := c.Call("GoedRpc.Open", []interface{}{cwd, loc}, &struct{}{})
	return err
}

func Edit(instance int64, cwd, loc string) error {
	c := getClient(instance)
	defer c.Close()
	err := c.Call("GoedRpc.Edit", []interface{}{cwd, loc}, &struct{}{})
	return err
}

func getClient(id int64) *rpc.Client {
	sock := core.GoedSocket(id)
	c, err := rpc.DialHTTP("unix", sock)
	if err != nil {
		panic(err)
	}
	return c
}

func GoedSocket(id int64) string {
	p := path.Join(core.GoedHome(), "instances", fmt.Sprintf("%d.sock", id))
	if _, err := os.Stat(p); os.IsNotExist(err) {
		p = path.Join(core.GoedHome()+"_test", "instances", fmt.Sprintf("%d.sock", id))
	}
	return p
}
