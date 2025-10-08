package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/nutchapon-m/web-server/foundation/logger"
)

var defaultClient = http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 15 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

type Header struct {
	Key   string
	Value string
}

type Query struct {
	Key   string
	Value string
}

type Response struct {
	StatusCode int
	Body       []byte
	Entries    *http.Response
}

func (r Response) Decode(v any) error {
	return json.Unmarshal(r.Body, v)
}

type Client struct {
	log     *logger.Logger
	url     string
	http    *http.Client
	headers []Header
	queries []Query
}

func NewHttpClient(log *logger.Logger, options ...func(cln *Client)) *Client {
	cln := Client{
		log:  log,
		http: &defaultClient,
	}

	for _, option := range options {
		option(&cln)
	}
	return &cln
}

// WithClient adds a custom client for processing requests. It's recommend
// to not use the default client and provide your own.
func WithClient(http *http.Client) func(cln *Client) {
	return func(cln *Client) {
		cln.http = http
	}
}

func (c Client) URL() string {
	return c.url
}

func (c *Client) Header(key, v string) {
	c.headers = append(c.headers, Header{Key: key, Value: v})
}

func (c *Client) QueryParam(key, v string) {
	c.queries = append(c.queries, Query{Key: key, Value: v})
}

func (c *Client) request(method, url string, body io.Reader) (*http.Request, error) {
	finalURL := url
	if c.url != "" {
		finalURL = c.url + url
	}

	req, err := http.NewRequest(method, finalURL, body)
	if err != nil {
		return nil, err
	}

	if len(c.headers) > 0 {
		for _, h := range c.headers {
			req.Header.Set(h.Key, h.Value)
		}
	}

	if len(c.queries) > 0 {
		q := req.URL.Query()
		for _, query := range c.queries {
			q.Set(query.Key, query.Value)
		}
		req.URL.RawQuery = q.Encode()
	}

	return req, nil
}

func (c *Client) Get(url string) (Response, error) {
	req, err := c.request(http.MethodGet, url, http.NoBody)
	if err != nil {
		return Response{}, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return Response{}, err
	}

	defer res.Body.Close()
	buff, err := io.ReadAll(res.Body)
	if err != nil {
		return Response{}, err
	}

	response := Response{
		StatusCode: res.StatusCode,
		Body:       buff,
		Entries:    res,
	}
	return response, nil
}

func (c *Client) Post(url string, body any) (Response, error) {
	buff, err := json.Marshal(body)
	if err != nil {
		return Response{}, err
	}

	req, err := c.request(http.MethodPost, url, bytes.NewBuffer(buff))
	if err != nil {
		return Response{}, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return Response{}, err
	}

	defer res.Body.Close()

	buff, err = io.ReadAll(res.Body)
	if err != nil {
		return Response{}, err
	}

	response := Response{
		StatusCode: res.StatusCode,
		Body:       buff,
		Entries:    res,
	}
	return response, nil
}
