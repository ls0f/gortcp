package main

import (
	"flag"
	"gortcp/src"
	"log"
	"os"
)

func main() {

	addr := flag.String("addr", ":33456", "listen addr")
	auth := flag.String("auth", "123456", "auth")
	id := flag.String("id", "", "node id")
	action := flag.String("action", "", "action [list|exec]")
	cmd := flag.String("cmd", "", "cmd, required when action is exec")
	src := flag.String("src", "", "src path, required when action is upload")
	dst := flag.String("dst", "", "dst path, required when action is upload")
	flag.Parse()
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
	flag.PrintDefaults()
	os.Exit(1)
}
