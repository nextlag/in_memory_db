package test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/nextlag/in_memory_db/internal/errs"
	"github.com/nextlag/in_memory_db/internal/server/repository"
	"github.com/nextlag/in_memory_db/internal/server/usecase/storage/engine/in_memory"
)

func TestInMemoryRepository(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	engine, err := in_memory.NewEngine(log)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	repo, err := repository.New(engine, log)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	ctx := context.Background()

	// Test Set/Get
	t.Run("Set and Get", func(t *testing.T) {
		key := "test-key"
		value := "test-value"

		// Set value
		err := repo.Set(ctx, key, value)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Get value
		retrieved, err := repo.Get(ctx, key)
		if err != nil {
			t.Fatalf("Failed to get value: %v", err)
		}

		if retrieved != value {
			t.Errorf("Expected value %q, got %q", value, retrieved)
		}
	})

	// Test Get non-existent key
	t.Run("Get non-existent key", func(t *testing.T) {
		_, err := repo.Get(ctx, "non-existent-key")
		if err == nil {
			t.Error("Expected error for non-existent key, got nil")
		}
		if !errs.IsNotFound(err) {
			t.Errorf("Expected NotFound error, got %v", err)
		}
	})

	// Test Del
	t.Run("Delete key", func(t *testing.T) {
		key := "delete-test-key"
		value := "delete-test-value"

		// Set value
		err := repo.Set(ctx, key, value)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Delete key
		err = repo.Del(ctx, key)
		if err != nil {
			t.Fatalf("Failed to delete key: %v", err)
		}

		// Verify key is deleted
		_, err = repo.Get(ctx, key)
		if err == nil {
			t.Error("Expected error for deleted key, got nil")
		}
		if !errs.IsNotFound(err) {
			t.Errorf("Expected NotFound error, got %v", err)
		}
	})

	// Test Del non-existent key
	t.Run("Delete non-existent key", func(t *testing.T) {
		err := repo.Del(ctx, "non-existent-key")
		if err == nil {
			t.Error("Expected error for deleting non-existent key, got nil")
		}
		if !errs.IsNotFound(err) {
			t.Errorf("Expected NotFound error, got %v", err)
		}
	})

	// Test with canceled context
	t.Run("Operations with canceled context", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		err := repo.Set(canceledCtx, "key", "value")
		if err == nil {
			t.Error("Expected error for Set with canceled context, got nil")
		}

		_, err = repo.Get(canceledCtx, "key")
		if err == nil {
			t.Error("Expected error for Get with canceled context, got nil")
		}

		err = repo.Del(canceledCtx, "key")
		if err == nil {
			t.Error("Expected error for Del with canceled context, got nil")
		}
	})
}
