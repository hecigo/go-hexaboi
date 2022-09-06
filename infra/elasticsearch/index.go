package elasticsearch

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v7/esapi"

	log "github.com/sirupsen/logrus"
	"hoangphuc.tech/go-hexaboi/infra/core"
)

type bulkResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			ID     string `json:"_id"`
			Result string `json:"result"`
			Status int    `json:"status"`
			Error  struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
				Cause  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
				} `json:"caused_by"`
			} `json:"error"`
		} `json:"index"`
	} `json:"items"`
}

func Index(indexName string, docIdField string, documents ...interface{}) error {
	var (
		buf bytes.Buffer
		res *esapi.Response
		raw map[string]interface{}
		blk *bulkResponse

		numItems   int
		numErrors  int
		numIndexed int
		currBatch  int
	)

	es := Client()
	cfg := GetConfig()
	count := len(documents)

	log.Infof(
		"\x1b[1mBulk\x1b[0m: documents [%s] batch size [%s]\n",
		humanize.Comma(int64(count)), humanize.Comma(int64(cfg.BatchIndexSize)))

	for i, d := range documents {
		numItems++

		currBatch = i / cfg.BatchIndexSize
		if i == count-1 {
			currBatch++
		}

		var doc map[string]interface{}
		if err := core.UnmarshalNoPanic(d, &doc); err != nil {
			return err
		}
		docId := core.Utils.ToStr(doc[docIdField])
		if docId == "" {
			return fmt.Errorf("document ID must not be empty")
		}

		// Prepare the metadata payload
		meta := []byte(fmt.Sprintf(`{"index": {"_id": %s}}%s`, docId, "\n"))
		if cfg.EnableDebugLogger {
			fmt.Printf("%s", meta)
		}

		// Prepare the data payload: encode article to JSON
		data, err := json.Marshal(d)
		if err != nil {
			return fmt.Errorf("cannot encode article %s: %s", docId, err)
		}

		// Append newline to the data payload
		data = append(data, "\n"...)
		if cfg.EnableDebugLogger {
			fmt.Printf("%s", data)
		}

		// Append payloads to the buffer (ignoring write errors)
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)

		// When a threshold is reached, execute the Bulk() request with body from buffer
		if i > 0 && i%cfg.BatchIndexSize == 0 || i == count-1 {

			res, err = es.Bulk(bytes.NewReader(buf.Bytes()), es.Bulk.WithIndex(indexName))
			if err != nil {
				return fmt.Errorf("failure indexing batch %d: %s", currBatch, err)
			}

			// If the whole request failed, print error and mark all documents as faileda
			if res.IsError() {
				numErrors += numItems
				if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
					return fmt.Errorf("failure to to parse response body: %s", err)
				} else {
					return &core.HPIResult{
						Status:    res.StatusCode,
						Message:   fmt.Sprintf("failure indexing batch %d", currBatch),
						Data:      raw["error"],
						ErrorCode: "ES_ERROR",
					}
				}
			} else {
				// A successful response might still contain errors for particular documents...

				if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
					return fmt.Errorf("failure to to parse response body: %s", err)
				} else {
					var errData []string
					for _, d := range blk.Items {
						// ... so for any HTTP status above 201 ...
						if d.Index.Status > 201 {
							// ... increment the error counter ...
							numErrors++

							// ... and print the response status and error information ...
							errData = append(errData, core.Utils.ToStr(d.Index))
						} else {
							// ... otherwise increase the success counter.
							numIndexed++
						}
					}

					if numErrors > 0 {
						return &core.HPIResult{
							Status:    http.StatusInternalServerError,
							Message:   fmt.Sprintf("There are [%d] failed documents", numErrors),
							Data:      errData,
							ErrorCode: "ES_ERROR",
						}
					}
				}
			}

			// Close the response body, to prevent reaching the limit for goroutines or file handles
			res.Body.Close()

			// Reset the buffer and items counter
			buf.Reset()
			numItems = 0
		}
	}

	return nil
}
