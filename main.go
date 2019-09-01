package main

import (
	"aviasearch/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	runServer()

}
func runServer() {
	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Println("Stop service")
		os.Exit(1)
	}()

	http.HandleFunc("/tickets", api.GetTicketsEndpoint)
	http.HandleFunc("/ticket", api.GetTicketEndpoint)

	log.Println("Start http server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
