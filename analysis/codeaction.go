package analysis

import "sunny-lsp/lsp"

func (s *State) TextCodeAction(id int, uri string, action_range lsp.Range) lsp.CodeActionResponse {
    actions := []lsp.CodeAction{}

    funcSnippet := `func FOO() {
}`
    actions = append(actions, lsp.CodeAction{
        Title: "Create function template",
        Edit: &lsp.WorkspaceEdit{
            Changes: map[string][]lsp.TextEdit{
                uri: {
                    {
                        Range:   action_range,
                        NewText: funcSnippet,
                    },
                },
            },
        },
    })

    forSnippet := `i32 n := 10;
for (mut i32 i := 0; i < n; i := i + 1) {
}`
    actions = append(actions, lsp.CodeAction{
        Title: "Create for loop",
        Edit: &lsp.WorkspaceEdit{
            Changes: map[string][]lsp.TextEdit{
                uri: {
                    {
                        Range:   action_range,
                        NewText: forSnippet,
                    },
                },
            },
        },
    })

    ifSnippet := `if (true) {
} else {
}`
    actions = append(actions, lsp.CodeAction{
        Title: "Create if-else block",
        Edit: &lsp.WorkspaceEdit{
            Changes: map[string][]lsp.TextEdit{
                uri: {
                    {
                        Range:   action_range,
                        NewText: ifSnippet,
                    },
                },
            },
        },
    })

    return lsp.CodeActionResponse{
        Response: lsp.Response{
            RPC: "2.0",
            ID:  &id,
        },
        Result: actions,
    }
}
