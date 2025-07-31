package main

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/emicklei/melrose-mcp/mcpserver"
	"github.com/emicklei/melrose/system"
	"github.com/modelcontextprotocol/go-sdk/mcp"
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

	req := &mcp.CallToolParamsFor[mcpserver.PlayParams]{
		Name: "play-melrose",
		Arguments: mcpserver.PlayParams{
			Expression: `
a=note('c')
b=a+a`,
		},
	}
	result, err := playServer.HandlePlay(context.Background(), nil, req)
	if err != nil {
		t.Fatalf("Handle failed: %v", err)
	}
	if result.IsError {
		t.Fatalf("Handle returned error: %v", result)
	}
	t.Log("Handle result:", result)
}
