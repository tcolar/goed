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

// Execute an external script, meant to be ran within a routine.
func ExecScript(script string) {
	vid := Ar.EdCurView()
	loc := core.FindResource(path.Join("actions", script))
	if _, err := os.Stat(loc); os.IsNotExist(err) {
		loc = script // assume a system wide command
	}
	env := os.Environ()
	env = append(env, fmt.Sprintf("GOED_INSTANCE=%d", core.InstanceId))
	env = append(env, fmt.Sprintf("GOED_VIEW=%d", vid))
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
		errv := Ar.EdViewsByLoc(fp)
		vid := int64(-1)
		if len(errv) > 0 {
			vid = errv[0]
		}
		Ar.EdOpen(fp, vid, "", true)
	} else {
		// no error
		errv := Ar.EdViewsByLoc(fp)
		if len(errv) > 0 {
			Ar.EdDelView(errv[0], false)
		}
	}
	Ar.EdRender()
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
