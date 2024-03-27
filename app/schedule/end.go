package schedule

import (
	"fmt"
	"time"
	"vote/app/model"
)

func Start() {
	go EndVote()
}

func EndVote() {
	t := time.NewTicker(5 * time.Second)
	defer func() {
		t.Stop()
	}()

	for {
		select {
		case <-t.C:
			fmt.Println("EndVote启动")
			model.EndVoteV1()
			fmt.Println("EndVote运行完毕")
		}

	}
}
