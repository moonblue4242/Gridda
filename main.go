//
//   Copyright (C) 2019 moonblue4242@gmail.com
//

package main

import (
	"bufio"
	"log"
	"os"
)

// Build with -ldflags "-H windowsgui" to create a gui element
func main() {

	f, w := setup()
	defer teardown(f, w)
	// create business logic controller
	logic := NewLogic()
	logic.Loop()

}

func setup() (f *os.File, w *bufio.Writer) {
	f, err := os.Create("./gridda.log")
	if err == nil {
		w = bufio.NewWriter(f)
		log.SetOutput(w)
	} else {
		panic(err)
	}
	return
}

func teardown(logFile *os.File, writer *bufio.Writer) {
	writer.Flush()
	logFile.Close()
}
