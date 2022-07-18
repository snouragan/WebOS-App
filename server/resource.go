package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type resource struct {
	PreparedFit     [spsz]bool `json:"preparedfit"`
	PreparedStretch [spsz]bool `json:"preparedstretch"`
	InProgress      bool       `json:"inprogress"`

	// TODO:
	// thumbnail string
	// length    int
}

var resourcedb struct {
	*sync.RWMutex

	Resources map[string]*resource `json:"resources"`
}

func initResourcedb() {
	content, err := ioutil.ReadFile(filepath.Join(options.dir, "resourcedb.json"))
	if err != nil {
		log.Fatal(err.Error())
	}

	err = json.Unmarshal(content, &resourcedb)
	if err != nil {
		log.Fatal(err.Error())
	}

	if resourcedb.Resources == nil {
		resourcedb.Resources = make(map[string]*resource)
	}

	for _, v := range resourcedb.Resources {
		v.InProgress = false
	}

	resourcedb.RWMutex = &sync.RWMutex{}

	go updateDB()
}

func updateDB() {
	for {
		resourcedb.RLock()
		b, _ := json.Marshal(resourcedb)
		resourcedb.RUnlock()

		ioutil.WriteFile(filepath.Join(options.dir, "resourcedb.json.tmp"), b, 0644)

		os.Rename(filepath.Join(options.dir, "resourcedb.json.tmp"), filepath.Join(options.dir, "resourcedb.json"))

		time.Sleep(1 * time.Minute)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	logrq(r, "upload")

	errResponse := func(err error, response string, statusCode int) {
		log.Error("Upload error", zap.String("Err", err.Error()))

		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(statusCode)
		io.WriteString(w, response)
	}

	if r.Method != "POST" {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest) // TODO: do method verification everywhere
		io.WriteString(w, "POST only")

		return
	}

	err := r.ParseMultipartForm(10 << 24) // TODO: Config this

	if err != nil {
		errResponse(err, err.Error(), http.StatusBadRequest)

		return
	}

	wf, header, err := r.FormFile("file")

	if err != nil {
		errResponse(err, err.Error(), http.StatusBadRequest)

		return
	}

	defer wf.Close()

	f, err := os.Create(filepath.Join(options.dir, "raw", header.Filename))

	if err != nil {
		errResponse(err, "Internal server error", http.StatusInternalServerError)

		return
	}

	defer f.Close()

	if _, err := io.Copy(f, wf); err != nil {
		errResponse(err, "Internal server error", http.StatusInternalServerError)

		return
	}

	resourcedb.Lock()
	resourcedb.Resources[header.Filename] = &resource{}
	resourcedb.Unlock()
}

func list(w http.ResponseWriter, r *http.Request) {
	logrq(r, "list")

	resourcedb.RLock()
	defer resourcedb.RUnlock()

	b, _ := json.Marshal(resourcedb)

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func prepare(w http.ResponseWriter, r *http.Request) {
	logrq(r, "prepare")

	s := strings.Split(r.URL.Path, "/")

	fmt.Printf("%#v%#v\n", s, len(s))

	if len(s) != 6 {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)

		io.WriteString(w, "n must be a number between 1 and "+strconv.Itoa(spsz)+"\n")

		return
	}

	resource := s[3]

	n, err := strconv.Atoi(s[4])

	if err != nil || n < 1 || n > spsz {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)

		io.WriteString(w, "n must be a number between 1 and "+strconv.Itoa(spsz)+"\n")

		return
	}

	var resolution string

	switch n {
	case 1:
		resolution = "1080:1920"
	case 2:
		resolution = "2160:1920"
	case 3:
		resolution = "3240:1920"
	case 4:
		resolution = "4320:1920"
	case 5:
		resolution = "5400:1920"
	case 6:
		resolution = "6480:1920"
	default:
		panic(errors.New("spsz more than 6"))
	}

	sf := s[5]

	if sf != "stretch" && sf != "fit" {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)

		io.WriteString(w, "sf must be stretch or fit\n")

		return
	}

	resourcedb.RLock()
	if resourcedb.Resources[resource].InProgress {
		w.WriteHeader(http.StatusConflict)

		return
	}
	resourcedb.RUnlock()

	resourcedb.Lock()
	resourcedb.Resources[resource].InProgress = true
	resourcedb.Unlock()

	ffmpeg := exec.Command(filepath.Join(options.dir, "ffmpeg_script.sh"), strconv.Itoa(n), resolution, sf, resource)

	ffmpeg.Dir = options.dir

	go func() {
		err := ffmpeg.Run()

		resourcedb.Lock()

		defer func() {
			resourcedb.Resources[resource].InProgress = false

			resourcedb.Unlock()
		}()

		exitError := &exec.ExitError{}

		if errors.As(err, &exitError) {
			log.Error("Ffmpeg error", zap.Int("ExitCode", exitError.ExitCode()), zap.ByteString("Output", exitError.Stderr))

			return
		} else if err != nil {
			log.Error("Ffmpeg run error", zap.String("Err", err.Error()))

			return
		}

		if sf == "fit" {
			resourcedb.Resources[resource].PreparedFit[n-1] = true
		} else {
			resourcedb.Resources[resource].PreparedStretch[n-1] = true
		}
	}()

	w.Header().Add("Content-Type", "text/plain")
	fmt.Println("ok")
}

// TODO: func rm
