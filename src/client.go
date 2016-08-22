package gortcp

import (
	"net"
	"time"
)

type Client struct{
	addr string;
}

func (c *Client)Connect() {

	for {
		conn, err := net.Dial("tcp", c.addr)
		if err != nil {
			logger.Errorf(err.Error())
			time.Sleep(20 * time.Second)
		}else {
			c.handlerConn(conn)
		}
	}
}


func (c *Client)handlerConn(conn net.Conn){
}
