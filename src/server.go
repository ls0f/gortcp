package gortcp

import (
	"net"
)

type Server struct {
	addr string
}

func (s * Server) Listen() {

	server, err := net.Listen("tcp", s.addr)
	if err != nil {
		logger.Error(err.Error())
	}
	for {
		conn, err := server.Accept()
		if err == nil {
			logger.Debugf("%s connect server", conn.RemoteAddr().String())
			go s.handlerConn(conn)
		}

	}

}

func (s *Server) handlerConn(conn net.Conn){
}

