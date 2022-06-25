package elasticsearch

import (
	"bytes"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/goccy/go-json"

	"hoangphuc.tech/go-hexaboi/infra/core"
)

func Search(index string, query interface{}, result interface{}) (uint64, error) {
	client := Client()
	var buf bytes.Buffer
	var reqBody map[string]interface{}

	reqBody, ok := query.(map[string]interface{})
	if !ok {
		if err := core.UnmarshalNoPanic(query, &reqBody); err != nil {
			return 0, err
		}
	}
	if err := json.NewEncoder(&buf).Encode(reqBody); err != nil {
		return 0, err
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
		return 0, err
	}
	defer resp.Body.Close()

	var respBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return 0, err
	}

	if resp.IsError() {
		return 0, &core.HPIResult{
			Status:    resp.StatusCode,
			Message:   resp.String(),
			Data:      respBody,
			ErrorCode: "ES_ERROR",
		}
	}

	esResult := respBody["hits"].(map[string]interface{})
	totalHits := uint64(esResult["total"].(map[string]interface{})["value"].(float64))
	var tmpResult []map[string]interface{}
	for _, h := range esResult["hits"].([]interface{}) {
		var m map[string]interface{}
		if err := core.UnmarshalNoPanic(h, &m); err != nil {
			return 0, err
		}
		tmpResult = append(tmpResult, m["_source"].(map[string]interface{}))
	}

	if err := core.UnmarshalNoPanic(tmpResult, result); err != nil {
		return 0, err
	}

	return totalHits, nil
}

func withTimeout() func(*esapi.SearchRequest) {
	return func(r *esapi.SearchRequest) {
		r.Timeout = GetConfig().SearchTimeout
	}
}

func withPretty() func(*esapi.SearchRequest) {
	return func(r *esapi.SearchRequest) {
		r.Pretty = GetConfig().EnableDebugLogger
	}
}

func withErrorTrace() func(*esapi.SearchRequest) {
	return func(r *esapi.SearchRequest) {
		r.ErrorTrace = GetConfig().EnableDebugLogger
	}
}
