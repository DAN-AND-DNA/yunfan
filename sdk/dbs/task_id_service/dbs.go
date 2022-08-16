package task_id_service

type Segments struct {
	BizTag     string `json:"biz_tag" gorm:"type:VARCHAR(128)"`
	MaxId      int64  `json:"max_id" gorm:"type:bigint"`
	Step       int64  `json:"step" gorm:"int"`
	CreateTime int64  `json:"create_time" gorm:"bigint"`
	UpdateTime int64  `json:"update_time" gorm:"bigint"`
}

func (tb *Segments) TableName() string {
	return "segments"
}
