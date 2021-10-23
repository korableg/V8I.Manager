package v8imockwriter

type V8IMockWriter struct {
	W chan struct{}
}

func New() *V8IMockWriter {
	return &V8IMockWriter{
		W: make(chan struct{}, 1),
	}
}

func (w *V8IMockWriter) Write(data []byte) (int, error) {
	w.W <- struct{}{}
	return len(data), nil
}

func (w *V8IMockWriter) String() string {
	return ""
}
