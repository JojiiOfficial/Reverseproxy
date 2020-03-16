package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// Get 403 forbidden response
func getForbiddenResponse(req *http.Request) *http.Response {
	return buildResponse(req, http.StatusForbidden, "403 Forbidden", "403 Forbidden", nil)
}

// Build http response
func buildResponse(req *http.Request, statusCode int, body, status string, header http.Header) *http.Response {
	response := http.Response{
		Status:        status,
		StatusCode:    statusCode,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Request:       req,
		Header:        header,
	}

	if header != nil {
		response.Header = header
	}

	return &response
}
