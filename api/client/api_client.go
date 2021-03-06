// Package client provides client functions to the Goed API server
// via RPC over a socket.
package client

import (
	"net/rpc"

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
