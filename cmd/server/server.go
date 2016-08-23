package main

import (
	"flag"
	"gortcp/src"
)

func main() {

	addr := flag.String("addr", ":33456", "listen addr")
	auth := flag.String("auth", "123456", "auth")
	s := &gortcp.Server{Addr: *addr, Auth: *auth}
	s.Listen()
}
