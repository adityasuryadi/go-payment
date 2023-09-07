package service

import (
	"payment/model"

	"github.com/go-resty/resty/v2"
)

type QueueService interface {
	CreateTicket(request model.GenerateTicketRequest) (*resty.Response, error)
}
