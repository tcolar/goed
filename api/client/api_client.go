// Package client provides client functions to the Goed API server
// via RPC over a socket.
package client

import (
	"fmt"
	"net/rpc"
	"path"

	"github.com/tcolar/goed/core"
)

func ApiVersion(instance int64) (version string, err error) {
	c := getClient(instance)
	defer c.Close()
	err = c.Call("GoedRpc.ApiVersion", struct{}{}, &version)
	return version, err
}

func ViewReload(instance, viewId int64) (err error) {
	c := getClient(instance)
	defer c.Close()
	err = c.Call("GoedRpc.ViewReload", viewId, &struct{}{})
	return err
}

func ViewCwd(instance, viewId int64, loc string) (err error) {
	c := getClient(instance)
	defer c.Close()
	err = c.Call("GoedRpc.ViewCwd", []interface{}{viewId, loc}, &struct{}{})
	return err
}

func ViewSave(instance, viewId int64) (err error) {
	c := getClient(instance)
	defer c.Close()
	err = c.Call("GoedRpc.ViewSave", viewId, &struct{}{})
	return err
}

func ViewSrcLoc(instance, viewId int64) (loc string, err error) {
	c := getClient(instance)
	defer c.Close()
	err = c.Call("GoedRpc.ViewSrcLoc", viewId, &loc)
	return loc, err
}

func ViewRows(instance, viewId int64) (rows int, err error) {
	c := getClient(instance)
	defer c.Close()
	err = c.Call("GoedRpc.ViewRows", viewId, &rows)
	return rows, err
}

func ViewCols(instance, viewId int64) (cols int, err error) {
	c := getClient(instance)
	defer c.Close()
	err = c.Call("GoedRpc.ViewCols", viewId, &cols)
	return cols, err
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
	return path.Join(core.GoedHome(), "instances", fmt.Sprintf("%d.sock", id))
}
