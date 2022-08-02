package main

import (
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var log *zap.Logger

var runState struct {
	*sync.RWMutex

	tvs int
}

func init() {
	var err error

	log, err = zap.NewDevelopment()

	if err != nil {
		panic(err)
	}
}

func logrq(r *http.Request, handler string) {
	log.Info("Request", zap.String("Resourse", r.URL.Path), zap.String("Method", r.Method), zap.String("Handler", handler), zap.String("From", fmtip(r.RemoteAddr)))
}

func split(w http.ResponseWriter, r *http.Request) {
	logrq(r, "split")

	tv := tvidstring(r.RemoteAddr)

	if tv == "0" {
		tv = "1"
	}

	res := strings.TrimPrefix(r.URL.Path, "/split")
	ext := path.Ext(res)
	res = strings.TrimSuffix(res, ext)

	res = "/sdata" + res + "." + tv + ext

	http.Redirect(w, r, res, http.StatusFound)
}

func cors(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrq(r, "cors")

		fs.ServeHTTP(w, r)
	}
}

func addAccessControlAllowOrigin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if options.debug && len(r.URL.Path) > 2 {
			if n, err := strconv.Atoi(r.URL.Path[1:2]); err != nil && n >= 1 && n <= spsz {
				r.RemoteAddr = idtv(n - 1)
			}
		}

		w.Header().Add("Access-Control-Allow-Origin", "*")

		http.DefaultServeMux.ServeHTTP(w, r)
	}
}

func runServer() {
	runState.tvs = 0b11111

	s := newsyncplay()

	http.Handle("/sync", s)

	ps := newsyncplay()

	http.Handle("/syncpause", ps)

	c := newcontroller()

	http.Handle("/poll", c)
	http.Handle("/ctl/", c)

	http.HandleFunc("/ctl/upload", upload)
	http.HandleFunc("/ctl/list", list)
	http.HandleFunc("/ctl/prepare/", prepare)
	http.HandleFunc("/ctl/rm/", rm)

	http.HandleFunc("/split/", split)

	fs := http.FileServer(http.Dir(options.dir))
	fsh := cors(fs)
	http.Handle("/data/", fsh)
	http.Handle("/sdata/", fsh)

	err := http.ListenAndServe(":8069", addAccessControlAllowOrigin())
	log.Fatal("Fatal server error", zap.String("Error", err.Error()))
}
