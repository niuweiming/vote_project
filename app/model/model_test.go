package model

import (
	"fmt"
	"testing"
	"time"
)

func TestGetVotes(t *testing.T) {
	NewMysql()
	//测试用例
	r := GetVotes()
	fmt.Printf("ret:%+v", r)
	Close()
}

func TestGetVote(t *testing.T) {
	NewMysql()
	//测试用例
	r := GetVote(3)
	fmt.Printf("%v", r)
	Close()
}

func TestDoVote(t *testing.T) {
	NewMysql()
	//测试用例
	r := DoVote(8, 3, []int64{1, 2})
	fmt.Printf("%v", r)
	Close()

}

func TestAddvote(t *testing.T) {
	NewMysql()
	//测试用例
	vote := Vote{
		Title:       "测试用例",
		Type:        0,
		Status:      0,
		Time:        0,
		UserId:      1,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	opt := make([]VoteOpt, 0)

	opt = append(opt, VoteOpt{
		Name: "测试用例1",
	})
	opt = append(opt, VoteOpt{
		Name: "测试用例2",
	})
	opt = append(opt, VoteOpt{
		Name: "测试用例3",
	})
	opt = append(opt, VoteOpt{
		Name: "测试用例4",
	})

	r := Addvote(vote, opt)
	fmt.Printf("%v", r)
	Close()

}

func TestGetUserV1(t *testing.T) {
	NewMysql()
	a := GetUserV1("admin")
	fmt.Printf("a:%+v", a)

}
