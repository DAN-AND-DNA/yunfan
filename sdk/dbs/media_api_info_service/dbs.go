package media_api_info_service

import "time"

// 事务表
type Transaction_info struct {
	Id               uint64 `gorm:"primaryKey;column:id;autoIncrement"`
	Tid              string `gorm:"uniqueIndex:t_key;column:tid;not null"`
	Sid              string `gorm:"uniqueIndex:t_key;column:sid;null"`
	Sname            string `gorm:"uniqueIndex:t_key;column:sname;null"`
	Status           string `gorm:"column:status;null"`
	Create_timestamp uint64 `gorm:"column:create_timestamp;not null"`
}

// 头条api信息表
type Toutiao_api_info struct {
	Tid                     string    `gorm:"index;column:tid;not null"`
	App_id                  uint64    `gorm:"primaryKey;column:app_id;not null"`
	Secret                  string    `gorm:"column:secret;not null"`
	Create_time             time.Time `gorm:"column:create_time;type:time;not null"`
	Token_update_timestamp  uint64    `gorm:"column:token_update_timestamp;null"`
	Token_expired_timestamp uint64    `gorm:"column:token_expired_timestamp;null"`
	Access_token            string    `gorm:"column:access_token;null"`
	Refresh_token           string    `gorm:"column:refresh_token;null"`
}
