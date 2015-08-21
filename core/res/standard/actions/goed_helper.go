package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// goed scripts utilites

// Returns the ViewId passed by Goed when calling the script
// Exit(1) if missing or invalid
func goedView() int64 {
	view := os.Getenv("GOED_VIEW")
	if len(view) == 0 {
		fmt.Println("GOED_VIEW missing")
		os.Exit(1)
	}
	v, err := strconv.ParseInt(view, 10, 64)
	if err != nil {
		log.Fatalf("Could not parse GOED_VIEW : %v", err)
	}
	return v
}

// Returns the InstanceId passed by Goed when calling the script
// Exit(1) if missing or invalid
func goedInstance() int64 {
	inst := os.Getenv("GOED_INSTANCE")
	if len(inst) == 0 {
		fmt.Println("GOED_INSTANCE missing")
		os.Exit(1)
	}
	i, err := strconv.ParseInt(inst, 10, 64)
	if err != nil {
		log.Fatalf("Could not parse GOED_INSTANCE : %v", err)
	}
	return i
}

func main() {
}
