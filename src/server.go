package gortcp

import (
	"io"
	"net"
	"strconv"
	"time"
)

type Server struct {
	Addr string
	Auth string
	Pool *NodeMap
	wrap *MessageWrap
}

func (s *Server) Listen() {

	server, err := net.Listen("tcp", s.Addr)
	if err != nil {
		logger.Error(err.Error())
	}
	s.Pool = new(NodeMap)
	s.Pool.pool = make(map[uint32]*Node)
	for {
		conn, err := server.Accept()
		if err == nil {
			logger.Debugf("%s connect server", conn.RemoteAddr().String())
			go s.handlerConn(conn)
		}

	}

}

// send node list info
func (s *Server) listNode() {
	msg := &Message{msgType: listNodeResultMessage, content: s.Pool.Bytes()}
	s.wrap.SendOneMessage(msg)
}

func (s *Server) forward(n *Node) {
	done := make(chan struct{})
	go func() {
		io.Copy(n.n2, n.n1)
		done <- struct{}{}
	}()
	io.Copy(n.n1, n.n2)
	<-done
}

func (s *Server) handler(conn net.Conn) {
	m, err := s.wrap.ReadOneMessage()
	if err != nil {
		return
	}
	if m.msgType == connectMessage {
		id := s.Pool.addNewNode(&Node{
			n1: conn, n2: nil,
			ConnectTime: time.Now()})
		defer s.Pool.removeNode(id)
		done := make(chan bool)
		<-done
	} else if m.msgType == authMessage {
		if string(m.content) != s.Auth {
			logger.Debugf("auth error.expected:%s, actual:%s", s.Auth, m.content)
			return
		}
		m, err := s.wrap.ReadOneMessage()
		if err != nil {
			return
		}
		if m.msgType == listNodeMessage {
			s.listNode()
			return
		} else if m.msgType == matchNodeMessage {
			i, err := strconv.ParseUint(string(m.content), 10, 32)
			if err != nil {
				logger.Error(err)
				return
			}
			node := s.Pool.getNode(uint32(i))
			if node == nil {
				logger.Debugf("id:%s is not found in the server", m.content)
				return
			}
			node.n2 = conn
			s.forward(node)
			node.n2 = nil
			return
		}
	}
}

func (s *Server) handlerConn(conn net.Conn) {
	defer conn.Close()
	defer logger.Debugf("%s disconnect server", conn.RemoteAddr().String())
	s.wrap = &MessageWrap{rw: conn}
	s.handler(conn)
}
