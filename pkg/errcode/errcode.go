package errcode

type Errcode struct {
	Msg   string `json:"msg" example:"错误消息"` // 错误消息
	From1 int    `json:"from"`               // 来源 {3000: user_service}
	Code1 int    `json:"code"`               // 错误码
}

func New(msg string, from, code int) *Errcode {
	return &Errcode{msg, from, code}
}

func (this *Errcode) Error() string {
	return this.Msg
}

func (this *Errcode) From() int {
	return this.From1
}

func (this *Errcode) Code() int {
	return this.Code1
}
