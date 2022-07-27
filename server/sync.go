package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

const spsz = 6

type arrivalmsg struct {
	tv   int
	ct   string
	body []byte
}

type syncplay struct {
	arrival chan arrivalmsg
	play    [spsz]chan arrivalmsg
}

func newsyncplay() *syncplay {
	s := &syncplay{arrival: make(chan arrivalmsg)}

	for i := 0; i < spsz; i++ {
		s.play[i] = make(chan arrivalmsg)
	}

	go s.syncd()

	return s
}

func (s *syncplay) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrq(r, "syncplay")

	tv := tvid(r.RemoteAddr)

	if tv == -1 {
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintln(w, "not a tv")

		return

	}

	var msg arrivalmsg

	if r.Method == "POST" {
		ct := r.Header.Get("Content-Type")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		msg.tv = tv
		msg.body = body
		msg.ct = ct
	} else {
		msg.tv = tv
		msg.body = []byte("sync")
		msg.ct = "text/plain"
	}

	s.arrival <- msg

	msg = <-s.play[tv]

	w.Header().Add("Content-Type", msg.ct)
	w.Write(msg.body)

	log.Info("Sync Response", zap.String("From", fmtip(r.RemoteAddr)), zap.String("Content-Type", msg.ct), zap.String("Response", string(msg.body)))
}

func (s *syncplay) syncd() {
	waitvar := 0

	var msgtosend arrivalmsg
	for {
		msg := <-s.arrival

		if msg.tv == 0 { // TODO: generalize
			msgtosend = msg
		}

		waitvar = waitvar | (1 << msg.tv)

		runState.RLock()
		if waitvar == runState.tvs {
			for _, c := range s.play {
				c <- msgtosend
			}

			waitvar = 0
		}
		runState.RUnlock()
	}
}
