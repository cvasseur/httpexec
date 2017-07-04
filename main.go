package main

import (
	"net/http"

	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var currentBatch *Batch = &Batch{id: 1, running: false}
var nextBatch *Batch = &Batch{id: 2, running: false}

func runNextBatch(cmd string) {
	currentBatch = &Batch{running: false}
	if len(nextBatch.requests) > 0 {
		currentBatch = nextBatch
		nextBatch = &Batch{running: false}
		go currentBatch.process(cmd)
	}
}

func requestHandler(cmd string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		execRequest := ExecRequest{processed: make(chan bool, 1)}

		if currentBatch.enqueueRequest(&execRequest) {
			log.Print("Incomming request enqueued in currentBatch")
			currentBatch.process(cmd)
			runNextBatch(cmd)
		} else {
			log.Print("Incomming request enqueued in nextbatch")
			nextBatch.enqueueRequest(&execRequest)
			<-execRequest.processed
		}

		if execRequest.ok {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(execRequest.output)
	})
}

func portNumber(portString string) uint16 {
	var port int
	port, err := strconv.Atoi(portString)
	if err != nil {
		log.Fatal("Error parsing port number")
	}

	if port < 0 || port > 65535 {
		log.Fatal("Wrong port number, must be between 0 and 65535")
	}

	return uint16(port)
}

func main() {
	var RootCmd = &cobra.Command{
		Use: "httpexec",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}

	var cmdStart = &cobra.Command{
		Use:   "start port exec",
		Short: "Listen to port, execute exec",
		Run: func(cmd *cobra.Command, args []string) {
			port := portNumber(args[0])
			execCmd := args[1]

			mux := http.NewServeMux()
			mux.Handle("/", requestHandler(execCmd))

			addr := fmt.Sprintf(":%d", port)
			log.Printf("Server started, listening on %s", addr)
			err := http.ListenAndServe(addr, mux)
			if err != nil {
				log.Fatal("Error starting http server")
			}

		},
	}

	RootCmd.AddCommand(cmdStart)

	RootCmd.Execute()
}
