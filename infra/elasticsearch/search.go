package elasticsearch

import (
	"bytes"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/goccy/go-json"

	"hoangphuc.tech/go-hexaboi/infra/core"
)

var (
	esDefaultConfig = GetConfig()
)

func Search(index string, query interface{}, result interface{}) error {
	client := Client()
	var buf bytes.Buffer
	var reqBody map[string]interface{}

	switch query.(type) {
	case string:
		core.UnmarshalNoPanic(query, &reqBody)
		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			return err
		}
	case map[string]interface{}:
		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			return err
		}
	default:
		return fmt.Errorf("type of query must be JSON string or map[string]interface{}")
	}

	resp, err := client.Search(
		client.Search.WithIndex(index),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		withTimeout(),
		withPretty(),
		withErrorTrace(),
	)
	if err != nil {
		return err
	}

	if resp.IsError() {
		var errBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
			return err
		}
		return &core.HPIResult{
			Status:    resp.StatusCode,
			Message:   resp.String(),
			Data:      errBody,
			ErrorCode: "ES_ERROR",
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	return nil
}

func withTimeout() func(*esapi.SearchRequest) {
	return func(r *esapi.SearchRequest) {
		r.Timeout = esDefaultConfig.SearchTimeout
	}
}

func withPretty() func(*esapi.SearchRequest) {
	return func(r *esapi.SearchRequest) {
		r.Pretty = esDefaultConfig.EnableDebugLogger
	}
}

func withErrorTrace() func(*esapi.SearchRequest) {
	return func(r *esapi.SearchRequest) {
		r.ErrorTrace = esDefaultConfig.EnableDebugLogger
	}
}
