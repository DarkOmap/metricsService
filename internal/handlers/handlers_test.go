package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DarkOmap/metricsService/internal/models"
	"github.com/DarkOmap/metricsService/internal/storage"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	textCT = "text/plain; charset=utf-8"
	jsonCT = "application/json; charset=utf-8"
)

type StorageMockedObject struct {
	mock.Mock
}

func (sm *StorageMockedObject) UpdateByMetrics(ctx context.Context, m models.Metrics) (*models.Metrics, error) {
	args := sm.Called(m)

	return args.Get(0).(*models.Metrics), args.Error(1)
}

func (sm *StorageMockedObject) ValueByMetrics(ctx context.Context, m models.Metrics) (*models.Metrics, error) {
	args := sm.Called(m)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Metrics), args.Error(1)
}

func (sm *StorageMockedObject) GetAllGauge(ctx context.Context) (map[string]storage.Gauge, error) {
	args := sm.Called()

	return args.Get(0).(map[string]storage.Gauge), args.Error(1)
}

func (sm *StorageMockedObject) GetAllCounter(ctx context.Context) (map[string]storage.Counter, error) {
	args := sm.Called()

	return args.Get(0).(map[string]storage.Counter), args.Error(1)
}

func (sm *StorageMockedObject) PingDB(ctx context.Context) error {
	args := sm.Called()

	return args.Error(0)
}

func (sm *StorageMockedObject) Updates(ctx context.Context, metrics []models.Metrics) error {
	args := sm.Called()

	return args.Error(0)
}

func TestServiceHandlers_updateByJSON(t *testing.T) {
	ms := new(StorageMockedObject)
	ms.On("UpdateByMetrics", *models.NewMetricsForGauge("test", 1.1)).Return(models.NewMetricsForGauge("test", 1.1), nil)
	ms.On("UpdateByMetrics", *models.NewMetricsForCounter("test", 1)).Return(models.NewMetricsForCounter("test", 1), nil)
	sh := NewServiceHandlers(ms)
	r := ServiceRouter(sh, "")

	srv := httptest.NewServer(r)
	defer srv.Close()

	type want struct {
		code              int
		contentType, body string
	}
	type param struct {
		method, body string
	}
	tests := []struct {
		name  string
		param param
		want  want
	}{
		{
			name:  "method get",
			param: param{method: http.MethodGet},
			want:  want{code: http.StatusMethodNotAllowed, contentType: ""},
		},
		{
			name: "wrong type value",
			param: param{method: http.MethodPost, body: `{
				"id": "test",
				"type": "type",
				"value": 0
			}`},
			want: want{code: http.StatusBadRequest, contentType: textCT},
		},
		{
			name: "wrong gauge value",
			param: param{method: http.MethodPost, body: `{
				"id": "test",
				"type": "gauge",
				"value": "0"
			}`},
			want: want{code: http.StatusBadRequest, contentType: textCT},
		},
		{
			name: "wrong counter value",
			param: param{method: http.MethodPost, body: `{
				"id": "test",
				"type": "counter",
				"delta": "0"
			}`},
			want: want{code: http.StatusBadRequest, contentType: textCT},
		},
		{
			name: "positive gauge",
			param: param{method: http.MethodPost, body: `{
				"id": "test",
				"type": "gauge",
				"value": 1.1
			}`},
			want: want{code: http.StatusOK, contentType: jsonCT, body: `{
				"id": "test",
				"type": "gauge",
				"value": 1.1
			}`},
		},
		{
			name: "positive counter",
			param: param{method: http.MethodPost, body: `{
				"id": "test",
				"type": "counter",
				"delta": 1
			}`},
			want: want{code: http.StatusOK, contentType: jsonCT, body: `{
				"id": "test",
				"type": "counter",
				"delta": 1
			}`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := testRequest(t, srv, tt.param.method, "/update", tt.param.body)
			assert.Equal(t, tt.want.code, res.StatusCode())
			assert.Equal(t, tt.want.contentType, strings.Join(res.Header().Values("Content-Type"), "; "))
			if tt.want.body != "" {
				assert.JSONEq(t, tt.want.body, string(res.Body()))
			}
		})
	}

	ms.AssertExpectations(t)
}

