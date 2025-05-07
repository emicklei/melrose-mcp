## installation

This mcp server requires the [Go SDK](https://go.dev) for installation.
To produce music, you need to setup a MIDI receiver such as a DAW or a hardware synthesizer.

    go install github.com/emicklei/melrose-mcp@latest

## session

>   change the tempo in melrose to a lower value

`melrose_bmp` request:
```javascript
{
  `bpm`: `80`
}
```

`llm` response:
```
I've reduced the tempo to 80 beats per minute from the default of 120 BPM. This creates a more relaxed and slower pace for any musical expressions we play.
```

>   play a tune from Claude Debussy

`llm` response:
```
I'll create a small segment inspired by Claude Debussy's impressionist style using the Melrose language at our current tempo of 80 BPM.
```

`melrose_play` request:
```javascript
{
  `expression`: `sequence('(F4 A4 D5) 4(F4 A4 D5) (G4 B4 E5) 4(G4 B4 E5) (A4 C5 F5) 4(A4 C5 F5) (G4 B4 E5) 4(G4 B4 E5) (F4 A4 D5) 8= 8(E4 G4 C5) 4(F4 A4 D5)')`
}
```
`melrose_play` response:
```javascript
12.749963417s
```

## configuration in your client

Below an example to configure it in Claude Desktop.

    {
        "mcpServers": {
            "melrōse": {
              "command": "<YOUR DIRECTORY TO THE BINARY>/melrose-mcp",
              "args": []
            }
        }
    }

## example prompts

- play the music sequence  C D E
- play the first bar of fur elise
- play the famous tune from Mike Oldfield called Tubular Bells
- play the notes of Frere Jacques
- play the C# chord

To help your LLM understand the Melrōse language better, you can use this [context](resources/melrose_llm_context.txt).

### Contributions

Fixes, suggestions, documentation improvements are all welcome.
Fork this project and submit small Pull requests. 
Discuss larger ones in the Issues list.
You can also sponsor Melrōse via [Github Sponsors](https://github.com/sponsors/emicklei).

Software is licensed under [MIT](LICENSE).
&copy; 2025 [ernestmicklei.com](http://ernestmicklei.com)