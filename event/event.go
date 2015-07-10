package event

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/tcolar/goed/core"
)

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
