package gortcp

import (
	"os"
	"testing"

	"log"

	"github.com/stretchr/testify/assert"
)

func TestControl_UploadFile(t *testing.T) {
	f, err := os.Open("./control.go")
	assert.NoError(t, err)
	buf := make([]byte, 2497)
	n, err := f.Read(buf)
	log.Println(n, err)
	n2, err2 := f.Read(buf)
	log.Println(n2, err2)

}
