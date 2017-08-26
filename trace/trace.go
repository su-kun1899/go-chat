package trace

import "io"

// Tracer はコード内での出来事を記録できるオブジェクトを表す
type Tracer interface {
	// 任意の型の引数を何個でも受け取れる
	Trace(...interface{})
}

func New(w io.Writer) Tracer {
	return nil
}
