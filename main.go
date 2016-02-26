package main

import (
	// "bytes"
	"fmt"
	// "io"
	"net"
	"net/http"
	"os"
)

type ConnectStruct struct {
	id        string
	hostname  string
	conn1     net.Conn
	conn2     net.Conn
	cache_buf []byte
}

type MsgChannel struct {
	id  string
	msg []byte
}

var DEBUG bool = true

var ConnMap map[string]*ConnectStruct
var SocketChannel chan string
var ReadMsgChannel1 chan MsgChannel
var ReadMsgChannel2 chan MsgChannel

func usage() {
	doc := "usage: rtcp l|c addr1 [addr2] [addr3](addr2、addr3 is need when mode is listen]"
	fmt.Println(doc)
}

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}

	mode := os.Args[1]
	addr1 := os.Args[2]

	if mode == "l" || mode == "L" {
		if len(os.Args) < 5 {
			usage()
			os.Exit(1)
		}
		addr2 := os.Args[3]
		addr3 := os.Args[4]
		ConnMap = make(map[string]*ConnectStruct)
		SocketChannel = make(chan string)
		ReadMsgChannel1 = make(chan MsgChannel)
		ReadMsgChannel2 = make(chan MsgChannel)

		go HttpServer(addr3)
		go Server1(addr1)
		go Server2(addr2)
		for {
			select {
			case rec := <-ReadMsgChannel1:

				conn_struct, ok := ConnMap[rec.id]
				if DEBUG {
					fmt.Println("recieve msg from channel1", conn_struct, ok)
				}
				if ok && conn_struct.conn2 != nil {
					n, err := conn_struct.conn2.Write(rec.msg)
					if DEBUG {
						fmt.Println("channel1:", n, err)
					}

				} else if ok {
					if conn_struct.hostname == "" {
						conn_struct.hostname = string(rec.msg)
					} else {
						conn_struct.cache_buf = rec.msg
					}
				} else {

				}
			case rec := <-ReadMsgChannel2:
				conn_struct, ok := ConnMap[rec.id]
				if DEBUG {
					fmt.Println("recieve msg from channel2", conn_struct, ok)
				}
				if ok {
					n, err := conn_struct.conn1.Write(rec.msg)
					if DEBUG {
						fmt.Println("channel2:", n, err)
					}

				}

			}
		}
	} else if mode == "c" || mode == "c" {
		Connect(addr1)
	} else {
		usage()
		os.Exit(1)
	}

}

func Server1(addr string) {

	server, err := net.Listen("tcp", addr)
	if err != nil {
		fatal("cannot listen: %s", err)
	}

	for {
		conn, err := server.Accept()
		id := conn.RemoteAddr().String()

		if err == nil {
			fmt.Println(conn.RemoteAddr().String(), "connect server1")
		}
		go func() {
			ConnMap[id] = &ConnectStruct{id: id, conn1: conn, conn2: nil, cache_buf: nil, hostname: ""}
			defer conn.Close()
			for {
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					break
				}
				ReadMsgChannel1 <- MsgChannel{id: id, msg: buf[:n]}
			}
			delete(ConnMap, id)

		}()
	}

}

func Server2(addr string) {
	server, err := net.Listen("tcp", addr)
	if err != nil {
		fatal("cannot listen: %s", err)
	}

	for {
		if DEBUG {
			fmt.Println("server2 accept a connection before notice me i can.")
		}
		id := <-SocketChannel
		if DEBUG {
			fmt.Println("receive id:", id)
		}
		conn_struct, ok := ConnMap[id]
		if !ok {
			fmt.Println("id not exists")
			continue
		}
		if conn_struct.conn2 != nil {
			fmt.Println("id has a connection.....")
			continue
		}
		if DEBUG {
			fmt.Println("will accept a connection")
		}
		conn, err := server.Accept()
		if err == nil {
			if DEBUG {
				fmt.Println(conn.RemoteAddr().String(), "connect server2")
			}
			conn_struct, ok := ConnMap[id]
			if !ok {
				conn.Close()
				fmt.Println("id not exists")
				continue
			}
			conn_struct.conn2 = conn
			go func() {
				defer conn.Close()
				if conn_struct.cache_buf != nil {
					conn.Write(conn_struct.cache_buf)
				}
				// conn_struct.cache_buf = nil
				for {
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						break
					}

					ReadMsgChannel2 <- MsgChannel{id: id, msg: buf[:n]}
				}
				// may be crash...
				// conn_struct.conn2 = nil
			}()
		}
	}

}

func HttpServer(addr string) {
	http.HandleFunc("/list", ListConn)
	http.HandleFunc("/conn", ConnectConn)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fatal("cannot listen: %s", err)
	}
}

func ListConn(w http.ResponseWriter, req *http.Request) {
	var html string = ""
	for _, v := range ConnMap {

		html += (v.hostname + " " + v.conn1.RemoteAddr().String())
		if v.conn2 != nil {
			html += (" ---> " + v.conn2.RemoteAddr().String())
		}
		html += "\n"
	}
	if html == "" {
		html = "no host connect.\n"
	}
	fmt.Fprintf(w, html)
}

func ConnectConn(w http.ResponseWriter, req *http.Request) {

	id := req.URL.Query().Get("id")
	addr := req.URL.Query().Get("addr")
	if id == "" {
		fmt.Fprintf(w, "need id arg.\n")
		return
	}
	if addr == "" {
		fmt.Fprintf(w, "need addr arg.\n")
		return
	}
	conn_struct, ok := ConnMap[id]
	if !ok {
		fmt.Fprintf(w, "id not exists.\n")
		return
	}
	conn_struct.conn1.Write([]byte(addr))
	SocketChannel <- id
	fmt.Fprintf(w, "success.\n")

}

func Connect(addr1 string) {

	for {
		conn1, err := net.Dial("tcp", addr1)
		if err != nil {
			fatal("remote dial failed: %v", err)
		}
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}
		//upload hostname
		conn1.Write([]byte(hostname))
		// read the addr2
		buf := make([]byte, 100)
		n, err := conn1.Read(buf)
		if err != nil {
			fmt.Println("remote read error")
			conn1.Close()
			break
		}
		addr2 := string(buf[:n])
		conn2, err := net.Dial("tcp", addr2)
		if err != nil {
			fatal("remote dial failed: %v", err)
		}
		finish := make(chan bool, 1)
		forward(conn1, conn2, finish)
		if DEBUG {
			fmt.Println("local connect finish")
		}
		<-finish
		conn1.Close()
		conn2.Close()
		if DEBUG {
			fmt.Println("local will rebuid connection")
		}
	}

}

func forward(conn1, conn2 net.Conn, ch chan bool) {
	go _forward(conn1, conn2, ch)
	go _forward(conn2, conn1, ch)

}

func _forward(src, dst net.Conn, ch chan bool) {

	for {
		data := make([]byte, 1024)
		n, err := src.Read(data)
		if err != nil {
			if DEBUG {
				fmt.Println(err)
			}
			ch <- true
			break
		}
		_, err = dst.Write(data[:n])
		if err != nil {
			if DEBUG {
				fmt.Println(err)
			}
			ch <- true
			break
		}
	}
}

func fatal(s string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "netfwd: %s\n", fmt.Sprintf(s, a))
	os.Exit(2)
}
