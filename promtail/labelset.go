package promtail

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
)

type LabelSet map[string]string

func (ls LabelSet) Append(label string, value string) LabelSet {
	ls[label] = value
	return ls
}

func (ls LabelSet) Copy() LabelSet {
	c := make(LabelSet, len(ls))
	for k, v := range ls {
		c[k] = v
	}
	return c
}

func (ls LabelSet) String() string {
	d := make([]string, 0, len(ls))
	for k, v := range ls {
		d = append(d, fmt.Sprintf("%s=%q", k, v))
	}

	sort.Strings(d)
	return fmt.Sprintf("{%s}", strings.Join(d, ", "))
}

func (ls LabelSet) WithExtras(extra LabelSet) LabelSet {
	d := make(LabelSet, len(ls))
	for k, v := range ls {
		d[k] = v
	}
	for k, v := range extra {
		d[k] = v
	}
	return d
}

func hash(arr []string) uint64 {
	h := fnv.New64a()
	for _, s := range arr {
		_, _ = h.Write([]byte(s))
	}

	return h.Sum64()
}

func (ls LabelSet) Fingerprint() uint64 {
	if len(ls) == 0 {
		return 0
	}

	var labels []string
	for k, v := range ls {
		labels = append(labels, fmt.Sprintf("%s:%s", k, v))
	}
	sort.Strings(labels)
	hash := hash(labels)
	return hash
}
