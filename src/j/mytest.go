package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println(err)
		return
	}
	seq := 0
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go echoFunc(conn, seq)
		seq += 1
	}
}

func echoFunc(c net.Conn, seq int) {
	buf := make([]byte, 1024)

	for {
		_, err := c.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		c.Write([]byte("jjjj"))
		fmt.Println("i am", seq)
	}
}
