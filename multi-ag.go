package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s PATTERN\n", os.Args[0])
}

type Config struct {
	Directory []string `yaml:"Directory"`
}

func search(query string, directories []string) string {
	var result []byte
	for _, directory := range directories {
		out, _ := exec.Command("ag", query, directory).Output()
		result = append(result, out...)
	}
	return string(result)
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
	result := search(args[0], config.Directory)
	fmt.Println(result)
}
