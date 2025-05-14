package lsp

type Request struct {
    RPC string `json:"jsonrpc"` // always 2.0
    ID int `json:"id"`
    Method string `json:"method"`

    // Params
}

type Response struct {
    RPC string `json:"jsonrpc"` // always 2.0
    ID *int `json:"id,omitempty"`

    // Result, Error
}

type Notification struct {
    RPC string `json:"jsonrpc"`
    Method string `json:"method"`
}
