package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"

	"gopkg.in/yaml.v2"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s PATTERN\n", os.Args[0])
}

type Config struct {
	Directory []string `yaml:"Directory"`
}

func search(query string, directory string, wg *sync.WaitGroup) {
	defer wg.Done()
	out, _ := exec.Command("ag", query, directory).Output()
	fmt.Println(string(out))
}

func readConfigFile() Config {
	configFile := os.Getenv("HOME") + "/.multi-ag.yml"

	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	var parsedMap Config
	if err = yaml.Unmarshal(buf, &parsedMap); err != nil {
		panic(err)
	}

	return parsedMap
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		usage()
		os.Exit(2)
	}

	config := readConfigFile()

	var wg sync.WaitGroup
	for _, directory := range config.Directory {
		wg.Add(1)
		go search(args[0], directory, &wg)
	}
	wg.Wait()
}
