package sseserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type hubStatus struct {
	Node        string             `json:"node"`
	Status      string             `json:"status"`
	Reported    int64              `json:"reported_at"`
	Connections []connectionStatus `json:"connections"`
}

// Status returns the status struct for a given server.
// This mirrors what is made available on the `/admin/status.json` endpoint.
//
// Mostly useful for reporting.
func (s *Server) Status() hubStatus {
	return s.hub.Status()
}

// Status returns the status struct for a given connection hub.
// This hub is the real source of truth and Server is just a wrapper but people
// don't know that...
func (h *hub) Status() hubStatus {

	stat := hubStatus{
		Node:     fmt.Sprintf("%s-%s-%s", platform(), env(), dyno()),
		Status:   "OK",
		Reported: time.Now().Unix(),
	}

	stat.Connections = []connectionStatus{}
	for k := range h.connections {
		stat.Connections = append(stat.Connections, k.Status())
	}

	return stat
}

func platform() string {
	return "go"
}

func dyno() string {
	dyno := os.Getenv("DYNO")
	if dyno != "" {
		return dyno
	} else {
		return "dev.1"
	}
}

func env() string {
	env := os.Getenv("GO_ENV")
	if env != "" {
		return env
	} else {
		return "development"
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request, h *hub) {
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.MarshalIndent(h.Status(), "", "  ")
	fmt.Fprint(w, string(b))
	return
}