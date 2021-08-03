package promtail

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"sync"
	"time"
)

type jsonLogEntry struct {
	Ts    time.Time `json:"ts"`
	Line  string    `json:"line"`
	level LogLevel  // not used in JSON
}

type promtailStream struct {
	Labels  string          `json:"labels"`
	Entries []*jsonLogEntry `json:"entries"`
}

type promtailMsg struct {
	Streams []promtailStream `json:"streams"`
}

type clientJson struct {
	config    *ClientConfig
	quit      chan struct{}
	entries   chan *jsonLogEntry
	waitGroup sync.WaitGroup
	client    httpClient
	token     *oauth2.Token
	sent      int
	buffered  int
}

func NewClientJson(conf ClientConfig, parent *http.Client) (Client, error) {
	client := clientJson{
		config:  &conf,
		quit:    make(chan struct{}),
		entries: make(chan *jsonLogEntry, LOG_ENTRIES_CHAN_SIZE),
		client: httpClient{
			parent: *parent,
		},
	}

	client.waitGroup.Add(1)
	go client.run()

	return &client, nil
}

func (c *clientJson) Debugf(format string, args ...interface{}) {
	c.log(format, DEBUG, "Debug: ", args...)
}

func (c *clientJson) Infof(format string, args ...interface{}) {
	c.log(format, INFO, "Info: ", args...)
}

func (c *clientJson) Warnf(format string, args ...interface{}) {
	c.log(format, WARN, "Warn: ", args...)
}

func (c *clientJson) Errorf(format string, args ...interface{}) {
	c.log(format, ERROR, "Error: ", args...)
}

func (c *clientJson) log(format string, level LogLevel, prefix string, args ...interface{}) {
	if (level >= c.config.SendLevel) || (level >= c.config.PrintLevel) {
		c.entries <- &jsonLogEntry{
			Ts:    time.Now(),
			Line:  fmt.Sprintf(prefix+format, args...),
			level: level,
		}
	}
}

func (c *clientJson) Shutdown() {
	close(c.quit)
	c.waitGroup.Wait()
}

func (c *clientJson) Sent() int {
	return c.sent
}

func (c *clientJson) Buffered() int {
	return c.buffered
}

func (c *clientJson) run() {
	var batch []*jsonLogEntry
	batchSize := 0
	maxWait := time.NewTimer(c.config.BatchWait)

	defer func() {
		if batchSize > 0 {
			c.send(batch)
		}

		c.waitGroup.Done()
	}()

	for {
		select {
		case <-c.quit:
			return
		case entry := <-c.entries:
			c.buffered++
			if entry.level >= c.config.PrintLevel {
				log.Printf(entry.Line)
			}

			if entry.level >= c.config.SendLevel {
				batch = append(batch, entry)
				batchSize++
				if batchSize >= c.config.BatchEntriesNumber {
					c.send(batch)
					batch = []*jsonLogEntry{}
					batchSize = 0
					maxWait.Reset(c.config.BatchWait)
				}
			}
		case <-maxWait.C:
			if batchSize > 0 {
				c.send(batch)
				batch = []*jsonLogEntry{}
				batchSize = 0
			}
			maxWait.Reset(c.config.BatchWait)
		}
	}
}

func (c *clientJson) send(entries []*jsonLogEntry) {
	var streams []promtailStream
	batchSize := len(entries)
	defer func() {
		c.buffered = c.buffered - batchSize
		c.sent = c.sent + batchSize
	}()

	streams = append(streams, promtailStream{
		Labels:  c.config.Labels,
		Entries: entries,
	})

	msg := promtailMsg{Streams: streams}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("unable to marshal a JSON document")
		return
	}

	resp, body, err := c.client.sendJsonReq("POST", c.config.PushURL, "application/json", jsonMsg)
	if err != nil {
		log.Printf("unable to send an HTTP request: %v", err)
		return
	}

	if resp.StatusCode != 204 {
		log.Printf("unexpected HTTP status code: %d, message: %s\n", resp.StatusCode, body)
		return
	}
}
