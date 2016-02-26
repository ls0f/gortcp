package main

import (
	"fmt"
	"net"
	"os"
	// "strconv"
	// "time"
)

type ConnRelation struct {
	c1 net.Conn
	c2 net.Conn
}

type ConnChannel struct {
	id      string
	content []byte
}

var ConnMap = make(map[string]ConnRelation)
var channel1, channel2 chan ConnChannel = nil, nil

func main() {

	ConnMap = make(map[string]ConnRelation)
	channel1 = make(chan ConnChannel)
	channel2 = make(chan ConnChannel)
	server1(":12000")
	server2(":13000")
	for {
		select {
		case t := <-channel1:
			fmt.Println("recieve lan package, %i", len(t.content))
			relation, ok := ConnMap[t.id]
			fmt.Println(relation, ok)
			if ok {
				if relation.c2 != nil {
					relation.c2.Write(t.content)
					// fmt.Println(t.content)
				}
			}
		case t := <-channel2:
			fmt.Println("recieve wan package, %i", len(t.content))
			relation, ok := ConnMap[t.id]
			fmt.Println(relation, ok)
			if ok {
				if relation.c1 != nil {
					relation.c1.Write(t.content)
				}
			}
		}
	}

}

func server1(service string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go handleClient1(conn)
		}
	}()
}

func server2(service string) {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			id := "1"
			go handleClient2(conn, id)
		}
	}()
}

func handleClient1(conn net.Conn) {
	id := "1"
	ConnMap[id] = ConnRelation{c1: conn, c2: nil}
	fmt.Println("lan host %s connect ...", conn.RemoteAddr().String())
	defer conn.Close()
	defer fmt.Println("lan host %s disconnect ...", conn.RemoteAddr().String())
	defer delete(ConnMap, id)

	for {
		buf := make([]byte, 1024)
		read_len, err := conn.Read(buf)

		if err != nil {
			fmt.Println(err)
			break
		}
		if read_len == 0 {
			continue // connection already closed by client

		}
		fmt.Println(string(buf))
		channel1 <- ConnChannel{id: id, content: buf}

	}
}

func handleClient2(conn net.Conn, id string) {

	fmt.Println("wan host %s connect ...", conn.RemoteAddr().String())
	defer conn.Close()
	defer fmt.Println("wan host %s disconnect ...", conn.RemoteAddr().String())

	relation, ok := ConnMap[id]
	if !ok {
		return
	}
	relation.c2 = conn

	for {
		buf := make([]byte, 1024)
		read_len, err := conn.Read(buf)

		if err != nil {
			fmt.Println(err)
			break
		}
		if read_len == 0 {
			continue // connection already closed by client

		}
		channel2 <- ConnChannel{id: id, content: buf}

	}
	relation.c2 = nil
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
