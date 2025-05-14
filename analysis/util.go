package analysis

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"slices"
	"sunny-lsp/lsp"
)

type State struct {
    // map of file name to content
	Documents map[string]string
    Logger *log.Logger
}

type SymbolNode struct {
    Name string `json:"name"`
    ReachableScopes []int `json:"reachable_scopes"`
    Type string `json:"type"`
    Range lsp.Range `json:"range"`
}

type ASTNode struct {
    Name string `json:"name"`
    Scope int `json:"scope"`
    Range lsp.Range `json:"range"`

    LiteralType string `json:"literalType,omitempty"`
}

type CompilerContext struct {
	SymbolTable []SymbolNode `json:"symbols"`
	AST []ASTNode `json:"ast"`
	Diagnostics []lsp.Diagnostic `json:"diagnostics"`
}

func logCompilerOutput(output []byte, logger *log.Logger) {
    prettyPath := "/tmp/compiler_output.pretty.json"
    var prettyJSON bytes.Buffer
    if err := json.Indent(&prettyJSON, output, "", "  "); err != nil {
        logger.Printf("Couldn't format JSON: %v", err)
        return
    }
    if err := os.WriteFile(prettyPath, prettyJSON.Bytes(), 0644); err != nil {
        logger.Printf("Error saving pretty JSON: %v", err)
        return
    }
    logger.Printf("Saved compiler output to: %s", prettyPath)
}

func findSymbolDefinition(ctx *CompilerContext, pos lsp.Position) (*ASTNode, *SymbolNode) {
    // find the containing AST node to get the name and current scope
    containingNode := findContainingNode(ctx.AST, pos)
    if containingNode == nil {
        return nil, nil
    }

    name := containingNode.Name
    currentScope := containingNode.Scope

    // collect all symbols with the same name
    var candidates []SymbolNode
    for _, symbol := range ctx.SymbolTable {
        if symbol.Name == name {
            candidates = append(candidates, symbol)
        }
    }

    // find viable symbols where currentScope is in their ReachableScopes
    var viable []SymbolNode
    for _, symbol := range candidates {
        if containsScope(symbol.ReachableScopes, currentScope) {
            viable = append(viable, symbol)
        }
    }

    if len(viable) == 0 {
        return nil, nil
    }

    // find the declaration scope for each viable symbol
    type SymbolWithScope struct {
        Symbol SymbolNode
        Scope int
    }

    var symbolsWithScope []SymbolWithScope
    for _, symbol := range viable {
        declScope := findDeclarationScope(ctx, symbol.Range)
        if declScope != -1 {
            symbolsWithScope = append(symbolsWithScope, SymbolWithScope{symbol, declScope})
        }
    }

    if len(symbolsWithScope) == 0 {
        return nil, nil
    }

    // select the symbol with the highest declaration scope (innermost)
    maxScope := -1
    var best *SymbolNode
    for _, s := range symbolsWithScope {
        if s.Scope > maxScope {
            maxScope = s.Scope
            best = &s.Symbol
        }
    }

    return containingNode, best
}

// helper function to find the declaration scope of a symbol based on its Range
func findDeclarationScope(ctx *CompilerContext, r lsp.Range) int {
    for _, astNode := range ctx.AST {
        if rangesEqual(astNode.Range, r) {
            return astNode.Scope
        }
    }
    return -1
}

// helper function to check if two ranges are equal
func rangesEqual(a, b lsp.Range) bool {
    return a.Start.Line == b.Start.Line &&
        a.Start.Character == b.Start.Character &&
        a.End.Line == b.End.Line &&
        a.End.Character == b.End.Character
}

func containsScope(scopes []int, target int) bool {
    return slices.Contains(scopes, target)
}

func findContainingNode(nodes []ASTNode, pos lsp.Position) *ASTNode {
	for i := range nodes {
		node := &nodes[i]
		if positionInRange(pos, node.Range) {
            return node
		}
	}
	return nil
}

func positionInRange(pos lsp.Position, r lsp.Range) bool {
	if pos.Line < r.Start.Line || pos.Line > r.End.Line {
		return false
	}
	if pos.Line == r.Start.Line && pos.Character < r.Start.Character {
		return false
	}
	if pos.Line == r.End.Line && pos.Character > r.End.Character {
		return false
	}
	return true
}

func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}
