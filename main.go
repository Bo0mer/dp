package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/Bo0mer/flagvar"
	"github.com/fatih/color"
)

var (
	addr     string
	target   string
	format   string
	insecure bool
	headers  flagvar.Map
)

func init() {
	flag.StringVar(&addr, "addr", "localhost:8080", "Address to bind to.")
	flag.StringVar(&target, "target", "https://example.com", "Target to proxy to.")
	flag.StringVar(&format, "format", "auto", "Attempt to format payloads as.")
	flag.BoolVar(&insecure, "insecure", false, "Please do not!")
	flag.Var(&headers, "header", "Header to add. Must be in Name:value format.")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	t, err := url.Parse(target)
	if err != nil {
		log.Fatalf("dp: error parsing target url: %v\n", err)
	}

	fmts := make(map[string]Formatter)
	switch format {
	case "json":
		fmts[""] = &JSONFormatter{}
	case "plain":
		fmts[""] = NopFormatter()
	case "auto":
		fmts["application/json"] = &JSONFormatter{}
		fmts["text/plain"] = NopFormatter()
		fmts[""] = NopFormatter()
	default:
		log.Fatalf("dp: unknown format: %q\n", format)
	}

	p := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = t.Host
			req.URL.Scheme = t.Scheme
			req.URL.Host = t.Host
			for k, v := range headers {
				req.Header.Add(k, v)
			}
		},
		Transport: &SniffTransport{
			RoundTripper: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecure,
				},
			},
			Formatters: fmts,
			In:         color.New(color.FgCyan),
			Out:        color.New(color.FgMagenta),
		},
	}

	http.Handle("/", p)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("dp: error listening: %v\n", err)
	}
}

type SniffTransport struct {
	Formatters   map[string]Formatter
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
	var err error
	var save io.ReadCloser
	savecl := req.ContentLength

	reqHead, err := httputil.DumpRequest(req, false)
	if err != nil {
		return err
	}
	s.In.Printf(string(reqHead))

	save, req.Body, err = drainBody(req.Body)
	if err != nil {
		return err
	}
	req.ContentLength = savecl

	orig, err := ioutil.ReadAll(save)
	if err != nil {
		return err
	}

	ct := detectContentType(req.Header, orig)
	fmted, err := s.formatBody(orig, ct)
	s.In.Println(fmted)
	return err
}

func (s *SniffTransport) dumpResponse(resp *http.Response) error {
	var err error
	var save io.ReadCloser
	savecl := resp.ContentLength

	respHead, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return err
	}
	s.Out.Printf(string(respHead))

	save, resp.Body, err = drainBody(resp.Body)
	if err != nil {
		return err
	}
	resp.ContentLength = savecl

	orig, err := ioutil.ReadAll(save)
	if err != nil {
		return err
	}

	ct := detectContentType(resp.Header, orig)
	fmted, err := s.formatBody(orig, ct)
	s.Out.Println(fmted)
	return err
}

func (s *SniffTransport) formatBody(body []byte, contentType string) (string, error) {
	var f Formatter
	var ok bool
	if f, ok = s.Formatters[contentType]; !ok {
		f = s.Formatters[""]
	}

	fmted, err := f.Format(body)
	lf := byte('\n')
	if !bytes.HasSuffix(fmted, []byte{lf}) {
		fmted = append(fmted, lf)
	}
	return string(fmted), err
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err := b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewBuffer(buf.Bytes())), nil
}

func detectContentType(headers http.Header, body []byte) string {
	ct := headers.Get("Content-Type")
	if ct == "" {
		ct = http.DetectContentType(body)
	}
	vals := strings.Split(ct, ";")
	if len(vals) > 1 {
		ct = vals[0]
	}

	return ct
}
