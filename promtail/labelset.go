package promtail

import (
	"fmt"
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
