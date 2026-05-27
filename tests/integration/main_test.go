//go:build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/nimbleflux/clickhouse-query-analyzer/internal/clickhouse"
)

func clickhouseURL() string {
	if v := os.Getenv("CLICKHOUSE_URL"); v != "" {
		return v
	}
	return "clickhouse://localhost:19000"
}

func newPool(t *testing.T) *clickhouse.Pool {
	t.Helper()
	pool := clickhouse.NewPool()
	t.Cleanup(pool.CloseAll)
	return pool
}

func newClient(t *testing.T) *clickhouse.Client {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	pool := newPool(t)
	c, err := pool.Get(ctx, clickhouse.ConnParams{
		URL:      clickhouseURL(),
		User:     "default",
		Database: "system",
	})
	if err != nil {
		t.Fatalf("failed to connect to ClickHouse: %v", err)
	}
	return c
}

func TestConnect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool := newPool(t)
	_, err := pool.Get(ctx, clickhouse.ConnParams{
		URL:      clickhouseURL(),
		User:     "default",
		Database: "system",
	})
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
}

func TestConnect_PoolCaching(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool := newPool(t)
	params := clickhouse.ConnParams{
		URL:      clickhouseURL(),
		User:     "default",
		Database: "system",
	}

	c1, err := pool.Get(ctx, params)
	if err != nil {
		t.Fatalf("first connect failed: %v", err)
	}
	c2, err := pool.Get(ctx, params)
	if err != nil {
		t.Fatalf("second connect failed: %v", err)
	}
	if c1 != c2 {
		t.Error("expected same client from pool cache")
	}
}
