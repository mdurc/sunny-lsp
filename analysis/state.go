package analysis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sunny-lsp/lsp"
)

func NewState(logger *log.Logger) State {
    return State {
        Documents: map[string]string{},
        Logger: logger,
    }
}

const (
	CompilerPath = "/Users/mdurcan/personal/git_projects/tools/lang-dev/sunny-lang/compile.out"
)

func (s *State) RunCompiler(uri string) (*CompilerContext, error) {
	content, exists := s.Documents[uri]
	if !exists {
		return nil, fmt.Errorf("document not found: %s", uri)
	}

	tmpFile, err := os.CreateTemp("", "lsp-*.code")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		return nil, err
	}
	tmpFile.Close()

    cmd := exec.Command(CompilerPath, "--export-json", tmpFile.Name())
    var stderr bytes.Buffer
    cmd.Stderr = &stderr  // stderr separately
    output, err := cmd.Output()  // stdout only

	if err != nil {
		return nil, fmt.Errorf("compilation failed: %v\n%s", err, string(output))
	}

    //logCompilerOutput(output, s.Logger)

	var ctx CompilerContext
	if err := json.Unmarshal(output, &ctx); err != nil {
		return nil, fmt.Errorf("JSON parse error: %v", err)
	}

	return &ctx, nil
}

func (s *State) GetDiagnostics(uri string) []lsp.Diagnostic {
	ctx, err := s.RunCompiler(uri)
	if err != nil {
		return []lsp.Diagnostic{{
            Range:    LineRange(0,0,0),
            Severity: 1,
            Source:   "sunny-lsp:compiler",
			Message:  err.Error(),
		}}
	}
	return ctx.Diagnostics
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
    s.Documents[uri] = text
    return s.GetDiagnostics(uri)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
    s.Documents[uri] = text
    return s.GetDiagnostics(uri)
}

func (s *State) Hover(id int, uri string, pos lsp.Position) lsp.HoverResponse {
	ctx, err := s.RunCompiler(uri)
	if err != nil {
		return lsp.HoverResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &id,
			},
			Result: lsp.HoverResult{
				Contents: "Error: " + err.Error(),
			},
		}
	}

	// check diagnostics first
	for _, diag := range ctx.Diagnostics {
		if positionInRange(pos, diag.Range) {
			return lsp.HoverResponse{
				Response: lsp.Response{
					RPC: "2.0",
					ID:  &id,
				},
				Result: lsp.HoverResult{
					Contents: diag.Message,
				},
			}
		}
	}

	// find AST node and its associated symbol
	node, symbol := findSymbolDefinition(ctx, pos)
	if node == nil {
		return lsp.HoverResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &id,
			},
			Result: lsp.HoverResult{
				Contents: "No information found at position",
			},
		}
	}

	// build hover content
	var content strings.Builder
    content.WriteString(fmt.Sprintf("**%s**", node.Name))

	if symbol != nil {
		content.WriteString(fmt.Sprintf(" : *%s*", symbol.Type))
	} else if node.LiteralType != "" {
		content.WriteString(fmt.Sprintf(" : *%s*", node.LiteralType))
	}

	if symbol != nil && len(symbol.ReachableScopes) > 0 {
		content.WriteString(fmt.Sprintf("\nVisible in %d scopes", len(symbol.ReachableScopes)))
	}

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: content.String(),
		},
	}
}

// Jump to Definition gd
func (s *State) Definition(id int, uri string, pos lsp.Position) lsp.DefinitionResponse {
	ctx, err := s.RunCompiler(uri)
	if err != nil {
		return lsp.DefinitionResponse{
            Response: lsp.Response {
                RPC: "2.0",
                ID: &id,
            },
			Result: lsp.Location{
				URI:   uri,
                // do not move character at all
                Range: lsp.Range{
                    Start: pos,
                    End: pos,
                },
            },
        }
	}

	if _, symbol := findSymbolDefinition(ctx, pos); symbol != nil {
		return lsp.DefinitionResponse{
            Response: lsp.Response {
                RPC: "2.0",
                ID: &id,
            },
			Result: lsp.Location{
				URI:   uri,
				Range: symbol.Range,
			},
		}
	}

	return lsp.DefinitionResponse{
        Response: lsp.Response {
            RPC: "2.0",
            ID: &id,
        },
		Result: lsp.Location{
			URI:   uri,
			Range: lsp.Range{
                Start: pos,
                End: pos,
            },
		},
	}
}
