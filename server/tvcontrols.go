package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

const spsz = 6

type ctlCommand struct {
	Command string `json:"command"`

	Src     string `json:"src,omitempty"`     // for load command
	Tvs     []int  `json:"tvs,omitempty"`     // for load command -- from webapp
	SyncID  string `json:"playid,omitempty"`  // for load coomand -- to tv
	PauseID string `json:"pauseid,omitempty"` // for load command -- to tv

	Pos string `json:"pos,omitempty"` // for seek command
}

var nocomm []byte

func init() {
	nocomm, _ = json.Marshal(&ctlCommand{Command: "none"})
}

var idGenerator = make(chan string)

func initIDGenerator() {
	go func() {
		i := 0

		for {
			s := strconv.Itoa(i)

			idGenerator <- s

			i++
		}
	}()
}

var tvArray struct {
	*sync.RWMutex

	playSessions [spsz]*playSession
	main         *playSession

	commands chan *ctlCommand
}

func init() {
	tvArray.RWMutex = &sync.RWMutex{}
	tvArray.commands = make(chan *ctlCommand)
}

func tvArrayServeCtl(w http.ResponseWriter, r *http.Request) {
	logrq(r, "tvArrayServeCtl")

	if !requireMethod(w, r, http.MethodGet) {
		return
	}

	switch r.URL.Path {
	case "/ctl/pause":
		tvArray.RLock()
		defer tvArray.RUnlock()

		tvArrayExecCommand(&ctlCommand{Command: "pause"})
	case "/ctl/play":
		tvArray.RLock()
		defer tvArray.RUnlock()

		tvArrayExecCommand(&ctlCommand{Command: "play"})

	case "/ctl/seek":
		tvArray.RLock()
		defer tvArray.RUnlock()

		tvArrayExecCommand(&ctlCommand{Command: "seek", Pos: r.URL.Query().Get("pos")})

	case "/ctl/load":
		tvArray.Lock()
		defer tvArray.Unlock()

		var tvs []int

		for _, tv := range r.URL.Query()["tv"] {
			tvi, err := strconv.Atoi(tv)

			if err != nil {
				return
			}

			tvs = append(tvs, tvi-1)
		}

		tvArrayExecCommand(&ctlCommand{Command: "load", Tvs: tvs, Src: r.URL.Query().Get("src")})
	}
}

func tvArrayExecCommand(cmd *ctlCommand) {
	switch cmd.Command {
	case "play", "pause", "seek":
		if tvArray.main != nil {
			tvArray.main.muxCommand(cmd)
		}
	case "load":
		resourcedb.RLock()
		resource := resourcedb.Resources[cmd.Src]
		resourcedb.RUnlock()

		if resource == nil {
			return // TODO: Proper handling
		}

		for _, tv := range cmd.Tvs {
			if ps := tvArray.playSessions[tv]; ps != nil {
				ps.die <- struct{}{}

				for _, i := range ps.tvs.toSlice() {
					tvArray.playSessions[i] = nil
				}
			}
		}

		if resource.Type == resourceTypeVideo {
			if tvArray.main != nil {
				tvArray.main.die <- struct{}{}

				for _, i := range tvArray.main.tvs.toSlice() {
					tvArray.playSessions[i] = nil
				}
			}
		}

		ps := newPlaySession(cmd)

		for _, i := range cmd.Tvs {
			tvArray.playSessions[i] = ps
		}

		if resource.Type == resourceTypeVideo {
			tvArray.main = ps
		}
	}
}

func tvArrayServePoll(w http.ResponseWriter, r *http.Request) {
	logrq(r, "tvArrayServePoll")

	if !requireMethod(w, r, http.MethodGet) {
		return
	}

	tvArray.RLock()
	defer tvArray.RUnlock()

	tv := tvid(r.RemoteAddr)

	if tvArray.playSessions[tv] != nil {
		tvArray.playSessions[tv].poll(tv, w, r)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.Write(nocomm)
	}
}

