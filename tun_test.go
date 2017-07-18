package pkt

import (
	"bytes"
	"syscall"
	"testing"
)

func TestTUNPacket(t *testing.T) {
	p := make(TUNPacket, 10)
	p.SetProto(syscall.AF_INET)
	if p[3] != syscall.AF_INET {
		t.Errorf("failed to set proto")
	}
	b := []byte{1, 2, 3, 4, 5, 6}
	l, err := p.CopyPayload(b)
	if err != nil {
		t.Errorf("failed to set payload")
	}
	if l != 6 {
		t.Errorf("set payload returns wrong count")
	}
	nb, _ := p.Payload()
	if !bytes.Equal(b, nb) {
		t.Errorf("not equal")
	}
	if !bytes.Equal(p, []byte{0, 0, 0, syscall.AF_INET, 1, 2, 3, 4, 5, 6}) {
		t.Errorf("not equal")
	}
}
