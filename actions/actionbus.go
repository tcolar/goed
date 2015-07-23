package actions

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/tcolar/goed/core"
)

type actionBus struct {
	actionChan chan (core.Action)
	quitc      chan (struct{})
}

func NewActionBus() core.ActionDispatcher {
	return actionBus{
		actionChan: make(chan (core.Action), 1000),
		quitc:      make(chan (struct{})),
	}
}

func (e actionBus) Dispatch(action core.Action) {
	e.actionChan <- action
}

func (e actionBus) Start() {
	for {
		select {
		case action := <-e.actionChan:
			err := action.Run()
			if err != nil {
				core.Ed.SetStatusErr(err.Error())
				log.Println(err.Error())
			}
		case <-e.quitc:
			break
		}
	}
}

func (e actionBus) Shutdown() {
	e.quitc <- struct{}{}
}

func RunAction(name string) error {
	e := core.Ed
	v := e.CurView()
	loc := core.FindResource(path.Join("actions", name))
	if _, err := os.Stat(loc); os.IsNotExist(err) {
		return fmt.Errorf("Action not found : %s", name)
	}
	env := os.Environ()
	env = append(env, fmt.Sprintf("GOED_INSTANCE=%d", core.InstanceId))
	env = append(env, fmt.Sprintf("GOED_VIEW=%d", v.Id()))
	cmd := exec.Command(loc)
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	fp := path.Join(core.Home, "errors.txt")
	if err != nil {
		file, _ := os.Create(fp)
		file.Write([]byte(err.Error()))
		file.Write([]byte{'\n'})
		file.Write(out)
		file.Close()
		errv := e.ViewByLoc(fp)
		errv, err = e.Open(fp, errv, "Errors")
		if err != nil {
			e.SetStatusErr(err.Error())
		}
		return fmt.Errorf("%s failed", name)
	}
	errv := e.ViewByLoc(fp)
	if errv != nil {
		e.DelView(errv, true)
	}
	return nil
}
