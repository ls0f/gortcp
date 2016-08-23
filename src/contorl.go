package gortcp

import (
	"net"
	"os"
)

type Control struct {
	Addr string
	Auth string
	wrap *MessageWrap
}

func (c *Control) auth() {
	msg := &Message{msgType: authMessage, content: []byte(c.Auth)}
	err := c.wrap.SendOneMessage(msg)
	if err != nil {
		logger.Panic(err)
	}
}

func (c *Control) connect() net.Conn {
	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		logger.Panic(err)
	}
	c.wrap = &MessageWrap{rw: conn}
	return conn

}

func (c *Control) print() {

	rm, err := c.wrap.ReadOneMessage()
	if err != nil {
		logger.Panic(err)
	}
	os.Stdout.Write([]byte("#######################\n"))
	os.Stdout.Write(rm.content)
	os.Stdout.Write([]byte("#######################\n"))
}

func (c *Control) matchNode(id string) {
	msg := &Message{msgType: matchNodeMessage, content: []byte(id)}
	err := c.wrap.SendOneMessage(msg)
	if err != nil {
		logger.Panic(err)
	}
}

func (c *Control) exec(cmd string) {

	msg := &Message{msgType: execCmdMessage, content: []byte(cmd)}
	err := c.wrap.SendOneMessage(msg)
	if err != nil {
		logger.Panic(err)
	}
}

func (c *Control) ListNode() {
	conn := c.connect()
	defer conn.Close()
	c.auth()
	msg := &Message{msgType: listNodeMessage}
	if err := c.wrap.SendOneMessage(msg); err != nil {
		logger.Panic(err)
	}
	c.print()
}

func (c *Control) ExecCommand(id, cmd string) {

	conn := c.connect()
	defer conn.Close()
	c.auth()
	c.matchNode(id)
	c.exec(cmd)
	c.print()
}
