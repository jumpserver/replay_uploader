package httplib

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type AuthSign interface {
	Sign(req *http.Request)
}

func NewClient(baseUrl string, timeout time.Duration) (*Client, error) {
	_, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	jar := &customCookieJar{
		data: map[string]string{},
	}
	con := http.Client{
		Timeout: timeout,
		Jar:     jar,
	}
	return &Client{
		Timeout: 0,
		baseUrl: baseUrl,
		cookies: make(map[string]string),
		headers: make(map[string]string),
		http:    &con,
	}, nil
}

type Client struct {
	Timeout  time.Duration
	baseUrl  string
	cookies  map[string]string
	headers  map[string]string
	http     *http.Client
	authSign AuthSign
}

func (c *Client) SetCookie(key string, value string) {
	c.cookies[key] = value
}

func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

func (c *Client) SetAuthSign(auth AuthSign) {
	c.authSign = auth
}

func (c *Client) setAuthHeader(r *http.Request) {
	if len(c.cookies) != 0 {
		for k, v := range c.cookies {
			co := http.Cookie{Name: k, Value: v}
			r.AddCookie(&co)
		}
	}
	if c.authSign != nil {
		c.authSign.Sign(r)
	}
}

func (c *Client) setReqHeaders(req *http.Request, params []map[string]string) {
	if len(c.headers) != 0 {
		for k, v := range c.headers {
			req.Header.Set(k, v)
		}
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", "koko-client")
	c.setAuthHeader(req)
	if len(params) >= 2 {
		for k, v := range params[1] {
			req.Header.Set(k, v)
		}
	}
}

func (c *Client) parseUrl(reqUrl string, params []map[string]string) string {
	if len(params) < 1 {
		return reqUrl
	}
	query := url.Values{}
	for _, item := range params {
		for k, v := range item {
			query.Add(k, v)
		}
	}
	if strings.Contains(reqUrl, "?") {
		reqUrl += "&" + query.Encode()
	} else {
		reqUrl += "?" + query.Encode()
	}
	return reqUrl
}

func (c *Client) newRequest(method, reqUrl string, data interface{}, params []map[string]string) (*http.Request, error) {
	reqUrl = c.parseUrl(reqUrl, params)
	if c.baseUrl != "" {
		reqUrl = strings.TrimRight(c.baseUrl, "/") + reqUrl
	}
	dataRaw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(dataRaw)
	req, err := http.NewRequest(method, reqUrl, reader)
	if err != nil {
		return req, err
	}
	c.setReqHeaders(req, params)
	return req, nil
}

func (c *Client) Do(method, reqUrl string, data, res interface{}, params ...map[string]string) (resp *http.Response, err error) {
	req, err := c.newRequest(method, reqUrl, data, params)
	if err != nil {
		return
	}
	resp, err = c.http.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		msg := fmt.Sprintf("%s %s failed, get code: %d, %s", req.Method, req.URL, resp.StatusCode, body)
		err = errors.New(msg)
		return
	}

	// If is buffer return the raw response body
	if buf, ok := res.(*bytes.Buffer); ok {
		buf.Write(body)
		return
	}
	// Unmarshal response body to result struct
	if res != nil {
		switch {
		case strings.Contains(resp.Header.Get("Content-Type"), "application/json"):
			err = json.Unmarshal(body, res)
			if err != nil {
				msg := fmt.Sprintf("%s %s failed, unmarshal '%s' response failed: %s", req.Method, req.URL, body[:12], err)
				err = errors.New(msg)
				return
			}
		}
	}
	return
}

func (c *Client) Get(url string, res interface{}, params ...map[string]string) (resp *http.Response, err error) {
	return c.Do("GET", url, nil, res, params...)
}

func (c *Client) Post(url string, data interface{}, res interface{}, params ...map[string]string) (resp *http.Response, err error) {
	return c.Do("POST", url, data, res, params...)
}

