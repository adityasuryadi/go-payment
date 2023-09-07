package service

import (
	"os"
	"payment/model"

	"github.com/go-resty/resty/v2"
)

func NewQueueService() QueueService {
	return &QueueServiceImpl{}
}

type QueueServiceImpl struct {
}

// CreateTicket implements QueueService.
func (service *QueueServiceImpl) CreateTicket(request model.GenerateTicketRequest) (*resty.Response, error) {
	url := os.Getenv("CREATE_TICKET_URL")
	client := resty.New()
	resp, err := client.R().
		SetFormData(map[string]string{
			"user_id":      string(request.UserId),
			"service_code": request.ServiceCode,
			"queue_name":   request.QueueName,
			"service_time": request.ServiceTime,
			"note":         request.Note,
			"queue_hp":     request.QueueHp,
		}).
		SetHeader("Accept", "application/json").
		Post(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
