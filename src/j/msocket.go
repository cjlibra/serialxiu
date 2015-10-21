package main

import (
	"fmt"
	"io/ioutil"
	//"bufio"
	//"fmt"
	"log"
	"net"
)

func handleConnection(conn net.Conn, connmain net.Conn) {
	for {
		data, err := ioutil.ReadAll(conn)
		fmt.Println(data)
		if err != nil {
			log.Fatal("get client data  from mutil socket error: ", err)
		}
		_, err = connmain.Write(data)
		if err != nil {
			log.Fatal("can not write data error to mainsocket: ", err)
		}

		data, err = ioutil.ReadAll(connmain)
		if err != nil {
			log.Fatal("can not read data from main socket error: ", err)
		}
		_, err = conn.Write(data)
		if err != nil {
			log.Fatal("can not write back data to mutil socket error: ", err)
		}

	}
}
func setMainsock() net.Conn {
	conn, err := net.Dial("tcp", "www.sina.com.cn:80")
	if err != nil {
		panic(err)
	}

	return conn

}

func main() {
	ln, err := net.Listen("tcp", ":6010")
	if err != nil {
		panic(err)
	}
	connmain := setMainsock()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("get client connection error: ", err)
		}

		go handleConnection(conn, connmain)
	}
}
