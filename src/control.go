package gortcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Control struct {
	Addr    string
	Auth    string
	wrap    *MessageWrap
	curConn *net.TCPConn
}

func (c *Control) auth() {
	msg := &Message{msgType: authMessage, content: []byte(c.Auth)}
	err := c.wrap.SendOneMessage(msg)
	if err != nil {
		logger.Panic(err)
	}
}

func (c *Control) checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: %s", err.Error())
		os.Exit(1)
	}
}

func (c *Control) connect() *net.TCPConn {
	addr, err := net.ResolveTCPAddr("tcp", c.Addr)
	c.checkError(err)
	conn, err := net.DialTCP("tcp", nil, addr)
	c.checkError(err)
	c.wrap = &MessageWrap{rw: conn}
	return conn

}

func (c *Control) print(w io.Writer) {

	rm, err := c.wrap.ReadOneMessage()
	c.checkError(err)
	w.Write(rm.content)
}

func (c *Control) matchNode(id string) {
	msg := &Message{msgType: matchNodeMessage, content: []byte(id)}
	err := c.wrap.SendOneMessage(msg)
	c.checkError(err)
	m, err := c.wrap.ReadOneMessage()
	c.checkError(err)
	if m.msgType != matchOKMessage {
		os.Stdout.Write(m.content)
		os.Exit(1)
	}
}

func (c *Control) createTunnel(remoteAddr string) {

	msg := &Message{msgType: tunnelMessage, content: []byte(remoteAddr)}
	err := c.wrap.SendOneMessage(msg)
	c.checkError(err)
	m, err := c.wrap.ReadOneMessage()
	c.checkError(err)
	if m.msgType != tunnelOKMessage {
		os.Stdout.Write(m.content)
		os.Exit(1)
	}
}

func (c *Control) exec(cmd string) {

	msg := &Message{msgType: execCmdMessage, content: []byte(cmd)}
	err := c.wrap.SendOneMessage(msg)
	c.checkError(err)
}

func (c *Control) ListNode() {
	conn := c.connect()
	defer conn.Close()
	c.auth()
	msg := &Message{msgType: listNodeMessage}
	if err := c.wrap.SendOneMessage(msg); err != nil {
		c.checkError(err)
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
	c.checkError(err)
	if s.IsDir() {
		err = errors.New(fmt.Sprintf("%s is not a file", file))
		c.checkError(err)
	}
}

func (c *Control) upload(srcPath, dstPath string) {
	f, err := os.Open(srcPath)
	c.checkError(err)
	defer f.Close()

	fm := &FileMsg{dstPath: dstPath}
	content, err := fm.Bytes(srcPath)
	c.checkError(err)
	m := &Message{msgType: fileInfoMessage, content: content}
	if err = c.wrap.SendOneMessage(m); err != nil {
		c.checkError(err)
	}

	m.msgType = uploadFileMessage
	buf := make([]byte, 1024*4)
	size := 0
	start := time.Now()
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
			c.checkError(err)
		}
		size += n
		spendTime := time.Since(start)
		speed := time.Second.Nanoseconds() * int64(size) / 1024 / spendTime.Nanoseconds()
		os.Stdout.WriteString(fmt.Sprintf("send: %dKB | time: %.2fS | speed: %dKB/S\r", size/1024, spendTime.Seconds(), speed))
	}
	m.msgType = uploadDoneMessage
	m.content = []byte{}
	if err = c.wrap.SendOneMessage(m); err != nil {
		c.checkError(err)
	}
}

func (c *Control) UploadFile(id, srcPath, dstPath string) {
	c.checkFile(srcPath)
	conn := c.connect()
	defer conn.Close()
	c.auth()
	c.matchNode(id)
	c.wrap.disableReadTimeOut = true
	go c.upload(srcPath, dstPath)
	c.print(os.Stdout)
}

func (c *Control) listen(addr string) *net.TCPListener {
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	c.checkError(err)
	server, err := net.ListenTCP("tcp4", laddr)
	c.checkError(err)
	logger.Infof("[tcp] listen on local %v", addr)
	return server
}

func (c *Control) forwardData(id, remoteAddr string) {
	conn := c.connect()
	c.auth()
	c.matchNode(id)
	c.createTunnel(remoteAddr)
	c.wrap.disableReadTimeOut = true
	go func() {
		for {
			if m, err := c.wrap.ReadOneMessage(); err == nil && m.msgType == tunnelForwardMessage {
				if _, err := c.curConn.Write(m.content); err != nil {
					return
				}
			} else if err != nil {
				return
			}
		}
	}()
	buf := make([]byte, 1024)
	for {
		n, err := c.curConn.Read(buf)
		if err != nil {
			break
		}
		if err := c.wrap.SendOneMessage(&Message{msgType: tunnelForwardMessage, content: buf[:n]}); err != nil {
			break

		}
	}
	conn.Close()
	logger.Debugf("%s disconnect server", c.curConn.RemoteAddr().String())
	c.curConn.Close()
	c.curConn = nil
}

func (c *Control) Forward(id, localAddr, remoteAddr string) {
	server := c.listen(localAddr)
	for {
		conn, err := server.AcceptTCP()
		if err == nil {
			logger.Debugf("%s connect server", conn.RemoteAddr().String())
		}
		if c.curConn == nil {
			c.curConn = conn
			go c.forwardData(id, remoteAddr)
		} else {
			conn.Close()
			logger.Debugf("now it's forwarding conn %s", c.curConn.RemoteAddr().String())
			continue
		}
	}

}
