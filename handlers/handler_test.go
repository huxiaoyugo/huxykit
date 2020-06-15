package handlers

import "testing"

func TestHandler(t *testing.T) {

	ctx := NewEchoContext()

	err := DoneHandler(OneHandler(TwoHandler(ThreeHandler(Default))))(ctx)
	if err != nil {
		t.Error(err)
	}
}
