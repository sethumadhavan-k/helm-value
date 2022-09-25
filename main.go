package main

import (
	"flag"
	"fmt"
	"os/exec"

	"github.com/sethumadhavan-k/helm-values/cmd"
)

var searchPtr = flag.String("search", "", "search for properties ex: -search image")

func main() {

	flag.Parse()
	args := flag.Args()
	fmt.Println("args", args)
	_, isCommandAvailable := exec.Command("helm").Output()
	if isCommandAvailable != nil {
		fmt.Println("helm is not installed ")
		return
	}
	if len(args) == 0 {
		fmt.Println("repo name is missing")
		return
	}
	repo := args[0]
	cmd.App(repo, *searchPtr)
}
