package gortcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

type Server struct {
	Addr string
	Auth string
	Pool *NodeMap
}

func (s *Server) Listen() {

	addr, err := net.ResolveTCPAddr("tcp", s.Addr)
	if err != nil {
		logger.Panic(err)
	}
	server, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		logger.Panic(err.Error())
	}
	logger.Infof("[tcp] listen on %v", addr)
	s.Pool = new(NodeMap)
	s.Pool.pool = make(map[uint32]*Node)
	for {
		conn, err := server.AcceptTCP()
		if err == nil {
			logger.Debugf("%s connect server", conn.RemoteAddr().String())
			go s.handlerConn(conn)
		}

	}

}

// send node list info
func (s *Server) listNode(wrap *MessageWrap) {
	msg := &Message{msgType: listNodeResultMessage, content: s.Pool.Bytes()}
	wrap.SendOneMessage(msg)
}

func (s *Server) forward(n *Node) {
	done := make(chan struct{})
	go func() {
		_, err := io.Copy(n.n1, n.n2)
		logger.Debug(err)
		//n.n2.CloseWrite()
		done <- struct{}{}
	}()
	_, err := io.Copy(n.n2, n.n1)
	//n.n2.CloseRead()
	logger.Debug(err)
	<-done
}

func (s *Server) handlerClientMessage(node *Node) {
	w1 := &MessageWrap{rw: node.n1}
	w2 := &MessageWrap{}
	for {
		m, err := w1.ReadOneMessage()
		if err != nil {
			return
		}
		if node.n2 != nil {
			w2.rw = node.n2
			w2.SendOneMessage(m)
		}

	}
}

func (s *Server) writeErrorMessage(wrap *MessageWrap, err error) {
	str := fmt.Sprintf("SERVER ERROR: %s\n", err.Error())
	msg := &Message{msgType: errorMessage, content: []byte(str)}
	wrap.SendOneMessage(msg)

}

func (s *Server) writeMatchOkMessage(wrap *MessageWrap) {
	msg := &Message{msgType: matchOKMessage}
	wrap.SendOneMessage(msg)
}

func (s *Server) handler(conn *net.TCPConn) {
	wrap := &MessageWrap{rw: conn}
	m, err := wrap.ReadOneMessage()
	if err != nil {
		return
	}
	if m.msgType == connectMessage {
		node := &Node{n1: conn, n2: nil, ConnectTime: time.Now()}
		id := s.Pool.addNewNode(node)
		defer s.Pool.removeNode(id)
		//done := make(chan bool)
		//<-done
		s.handlerClientMessage(node)
	} else if m.msgType == authMessage {
		if string(m.content) != s.Auth {
			logger.Debugf("auth error expected %s, actual: %s", s.Auth, m.content)
			s.writeErrorMessage(wrap, errors.New("auth error"))
			return
		}
		m, err := wrap.ReadOneMessage()
		if err != nil {
			return
		}
		if m.msgType == listNodeMessage {
			s.listNode(wrap)
			return
		}
		if m.msgType == matchNodeMessage {
			i, err := strconv.ParseUint(string(m.content), 10, 32)
			if err != nil {
				logger.Error(err)
				s.writeErrorMessage(wrap, err)
				return
			}
			node := s.Pool.getNode(uint32(i))
			if node == nil {
				errStr := fmt.Sprintf("conn id %s is not found in the server", m.content)
				logger.Error(errStr)
				s.writeErrorMessage(wrap, errors.New(errStr))
				return
			}
			node.n2 = conn
			logger.Debugf("%s match %s", node.n2.RemoteAddr(), node.n1.RemoteAddr())
			s.writeMatchOkMessage(wrap)
			io.Copy(node.n1, node.n2)
			//s.forward(node)
			node.n2 = nil
			return
		}
	}
}

func (s *Server) handlerConn(conn *net.TCPConn) {
	defer conn.Close()
	defer logger.Debugf("%s disconnect server", conn.RemoteAddr().String())
	s.handler(conn)
}
