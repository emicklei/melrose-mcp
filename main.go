package main

import (
	"context"
	"fmt"
	"os"

	"github.com/emicklei/melrose-mcp/mcpserver"
	"github.com/emicklei/melrose/notify"
	"github.com/emicklei/melrose/system"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	_ "embed"
)

var BuildTag = "dev"

//go:embed resources/melrose_note_syntax.txt
var noteSyntaxContent string

//go:embed resources/melrose_llm_context.txt
var playLLMContext string

func main() {
	notify.SetANSIColorsEnabled(false) // error messages cannot be colored

	ctx, err := system.Setup(BuildTag)
	if err != nil {
		notify.Errorf("setup failed: %v", err)
		os.Exit(1)
	}

	ioServer := server.NewMCPServer(
		"melrose",
		"v0.56.0",
	)
	playServer := mcpserver.NewMCPServer(ctx)

	// Add resource for syntax
	// syntax := mcp.NewResource("file://melrose/note/syntax", "melrose note syntax", mcp.WithMIMEType("text/plain"))
	// ioServer.AddResource(syntax, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// 	return []mcp.ResourceContents{
	// 		mcp.TextResourceContents{
	// 			URI:      "file://melrose/note/syntax",
	// 			MIMEType: "text/plain",
	// 			Text:     noteSyntaxContent,
	// 		},
	// 	}, nil
	// })

	// Add resource for context
	syntax := mcp.NewResource("docs://melrose_play", "melrose expressions llm system context", mcp.WithMIMEType("text/plain"))
	ioServer.AddResource(syntax, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://melrose_play",
				MIMEType: "text/plain",
				Text:     playLLMContext,
			},
		}, nil
	})

	// Add play tool
	tool1 := mcp.NewTool("melrose_play",
		mcp.WithDescription(`Melrōse is a language to create music by programming expressions.
		 The language uses musical primitives (note, sequence, chord) and many functions (map, group, transpose).
		 See docs://melrose_play for more information.`),
		mcp.WithString("expression",
			mcp.Required(),
			mcp.Description("functional expression using the syntax rules docs://melrose_play"),
		),
	)
	ioServer.AddTool(tool1, playServer.HandlePlay)

	// Add bpm tool
	tool2 := mcp.NewTool("melrose_bpm",
		mcp.WithDescription(`Changes the beats per minutes setting. Default is 120.`),
		mcp.WithString("bpm",
			mcp.Required(),
			mcp.Description("number representing beats per minute, must be between 1 and 300"),
		),
	)
	ioServer.AddTool(tool2, playServer.HandleBPM)

	// Add device listing
	tool3 := mcp.NewTool("melrose_devices",
		mcp.WithDescription(`List all available input and output MIDI devices.`),
	)
	ioServer.AddTool(tool3, playServer.HandleListDevices)

	// Add device selector
	tool4 := mcp.NewTool("melrose_change_output_device",
		mcp.WithDescription(`Change the output device to the one specified by the device name.`),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("id me of the output device")),
		mcp.WithNumber("channel",
			mcp.Description("default channel number for this device. Must be between 1 and 16")))
	ioServer.AddTool(tool4, playServer.HandleChangeOutputDevice)

	// Add chord prompt
	chordHander := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		note := request.Params.Arguments["ground"]
		if note == "" {
			note = "C"
		}
		fraction := request.Params.Arguments["fraction"]
		if fraction == "" {
			fraction = "4"
		}
		octave := request.Params.Arguments["octave"]
		if octave == "" {
			octave = "4"
		}
		return mcp.NewGetPromptResult(
			"playing a chord",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent(fmt.Sprintf("chord('%s%s%s')", fraction, note, octave)),
				),
			},
		), nil
	}
	ioServer.AddPrompt(mcp.NewPrompt("play-chord",
		mcp.WithPromptDescription("play the notes of a chord")), chordHander)

	// Start the stdio server
	if err := server.ServeStdio(ioServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
