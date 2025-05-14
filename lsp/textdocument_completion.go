package lsp

type CompletionRequest struct {
    Request
    Params CompletionParams `json:"params"`
}

type CompletionParams struct {
    TextDocumentPositionParam
}

type CompletionResponse struct {
    Response
    Result []CompletionItem `json:"result"`
}

type CompletionItem struct {
    Label string `json:"label"`
    Detail string `json:"detail"`
    Documentation string `json:"documentation"`
}
