package promtail

import "testing"

func TestFingerprint(t *testing.T) {

	ls1 := make(LabelSet, 2)
	ls1.Append("set", "1").
		Append("value", "batch_label1")
	ls2 := make(LabelSet, 2)
	ls2.Append("set", "2").
		Append("value", "batch_label2")
	ls2b := make(LabelSet, 2)
	ls2b.Append("set", "2").
		Append("value", "batch_label2").
		Append("value", "batch_label2")

	if ls1.Fingerprint() == ls2.Fingerprint() {
		t.Fatalf("finterprints should not match when labels differ")
	}

	if ls2.Fingerprint() != ls2b.Fingerprint() {
		t.Fatalf("fingerprints should match")
	}
}
