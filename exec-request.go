package main

type ExecRequest struct {
	ok        bool
	output    []byte
	processed chan bool
}

func (request *ExecRequest) setBatchResult(batch *Batch) {
	request.ok = batch.ok
	request.output = batch.output
	request.processed <- true
}
