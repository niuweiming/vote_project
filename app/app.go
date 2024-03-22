package app

import (
	"vote/app/model"
	"vote/app/router"
	"vote/app/tools"
)

// /启动器方法
func Start() {
	model.NewMysql()
	defer func() {
		model.Close()
	}()
	tools.NewLogger()
	//schedule.Start()
	router.New()
}
