package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffe
	tracer := New(&buf)
	if tracer == nil {
		t.Error("New からの戻り値がnilです")
	} else {
		tracer.Trace("こんにちは、traceパッケージ")
		if buf.String() != "こんにちは、traceパッケージ\n" {
			t.Errorf("'%s' という誤った文字が出力されました", buf.String())
		}
	}
}
