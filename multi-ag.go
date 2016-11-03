package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	logger *log.Logger
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s PATTERN\n", os.Args[0])
}

// Config type
type Config struct {
	Directory []string `yaml:"Directory"`
}

func search(query string, directory string, wg *sync.WaitGroup) {
	defer wg.Done()
	out, _ := exec.Command("ag", query, directory).Output()
	logger.Println(string(out))
}

func readConfigFile() (Config, error) {
	var config Config
	configFile := os.Getenv("HOME") + "/.multi-ag.yml"

	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, err
	}

	if err = yaml.Unmarshal(buf, &config); err != nil {
		return config, err
	}

	return config, nil
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	logger = log.New(os.Stdout, "", 0)
	config, err := readConfigFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for _, directory := range config.Directory {
		wg.Add(1)
		go search(args[0], directory, &wg)
	}
	wg.Wait()
}
