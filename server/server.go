package main

import (
	"net/http"
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
)

var log *zap.Logger

func init() {
	var err error

	log, err = zap.NewDevelopment()

	if err != nil {
		panic(err)
	}
}

func logrq(r *http.Request) {
	log.Info("Request", zap.String("Resourse", r.URL.Path), zap.String("From", fmtip(r.RemoteAddr)))
}

func splitdata(w http.ResponseWriter, r *http.Request) {
	logrq(r)

	tv := tvidstring(r.RemoteAddr)

	if tv == "0" {
		tv = "1"
	}

	res := strings.TrimPrefix(r.URL.Path, "/sdata")
	ext := path.Ext(res)
	res = strings.TrimSuffix(res, ext)

	res = "/data" + res + tv + ext

	w.Header().Add("Access-Control-Allow-Origin", "*")
	http.Redirect(w, r, res, http.StatusFound)
}

func cors(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrq(r)

		w.Header().Add("Access-Control-Allow-Origin", "*")
		fs.ServeHTTP(w, r)
	}
}

func main() {
	s := newsyncplay()

	go s.syncd()

	http.Handle("/sync", s)
	http.Handle("/sync/", s)

	home, _ := os.UserHomeDir()

	fs := http.FileServer(http.Dir(home + "/httpdata"))
	http.Handle("/data/", cors(fs))

	http.HandleFunc("/sdata/", splitdata)

	log.Fatal("Fatal server error", zap.String("Error", http.ListenAndServe(":8069", nil).Error()))
}
