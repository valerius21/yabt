package client

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

type HttpClient struct {
	client    *fasthttp.Client
	targetURL string
}

func NewClient(baseURL string) *HttpClient {
	client := &fasthttp.Client{}
	return &HttpClient{targetURL: baseURL, client: client}
}

func (c *HttpClient) Get() *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(c.targetURL)
	resp := fasthttp.AcquireResponse()
	start := time.Now()
	err := c.client.Do(req, resp)
	end := time.Since(start)
	log.Info().Msgf("POST request to %s completed - Duration: %v", c.targetURL, end)
	if err != nil {
		panic(err)
	}
	return resp
}

func (c *HttpClient) SendFile(file *os.File) *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(c.targetURL)
	req.Header.SetMethod(fasthttp.MethodPost)
	body, contentType := readFile(file)
	req.Header.SetContentType(contentType)
	req.SetBodyString(body.String())

	resp := fasthttp.AcquireResponse()
	start := time.Now()
	err := c.client.Do(req, resp)
	end := time.Since(start)
	log.Info().Msgf("POST request to %s completed - Duration: %v", c.targetURL, end)
	if err != nil {
		panic(err)
	}
	return resp
}

func readFile(file *os.File) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	formFile, err := writer.CreateFormFile("image", file.Name())
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(formFile, file)
	if err != nil {
		panic(err)
	}
	writer.Close()
	contentType := writer.FormDataContentType()
	return body, contentType
}
