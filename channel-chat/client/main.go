package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	go func() {
		io.Copy(os.Stdout, conn)
		log.Println("done")
		done <- true
	}()

	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	conn.Close()
	<-done
}
