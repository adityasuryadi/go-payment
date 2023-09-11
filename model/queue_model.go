package model

type GenerateTicketRequest struct {
	UserId      int    `json:"user_id"`
	ServiceCode string `json:"string"`
	QueueName   string `json:"queue_name"`
	ServiceTime string `json:"service_time"`
	QueueHp     string `json:"queue_hp"`
	Note        string `json:"note"`
}
