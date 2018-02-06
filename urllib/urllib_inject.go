// +build unittest

package urllib

import (
	"fmt"
	"io"
	"net/http"
)

// MockBody ...
type MockBody struct {
	n        int
	body     []byte
	readErr  error
	closeErr error
}

// Read ...
func (fb *MockBody) Read(p []byte) (n int, err error) {
	n = copy(p[fb.n:], fb.body)
	fb.n += n
	return n, fb.readErr
}

// Close ...
func (fb *MockBody) Close() error {
	return fb.closeErr
}

type MockResponse struct {
	Count    int
	MockFunc func(m *MockResponse, b *HttpRequest) (*http.Response, error)
}

var MockResp *MockResponse = &MockResponse{}

func (m *MockResponse) MockError(err error) (*http.Response, error) {
	return nil, err
}

func (m *MockResponse) Reset() {
	m.Count = 0
}

func (m *MockResponse) MockNoBody() (*http.Response, error) {
	resp := &http.Response{
		Body:          nil,
		Close:         true,
		ContentLength: 0,
		Header:        http.Header{},
		Proto:         "",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		Uncompressed:  true,
	}
	return resp, nil
}

func (m *MockResponse) MockReadBodyError() (*http.Response, error) {
	resp, _ := m.MockNoBody()
	resp.Body = &MockBody{
		readErr: fmt.Errorf("read error on mock body"),
	}
	return resp, nil
}

func (m *MockResponse) Mock5xxError() (*http.Response, error) {
	resp, _ := m.MockNoBody()
	resp.Body = &MockBody{
		readErr: io.EOF,
	}
	resp.StatusCode = http.StatusInternalServerError
	return resp, nil
}

func (b *HttpRequest) GetResponse() (*http.Response, error) {
	return b.getResponse()
}

// Response executes request client gets response mannually.
func (b *HttpRequest) Response() (*http.Response, error) {
	defer func() {
		MockResp.Count += 1
	}()
	if MockResp.MockFunc != nil {
		resp, err := MockResp.MockFunc(MockResp, b)
		return resp, err
	}
	return b.getResponse()
}
