package main

import (
	"fmt"
	"github.com/habakke/promtail-client/promtail"
	"log"
	"net/http"
	"os"
	"time"
)

func displayUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s proto|json source-name job-name\n", os.Args[0])
	os.Exit(1)
}

func displayInvalidName(arg string) {
	fmt.Fprintf(os.Stderr, "Invalid %s: allowed characters are a-zA-Z0-9_-\n", arg)
	os.Exit(1)
}

func nameIsValid(name string) bool {
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			(c == '-') || (c == '_')) {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 4 {
		displayUsage()
	}

	format := os.Args[1]
	source_name := os.Args[2]
	job_name := os.Args[3]
	if format != "proto" && format != "json" {
		displayUsage()
	}

	if !nameIsValid(source_name) {
		displayInvalidName("source-name")
	}

	if !nameIsValid(job_name) {
		displayInvalidName("job-name")
	}

	labels := make(promtail.LabelSet)
	labels = labels.
		Append("source", source_name).
		Append("job", job_name)

	conf := promtail.ClientConfig{
		PushURL:            "http://localhost:3100/api/prom/push",
		Labels:             labels,
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel:          promtail.DEBUG,
	}

	var (
		loki promtail.Client
		err  error
	)

	c := &http.Client{Timeout: time.Duration(1) * time.Second}

	loki, err = promtail.NewClientProto(conf, c)

	if err != nil {
		log.Printf("promtail.NewClient: %s\n", err)
		os.Exit(1)
	}

	extralabels := promtail.LabelSet{}
	for i := 1; i < 5; i++ {
		tstamp := time.Now().String()
		extralabels.Append("seq", fmt.Sprintf("%d", i))
		loki.Log(fmt.Sprintf("source = %s time = %s, i = %d\n", source_name, tstamp), promtail.DEBUG, extralabels)
		loki.Log(fmt.Sprintf("source = %s time = %s, i = %d\n", source_name, tstamp), promtail.INFO, extralabels)
		loki.Log(fmt.Sprintf("source = %s time = %s, i = %d\n", source_name, tstamp), promtail.WARNING, extralabels)
		loki.Log(fmt.Sprintf("source = %s time = %s, i = %d\n", source_name, tstamp), promtail.ERROR, extralabels)
		time.Sleep(1 * time.Second)
	}

	loki.Shutdown()
}