func (c *Client) Delete(url string, res interface{}, params ...map[string]string) (resp *http.Response, err error) {
	return c.Do("DELETE", url, nil, res, params...)
}

func (c *Client) Put(url string, data interface{}, res interface{}, params ...map[string]string) (resp *http.Response, err error) {
	return c.Do("PUT", url, data, res, params...)
}

func (c *Client) Patch(url string, data interface{}, res interface{}, params ...map[string]string) (resp *http.Response, err error) {
	return c.Do("PATCH", url, data, res, params...)
}

func (c *Client) UploadFile(url string, gFile string, res interface{}, params ...map[string]string) (err error) {
	f, err := os.Open(gFile)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(buf)
	gName := filepath.Base(gFile)
	part, err := bodyWriter.CreateFormFile("file", gName)
	if err != nil {
		return err
	}
	if _, err = io.Copy(part, f); err != nil {
		return err
	}
	err = bodyWriter.Close()
	if err != nil {
		return err
	}
	url = c.parseUrl(url, params)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	c.setReqHeaders(req, params)
	/*
		上传文件时，取消 timeout
		A Timeout of zero means no timeout.
	*/
	client := http.Client{
		Jar: c.http.Jar,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		msg := fmt.Sprintf("%s %s failed, get code: %d, %s", req.Method, req.URL, resp.StatusCode, string(body))
		err = errors.New(msg)
		return
	}

	// If is buffer return the raw response body
	if buf, ok := res.(*bytes.Buffer); ok {
		buf.Write(body)
		return
	}
	// Unmarshal response body to result struct
	if res != nil {
		err = json.Unmarshal(body, res)
		if err != nil {
			msg := fmt.Sprintf("%s %s failed, unmarshal '%s' response failed: %s", req.Method, req.URL, body, err)
			err = errors.New(msg)
			return
		}
	}
	return
}

func (c *Client) UploadMultiPartFile(url string, gFile string, res interface{}, params ...map[string]string) (err error) {
	f, err := os.Open(gFile)
	if err != nil {
		return err
	}
	defer f.Close()
	bufferedFileReader := bufio.NewReader(f)
	bodyReader, bodyWriter := io.Pipe()
	formWriter := multipart.NewWriter(bodyWriter)

	// Store the first write error in writeErr.
	var (
		writeErr error
		errOnce  sync.Once
	)
	setErr := func(err error) {
		if err != nil {
			errOnce.Do(func() { writeErr = err })
		}
	}
	go func() {
		gName := filepath.Base(gFile)
		partWriter, err := formWriter.CreateFormFile("file", gName)
		if err != nil {
			setErr(err)
			return
		}
		defer bodyWriter.Close()
		defer formWriter.Close()
		if _, err := io.Copy(partWriter, bufferedFileReader); err != nil {
			setErr(err)
		}

	}()
	url = c.parseUrl(url, params)
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return err
	}
	c.setReqHeaders(req, params)
	req.Header.Add("Content-Type", formWriter.FormDataContentType())

	// This operation will block until both the formWriter
	// and bodyWriter have been closed by the goroutine,
	// or in the event of a HTTP error.
	client := http.Client{
		Jar: c.http.Jar,
	}
	resp, err := client.Do(req)
	if writeErr != nil {
		return writeErr
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		msg := fmt.Sprintf("%s %s failed, get code: %d, %s", req.Method, req.URL, resp.StatusCode, string(body))
		err = errors.New(msg)
		return
	}

	// If is buffer return the raw response body
	if buf, ok := res.(*bytes.Buffer); ok {
		buf.Write(body)
		return
	}
	// Unmarshal response body to result struct
	if res != nil {
		err = json.Unmarshal(body, res)
		if err != nil {
			msg := fmt.Sprintf("%s %s failed, unmarshal '%s' response failed: %s", req.Method, req.URL, body, err)
			err = errors.New(msg)
			return
		}
	}
	return
}
