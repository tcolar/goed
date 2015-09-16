package actions

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/mattn/anko/parser"
	"github.com/mattn/anko/vm"
	"github.com/tcolar/goed/core"
)

func External(script string) {
	d(external{script: script})
}

type external struct {
	script string
}

func (a external) Run() error {
	e := core.Ed
	v := e.ViewById(e.CurViewId())
	loc := core.FindResource(path.Join("actions", a.script))
	if _, err := os.Stat(loc); os.IsNotExist(err) {
		loc = a.script // assume a system wide command
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
		errv, err = e.Open(fp, errv, "", true)
		if err != nil {
			e.SetStatusErr(err.Error())
		}
		e.Render()
		return fmt.Errorf("%s failed", a.script)
	}
	errv := e.ViewByLoc(fp)
	e.DelView(errv, true)
	return nil
}

func runAnko() {
	env := vm.NewEnv()
	scanner := new(parser.Scanner)
	scanner.Init(`foo + 1`)
	stmts, err := parser.Parse(scanner)
	if err != nil {
		log.Fatal(err)
	}
	v, err := vm.Run(stmts, env)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(v.Interface())
}
