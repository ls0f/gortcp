package main

import (
	"flag"
	"log"
	"os"

	"github.com/lovedboy/gortcp/src"
)

func main() {

	addr := flag.String("addr", ":33456", "forward server addr")
	auth := flag.String("auth", "123456", "forward server auth")
	id := flag.String("id", "", "client node id")
	action := flag.String("action", "", "action [list|exec|upload|forward]")
	cmd := flag.String("cmd", "", "cmd, required when action is exec")
	src := flag.String("src", "", "src path, required when action is upload")
	dst := flag.String("dst", "", "dst path, required when action is upload")
	laddr := flag.String("laddr", "", "local listen addr, required when action is forward")
	raddr := flag.String("raddr", "", "remote connet addr, required when action is forward")
	debug := flag.Bool("debug", false, "debug mode(true or false)")
	flag.Parse()
	gortcp.InitLogger(*debug)
	c := &gortcp.Control{Addr: *addr, Auth: *auth}
	if *action == "list" {
		c.ListNode()
		return
	}
	if *action == "exec" {
		if *id == "" {
			log.Panic("id is required")
		}
		if *cmd == "" {
			log.Panic("cmd is required")

		}
		c.ExecCommand(*id, *cmd)
		return
	}
	if *action == "upload" {
		if *id == "" {
			log.Panic("id is required")
		}
		if *src == "" {
			log.Panic("src path is required")
		}
		if *dst == "" {
			log.Panic("dst path is required")
		}
		c.UploadFile(*id, *src, *dst)
		return
	}
	if *action == "forward" {
		if *id == "" {
			log.Panic("id is required")
		}
		if *laddr == "" {
			log.Panic("laddr is required")
		}
		if *raddr == "" {
			log.Panic("raddr is required")
		}
		c.Forward(*id, *laddr, *raddr)
	}
	flag.PrintDefaults()
	os.Exit(1)
}
