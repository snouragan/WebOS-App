package main

import (
	"encoding/json"
	"net/http"
	"path"
)

type command struct {
	Command string `json:"command"`
	Src     string `json:"src,omitempty"`
	Pos     string `json:"pos,omitempty"`
}

var nocomm []byte

func init() {
	nocomm, _ = json.Marshal(&command{Command: "none"})
}

type controller struct {
	command [spsz]chan []byte // JSONed command
}

func newcontroller() *controller {
	c := &controller{}

	for i := 0; i < spsz; i++ {
		c.command[i] = make(chan []byte, 10)
	}

	return c
}

func (c *controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrq(r, "controller")

	tv := tvid(r.RemoteAddr)

	switch r.URL.Path {
	case "/poll":
		if tv == -1 {
			tv = 0
		}

		var comm []byte

		select {
		case comm = <-c.command[tv]:
		default:
			comm = nocomm
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(comm)

	case "/ctl/pause":
		c.ctl(&command{Command: "pause"})
	case "/ctl/play":
		c.ctl(&command{Command: "play"})
	}

	if path.Dir(r.URL.Path) == "/ctl/seek" {
		c.ctl(&command{Command: "seek", Pos: path.Base(r.URL.Path)})
	}

	// TODO: standardise video selection/prepare etc.
	if path.Dir(r.URL.Path) == "/ctl/load" {
		c.ctl(&command{Command: "load", Src: "/split/" + path.Base(r.URL.Path)})
	}
}

func (c *controller) ctl(cmd *command) {
	comm, err := json.Marshal(cmd)

	if err != nil {
		panic(err)
	}

	for i := 0; i < spsz; i++ {
		c.command[i] <- comm
	}
}
