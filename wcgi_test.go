package wcgi_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cgi"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/hashicorp/yamux"
	"github.com/shynome/err0/try"
	"github.com/shynome/wcgi"
)

func TestWCGI(t *testing.T) {
	stdio := &wcgi.Stdio{}
	cmd := exec.Command("go", "run", "./example")
	cmd.Env = append(os.Environ(), "WAGI_WCGI=true")
	cmd.Stderr = os.Stderr
	cmd.Stdin, stdio.Writer = io.Pipe()
	stdio.Reader, cmd.Stdout = io.Pipe()

	sess := try.To1(yamux.Client(stdio, nil))
	defer sess.Close()

	try.To(cmd.Start())
	defer cmd.Process.Kill()

	endpoint := fmt.Sprintf("http://yamux.proxy/")
	target := try.To1(url.Parse(endpoint))
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return sess.Open()
		},
	}

	l := try.To1(net.Listen("tcp", "127.0.0.1:0"))
	defer l.Close()
	go http.Serve(l, proxy)
	getBody := func(path string) string {
		link := "http://" + l.Addr().String() + path
		resp := try.To1(http.Get(link))
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			try.To(fmt.Errorf(""))
		}
		body := try.To1(io.ReadAll(resp.Body))
		return string(body)
	}
	if body := getBody("/"); body != "ok" {
		t.Error(body)
	} else {
		t.Log(body)
	}
	if body := getBody("/hello"); body != "world" {
		t.Error(body)
	} else {
		t.Log(body)
	}
}

func TestCGI(t *testing.T) {
	pwd := try.To1(os.Getwd())
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := &cgi.Handler{
			Path:   runtime.GOROOT() + "/bin/go",
			Dir:    pwd,
			Args:   []string{"run", "./example"},
			Env:    os.Environ(),
			Stderr: os.Stderr,
		}
		h.ServeHTTP(w, r)
	})

	l := try.To1(net.Listen("tcp", "127.0.0.1:0"))
	defer l.Close()
	go http.Serve(l, h)
	getBody := func(path string) string {
		link := "http://" + l.Addr().String() + path
		resp := try.To1(http.Get(link))
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			try.To(fmt.Errorf(""))
		}
		body := try.To1(io.ReadAll(resp.Body))
		return string(body)
	}
	if body := getBody("/"); body != "ok" {
		t.Error(body)
	} else {
		t.Log(body)
	}
	if body := getBody("/hello"); body != "world" {
		t.Error(body)
	} else {
		t.Log(body)
	}
}
