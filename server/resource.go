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

	"github.com/h2non/filetype"
	"go.uber.org/zap"
)

const (
	resourceTypeVideo = iota
	resourceTypeImage
)

type resource struct {
	Title string `json:"title"`
	Type  int    `json:"type"`

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

func newResource(r *http.Request, fname string, res *resource) {
	log.Info("New Resource", zap.String("From", fmtip(r.RemoteAddr)), zap.String("Resource", fname))

	resourcedb.Lock()
	resourcedb.Resources[fname] = res
	resourcedb.Unlock()
}

var tmp uint64

func upload(w http.ResponseWriter, r *http.Request) {
	logrq(r, "upload")

	if !requireMethod(w, r, http.MethodPost) {
		return
	}

	errResponse := func(err error, response string, statusCode int) {
		log.Error("Upload error", zap.String("Err", err.Error()))

		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(statusCode)
		io.WriteString(w, response)
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

	ext := path.Ext(header.Filename)

	tmpi := atomic.AddUint64(&tmp, 1)
	tmpFname := "tmp" + strconv.FormatUint(tmpi, 10) + ext

	f, err := os.Create(filepath.Join(options.dir, "raw", tmpFname))

	if err != nil {
		errResponse(err, "Internal server error (see logs)", http.StatusInternalServerError)

		return
	}

	sha := sha256.New()

	write := io.MultiWriter(f, sha)

	if _, err := io.Copy(write, formFile); err != nil {
		errResponse(err, "Internal server error (see logs)", http.StatusInternalServerError)

		return
	}

	formFile.Close()
	f.Close()

	fname := hex.EncodeToString(sha.Sum(nil))[:7] + ext

	resourcedb.RLock()
	_, exists := resourcedb.Resources[fname]
	resourcedb.RUnlock()

	if exists {
		os.Remove(filepath.Join(options.dir, "raw", tmpFname))
		errResponse(fmt.Errorf("resource %s already exists", fname), "resource already exists", http.StatusConflict)

		return
	}

	os.Rename(filepath.Join(options.dir, "raw", tmpFname), filepath.Join(options.dir, "raw", fname))

	t, err := filetype.MatchFile(filepath.Join(options.dir, "raw", fname))

	if err != nil {
		os.Remove(filepath.Join(options.dir, "raw", fname))
		errResponse(err, "Internal server error (see logs)", http.StatusInternalServerError)

		return
	}

	fmt.Print(t.MIME.Type)

	switch t.MIME.Type {
	default:
		os.Remove(filepath.Join(options.dir, "raw", fname))

		err := errors.New("bad mime type")
		errResponse(err, err.Error(), http.StatusBadRequest)
	case "image":
		fmt.Println("\n\nere " + fname + "\n\n")
		err := os.Rename(filepath.Join(options.dir, "raw", fname), filepath.Join(options.dir, "data", fname))

		if err != nil {
			panic(err)
		}

		newResource(r, fname, &resource{Title: title, Type: resourceTypeImage, Prepared: true, Thumbnail: "/data/" + fname})
	case "video":
		n, err := strconv.Atoi(r.FormValue("n"))

		if err != nil || n < 1 || n > spsz {
			os.Remove(filepath.Join(options.dir, "raw", fname))

			err := errors.New("n must be a number between 1 and " + strconv.Itoa(spsz))
			errResponse(err, err.Error(), http.StatusBadRequest)

			return
		}

		sf := r.FormValue("sf")

		if sf != "stretch" && sf != "fit" {
			os.Remove(filepath.Join(options.dir, "raw", fname))

			err := errors.New("sf must be stretch or fit")
			errResponse(err, err.Error(), http.StatusBadRequest)

			return
		}

		ffmpeg := exec.Command(filepath.Join(options.dir, "thumbnail.sh"), fname)
		ffmpeg.Dir = options.dir

		err = ffmpeg.Run()

		if err != nil {
			os.Remove(filepath.Join(options.dir, "raw", fname))
			errResponse(err, "Internal server error (see logs)", http.StatusInternalServerError)

			exitError := &exec.ExitError{}

			if errors.As(err, &exitError) {

				log.Error("Ffmpeg error", zap.Int("ExitCode", exitError.ExitCode()), zap.ByteString("Output", exitError.Stderr))

				return
			}

			return
		}

		newResource(r, fname, &resource{Title: title, Type: resourceTypeVideo, NumMonitors: n, SF: sf, Thumbnail: "/data/" + fname + ".thumb.jpg"})

		go func() {
			ffmpeg := exec.Command(filepath.Join(options.dir, "split.sh"), strconv.Itoa(n), sf, fname)

			ffmpeg.Dir = options.dir

			err := ffmpeg.Run()

			resourcedb.Lock()
			defer resourcedb.Unlock()

			exitError := &exec.ExitError{}

			if errors.As(err, &exitError) {
				log.Error("Ffmpeg error", zap.Int("ExitCode", exitError.ExitCode()), zap.ByteString("Output", exitError.Stderr))

				return
			} else if err != nil {
				log.Error("Ffmpeg run error", zap.String("Err", err.Error()))

				return
			}

			resourcedb.Resources[fname].Prepared = true
		}()
	}

}

func rm(w http.ResponseWriter, r *http.Request) {
	logrq(r, "rm")

	if !requireMethod(w, r, http.MethodGet) {
		return
	}

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

	if !requireMethod(w, r, http.MethodGet) {
		return
	}

	resourcedb.RLock()
	defer resourcedb.RUnlock()

	b, _ := json.Marshal(resourcedb)

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