func tvArrayServeSync(w http.ResponseWriter, r *http.Request) {
	logrq(r, "tvArrayServeSync")

	if !requireMethod(w, r, http.MethodGet) {
		return
	}

	tvArray.RLock()
	defer tvArray.RUnlock()

	tv := tvid(r.RemoteAddr)

	if tv == -1 {
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintln(w, "not a tv")

		return

	}

	if tvArray.playSessions[tv] == nil {
		return
	}

	var msg syncArrivalMsg

	if r.Method == "POST" {
		ct := r.Header.Get("Content-Type")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		msg.msg.body = body
		msg.msg.contentType = ct
	} else {
		msg.msg = defaultSyncMsg
	}

	msg.tv = tv
	msg.id = r.URL.Query().Get("id")
	msg.returnChan = make(chan syncMsg)

	tvArray.playSessions[tv].syncArrival <- msg

	var rcvMsg syncMsg

	timeoutChan := make(chan syncMsg)

	// in certain cases, syncd will refuse to respond (e.g. syncID in pause state)
	// in those cases, the tv should not be expecting a response at all
	// but we do not want to be leaving a dead goroutine here, so we send a timeout response
	go func() {
		time.Sleep(20 * time.Minute)
		timeoutChan <- timeoutSyncMsg
	}()

	select {
	case rcvMsg = <-msg.returnChan:
	case rcvMsg = <-timeoutChan:
	}

	w.Header().Add("Content-Type", rcvMsg.contentType)
	w.Write(rcvMsg.body)

	log.Info("Sync Response", zap.String("From", fmtip(r.RemoteAddr)), zap.String("Content-Type", rcvMsg.contentType), zap.String("Response", string(rcvMsg.body)))
}

type playSession struct {
	tvs tvCollection

	syncArrival chan syncArrivalMsg
	die         chan struct{}

	syncID  string
	pauseID string

	commandBuffer [spsz]chan []byte // Marshalled ctlCommand
}

func newPlaySession(load *ctlCommand) *playSession {
	ps := &playSession{
		tvs: newTVCollection(load.Tvs),

		syncArrival: make(chan syncArrivalMsg),
		die:         make(chan struct{}, spsz), // see load for why this is buffered

		syncID:  <-idGenerator,
		pauseID: <-idGenerator,
	}

	for _, tv := range load.Tvs {
		ps.commandBuffer[tv] = make(chan []byte, 3)
	}

	load.SyncID = ps.syncID
	load.PauseID = ps.pauseID
	load.Tvs = nil

	ps.muxCommand(load)

	go ps.syncd()

	return ps
}

func (ps *playSession) muxCommand(cmd *ctlCommand) {
	cdata, err := json.Marshal(cmd)
	if err != nil {
		panic(err)
	}

	for _, tv := range ps.tvs.toSlice() {
		ps.commandBuffer[tv] <- cdata
	}
}

func (ps *playSession) poll(tv int, w http.ResponseWriter, r *http.Request) {
	var comm []byte

	select {
	case comm = <-ps.commandBuffer[tv]:
		log.Info("Poll response", zap.ByteString("Command", comm), zap.String("To", fmtip(r.RemoteAddr)))
	default:
		comm = nocomm
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(comm)
}

func (ps *playSession) syncd() {
	waitvar := tvCollection(0)

	var msgToSend syncMsg
	var msg syncArrivalMsg

	returnChans := []chan syncMsg{}

	state := "sync"

	for {
		select {
		case msg = <-ps.syncArrival:
		case <-ps.die:
			return
		}

		returnChans = append(returnChans, msg.returnChan)

		if msg.tv == ps.tvs.masterTV() {
			msgToSend = msg.msg
		}

		switch {
		case state == "sync" && msg.id == ps.pauseID:
			state = "pause"

			waitvar = 0

			fallthrough
		case state == "sync" && msg.id == ps.syncID,
			state == "pause" && msg.id == ps.pauseID:

			waitvar = waitvar | (1 << msg.tv)

			if msg.tv == ps.tvs.masterTV() {
				msgToSend = msg.msg
			}

		case state == "pause" && msg.id == ps.syncID:
			// do nothing
		}

		if waitvar == ps.tvs {
			for _, c := range returnChans {
				c <- msgToSend
			}

			waitvar = 0
			state = "sync"
			returnChans = []chan syncMsg{}
		}
	}
}

type syncMsg struct {
	contentType string
	body        []byte
}

var (
	defaultSyncMsg = syncMsg{"text/plain", []byte("sync")}
	timeoutSyncMsg = syncMsg{"text/plain", []byte("sync timeout")}
)

type syncArrivalMsg struct {
	tv  int
	msg syncMsg

	id         string
	returnChan chan syncMsg
}

type tvCollection int

func newTVCollection(tvs []int) tvCollection {
	tvc := 0

	for _, tv := range tvs {
		tvc |= (1 << tv)
	}

	return tvCollection(tvc)
}

func (tvc tvCollection) masterTV() int {
	j := 1

	for i := 0; i < spsz; i++ {
		if j&int(tvc) != 0 {
			return i
		}

		j <<= 1
	}

	panic("Empty tv collection")
}

func (tvc tvCollection) toSlice() []int {
	j := 1
	sl := []int{}

	for i := 0; i < spsz; i++ {
		if j&int(tvc) != 0 {
			sl = append(sl, i)
		}

		j <<= 1
	}

	return sl
}
