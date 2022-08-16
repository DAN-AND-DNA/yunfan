package dbs

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

// 系统表
type System_info struct {
	Tid             string    `gorm:"index;column:tid;not null"`
	Sysid           uint64    `gorm:"primaryKey;column:sysid;autoIncrement"`
	System_name     string    `gorm:"uniqueIndex;column:system_name;not null"`
	System_describe string    `gorm:"column:system_describe;null"`
	Create_time     time.Time `gorm:"column:create_time;type:time;not null"`
}

// 系统和公司的映射表
type System_company_map struct {
	Sysid uint64 `gorm:"uniqueIndex:sc_key;column:sysid;not null"`
	Cid   uint64 `gorm:"uniqueIndex:sc_key;column:cid;not null"`
}

// 公司表id
type Company_info_id struct {
	Id               uint64 `gorm:"primaryKey;column:id;autoIncrement"`
	Create_timestamp uint64 `gorm:"column:create_timestamp;not null"`
}

// 公司表
type Company_info struct {
	Tid              string    `gorm:"index;column:tid;not null"`
	Cid              uint64    `gorm:"primaryKey;column:cid;not null"`
	Company_name     string    `gorm:"uniqueIndex;column:company_name;not null"`
	Company_describe string    `gorm:"column:company_describe;null"`
	Create_time      time.Time `gorm:"column:create_time;type:time;not null"`
}

// 公司和用户映射表
type Company_user_map struct {
	Cid uint64 `gorm:"uniqueIndex:cu_key;column:cid;not null"`
	Uid uint64 `gorm:"uniqueIndex:cu_key;column:uid;not null"`
}

// 用户表id
type User_info_id struct {
	Id               uint64 `gorm:"primaryKey;column:id;autoIncrement"`
	Create_timestamp uint64 `gorm:"column:create_timestamp;not null"`
}

// 用户表
type User_info struct {
	Tid         string    `gorm:"index;column:tid;not null"`
	Uid         uint64    `gorm:"primaryKey;column:uid;not null"`
	Username    string    `gorm:"uniqueIndex;column:username;not null"`
	Password    string    `gorm:"column:password;not null"`
	Create_time time.Time `gorm:"column:create_time;type:time;not null"`
}
