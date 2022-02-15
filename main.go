package main

import (
	"github.com/jwmatthews/case_watcher/cmd"
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("watcher.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	cmd.Execute()
}
