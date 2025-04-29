package clients

import (
	"go.temporal.io/sdk/client"
)

type TemporalClient struct {
	Client client.Client
}

func NewTemporalClient() (*TemporalClient, error) {
	client, err := client.Dial(client.Options{
		HostPort: "temporal:7233",
	})
	if err != nil {
		return nil, err
	}
	return &TemporalClient{Client: client}, nil
}
