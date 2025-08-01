package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {

	ctx := context.Background()

	// Create a new client, with no features.
	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	// Connect to a server over stdin/stdout
	transport := mcp.NewCommandTransport(exec.Command("/Users/emicklei/go/bin/melrose-mcp"))
	session, err := client.Connect(ctx, transport)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		fmt.Println("Closing client...")
		time.Sleep(2 * time.Second)
		session.Close()
	}()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// List Tools
	fmt.Println("Listing available tools...")

	{
		params := &mcp.ListToolsParams{}
		res, err := session.ListTools(ctx, params)
		if err != nil {
			log.Fatalf("CallTool failed: %v", err)
		}
		for _, c := range res.Tools {
			log.Print(c.Name, " - ", c.Description)
		}
	}

	{
		// Call a tool on the server.
		params := &mcp.CallToolParams{
			Name:      "melrose_play",
			Arguments: map[string]any{"expression": "sequence('A B C')"},
		}
		res, err := session.CallTool(ctx, params)
		if err != nil {
			log.Fatalf("CallTool failed: %v", err)
		}
		if res.IsError {
			log.Fatal("tool failed")
		}
		for _, c := range res.Content {
			log.Print(c.(*mcp.TextContent).Text)
		}
	}
}
