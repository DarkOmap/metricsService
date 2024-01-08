package storage

import (
	"errors"
	"strconv"
)

type Gauge float64
type Counter int64

type Repositories interface {
	SetGauge(value Gauge, name string)
	GetGauge(name string) (Gauge, error)
	AddCounter(value Counter, name string)
	GetCounter(name string) (Counter, error)
	GetAllGauge() map[string]Gauge
	GetAllCounter() map[string]Counter
}

type MemStorage struct {
	gauges   map[string]Gauge
	counters map[string]Counter
}

func NewMemStorage() *MemStorage {
	ms := MemStorage{}
	ms.counters = make(map[string]Counter)
	ms.gauges = make(map[string]Gauge)

	return &ms
}

func (m *MemStorage) SetGauge(value Gauge, name string) {
	m.gauges[name] = value
}

func (m *MemStorage) GetGauge(name string) (Gauge, error) {
	v, ok := m.gauges[name]

	if !ok {
		return v, errors.New("value not found")
	}

	return v, nil
}

func (m *MemStorage) AddCounter(value Counter, name string) {
	m.counters[name] += value
}

func (m *MemStorage) GetCounter(name string) (Counter, error) {
	v, ok := m.counters[name]

	if !ok {
		return v, errors.New("value not found")
	}

	return v, nil
}

func (m *MemStorage) GetAllGauge() map[string]Gauge {
	return m.gauges
}

func (m *MemStorage) GetAllCounter() map[string]Counter {
	return m.counters
}

func ParseGauge(g string) (Gauge, error) {
	v, err := strconv.ParseFloat(g, 64)

	return Gauge(v), err
}

func ParseCounter(c string) (Counter, error) {
	v, err := strconv.ParseInt(c, 10, 64)

	return Counter(v), err
}
