package gortcp

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"
)

type Node struct {
	n1                *net.TCPConn
	n2                *net.TCPConn
	ConnectTime       time.Time
	LastHeartBeatTime time.Time
	sync.Mutex
}

func (n *Node) updateLastHeartBeat() {

	n.LastHeartBeatTime = time.Now()
}

func (n *Node) setN2(conn *net.TCPConn) {
	n.Lock()
	n.n2 = conn
	defer n.Lock()

}

type NodeMap struct {
	sync.Mutex
	pool   map[uint32]*Node
	cur_id uint32
}

func (m *NodeMap) addNewNode(n *Node) uint32 {
	m.Lock()
	defer m.Unlock()
	m.cur_id += 1
	m.pool[m.cur_id] = n
	return m.cur_id
}

func (m *NodeMap) removeNode(id uint32) {
	delete(m.pool, id)
}

func (m *NodeMap) getNode(id uint32) (n *Node) {
	n, _ = m.pool[id]
	return
}

func (m *NodeMap) Bytes() []byte {

	buf := new(bytes.Buffer)
	buf.WriteString("ID   ADDRESS\n")
	for k, v := range m.pool {
		str := fmt.Sprintf("%d   %s\n", k, v.n1.RemoteAddr())
		buf.WriteString(str)
	}
	return buf.Bytes()

}
