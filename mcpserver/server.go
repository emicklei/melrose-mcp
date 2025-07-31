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
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPServer struct {
	service api.Service
}

func NewMCPServer(ctx core.Context) *MCPServer {
	return &MCPServer{service: api.NewService(ctx)}
}

type ChangeOutputDeviceParams struct {
	ID      int `json:"id" jsonschema:"the output device id"`
	Channel int `json:"channel" jsonschema:"the output channel, between 1 and 16"`
}

func (s *MCPServer) HandleChangeOutputDevice(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[ChangeOutputDeviceParams]) (*mcp.CallToolResultFor[any], error) {
	id := params.Arguments.ID
	channel := params.Arguments.Channel
	if id < 1 || id > 16 {
		return nil, fmt.Errorf("id must be a number between 1 and 16")
	}
	err := s.service.ChangeDefaultDeviceAndChannel(false, id, channel)
	toolResult := new(mcp.CallToolResult)
	if err != nil {
		toolResult.IsError = true
		toolResult.Content = []mcp.Content{
			&mcp.TextContent{
				Text: err.Error(),
			},
		}
	} else {
		toolResult.Content = []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Output device is set to %d", id),
			},
		}
	}
	return toolResult, nil
}

type BPMParams struct {
	BPM float64 `json:"bpm" jsonschema:"the beats per minute to set"`
}

func (p BPMParams) Get() float64 {
	if p.BPM <= 0 {
		return 120 // default
	}
	return p.BPM
}

func (s *MCPServer) HandleBPM(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[BPMParams]) (*mcp.CallToolResultFor[any], error) {
	bpm := params.Arguments.Get()
	if bpm < 1 || bpm > 300 {
		return nil, errors.New("parameter must be positive number between 1 and 300")
	}
	s.service.Context().Control().SetBPM(float64(bpm))
	toolResult := new(mcp.CallToolResult)
	toolResult.Content = []mcp.Content{
		&mcp.TextContent{
			Text: fmt.Sprintf("BPM set to %f", bpm),
		},
	}
	return toolResult, nil
}

type PlayParams struct {
	Expression string `json:"expression"  jsonschema:"the melrose expression to play"`
}

func (s *MCPServer) HandlePlay(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[PlayParams]) (*mcp.CallToolResultFor[any], error) {
	expression := params.Arguments.Expression
	toolResult := new(mcp.CallToolResultFor[any])

	// do not write to stdout as the MCP server is using that
	captured := new(bytes.Buffer)
	notify.Console.StandardOut = captured

	response, err := s.service.CommandPlay("melrose-mcp", 0, expression)
	if err != nil {
		fmt.Fprintf(os.Stderr, "play failed: %v\n", err)
		toolResult.IsError = true
		toolResult.Content = []mcp.Content{
			&mcp.TextContent{
				Text: expression,
			},
			&mcp.TextContent{
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
		&mcp.TextContent{
			Text: dur.String(),
		}}
	if p, ok := response.ExpressionResult.(core.Sequenceable); ok {
		ps := p.S()
		if len(ps.Notes) > 0 {
			content = append(content, &mcp.TextContent{
				Text: ps.Storex(),
			})
		}
	} else {
		content = append(content, &mcp.TextContent{
			Text: fmt.Sprintf("%v", response.ExpressionResult),
		})
	}
	toolResult.Content = content
	return toolResult, nil
}

type ListDevicesParams struct{}

func (s *MCPServer) HandleListDevices(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[ListDevicesParams]) (*mcp.CallToolResultFor[any], error) {
	list := s.service.ListDevices()
	toolResult := new(mcp.CallToolResult)
	for _, d := range list {
		kind := "input"
		if !d.IsInput {
			kind = "output"
		}
		toolResult.Content = append(toolResult.Content, &mcp.TextContent{
			Text: fmt.Sprintf("%s is available as %s with device id %d", d.Name, kind, d.ID),
		})
	}
	return toolResult, nil
}
