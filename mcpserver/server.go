package mcpserver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/emicklei/melrose/api"
	"github.com/emicklei/melrose/core"
	"github.com/emicklei/melrose/notify"
	"github.com/mark3labs/mcp-go/mcp"
)

type MCPServer struct {
	service api.Service
}

func NewMCPServer(ctx core.Context) *MCPServer {
	return &MCPServer{service: api.NewService(ctx)}
}

func (s *MCPServer) HandleChangeOutputDevice(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetInt("id", 0)
	channel := request.GetInt("channel", 1)
	if id < 1 || id > 16 {
		return nil, fmt.Errorf("id must be a number between 1 and 16")
	}
	err := s.service.ChangeDefaultDeviceAndChannel(false, id, channel)
	toolResult := new(mcp.CallToolResult)
	if err != nil {
		toolResult.IsError = true
		toolResult.Content = []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: err.Error(),
			},
		}
	} else {
		toolResult.Content = []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Output device is set to %d", id),
			},
		}
	}
	return toolResult, nil
}

func (s *MCPServer) HandleBPM(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	bpm := request.GetFloat("bpm", 120)
	if bpm < 1 || bpm > 300 {
		return nil, errors.New("parameter must be positive number between 1 and 300")
	}
	s.service.Context().Control().SetBPM(float64(bpm))
	toolResult := new(mcp.CallToolResult)
	toolResult.Content = []mcp.Content{
		mcp.TextContent{
			Type: "text",
			Text: fmt.Sprintf("BPM set to %f", bpm),
		},
	}
	return toolResult, nil
}

func (s *MCPServer) HandlePlay(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	expression := request.GetString("expression", "")
	toolResult := new(mcp.CallToolResult)

	// do not write to stdout as the MCP server is using that
	captured := new(bytes.Buffer)
	notify.Console.StandardOut = captured

	response, err := s.service.CommandPlay("melrose-mcp", 0, expression)
	if err != nil {
		fmt.Fprintf(os.Stderr, "play failed: %v\n", err)
		toolResult.IsError = true
		toolResult.Content = []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: expression,
			},
			mcp.TextContent{
				Type: "text",
				Text: err.Error(),
			}}
		return toolResult, err
	}
	dur := max(time.Until(response.EndTime), 0) // not negative
	// wait until music has stopped playing or it is taking too long (2 min)
	if dur > 0 {
		time.Sleep(min(2*time.Minute, dur))
	}
	content := []mcp.Content{
		mcp.TextContent{
			Type: "text",
			Text: dur.String(),
		}}
	if p, ok := response.ExpressionResult.(core.Sequenceable); ok {
		ps := p.S()
		if len(ps.Notes) > 0 {
			content = append(content, mcp.TextContent{
				Type: "text",
				Text: ps.Storex(),
			})
		}
	} else {
		content = append(content, mcp.TextContent{
			Type: "text",
			Text: fmt.Sprintf("%v", response.ExpressionResult),
		})
	}
	toolResult.Content = content
	return toolResult, nil
}

func (s *MCPServer) HandleListDevices(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	list := s.service.ListDevices()
	toolResult := new(mcp.CallToolResult)
	for _, d := range list {
		kind := "input"
		if !d.IsInput {
			kind = "output"
		}
		toolResult.Content = append(toolResult.Content, mcp.TextContent{
			Type: "text",
			Text: fmt.Sprintf("%s is available as %s with device id %d", d.Name, kind, d.ID),
		})
	}
	return toolResult, nil
}
