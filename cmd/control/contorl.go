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
	cmd := flag.String("cmd", "", "cmd")
	flag.Parse()
	c := &gortcp.Control{Addr: *addr, Auth: *auth}
	if *action == "list"{
		c.ListNode()
		return
	}
	if *action == "exec"{
		if *id == "" {
			log.Panic("id is required")
		}
		if *cmd == ""{
			log.Panic("cmd is required")

		}
		c.ExecCommand(*id, *cmd)
		return
	}
	flag.PrintDefaults()
	os.Exit(1)
}
