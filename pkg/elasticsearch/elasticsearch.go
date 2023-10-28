package elasticsearch

import (
	esLib "github.com/elastic/go-elasticsearch/v8"
)

var (
	esClient *esLib.TypedClient
)

func InitES() error {
	var err error
	esClient, err = esLib.NewTypedClient(esLib.Config{})
	if err != nil {
		return err
	}
	return nil
}

func GetESClient() (*esLib.TypedClient, error) {
	var err error
	if esClient == nil {
		err = InitES()
	}
	return esClient, err
}
