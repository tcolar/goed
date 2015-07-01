// Package client provides client functions to the Goed API server
// via RPC over a socket.
package client

import (
	"net/rpc"

	"github.com/tcolar/goed/core"
)

func ApiVersion(id int64) (version string, err error) {
	c := getClient(id)
	defer c.Close()
	err = c.Call("GoedRpc.ApiVersion", struct{}{}, &version)
	return version, err
}

func getClient(id int64) *rpc.Client {
	sock := core.GoedSocket(id)
	c, err := rpc.DialHTTP("unix", sock)
	if err != nil {
		panic(err)
	}
	return c
}
