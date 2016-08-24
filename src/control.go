package gortcp

import (
	"fmt"
	"io"
	"log"
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
	//os.Stdout.Write([]byte("#######################\n"))
	os.Stdout.Write(rm.content)
	//os.Stdout.Write([]byte("\n"))
	//os.Stdout.Write([]byte("#######################\n"))
}

func (c *Control) matchNode(id string) {
	msg := &Message{msgType: matchNodeMessage, content: []byte(id)}
	err := c.wrap.SendOneMessage(msg)
	if err != nil {
		logger.Panic(err)
	}
	m, err := c.wrap.ReadOneMessage()
	if err != nil {
		log.Panic(err)
	}
	if m.msgType != matchOKMessage {
		os.Stdout.Write(m.content)
		os.Exit(1)
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

func (c *Control) checkFile(file string) {

	s, err := os.Stat(file)
	if err != nil {
		logger.Panic(err)
	}
	if s.IsDir() {
		logger.Panic("%s is not a file", file)
	}
}

func (c *Control) upload(srcPath, dstPath string) {
	f, err := os.Open(srcPath)
	if err != nil {
		logger.Panic(err)
	}
	defer f.Close()

	m := &Message{msgType: fileInfoMessage, content: []byte(dstPath)}
	if err = c.wrap.SendOneMessage(m); err != nil {
		logger.Panic(err)
	}

	m.msgType = uploadFileMessage
	buf := make([]byte, 1024*64)
	size := 0
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			logger.Panic(err)
		}
		if err == io.EOF {
			os.Stdout.WriteString("\n")
			break
		}
		m.content = buf[:n]
		if err := c.wrap.SendOneMessage(m); err != nil {
			logger.Panic(err)
		}
		size += n
		os.Stdout.WriteString(fmt.Sprintf("\rsend %d KB", size/1024))
	}
	m.msgType = uploadDoneMessage
	m.content = []byte{}
	err = c.wrap.SendOneMessage(m)
	if err != nil {
		logger.Panic(err)
	}
}

func (c *Control) UploadFile(id, srcPath, dstPath string) {
	c.checkFile(srcPath)
	conn := c.connect()
	defer conn.Close()
	c.auth()
	c.matchNode(id)
	go c.upload(srcPath, dstPath)
	c.print()
}
