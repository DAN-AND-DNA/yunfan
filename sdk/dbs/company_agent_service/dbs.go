package company_agent_service

// 事务表
type Transaction_info struct {
	Id               uint64 `gorm:"primaryKey;column:id;autoIncrement"`
	Tid              string `gorm:"uniqueIndex:t_key;column:tid;not null"`
	Sid              string `gorm:"uniqueIndex:t_key;column:sid;null"`
	Sname            string `gorm:"uniqueIndex:t_key;column:sname;null"`
	Status           string `gorm:"column:status;null"`
	Create_timestamp uint64 `gorm:"column:create_timestamp;not null"`
}

// 公司层级数据
type Company_data_info struct {
	Id   uint64  `gorm:"primaryKey;column:id`
	Name string  `gorm:"index;column:name;not null"`
	Cost float64 `gorm:"column:cost"`
}
