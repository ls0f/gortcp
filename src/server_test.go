package gortcp

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestServer_Listen(t *testing.T) {

	defer func(){
		r:=recover()
		assert.True(t, r != nil)
		_, ok := r.(error)
		assert.False(t, ok)
	}()
	s := Server{Addr: ":12345", Auth:"123456"}
	go s.Listen()
	time.Sleep(500*time.Millisecond)
	s.Listen()

}