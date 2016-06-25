package main

import (
	"fmt"
	"github.com/monder/wain/wain"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Fprintf(os.Stderr, "usage: %s [config]\n", os.Args[0])
		os.Exit(1)
	}

	config, err := wain.ReadConfig(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	server, err := wain.CreateHTTPServer(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	log.Println("Listening...")
	server.ListenAndServe()
}
