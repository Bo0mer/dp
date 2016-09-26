package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/fatih/color"
)

var (
	addr     string
	target   string
	format   string
	insecure bool
)

func init() {
	flag.StringVar(&addr, "addr", "localhost:8080", "Address to bind to.")
	flag.StringVar(&target, "target", "https://google.com", "Target to proxy to.")
	flag.StringVar(&format, "format", "none", "Attempt to format payloads as.")
	flag.BoolVar(&insecure, "insecure", false, "Please do not!")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	t, err := url.Parse(target)
	if err != nil {
		log.Fatalf("dp: error parsing target url: %v\n", err)
	}

	var f Formatter
	switch format {
	case "json":
		f = &JSONFormatter{
			Prefix: "",
			Indent: "  ",
		}
	case "none":
		f = NopFormatter()
	default:
		log.Fatalf("dp: unknown format: %q\n", format)
	}

	p := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = t.Host
			req.URL.Scheme = t.Scheme
			req.URL.Host = t.Host
		},
		Transport: &SniffTransport{
			RoundTripper: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecure,
				},
			},
			Formatter: f,
			In:        color.New(color.FgCyan),
			Out:       color.New(color.FgMagenta),
		},
	}

	http.Handle("/", p)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("dp: error listening: %v\n", err)
	}
}

type nopformatter struct{}

func (n nopformatter) Format(src []byte) ([]byte, error) { return src, nil }

func NopFormatter() Formatter {
	return nopformatter{}
}

type JSONFormatter struct {
	Prefix string
	Indent string
}

func (f *JSONFormatter) Format(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return []byte{}, nil
	}
	var buf bytes.Buffer
	err := json.Indent(&buf, src, f.Prefix, f.Indent)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type Formatter interface {
	Format([]byte) ([]byte, error)
}

type Printer interface {
	Printf(string, ...interface{}) (int, error)
	Println(...interface{}) (int, error)
}

type SniffTransport struct {
	Formatter    Formatter
	RoundTripper http.RoundTripper
	In           Printer
	Out          Printer
}

func (s *SniffTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if derr := s.dumpRequest(req); derr != nil {
		log.Printf("dp: error dumping request: %v\n", derr)
	}
	resp, err := s.RoundTripper.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	if derr := s.dumpResponse(resp); derr != nil {
		log.Printf("dp: error dumping response: %v\n", derr)
	}
	return resp, err
}

func (s *SniffTransport) dumpRequest(req *http.Request) error {
	reqHead, err := httputil.DumpRequest(req, false)
	if err != nil {
		return err
	}
	s.In.Printf(string(reqHead))

	buf := bytes.Buffer{}
	_, err = io.Copy(&buf, req.Body)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)

	orig := buf.Bytes()
	fmted, err := s.Formatter.Format(orig)
	if err != nil {
		log.Printf("dp: error formatting: %v\n", err)
		fmted = orig
	}
	fmted = bytes.TrimRight(fmted, "\n")
	s.In.Println(string(fmted))
	s.In.Println()
	return nil
}

func (s *SniffTransport) dumpResponse(resp *http.Response) error {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}

	bytes.TrimRight(dump, "\n")
	s.Out.Println(string(dump))
	s.Out.Println()
	return nil
}
