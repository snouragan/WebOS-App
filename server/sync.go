package main

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"go.uber.org/zap"
)

const spsz = 6

type syncplay struct {
	arrival chan int
	play    [spsz]chan struct{}
}

func newsyncplay() *syncplay {
	s := &syncplay{arrival: make(chan int)}

	for i := 0; i < spsz; i++ {
		s.play[i] = make(chan struct{})
	}

	go s.syncd()

	return s
}

func (s *syncplay) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrq(r, "syncplay")

	tv := tvid(r.RemoteAddr)

	if tv == -1 {
		var err error

		tv, err = strconv.Atoi(path.Base(r.URL.Path))

		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			fmt.Fprintln(w, "not a tv")

			return
		}
	}

	s.arrival <- tv

	<-s.play[tv]

	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintln(w, "sync")

	log.Info("Sync Response", zap.String("From", fmtip(r.RemoteAddr)))
}

func (s *syncplay) syncd() {
	waitvar := 0
	for {
		tv := <-s.arrival

		waitvar = waitvar | (1 << tv)

		if waitvar == 0b111111 {
			for _, c := range s.play {
				c <- struct{}{}
			}

			waitvar = 0
		}
	}
}
