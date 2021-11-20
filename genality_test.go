package genality

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/infiniteloopcloud/genality/migration"
)

const (
	testConnection = "postgres://genality:genality@localhost:5435/genality?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(driverName, testConnection)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	if err := migration.Up(ctx, conn); err != nil {
		log.Fatal(err)
	}
	exitVal := m.Run()
	//migration.Down(ctx, conn)

	os.Exit(exitVal)
}

func TestExample(t *testing.T) {
	random := uuid.New().String()
	m, err := newGenality(Opts{
		ConnectionString: testConnection,
	})
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
		if err := m.add(context.Background(), random, tim); err != nil {
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

	resp, err := m.GetCountBuckets(context.Background(), random, time.Now().UTC().Add(-10*time.Hour), time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
