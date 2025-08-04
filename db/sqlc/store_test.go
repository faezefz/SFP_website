package sqlc

import (
    "context"
    "testing"
)

func TestDummy(t *testing.T) {
    ctx := context.Background()
    if ctx == nil {
        t.Fatal("context is nil")
    }
}