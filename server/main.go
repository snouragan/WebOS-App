package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//go:embed ffmpeg_script.sh resourcedb.empty.json
var embedFS embed.FS

var options struct {
	dir string
}

func init() {
	flag.StringVar(&options.dir, "dir", "", "Server directory")
}

func createDirIfNotExist() {
	dir, err := os.Open(options.dir)
	if err == nil {
		dir.Close()
		return
	}

	if !errors.Is(err, os.ErrNotExist) {
		panic(fmt.Errorf("trying to open working directory: %w", err))
	}

	err = os.MkdirAll(options.dir, 0755)

	if err != nil {
		panic(fmt.Errorf("trying to create working directory: %w", err))
	}

	createdir := func(s string) {
		err = os.Mkdir(filepath.Join(options.dir, s), 0755)

		if err != nil {
			panic(fmt.Errorf("trying to create working directory: %w", err))
		}
	}

	createdir("data")
	createdir("sdata")
	createdir("raw")
	createdir("processed")

	ffmpegScript, _ := embedFS.ReadFile("ffmpeg_script.sh")

	err = ioutil.WriteFile(filepath.Join(options.dir, "ffmpeg_script.sh"), ffmpegScript, 0755)

	if err != nil {
		panic(fmt.Errorf("trying to create working directory: %w", err))
	}

	emptyRemotedb, _ := embedFS.ReadFile("resourcedb.empty.json")

	err = ioutil.WriteFile(filepath.Join(options.dir, "resourcedb.json"), emptyRemotedb, 0644)

	if err != nil {
		panic(fmt.Errorf("trying to create working directory: %w", err))
	}
}

func main() {
	flag.Parse()

	if options.dir == "" {
		os.Stdout.WriteString("Must specify a working directory (-dir)")
		os.Exit(1)
	}

	createDirIfNotExist()

	initResourcedb()

	runServer()
}
