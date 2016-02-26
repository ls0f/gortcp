package main

import (
	"flag"

	"github.com/lovedboy/gortcp/src"
)

func main() {

	addr := flag.String("addr", ":33456", "listen addr")
	auth := flag.String("auth", "123456", "auth")
	debug := flag.Bool("debug", false, "debug mode(true or false)")
	flag.Parse()
	gortcp.InitLogger(*debug)
	s := &gortcp.Server{Addr: *addr, Auth: *auth}
	s.Listen()
}
