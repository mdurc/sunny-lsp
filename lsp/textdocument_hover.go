package lsp

type HoverRequest struct {
    Request
    Params HoverParams `json:"params"`
}

type HoverParams struct {
    TextDocumentPositionParam
}

type HoverResponse struct {
    Response
    Result HoverResult `json:"result"`
}

type HoverResult struct {
    Contents string `json:"contents"`
}
