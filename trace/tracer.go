package trace

import "io"
import "fmt"

// Tracer はコード内での出来事を記録できるオブジェクトを表す
type Tracer interface {
	// 任意の型の引数を何個でも受け取れる
	Trace(...interface{})
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}
