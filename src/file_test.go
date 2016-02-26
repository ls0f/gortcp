package gortcp

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashFileMd5(t *testing.T) {

	f, err := os.Create("test")
	assert.NoError(t, err)
	f.WriteString("123456")
	f.Close()
	defer os.Remove("test")
	m, err := MD5sum("test")
	assert.NoError(t, err)
	assert.Equal(t, m, "e10adc3949ba59abbe56e057f20f883e")

}

func TestFileMsg_Bytes(t *testing.T) {
	fm := &FileMsg{dstPath: "/tmp/test"}
	f, err := os.Create("test")
	assert.NoError(t, err)
	f.WriteString("123456")
	defer os.Remove("test")
	b, err := fm.Bytes("test")
	assert.NoError(t, err)
	assert.Equal(t, string(b), "e10adc3949ba59abbe56e057f20f883e/tmp/test")
}

func TestDecodeFileMsg(t *testing.T) {
	f, err := DecodeFileMsg([]byte("e10adc3949ba59abbe56e057f20f883e/tmp/test"))
	assert.NoError(t, err)
	assert.Equal(t, f.md5, "e10adc3949ba59abbe56e057f20f883e")
	assert.Equal(t, f.dstPath, "/tmp/test")

}
