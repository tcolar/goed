// +build

package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/kr/pty"
)

func main() {
	c := exec.Command("bash", "-i")
	f, err := pty.Start(c)
	if err != nil {
		log.Fatal(err)
	}
	//c.Stdout = os.Stdout
	//c.Env = []string{"TERM=xterm"}
	//var i bytes.Buffer
	//c.Stdin = &i
	//c.Stdin = strings.NewReader("git commit\n")
	//c.Stderr = os.Stderr
	go c.Run()
	go copy(f)
	f.Write([]byte("ls s*\n"))
	time.Sleep(3 * time.Second)
	f.Write([]byte("top\n"))
	time.Sleep(3 * time.Second)
	f.Close()
}

func copy(f *os.File) {
	for i := 0; i != 100; i++ {
		io.Copy(os.Stdout, f)
		time.Sleep(1 * time.Millisecond)
	}
}
