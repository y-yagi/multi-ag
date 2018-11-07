package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/y-yagi/configure"
)

var (
	logger *log.Logger
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s GROUP_NAME PATTERNS\n", os.Args[0])
}

type config struct {
	Groups []Group `toml:"group"`
}

type Group struct {
	Name        string   `toml:"name"`
	Directories []string `toml:"directories"`
}

func msg(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return 1
	}
	return 0
}

func search(query string, directory string, wg *sync.WaitGroup) {
	defer wg.Done()
	out, _ := exec.Command("ag", query, directory).Output()
	if len(string(out)) > 0 {
		// NOTE: Logger is safe from multiple goroutines. Ref: https://golang.org/pkg/log/#Logger
		logger.Print(string(out))
	}
}

func cmdEdit() error {
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		editor = "vim"
	}

	return configure.Edit("multi-ag", editor)
}

func init() {
	if !configure.Exist("multi-ag") {
		var cfg config
		cfg.Groups = []Group{}
		configure.Save("multi-ag", cfg)
	}
}

func main() {
	var edit bool

	flag.BoolVar(&edit, "c", false, "edit config")
	flag.Parse()

	if edit {
		os.Exit(msg(cmdEdit()))
	}

	args := os.Args[1:]
	if len(args) < 2 {
		usage()
		os.Exit(1)
	}

	var cfg config
	err := configure.Load("multi-ag", &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var group Group
	for _, g := range cfg.Groups {
		if g.Name == args[0] {
			group = g
			break
		}
	}

	if len(group.Directories) == 0 {
		fmt.Fprintf(os.Stderr, "GROUP_NAME not found: '%v'\n", args[0])
		os.Exit(1)
	}

	logger = log.New(os.Stdout, "", 0)
	var wg sync.WaitGroup
	for _, directory := range group.Directories {
		wg.Add(1)
		go search(args[1], directory, &wg)
	}
	wg.Wait()
}
