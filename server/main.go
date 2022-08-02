package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/go-units"
	"github.com/h2non/filetype"
)

//go:embed split.sh thumbnail.sh resourcedb.empty.json
var embedFS embed.FS

var options struct {
	dir           string
	debug         bool
	maxUploadSize int64
}

var maxUploadSizeString string

func init() {
	flag.StringVar(&options.dir, "dir", "", "Server directory")
	flag.BoolVar(&options.debug, "debug", false, "Debug mode")
	flag.StringVar(&maxUploadSizeString, "maxuploadsize", "200MB", "Max upload file size")
}

func createDirIfNotExist() {
	dir, err := os.Open(options.dir)
	if err == nil {
		dir.Close()
		return
	}

	if !errors.Is(err, os.ErrNotExist) {
		panic(fmt.Errorf("failed trying to open working directory: %w", err))
	}

	err = os.MkdirAll(options.dir, 0755)

	if err != nil {
		panic(fmt.Errorf("failed trying to create working directory: %w", err))
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

	ffmpegScript, _ = embedFS.ReadFile("thumbnail.sh")

	err = ioutil.WriteFile(filepath.Join(options.dir, "thumbnail.sh"), ffmpegScript, 0755)

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

	var err error
	options.maxUploadSize, err = units.FromHumanSize(maxUploadSizeString)

	if err != nil {
		os.Stderr.WriteString(fmt.Errorf("maxuploadsize not a size: %w", err).Error() + "\n")
		return
	}

	if options.dir == "" {
		os.Stderr.WriteString("Must specify a working directory (-dir)" + "\n")
		os.Exit(1)
	}

	if !filepath.IsAbs(options.dir) {
		os.Stderr.WriteString("Working directory must be absolute path" + "\n")
		os.Exit(1)
	}

	createDirIfNotExist()

	if err := os.Chdir(options.dir); err != nil {
		os.Stderr.WriteString(fmt.Errorf("failed to change into working directory: %w", err).Error() + "\n")
		os.Exit(1)
	}

	a, b := filetype.MatchFile("/home/grffn/repos/WebOS-App/server/dir/data/38ea1517085ce4233b8251d343b70e5d1ba1d568746bcb86dbc8551dcff4ac5d.webm.thumb.jpg")
	fmt.Printf("%#v %#v\n", a, b)

	initResourcedb()

	runServer()
}
