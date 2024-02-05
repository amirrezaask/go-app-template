package storage

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticClient(addresses []string, username string, password string) (*elasticsearch.TypedClient, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
	}
	typedClient, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to elasticsearch: %w", err)
	}

	return typedClient, nil
}
