package main

import (
	_ "github.com/go-sql-driver/mysql"
	"vote/app"
)

func main() {
	app.Start()
}
