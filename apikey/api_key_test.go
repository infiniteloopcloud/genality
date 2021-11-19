package apikey

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestExample(t *testing.T) {
	random := uuid.New().String()
	connStr := "postgres://api_key:api_key@localhost:5435/api_key?sslmode=disable"
	m, err := NewMetrics(connStr)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		if err := m.Add(context.Background(), random); err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < 50; i++ {
		tim := time.Now().UTC().Add(time.Duration(-(i * 15)) * time.Minute)
		if err := m.add(context.Background(),random, tim); err != nil {
			t.Fatal(err)
		}
	}

	count, err := m.GetCountFrom(context.Background(), random, time.Now().UTC().Add(-1*time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	if count != 104 {
		t.Fatal("count should be 104")
	}

	resp, err := m.GetCountBuckets(context.Background(),random, time.Now().UTC().Add(-10*time.Hour),time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
