package elasticsearch

import (
	"bytes"
	"context"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/goccy/go-json"
	"hecigo.com/go-hexaboi/infra/core"
)

func Search(ctx context.Context, index string, query interface{}, result interface{}) (total uint64, extra map[string]interface{}, err error) {
	client := Client()
	var buf bytes.Buffer
	var reqBody map[string]interface{}

	reqBody, ok := query.(map[string]interface{})
	if !ok {
		if err := core.UnmarshalNoPanic(query, &reqBody); err != nil {
			return 0, nil, err
		}
	}
	if err := json.NewEncoder(&buf).Encode(reqBody); err != nil {
		return 0, nil, err
	}

	resp, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(index),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		withTimeout(),
		withPretty(),
		withErrorTrace(),
	)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	var respBody map[string]interface{}
	buf.Reset()
	buf.ReadFrom(resp.Body)
	if err := core.UnmarshalNoPanic(buf.String(), &respBody); err != nil {
		return 0, nil, err
	}

	if resp.IsError() {
		return 0, nil, &core.HPIResult{
			Status:    resp.StatusCode,
			Message:   resp.String(),
			Data:      respBody,
			ErrorCode: "ES_ERROR",
		}
	}

	esResult := respBody["hits"].(map[string]interface{})
	extra = make(map[string]interface{})

	total = uint64(esResult["total"].(map[string]interface{})["value"].(float64))
	if respBody["aggregations"] != nil {
		extra["aggs"] = respBody["aggregations"].(map[string]interface{})
	}

	var tmpResult []map[string]interface{}
	extraSorts := make(map[string]interface{})
	for _, h := range esResult["hits"].([]interface{}) {
		var m map[string]interface{}
		if err := core.UnmarshalNoPanic(h, &m); err != nil {
			return 0, nil, err
		}
		_source := m["_source"].(map[string]interface{})

		if m["sort"] != nil {
			extraSorts[m["_id"].(string)] = m["sort"]
		}

		tmpResult = append(tmpResult, _source)
	}
	extra["sorts"] = extraSorts

	if err := core.UnmarshalNoPanic(tmpResult, result); err != nil {
		return 0, nil, err
	}

	return total, extra, nil
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
