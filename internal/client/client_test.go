package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendGauge(t *testing.T) {
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "not OK test",
			args: args{"test", 1.1},
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "test error", http.StatusBadRequest)
			},
			wantErr: true,
		},
		{
			name: "OK test",
			args: args{"test", 1.1},
			handler: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(tt.handler))
			defer ts.Close()

			c := NewClient(strings.TrimPrefix(ts.URL, "http://"), "")
			err := c.SendGauge(context.Background(), tt.args.name, tt.args.value)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}

	t.Run("test brocken server", func(t *testing.T) {
		c := NewClient("test", "")
		err := c.SendGauge(context.Background(), "test", 1.1)

		assert.Error(t, err)
	})
}

func TestSendCounter(t *testing.T) {
	type args struct {
		name  string
		delta int64
	}
	tests := []struct {
		name    string
		args    args
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "not OK test",
			args: args{"test", 1},
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "test error", http.StatusBadRequest)
			},
			wantErr: true,
		},
		{
			name: "OK test",
			args: args{"test", 1},
			handler: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(tt.handler))
			defer ts.Close()

			c := NewClient(strings.TrimPrefix(ts.URL, "http://"), "")
			err := c.SendCounter(context.Background(), tt.args.name, tt.args.delta)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}

	t.Run("test brocken server", func(t *testing.T) {
		c := NewClient("test", "")
		err := c.SendCounter(context.Background(), "test", 1)

		assert.Error(t, err)
	})
}
