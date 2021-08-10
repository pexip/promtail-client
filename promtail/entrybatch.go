package promtail

import "github.com/pexip/promtail-client/logproto"

type BatchMap map[uint64]*EntryBatch

func NewBatchMap() BatchMap {
	return make(BatchMap)
}

type EntryBatch struct {
	entries []*logproto.Entry
	labels  LabelSet
}

func NewEntryBatch(labels LabelSet) *EntryBatch {
	return &EntryBatch{
		entries: []*logproto.Entry{},
		labels:  labels,
	}
}

func (b EntryBatch) Append(entry *logproto.Entry) *EntryBatch {
	b.entries = append(b.entries, entry)
	return &b
}

func (b BatchMap) Append(labels LabelSet, entry *logproto.Entry) *BatchMap {
	fingerprint := labels.Fingerprint()
	if _, ok := b[fingerprint]; !ok {
		b[fingerprint] = NewEntryBatch(labels)
	}
	b[fingerprint].entries = append(b[fingerprint].entries, entry)
	return &b
}
