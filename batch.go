package main

import (
	"os/exec"
	"strings"
)

type Batch struct {
	id       int
	running  bool
	output   []byte
	ok       bool
	requests []*ExecRequest
}

func (batch *Batch) process(cmd string) {
	batch.running = true

	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:]
	out, err := exec.Command(head, parts...).Output()

	if err != nil {
		batch.ok = false
		batch.output = []byte(err.Error())
	} else {
		batch.ok = true
		batch.output = out
	}

	for _, request := range batch.requests {
		request.setBatchResult(batch)
	}

	batch.running = false
}

func (batch *Batch) enqueueRequest(request *ExecRequest) bool {
	if !batch.running {
		batch.requests = append(batch.requests, request)
		return true
	}
	return false
}
