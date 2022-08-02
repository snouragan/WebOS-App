package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"

	"go.uber.org/zap"
)

type resource struct {
	Title string `json:"title"`

	InProgress  bool   `json:"inprogress"`
	Prepared    bool   `json:"prepared"`
	NumMonitors int    `json:"nmonitors,omitempty"`
	SF          string `json:"sf,omitempty"`

	Thumbnail string `json:"thumbnail"`
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

func validTitle(title string) bool {
	return true

	if title == "" {
		return false
	}

	for _, r := range title {
		if !unicode.IsOneOf([]*unicode.RangeTable{unicode.P, unicode.L, unicode.M, unicode.S}, r) && r != ' ' {
			return false
		}
	}

	return true
}

var tmp uint64

func upload(w http.ResponseWriter, r *http.Request) {
	logrq(r, "upload")

	log.Info("here1")

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

	err := r.ParseMultipartForm(options.maxUploadSize)

	if err != nil {
		errResponse(err, err.Error(), http.StatusBadRequest)

		return
	}

	title := r.FormValue("title")

	if !validTitle(title) {
		err := errors.New("invalid title")
		errResponse(err, err.Error(), http.StatusBadRequest)

		return
	}

	formFile, header, err := r.FormFile("file")

	if err != nil {
		errResponse(err, err.Error(), http.StatusBadRequest)

		return
	}

	defer formFile.Close()

	ext := path.Ext(header.Filename)

	tmpi := atomic.AddUint64(&tmp, 1)
	tmpFname := "tmp" + strconv.FormatUint(tmpi, 10) + ext

	f, err := os.Create(filepath.Join(options.dir, "raw", tmpFname))

	if err != nil {
		errResponse(err, "Internal server error", http.StatusInternalServerError)

		return
	}

	sha := sha256.New()

	write := io.MultiWriter(f, sha)

	if _, err := io.Copy(write, formFile); err != nil {
		errResponse(err, "Internal server error", http.StatusInternalServerError)

		return
	}

	f.Close()

	fname := hex.EncodeToString(sha.Sum(nil)) + ext

	os.Rename(filepath.Join(options.dir, "raw", tmpFname), filepath.Join(options.dir, "raw", fname))

	ffmpeg := exec.Command(filepath.Join(options.dir, "thumbnail.sh"), fname)
	ffmpeg.Dir = options.dir

	err = ffmpeg.Run()

	if err != nil {
		//TODO: Things
		exitError := &exec.ExitError{}

		if errors.As(err, &exitError) {
			log.Error("Ffmpeg error", zap.Int("ExitCode", exitError.ExitCode()), zap.ByteString("Output", exitError.Stderr))

			return
		}

		fmt.Println(err.Error())

		return
	}

	resourcedb.Lock()
	resourcedb.Resources[fname] = &resource{Title: title, Thumbnail: "/data/" + fname + ".thumb.jpg"}
	resourcedb.Unlock()
}

func prepare(w http.ResponseWriter, r *http.Request) {
	logrq(r, "prepare")

	q := r.URL.Query()

	resource := path.Base(r.URL.Path)

	n, err := strconv.Atoi(q.Get("n"))

	if err != nil || n < 1 || n > spsz {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)

		io.WriteString(w, "n must be a number between 1 and "+strconv.Itoa(spsz)+"\n")

		return
	}

	sf := q.Get("sf")

	if sf != "stretch" && sf != "fit" {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)

		io.WriteString(w, "sf must be stretch or fit\n")

		return
	}

	resourcedb.RLock()
	if resourcedb.Resources[resource].InProgress {
		w.WriteHeader(http.StatusConflict)

		resourcedb.RUnlock()
		return
	}
	resourcedb.RUnlock()

	resourcedb.Lock()
	resourcedb.Resources[resource].InProgress = true
	resourcedb.Unlock()

	ffmpeg := exec.Command(filepath.Join(options.dir, "split.sh"), strconv.Itoa(n), sf, resource)

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

		resourcedb.Resources[resource].Prepared = true
		resourcedb.Resources[resource].SF = sf
		resourcedb.Resources[resource].NumMonitors = n
	}()

	w.Header().Add("Content-Type", "text/plain")
	fmt.Println("ok")
}

func rm(w http.ResponseWriter, r *http.Request) {
	resource := path.Base(r.URL.Path)

	if _, ok := resourcedb.Resources[resource]; !ok {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))

		return
	}

	ext := path.Ext(resource)
	base := strings.TrimSuffix(resource, ext)

	resourcedb.Lock()
	defer resourcedb.Unlock()
	if resourcedb.Resources[resource].Prepared {
		os.Remove(filepath.Join(options.dir, "processed", resource))

		for i := 1; i <= resourcedb.Resources[resource].NumMonitors; i++ {
			os.Remove(filepath.Join(options.dir, "sdata", base+"."+strconv.Itoa(i)+ext))
		}
	}

	os.Remove(filepath.Join(options.dir, "raw", resource))

	delete(resourcedb.Resources, resource)
}

func list(w http.ResponseWriter, r *http.Request) {
	logrq(r, "list")

	resourcedb.RLock()
	defer resourcedb.RUnlock()

	b, _ := json.Marshal(resourcedb)

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
