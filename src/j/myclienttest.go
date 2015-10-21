package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"runtime"
	"time"
)

func ping(num int, count *int) int {
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", "localhost:8080")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	howmuch := 0
	for err != nil {
		time.Sleep(1 * time.Second)
		conn, err = net.DialTCP("tcp", nil, tcpAddr)
		if err == nil {
			break
		} else {
			howmuch += 1
		}

		if howmuch == 50 {
			ioutil.WriteFile("log.txt", []byte("i can not dial 50 "), 0644)
			fmt.Println("i can not dial after 50times")
			return -1
		}
	}
	*count += 1
	fmt.Println("i am ", num)
	for {
		_, _ = conn.Write([]byte("Ping"))
		var buff [4]byte
		_, _ = conn.Read(buff[0:])
		time.Sleep(5 * time.Second)
	}

	conn.Close()
	return 0
}

func main() {
	runtime.GOMAXPROCS(10)

	//var totalPings int = 1000000
	//var concurrentConnections int = 100
	//var pingsPerConnection int = totalPings / concurrentConnections
	//var actualTotalPings int = pingsPerConnection * concurrentConnections
	count := 0
	for i := 0; i < 10000; i++ {
		go ping(i, &count)
	}
	for {
		time.Sleep(20 * time.Second)
		fmt.Println("count is", count)
	}
	fmt.Println("count is", count)
	lockChan := make(chan bool)
	<-lockChan
}
