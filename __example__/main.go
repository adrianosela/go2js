package main

import (
	"fmt"
	"log"
	"time"

	"github.com/adrianosela/go2js"
)

func main() {
	conn, err := go2js.NewJsConn(
		go2js.WithOnReadHandler("onRead"),
		go2js.WithOnWriteHandler("onWrite"),
		go2js.WithOnCloseHandler("onClose"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize new JsConn object: %v", err)
	}
	defer conn.Close()

	// mock write to JS
	go func() {
		for {
			_, err := conn.Write([]byte("Hello from Go!"))
			if err != nil {
				log.Printf("Write error: %v", err)
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// mock read from JS
	for {
		buf := make([]byte, 1024)
		i, err := conn.Read(buf)
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}
		fmt.Printf("Received: %s\n", string(buf[:i]))
		time.Sleep(1 * time.Second)
	}
}
