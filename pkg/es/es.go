package es

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"microservices/pkg/meta"
	"time"
)

type ElasticSearchClient[T any] struct {
	client *elasticsearch.Client
}

func (e *ElasticSearchClient[T]) Search(ctx context.Context, index string, query map[string]any) (any, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := e.client.Search(
		e.client.Search.WithTimeout(4*time.Second),
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex(index),
		e.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	response, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var r meta.SearchResult[T]
	if err := json.Unmarshal(response, &r); err != nil {
		return nil, err
	}
	return r, nil
}

func (e *ElasticSearchClient[T]) Create(ctx context.Context, index string, id *string, record T) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(record); err != nil {
		return err
	}
	if id == nil {
		req := esapi.IndexRequest{
			Index:   index,
			Body:    &buf,
			Refresh: "true",
		}
		res, err := req.Do(ctx, e.client)
		if err != nil {
			return err
		}
		if res.IsError() {
			return errors.New(res.String())
		}
	} else {
		response, err := e.client.Create(
			index,
			id,
			&buf,
			e.client.Create.WithContext(ctx),
		)
		if err != nil {
			return err
		}
		if response.IsError() {
			return errors.New(response.String())
		}
	}
	return nil
}

func NewElasticSearch[T any](opts *Options) (*ElasticSearchClient[T], error) {
	elasticSearchClient := &ElasticSearchClient[T]{}
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{opts.IP},
		Username:  opts.Username,
		Password:  opts.Password,
	})
	if err != nil {
		return nil, err
	}
	elasticSearchClient.client = client
	return elasticSearchClient, nil
}
