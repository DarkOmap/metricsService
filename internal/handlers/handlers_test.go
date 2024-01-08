package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DarkOmap/metricsService/internal/storage"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

const defaultCT = "text/plain; charset=utf-8"

func TestUpdate(t *testing.T) {
	ms := storage.NewMemStorage()
	sh := NewServiceHandlers(ms)
	r := ServiceRouter(sh)

	srv := httptest.NewServer(r)
	defer srv.Close()

	type want struct {
		code        int
		contentType string
	}
	type param struct {
		method, url string
	}
	tests := []struct {
		name  string
		param param
		want  want
	}{
		{
			name:  "method get",
			param: param{http.MethodGet, "/update/gauge/test/123"},
			want:  want{http.StatusMethodNotAllowed, ""},
		},
		{
			name:  "short url",
			param: param{http.MethodPost, "/update/gauge/test"},
			want:  want{http.StatusNotFound, defaultCT},
		},
		{
			name:  "wrong gauge value",
			param: param{http.MethodPost, "/update/gauge/test/wrong"},
			want:  want{http.StatusBadRequest, defaultCT},
		},
		{
			name:  "wrong counter value",
			param: param{http.MethodPost, "/update/counter/test/wrong"},
			want:  want{http.StatusBadRequest, defaultCT},
		},
		{
			name:  "positive gauge",
			param: param{http.MethodPost, "/update/gauge/test/12"},
			want:  want{http.StatusOK, defaultCT},
		},
		{
			name:  "positive counter",
			param: param{http.MethodPost, "/update/counter/test/12"},
			want:  want{http.StatusOK, defaultCT},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.param.method
			req.URL = srv.URL + tt.param.url

			res := testRequest(t, srv, tt.param.method, tt.param.url)
			assert.Equal(t, tt.want.code, res.StatusCode())
			assert.Equal(t, tt.want.contentType, strings.Join(res.Header().Values("Content-Type"), "; "))
		})
	}
}

func testRequest(t *testing.T, srv *httptest.Server, method, url string) *resty.Response {
	req := resty.New().R()
	req.Method = method
	req.URL = srv.URL + url

	res, err := req.Send()
	assert.NoError(t, err)

	return res
}

func Test_value(t *testing.T) {
	ms := storage.NewMemStorage()
	sh := NewServiceHandlers(ms)
	r := ServiceRouter(sh)

	srv := httptest.NewServer(r)
	defer srv.Close()

	g, _ := storage.ParseGauge("1")
	sh.ms.SetGauge(g, "test")

	c, _ := storage.ParseCounter("1")
	sh.ms.AddCounter(c, "test")

	type want struct {
		code        int
		contentType string
		value       string
	}
	type param struct {
		method, url string
	}
	tests := []struct {
		name  string
		param param
		want  want
	}{
		{
			name:  "method post",
			param: param{http.MethodPost, "/value/gauge/test"},
			want:  want{http.StatusMethodNotAllowed, "", ""},
		},
		{
			name:  "short url",
			param: param{http.MethodGet, "/value/gauge"},
			want:  want{http.StatusNotFound, defaultCT, "404 page not found"},
		},
		{
			name:  "wrong type",
			param: param{http.MethodGet, "/value/wrong/test"},
			want:  want{http.StatusNotFound, defaultCT, "unknown type"},
		},
		{
			name:  "wrong name",
			param: param{http.MethodGet, "/value/counter/wrong"},
			want:  want{http.StatusNotFound, defaultCT, "value not found"},
		},
		{
			name:  "positive gauge",
			param: param{http.MethodGet, "/value/gauge/test"},
			want:  want{http.StatusOK, defaultCT, "1"},
		},
		{
			name:  "positive counter",
			param: param{http.MethodGet, "/value/counter/test"},
			want:  want{http.StatusOK, defaultCT, "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.param.method
			req.URL = srv.URL + tt.param.url

			res := testRequest(t, srv, tt.param.method, tt.param.url)
			assert.Equal(t, tt.want.code, res.StatusCode())
			assert.Equal(t, tt.want.contentType, strings.Join(res.Header().Values("Content-Type"), "; "))
			assert.Equal(t, tt.want.value, res.String())
		})
	}
}

func Test_all(t *testing.T) {
	ms := storage.NewMemStorage()
	sh := NewServiceHandlers(ms)
	r := ServiceRouter(sh)

	srv := httptest.NewServer(r)
	defer srv.Close()

	g, _ := storage.ParseGauge("1")
	sh.ms.SetGauge(g, "test")

	c, _ := storage.ParseCounter("1")
	sh.ms.AddCounter(c, "test")

	htmlText := `<!DOCTYPE html>
	<html>
	
	<head>
		<meta charset="UTF-8">
		<title>add data from service</title>
	</head>
	
	<body>
		<table>
			<tr>
				<th>name</th>
				<th>value</th>
			</tr>
			<tr>
				<td>test</td>
				<td>1</td>
			</tr>
			<tr>
				<td>test</td>
				<td>1</td>
			</tr>
		</table>
	</body>
	
	</html>`

	type want struct {
		code        int
		contentType string
		value       string
	}
	type param struct {
		method, url string
	}
	tests := []struct {
		name  string
		param param
		want  want
	}{
		{
			name:  "method post",
			param: param{http.MethodPost, "/"},
			want:  want{http.StatusMethodNotAllowed, "", ""},
		},
		{
			name:  "wrong name",
			param: param{http.MethodGet, "/value/counter/wrong"},
			want:  want{http.StatusNotFound, defaultCT, "value not found"},
		},
		{
			name:  "positive",
			param: param{http.MethodGet, "/"},
			want:  want{http.StatusOK, "text/html; charset=utf-8", htmlText},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.param.method
			req.URL = srv.URL + tt.param.url

			res := testRequest(t, srv, tt.param.method, tt.param.url)
			assert.Equal(t, tt.want.code, res.StatusCode())
			assert.Equal(t, tt.want.contentType, strings.Join(res.Header().Values("Content-Type"), "; "))
			assert.Equal(t, strings.ReplaceAll(strings.ReplaceAll(tt.want.value, "\t", ""), "\n", ""), strings.ReplaceAll(strings.ReplaceAll(res.String(), "\t", ""), "\n", ""))
		})
	}
}
