package main

import "gopkg.in/alecthomas/kingpin.v1"

var (
	Version = "0.0.1"
	app     = kingpin.New("goed", "A code editor")
	test    = kingpin.Flag("testterm", "Pints colors to the terminal to test it.").Bool()
	colors  = kingpin.Flag("c", "Number of colors(0,2,16,256). 0 means Detect.").Default("0").Int()

	Ed Editor
)

func main() {
	kingpin.Version("0.0.1")

	kingpin.Parse()
	if *test {
		testTerm()
		return
	}
	if *colors == 0 {
		*colors = detectColors()
	}
	if *colors != 256 && *colors != 16 {
		*colors = 2
	}

	/*a := []int{0, 1, 2, 3, 4, 5, 6}
	// 0,2,3,4,1,5,6
	i1 := 1
	i2 := 4
	v := a[i1]
	copy(a[i1:i2], a[i1+1:i2+1])
	a[i2] = v
	pretty.Print(a)
	return*/

	Ed = Editor{}
	Ed.Start()
}
