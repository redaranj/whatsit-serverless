package main

import (
	"context"
	"testing"
)

func TestHandler(t *testing.T) {
	ctx := context.Background()
	res := handler(ctx, nil)
}
