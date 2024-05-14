package wcgi

import (
	"io"
	"net/http"
	"net/http/cgi"
	"os"

	"github.com/hashicorp/yamux"
)

//go:wasmexport wagi_wcgi
func wagi_wcgi() {}

func Serve(h http.Handler) error {
	wcgi := os.Getenv("WAGI_WCGI") == "true"
	if !wcgi {
		return cgi.Serve(h)
	}
	stdio := &Stdio{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}
	sess, err := yamux.Server(stdio, nil)
	if err != nil {
		return err
	}
	return http.Serve(sess, h)
}

type Stdio struct {
	io.Reader
	io.Writer
}

var _ io.ReadWriteCloser = (*Stdio)(nil)

func (stdio *Stdio) Close() error {
	if reader, ok := stdio.Reader.(io.Closer); ok {
		if err := reader.Close(); err != nil {
			return err
		}
	}
	if writer, ok := stdio.Writer.(io.Closer); ok {
		if err := writer.Close(); err != nil {
			return err
		}
	}
	return nil
}
