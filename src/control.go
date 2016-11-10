package gortcp

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
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
		logger.Fatal(err)
	}
}

func (c *Control) connect() *net.TCPConn {
	addr, err := net.ResolveTCPAddr("tcp", c.Addr)
	if err != nil {
		logger.Fatal(err)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		logger.Fatal(err)
	}
	c.wrap = &MessageWrap{rw: conn}
	return conn

}

func (c *Control) print(w io.Writer) {

	rm, err := c.wrap.ReadOneMessage()
	if err != nil {
		logger.Fatal(err)
	}
	w.Write(rm.content)
}

func (c *Control) matchNode(id string) {
	msg := &Message{msgType: matchNodeMessage, content: []byte(id)}
	err := c.wrap.SendOneMessage(msg)
	if err != nil {
		logger.Fatal(err)
	}
	m, err := c.wrap.ReadOneMessage()
	if err != nil {
		log.Fatal(err)
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
		logger.Fatal(err)
	}
}

func (c *Control) ListNode() {
	conn := c.connect()
	defer conn.Close()
	c.auth()
	msg := &Message{msgType: listNodeMessage}
	if err := c.wrap.SendOneMessage(msg); err != nil {
		logger.Fatal(err)
	}
	c.print(os.Stdout)
}

func (c *Control) ExecCommand(id, cmd string) {

	conn := c.connect()
	defer conn.Close()
	c.auth()
	c.matchNode(id)
	c.exec(cmd)
	c.print(os.Stdout)
}

func (c *Control) checkFile(file string) {

	s, err := os.Stat(file)
	if err != nil {
		logger.Fatal(err)
	}
	if s.IsDir() {
		logger.Fatal("%s is not a file", file)
	}
}

func (c *Control) upload(srcPath, dstPath string) {
	f, err := os.Open(srcPath)
	if err != nil {
		logger.Fatal(err)
	}
	defer f.Close()

	fm := &FileMsg{dstPath: dstPath}
	content, err := fm.Bytes(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	m := &Message{msgType: fileInfoMessage, content: content}
	if err = c.wrap.SendOneMessage(m); err != nil {
		logger.Fatal(err)
	}

	m.msgType = uploadFileMessage
	buf := make([]byte, 1024*4)
	size := 0
	start := time.Now()
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			logger.Fatal(err)
		}
		if err == io.EOF {
			os.Stdout.WriteString("\n")
			break
		}
		m.content = buf[:n]
		if err := c.wrap.SendOneMessage(m); err != nil {
			logger.Fatal(err)
		}
		size += n
		spendTime := time.Since(start)
		speed := time.Second.Nanoseconds() * int64(size) / 1024 / spendTime.Nanoseconds()
		os.Stdout.WriteString(fmt.Sprintf("send: %dKB | time: %.2fS | speed: %dKB/S\r", size/1024, spendTime.Seconds(), speed))
	}
	m.msgType = uploadDoneMessage
	m.content = []byte{}
	err = c.wrap.SendOneMessage(m)
	if err != nil {
		logger.Fatal(err)
	}
}

func (c *Control) UploadFile(id, srcPath, dstPath string) {
	c.checkFile(srcPath)
	conn := c.connect()
	defer conn.Close()
	c.auth()
	c.matchNode(id)
	go c.upload(srcPath, dstPath)
	c.print(os.Stdout)
}