func TestServiceHandlers_updateByURL(t *testing.T) {
	ms := new(StorageMockedObject)
	ms.On("UpdateByMetrics", *models.NewMetricsForGauge("test", 12)).Return(models.NewMetricsForGauge("test", 12), nil)
	ms.On("UpdateByMetrics", *models.NewMetricsForCounter("test", 12)).Return(models.NewMetricsForCounter("test", 12), nil)
	sh := NewServiceHandlers(ms)
	r := ServiceRouter(sh, "")

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
			want:  want{http.StatusNotFound, textCT},
		},
		{
			name:  "wrong gauge value",
			param: param{http.MethodPost, "/update/gauge/test/wrong"},
			want:  want{http.StatusBadRequest, textCT},
		},
		{
			name:  "wrong counter value",
			param: param{http.MethodPost, "/update/counter/test/wrong"},
			want:  want{http.StatusBadRequest, textCT},
		},
		{
			name:  "positive gauge",
			param: param{http.MethodPost, "/update/gauge/test/12"},
			want:  want{http.StatusOK, textCT},
		},
		{
			name:  "positive counter",
			param: param{http.MethodPost, "/update/counter/test/12"},
			want:  want{http.StatusOK, textCT},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.param.method
			req.URL = srv.URL + tt.param.url

			res := testRequest(t, srv, tt.param.method, tt.param.url, "")
			assert.Equal(t, tt.want.code, res.StatusCode())
			assert.Equal(t, tt.want.contentType, strings.Join(res.Header().Values("Content-Type"), "; "))
		})
	}

	ms.AssertExpectations(t)
}

func testRequest(t *testing.T, srv *httptest.Server, method, url string, body string) *resty.Response {
	req := resty.New().R()
	req.Method = method
	req.URL = srv.URL + url
	req.SetBody(body)
	res, err := req.Send()
	assert.NoError(t, err)

	return res
}

func TestServiceHandlers_valueByURL(t *testing.T) {
	ms := new(StorageMockedObject)

	testWrong, err := models.NewMetrics("wrong", "counter")
	require.NoError(t, err)
	ms.On("ValueByMetrics", *testWrong).Return(nil, errors.New("unknown type wrong"))

	testGauge, err := models.NewMetrics("test", "gauge")
	require.NoError(t, err)
	ms.On("ValueByMetrics", *testGauge).Return(models.NewMetricsForGauge("test", 1.1), nil)

	testCounter, err := models.NewMetrics("test", "counter")
	require.NoError(t, err)
	ms.On("ValueByMetrics", *testCounter).Return(models.NewMetricsForCounter("test", 1), nil)

	sh := NewServiceHandlers(ms)
	r := ServiceRouter(sh, "")

	srv := httptest.NewServer(r)
	defer srv.Close()

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
			want:  want{code: http.StatusMethodNotAllowed},
		},
		{
			name:  "short url",
			param: param{http.MethodGet, "/value/gauge"},
			want:  want{code: http.StatusNotFound, contentType: textCT},
		},
		{
			name:  "wrong type",
			param: param{http.MethodGet, "/value/wrong/test"},
			want:  want{code: http.StatusNotFound, contentType: textCT},
		},
		{
			name:  "wrong name",
			param: param{http.MethodGet, "/value/counter/wrong"},
			want:  want{code: http.StatusNotFound, contentType: textCT},
		},
		{
			name:  "positive gauge",
			param: param{http.MethodGet, "/value/gauge/test"},
			want:  want{http.StatusOK, textCT, "1.1"},
		},
		{
			name:  "positive counter",
			param: param{http.MethodGet, "/value/counter/test"},
			want:  want{http.StatusOK, textCT, "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.param.method
			req.URL = srv.URL + tt.param.url

			res := testRequest(t, srv, tt.param.method, tt.param.url, "")
			assert.Equal(t, tt.want.code, res.StatusCode())

			if res.StatusCode() != http.StatusOK {
				return
			}

			assert.Equal(t, tt.want.contentType, strings.Join(res.Header().Values("Content-Type"), "; "))
			assert.Equal(t, tt.want.value, res.String())
		})
	}

	ms.AssertExpectations(t)
}

