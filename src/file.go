package gortcp

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
)

const bufferSize = 65536

type FileMsg struct {
	dstPath string
	md5     string
}

// MD5sum returns MD5 checksum of filename
func MD5sum(filename string) (string, error) {
	if info, err := os.Stat(filename); err != nil {
		return "", err
	} else if info.IsDir() {
		return "", errors.New(fmt.Sprintf("%s is a dir", filename))
	}

	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	for buf, reader := make([]byte, bufferSize), bufio.NewReader(file); ; {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		hash.Write(buf[:n])
	}

	checksum := fmt.Sprintf("%x", hash.Sum(nil))
	return checksum, nil
}

func (f *FileMsg) Bytes(srcPath string) (b []byte, err error) {
	md5, err := MD5sum(srcPath)
	if err != nil {
		return
	}
	f.md5 = md5
	buf := new(bytes.Buffer)
	buf.WriteString(md5)
	buf.WriteString(f.dstPath)
	return buf.Bytes(), nil

}

func DecodeFileMsg(b []byte) (f *FileMsg, err error) {
	if len(b) <= 32 {
		err = errors.New("invalid FileMsg struct")
		return
	}
	f = &FileMsg{}
	f.md5 = string(b[:32])
	f.dstPath = string(b[32:])
	return
}
