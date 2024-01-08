package client

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	memstats "github.com/DarkOmap/metricsService/internal/memstats"
	"github.com/stretchr/testify/assert"
)

func TestSendGauge(t *testing.T) {
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "not OK test",
			args: args{"test", "test"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "test error", http.StatusBadRequest)
			},
			wantErr: true,
		},
		{
			name: "OK test",
			args: args{"test", "test"},
			handler: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(tt.handler))
			defer ts.Close()

			serviceAddr = strings.TrimPrefix(ts.URL, "http://")

			err := sendGauge(tt.args.name, tt.args.value)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestSendCounter(t *testing.T) {
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "not OK test",
			args: args{"test", "test"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "test error", http.StatusBadRequest)
			},
			wantErr: true,
		},
		{
			name: "OK test",
			args: args{"test", "test"},
			handler: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(tt.handler))
			defer ts.Close()

			serviceAddr = strings.TrimPrefix(ts.URL, "http://")

			err := sendCounter(tt.args.name, tt.args.value)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestPushStats(t *testing.T) {
	type args struct {
		ms []memstats.StringMS
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "not OK test",
			args: args{[]memstats.StringMS{
				{Name: "test", Value: "1111"},
				{Name: "test2", Value: "2222"},
			}},
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "test error", http.StatusBadRequest)
			},
			wantErr: true,
		},
		{
			name: "OK test",
			args: args{[]memstats.StringMS{
				{Name: "test", Value: "1111"},
				{Name: "test2", Value: "2222"},
			}},
			handler: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(tt.handler))
			defer ts.Close()

			serviceAddr = strings.TrimPrefix(ts.URL, "http://")

			err := pushStats(tt.args.ms)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
