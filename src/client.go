package gortcp

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Client struct {
	Addr       string
	wrap       *MessageWrap
	fm         *FileMsg
	f          *os.File
	tunnelConn *net.TCPConn
}

func (c *Client) writeErrorMessage(err error) {
	str := fmt.Sprintf("CLIENT ERROR: %s\n", err.Error())
	msg := &Message{msgType: errorMessage, content: []byte(str)}
	c.wrap.SendOneMessage(msg)

}

func (c *Client) handlerMessage(m *Message) {
	switch m.msgType {
	case execCmdMessage:
		c.handlerCommand(m.content)
	case fileInfoMessage:
		c.handlerFileInfo(m.content)
	case uploadFileMessage:
		c.handlerFile(m.content)
	case uploadDoneMessage:
		c.receiveFileComplete()
	case pingOKMessage:
		logger.Debugf("receive ping reply message")
	case tunnelMessage:
		c.handlerTunnel(m.content)
	case tunnelForwardMessage:
		c.forwardTunnelMsg(m.content)
	default:
		logger.Debugf("unsupport msgType %d", m.msgType)
	}

}

func (c *Client) forwardTunnelMsg(body []byte) {
	_, err := c.tunnelConn.Write(body)
	if err != nil {
		logger.Debug(err)
	}
}

func (c *Client) handlerTunnel(addr []byte) {
	laddr, err := net.ResolveTCPAddr("tcp", string(addr))
	if err != nil {
		logger.Debug(err)
		c.writeErrorMessage(err)
		return
	}
	c.tunnelConn, err = net.DialTCP("tcp", nil, laddr)
	if err != nil {
		logger.Debug(err)
		c.writeErrorMessage(err)
		return
	}
	c.wrap.SendOneMessage(&Message{msgType: tunnelOKMessage})
	buf := make([]byte, 1024)
	go func() {
		for {
			n, err := c.tunnelConn.Read(buf)
			if err != nil {
				c.tunnelConn.Close()
				return
			}
			if err := c.wrap.SendOneMessage(&Message{msgType: tunnelForwardMessage, content: buf[:n]}); err != nil {
				c.tunnelConn.Close()
				return
			}
		}
	}()

}

func (c *Client) handlerFileInfo(p []byte) {
	var err error
	if c.fm, err = DecodeFileMsg(p); err != nil {
		c.writeErrorMessage(err)
		return
	}
	if c.f != nil {
		c.f.Close()
		c.f = nil
	}
	logger.Debugf("receive file info msg, path is %s", c.fm.dstPath)

}

func (c *Client) handlerFile(body []byte) {
	if c.fm == nil {
		err := errors.New("FileMsg is nil, do you send a file info msg?")
		c.writeErrorMessage(err)
		return
	}
	if c.f == nil {
		f, err := os.Create(c.fm.dstPath)
		if err != nil {
			c.writeErrorMessage(err)
			return
		}
		c.f = f

	}
	_, err := c.f.Write(body)
	if err != nil {
		c.writeErrorMessage(err)
		return
	}
}

func (c *Client) receiveFileComplete() {
	c.f.Sync()
	c.f.Close()
	defer func() {
		c.fm = nil
	}()

	md5, err := MD5sum(c.fm.dstPath)
	if err != nil {
		c.writeErrorMessage(err)
		return
	}
	if md5 != c.fm.md5 {
		err = errors.New(fmt.Sprintf("md5 verify not passed, expectd %s actual %s", c.fm.md5, md5))
		return
	}
	m := &Message{msgType: replyUploadDoneMessage, content: []byte("CLIENT: receive complete, md5 verify passed\n")}
	c.wrap.SendOneMessage(m)
}

func (c *Client) handlerCommand(cmd []byte) {
	content, err := c.execCommand(string(cmd))
	if err != nil {
		c.writeErrorMessage(err)
	} else {
		m := &Message{msgType: execCmdResultMessage, content: content}
		c.wrap.SendOneMessage(m)
	}

}

func (c *Client) execCommand(cmd string) (res []byte, err error) {
	logger.Debugf("exec cmd is :%s", cmd)
	if cmd == "" {
		err = errors.New("invalid cmd")
		return
	}
	cmd_list := strings.Fields(strings.TrimSpace(cmd))
	exp := exec.Command(cmd_list[0], cmd_list[1:]...)
	go func() {
		time.Sleep(time.Second * CommandTimeOut)
		exp.Process.Kill()
	}()
	res, err = exp.CombinedOutput()
	return

}

func (c *Client) HandlerMessage() {
	for {
		m, err := c.wrap.ReadOneMessage()
		if err != nil {
			return
		}
		c.handlerMessage(m)
	}
}

func (c *Client) Connect() {

	for {
		addr, err := net.ResolveTCPAddr("tcp", c.Addr)
		if err != nil {
			logger.Panic(err)
		}
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			logger.Errorf(err.Error())
			time.Sleep(20 * time.Second)
		} else {
			c.handlerConn(conn)
		}
	}
}

func (c *Client) ping() {
	go func() {
		for {
			if err := c.wrap.SendOneMessage(&Message{msgType: pingMessage}); err != nil {
				return
			}
			time.Sleep(PingInterval * time.Second)
		}
	}()
}

func (c *Client) handlerConn(conn *net.TCPConn) {
	defer conn.Close()
	c.wrap = &MessageWrap{rw: conn}
	msg := &Message{msgType: connectMessage, content: []byte("test")}
	if err := c.wrap.SendOneMessage(msg); err != nil {
		return
	}
	if _, err := c.wrap.ReadTheSpecialTypeMessage(connectOKmessage); err != nil {
		return
	}
	c.ping()
	c.HandlerMessage()

}
