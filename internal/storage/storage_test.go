package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_SetGauge(t *testing.T) {
	type fields struct {
		gauges   map[string]Gauge
		counters map[string]Counter
	}
	type args struct {
		value Gauge
		name  string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantFields fields
	}{
		{
			name:       "add gauge",
			fields:     fields{map[string]Gauge{"test": 0.12}, map[string]Counter{"test": 1}},
			args:       args{Gauge(1.11), "tg"},
			wantFields: fields{map[string]Gauge{"test": 0.12, "tg": 1.11}, map[string]Counter{"test": 1}},
		},
		{
			name:       "add gauge exchange",
			fields:     fields{map[string]Gauge{"test": 0.12}, map[string]Counter{"test": 1}},
			args:       args{Gauge(1.11), "test"},
			wantFields: fields{map[string]Gauge{"test": 1.11}, map[string]Counter{"test": 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				gauges:   tt.fields.gauges,
				counters: tt.fields.counters,
			}

			m.SetGauge(tt.args.value, tt.args.name)
			assert.Equal(t, tt.wantFields, fields{m.gauges, m.counters})
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	type fields struct {
		gauges   map[string]Gauge
		counters map[string]Counter
	}
	type args struct {
		value Counter
		name  string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantFields fields
	}{
		{
			name:       "add counter",
			fields:     fields{map[string]Gauge{"test": 0.12}, map[string]Counter{"test": 1}},
			args:       args{Counter(1), "tc"},
			wantFields: fields{map[string]Gauge{"test": 0.12}, map[string]Counter{"test": 1, "tc": 1}},
		},
		{
			name:       "add counter increment",
			fields:     fields{map[string]Gauge{"test": 0.12}, map[string]Counter{"test": 1}},
			args:       args{Counter(1), "test"},
			wantFields: fields{map[string]Gauge{"test": 0.12}, map[string]Counter{"test": 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				gauges:   tt.fields.gauges,
				counters: tt.fields.counters,
			}

			m.AddCounter(tt.args.value, tt.args.name)
			assert.Equal(t, tt.wantFields, fields{m.gauges, m.counters})
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type fields struct {
		gauges   map[string]Gauge
		counters map[string]Counter
	}
	tests := []struct {
		name    string
		fields  fields
		args    string
		want    Gauge
		wantErr bool
	}{
		{
			name:    "gauge not found",
			fields:  fields{map[string]Gauge{"test": 0.12}, map[string]Counter{"test": 1}},
			args:    "wrongName",
			wantErr: true,
		},
		{
			name:   "get gauge",
			fields: fields{map[string]Gauge{"test": 0.12}, map[string]Counter{"test": 1}},
			args:   "test",
			want:   Gauge(0.12),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				gauges:   tt.fields.gauges,
				counters: tt.fields.counters,
			}
			got, err := m.GetGauge(tt.args)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	type fields struct {
		gauges   map[string]Gauge
		counters map[string]Counter
	}
	tests := []struct {
		name    string
		fields  fields
		args    string
		want    Counter
		wantErr bool
	}{
		{
			name:    "counter not found",
			fields:  fields{map[string]Gauge{"test": 0.12}, map[string]Counter{"test": 1}},
			args:    "wrongName",
			wantErr: true,
		},
		{
			name:   "get counter",
			fields: fields{map[string]Gauge{"test": 0.12}, map[string]Counter{"test": 1}},
			args:   "test",
			want:   Counter(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				gauges:   tt.fields.gauges,
				counters: tt.fields.counters,
			}
			got, err := m.GetCounter(tt.args)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseGauge(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    Gauge
		wantErr bool
	}{
		{
			name:    "wront number",
			arg:     "test",
			wantErr: true,
		},
		{
			name: "test 0.001",
			arg:  "0.001",
			want: 0.001,
		},
		{
			name: "test 1.001",
			arg:  "1.001",
			want: 1.001,
		},
		{
			name: "test -1.001",
			arg:  "-1.001",
			want: -1.001,
		},
		{
			name: "test -0",
			arg:  "-0",
			want: -0,
		},
		{
			name: "test 0",
			arg:  "0",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGauge(tt.arg)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseCounter(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    Counter
		wantErr bool
	}{
		{
			name:    "wront number",
			arg:     "test",
			wantErr: true,
		},
		{
			name: "test 1",
			arg:  "1",
			want: 1,
		},
		{
			name: "test -1",
			arg:  "-1",
			want: -1,
		},
		{
			name: "test 11",
			arg:  "11",
			want: 11,
		},
		{
			name: "test -0",
			arg:  "-0",
			want: -0,
		},
		{
			name: "test 0",
			arg:  "0",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCounter(tt.arg)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
