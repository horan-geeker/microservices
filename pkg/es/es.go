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

// Search 搜索 api
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

// Create 新增 doc 当 docId 有值时并且 es 中存在 docId 则进行更新 upsert
func (e *ElasticSearchClient[T]) Create(ctx context.Context, index string, id *string, record T) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(record); err != nil {
		return err
	}
	req := esapi.IndexRequest{
		Index:   index,
		Body:    &buf,
		Refresh: "true",
	}
	if id == nil {
		req.DocumentID = *id
	}
	res, err := req.Do(ctx, e.client)
	if err != nil {
		return err
	}
	if res.IsError() {
		return errors.New(res.String())
	}
	return nil
}

// Delete .
func (e *ElasticSearchClient[T]) Delete(ctx context.Context, index, docId string) error {
	res, err := e.client.Delete(
		index,
		docId,
		e.client.Delete.WithRefresh("true"),
		e.client.Delete.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	if res.IsError() {
		return errors.New(res.String())
	}
	return nil
}

// GetByDocId .
func (e *ElasticSearchClient[T]) GetByDocId(ctx context.Context, index string, docId string) (t T, err error) {
	res, err := e.client.Get(index, docId, e.client.Get.WithContext(ctx))
	if err != nil {
		return t, err
	}
	if res.IsError() {
		return t, errors.New(res.String())
	}
	defer res.Body.Close()
	var r T
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return t, err
	}
	return r, nil
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
