package promtail

import (
	"bytes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"io/ioutil"
	"net/http"
	"time"
)

const LOG_ENTRIES_CHAN_SIZE = 5000

type LogLevel int

const (
	CRITICAL LogLevel = 2
	ERROR    LogLevel = 3
	WARNING  LogLevel = 4
	NOTICE   LogLevel = 5
	INFO     LogLevel = 6
	DEBUG    LogLevel = 7
)

type ClientConfig struct {
	// E.g. http://localhost:3100/api/prom/push
	PushURL string
	// E.g. "{job=\"somejob\"}"
	Labels             LabelSet
	BatchWait          time.Duration
	BatchEntriesNumber int
}

type Client interface {
	Log(Line string, level LogLevel, timestamp *timestamp.Timestamp, extraLabels LabelSet)
	Shutdown()
	Sent() int
	Buffered() int
}

// http.Client wrapper for adding new methods, particularly sendJsonReq
type httpClient struct {
	parent http.Client
}

// A bit more convenient method for sending requests to the HTTP server
func (client *httpClient) sendJsonReq(method, url string, ctype string, reqBody []byte) (resp *http.Response, resBody []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", ctype)

	resp, err = client.parent.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	resBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, resBody, nil
}
