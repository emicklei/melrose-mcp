package main

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/emicklei/melrose-mcp/mcpserver"
	"github.com/emicklei/melrose/system"
)

func TestHandleCDE(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("Skipping test in CI environment")
	}
	ctx, err := system.Setup("test")
	if err != nil {
		log.Fatalln(err)
	}
	defer ctx.Device().Close()
	playServer := mcpserver.NewMCPServer(ctx)

	result, _, err := playServer.HandlePlay(context.Background(), nil, mcpserver.PlayParams{
		Expression: `
a=note('c')
b=a+a`,
	})
	if err != nil {
		t.Fatalf("Handle failed: %v", err)
	}
	if result.IsError {
		t.Fatalf("Handle returned error: %v", result)
	}
	t.Log("Handle result:", result)
}
