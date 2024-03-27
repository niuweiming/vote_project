package main

import (
	_ "github.com/go-sql-driver/mysql"
	"vote/app"
)

// @title           vote
// @version         1.0
// @description     This is a sample server vote server.

//@contact.email 2695062156@qq.com

//@license.name Apache 2.0

func main() {
	app.Start()
}
