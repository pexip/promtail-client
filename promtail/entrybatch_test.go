package promtail

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pexip/promtail-client/logproto"
	"testing"
	"time"
)

func TestBatchMapGrouping(t *testing.T) {
	m := NewBatchMap()

	ls1 := make(LabelSet, 2)
	ls1.Append("set", "1").
		Append("value", "batch_label1")
	ls2 := make(LabelSet, 2)
	ls2.Append("set", "2").
		Append("value", "batch_label2")
	ls3 := make(LabelSet, 2)
	ls3.Append("set", "3").
		Append("value", "batch_label3")
	ls4 := make(LabelSet, 2)
	ls4.Append("set", "4").
		Append("value", "batch_label4")

	n := 10
	ln := 0
	for i := 1; i <= n; i++ {
		ln++
		if ln > 4 {
			ln = 1
		}

		now := time.Now()
		ts := &timestamp.Timestamp{
			Seconds: now.UnixNano() / int64(time.Second),
			Nanos:   int32(now.UnixNano() % int64(time.Second)),
		}

		e := &logproto.Entry{
			Timestamp: ts,
			Line:      fmt.Sprintf("This is log line #%d", i),
		}

		switch ln {
		case 1:
			m.Append(ls1, e)
		case 2:
			m.Append(ls2, e)
		case 3:
			m.Append(ls3, e)
		case 4:
			m.Append(ls4, e)
		default:
			t.Fatalf("unexpectcted value for 'n'")
		}
	}

	if len(m) != 4 {
		t.Fatalf("map does not contain the expected number of entries")
	}

	c := 0
	for _, v := range m {
		c += len(v.entries)
	}

	if c != n {
		t.Fatalf("map does not contain the expected number of entries")
	}
}
