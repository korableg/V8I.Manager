package v8ifilewriter

import "os"

type V8IFileWriter struct {
	filename string
}

func New(filename string) *V8IFileWriter {
	w := &V8IFileWriter{filename: filename}
	return w
}

func (w *V8IFileWriter) Write(data []byte) (int, error) {

	f, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return 0, err
	}

	defer f.Close()

	return f.Write(data)

}

func (w *V8IFileWriter) String() string {
	return w.filename
}
