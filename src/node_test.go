package gortcp

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNewNode(t *testing.T) {
	nm := &NodeMap{pool: make(map[uint32]*Node)}
	n := &Node{n1: nil, n2: nil}
	i := nm.addNewNode(n)
	assert.Equal(t, i, uint32(1))
}

func TestRemoveNode(t *testing.T) {
	nm := &NodeMap{pool: make(map[uint32]*Node)}
	n := &Node{n1: nil, n2: nil}
	i := nm.addNewNode(n)
	nm.removeNode(i)
	assert.Len(t, nm.pool, 0)

}

func TestGetNode(t *testing.T) {
	nm := &NodeMap{pool: make(map[uint32]*Node)}
	n := &Node{n1: nil, n2: nil}
	i := nm.addNewNode(n)
	n = nm.getNode(i)
	assert.Equal(t, n.n1, (*net.TCPConn)(nil))
}

func TestNodeMap_Bytes(t *testing.T) {
	nm := &NodeMap{pool: make(map[uint32]*Node)}
	buf := nm.Bytes()
	assert.Contains(t, string(buf), "ADDRESS")
}
