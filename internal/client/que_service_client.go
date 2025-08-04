package client

import (
	"scorer/internal/models"
)

type QueServiceClient struct {
	host string
}

func NewQueServiceClient(host string) QueServiceClient {
	return QueServiceClient{host: host}
}

func (s QueServiceClient) EnqueueRequest(m models.Request) error {
	//que logic here
	return nil
}