func TestServiceHandlers_all(t *testing.T) {
	ms := new(StorageMockedObject)
	ms.On("GetAllGauge").Return(map[string]storage.Gauge{"test": 1.1}, nil)
	ms.On("GetAllCounter").Return(map[string]storage.Counter{"test": 1}, nil)
	sh := NewServiceHandlers(ms)
	r := ServiceRouter(sh, "")

	srv := httptest.NewServer(r)
	defer srv.Close()

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
				<td>1.1</td>
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
			want:  want{code: http.StatusMethodNotAllowed},
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

			res := testRequest(t, srv, tt.param.method, tt.param.url, "")
			assert.Equal(t, tt.want.code, res.StatusCode())

			if res.StatusCode() != http.StatusOK {
				return
			}

			assert.Equal(t, tt.want.contentType, strings.Join(res.Header().Values("Content-Type"), "; "))
			assert.Equal(t, strings.ReplaceAll(strings.ReplaceAll(tt.want.value, "\t", ""), "\n", ""), strings.ReplaceAll(strings.ReplaceAll(res.String(), "\t", ""), "\n", ""))
		})
	}

	ms.AssertExpectations(t)
}

func TestServiceHandlers_valueByJSON(t *testing.T) {
	ms := new(StorageMockedObject)

	testCounterNotFound, err := models.NewMetrics("notFound", "counter")
	require.NoError(t, err)
	ms.On("ValueByMetrics", *testCounterNotFound).Return(nil, errors.New("error"))

	testGaugeNotFound, err := models.NewMetrics("notFound", "gauge")
	require.NoError(t, err)
	ms.On("ValueByMetrics", *testGaugeNotFound).Return(nil, errors.New("error"))

	testGauge, err := models.NewMetrics("test", "gauge")
	require.NoError(t, err)
	ms.On("ValueByMetrics", *testGauge).Return(models.NewMetricsForGauge("test", 1.1), nil)

	testCounter, err := models.NewMetrics("test", "counter")
	require.NoError(t, err)
	ms.On("ValueByMetrics", *testCounter).Return(models.NewMetricsForCounter("test", 1), nil)

	sh := NewServiceHandlers(ms)
	r := ServiceRouter(sh, "")

	srv := httptest.NewServer(r)
	defer srv.Close()

	type want struct {
		code              int
		contentType, body string
	}
	type param struct {
		method, body string
	}
	tests := []struct {
		name  string
		param param
		want  want
	}{
		{
			name:  "method get",
			param: param{method: http.MethodGet},
			want:  want{code: http.StatusMethodNotAllowed, contentType: ""},
		},
		{
			name: "gauge not found",
			param: param{method: http.MethodPost, body: `{
				"id": "notFound",
				"type": "gauge"
			}`},
			want: want{code: http.StatusNotFound, contentType: textCT},
		},
		{
			name: "counter not found",
			param: param{method: http.MethodPost, body: `{
				"id": "notFound",
				"type": "counter"
			}`},
			want: want{code: http.StatusNotFound, contentType: textCT},
		},
		{
			name: "positive gauge",
			param: param{method: http.MethodPost, body: `{
				"id": "test",
				"type": "gauge"
			}`},
			want: want{code: http.StatusOK, contentType: jsonCT, body: `{
				"id": "test",
				"type": "gauge",
				"value": 1.1
			}`},
		},
		{
			name: "positive counter",
			param: param{method: http.MethodPost, body: `{
				"id": "test",
				"type": "counter"
			}`},
			want: want{code: http.StatusOK, contentType: jsonCT, body: `{
				"id": "test",
				"type": "counter",
				"delta": 1
			}`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := testRequest(t, srv, tt.param.method, "/value", tt.param.body)
			assert.Equal(t, tt.want.code, res.StatusCode())
			assert.Equal(t, tt.want.contentType, strings.Join(res.Header().Values("Content-Type"), "; "))
			if tt.want.body != "" {
				assert.JSONEq(t, tt.want.body, string(res.Body()))
			}
		})
	}

	ms.AssertExpectations(t)
}

func TestServiceHandlers_ping(t *testing.T) {
	ms := new(StorageMockedObject)
	sh := NewServiceHandlers(ms)
	r := ServiceRouter(sh, "")

	srv := httptest.NewServer(r)
	defer srv.Close()

	ms.On("PingDB").Return(nil)

	t.Run("test positive", func(t *testing.T) {
		res := testRequest(t, srv, http.MethodGet, "/ping", "")
		assert.Equal(t, http.StatusOK, res.StatusCode())
	})
	ms.AssertExpectations(t)
	ms.On("PingDB").Unset()
	ms.On("PingDB").Return(errors.New("test"))

	t.Run("test negative", func(t *testing.T) {
		res := testRequest(t, srv, http.MethodGet, "/ping", "")
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode())
	})
	ms.AssertExpectations(t)
}
