package model

import "time"

// VoteOptUser 用于记录用户投票选项情况
type VoteOptUser struct {
	Id          int64     `gorm:"column:id;primary_key;NOT NULL"`
	UserId      int64     `gorm:"column:user_id;default:NULL"`
	VoteId      int64     `gorm:"column:vote_id;default:NULL"`
	VoteOptId   int64     `gorm:"column:vote_opt_id;default:NULL"`
	CreatedTime time.Time `gorm:"column:created_time;default:NULL"`
}

func (v *VoteOptUser) TableName() string {
	return "vote_opt_user"
}

// VoteOpt 用于记录投票选项
type VoteOpt struct {
	Id          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name        string    `gorm:"column:name;default:暂无选项"`
	VoteId      int64     `gorm:"column:vote_id;NOT NULL`
	Count       int64     `gorm:"column:count;default:0"`
	CreatedTime time.Time `gorm:"column:created_time;default:NULL"`
	UpdatedTime time.Time `gorm:"column:updated_time;default:NULL"`
}

func (v *VoteOpt) TableName() string {
	return "vote_opt"
}

// Vote 用于记录投票主题
type Vote struct {
	Id          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL" json:"id"`
	Title       string    `gorm:"column:title;default:NULL" json:"title"`
	Type        int32     `gorm:"column:type;default:0;comment:'0单选1多选'" json:"type"`
	Status      int32     `gorm:"column:status;default:1;comment:'0正常1超时'" json:"status"`
	Time        int64     `gorm:"column:time;default:NULL;comment:'有效时长'" json:"time"`
	UserId      int64     `gorm:"column:user_id;default:NULL;comment:'创建人'" json:"user_id"`
	CreatedTime time.Time `gorm:"column:create_time;default:NULL" json:"create_time"`
	UpdatedTime time.Time `gorm:"column:update_time;default:NULL" json:"update_time"`
}

func (v *Vote) TableName() string {
	return "vote"
}

// User 用于记录用户信息
type User struct {
	Id          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name        string    `gorm:"column:name;default:NULL"`
	Password    string    `gorm:"column:password;default:NULL"`
	CreatedTime time.Time `gorm:"column:create_time;default:NULL"`
	UpdatedTime time.Time `gorm:"column:update_time;default:NULL"`
	//uuid
	Uuid string `gorm:"column:uuid;default:NULL"`
}

func (v *User) TableName() string {
	return "user"
}

// VoteWithOpt 用于表示包含投票主题和选项
type VoteWithOpt struct {
	Vote Vote
	Opt  []VoteOpt
}

// VoteWithOptV1 用于表示包含投票主题和选项（版本1）
type VoteWithOptV1 struct {
	Vote Vote
	Opt  []VoteOpt `gorm : "foreignKey:VoteId"`
}
