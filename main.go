package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"sunny-lsp/analysis"
	"sunny-lsp/lsp"
	"sunny-lsp/rpc"
)

func main() {
    logger := getLogger("/Users/mdurcan/personal/git_projects/tools/lang-dev/sunny-lsp/log.txt")
    logger.Println("Logger Started")

    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(rpc.Split)

    state := analysis.NewState(logger)
    writer := os.Stdout

    for scanner.Scan() {
        msg := scanner.Bytes()
        method, contents, err := rpc.DecodeMessage(msg)
        if err != nil {
            logger.Printf("Got an error: %s", err)
            continue
        }

        handleMessage(logger, writer, state, method, contents)
    }
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, contents []byte) {
    logger.Printf("Received msg with method: %s", method)

    switch method {
    case "initialize":
        var request lsp.InitializeRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("We could not parse initialize request: %s", err)
        }

        logger.Printf("Connected to: %s %s",
            request.Params.ClientInfo.Name,
            request.Params.ClientInfo.Version)

        msg := lsp.NewInitializeResponse(request.ID)
        writeResponse(writer, msg)
    case "textDocument/didOpen":
        var request lsp.TextDocumentDidOpenNotification
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/didOpen: %s", err)
            return;
        }

        uri := request.Params.TextDocument.URI
        text := request.Params.TextDocument.Text
        logger.Printf("Opened: %s", uri)
        diagnostics := state.OpenDocument(uri, text)
        writeResponse(writer, lsp.PublishDiagnosticNotification {
            Notification: lsp.Notification {
                RPC: "2.0",
                Method: "textDocument/publishDiagnostics",
            },
            Params: lsp.PublishDiagnosticParams {
                URI: uri,
                Diagnostics: diagnostics,
            },
        })
    case "textDocument/didChange":
        var request lsp.TextDocumentDidChangeNotification
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/didChange: %s", err)
            return
        }

        uri := request.Params.TextDocument.URI
        changes := request.Params.ContentChanges
        logger.Printf("Changed: %s", uri)
        for _, change := range changes {
            diagnostics := state.UpdateDocument(uri, change.Text)
            writeResponse(writer, lsp.PublishDiagnosticNotification {
                Notification: lsp.Notification {
                    RPC: "2.0",
                    Method: "textDocument/publishDiagnostics",
                },
                Params: lsp.PublishDiagnosticParams {
                    URI: uri,
                    Diagnostics: diagnostics,
                },
            })
        }
    case "textDocument/hover":
        var request lsp.HoverRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/hover: %s", err)
            return
        }

        // create response
        uri := request.Params.TextDocument.URI
        pos := request.Params.Position
        response := state.Hover(request.ID, uri, pos)

        // send it to LSP
        writeResponse(writer, response)
    case "textDocument/definition":
        var request lsp.DefinitionRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/definition: %s", err)
            return
        }

        // create response
        uri := request.Params.TextDocument.URI
        pos := request.Params.Position
        response := state.Definition(request.ID, uri, pos)

        // send it to LSP
        writeResponse(writer, response)
    case "textDocument/codeAction":
        var request lsp.CodeActionRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/codeAction: %s", err)
            return
        }

        uri := request.Params.TextDocument.URI
        action_range := request.Params.Range
        response := state.TextCodeAction(request.ID, uri, action_range)

        writeResponse(writer, response)
    case "textDocument/completion":
        var request lsp.CompletionRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/completion: %s", err)
            return
        }

        uri := request.Params.TextDocument.URI

        // should really also be passing position here
        response := state.Completion(request.ID, uri)

        writeResponse(writer, response)
    }
}

func writeResponse(writer io.Writer, msg any) {
    reply := rpc.EncodeMessage(msg)
    writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
    logfile, err := os.OpenFile(filename, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0666)
    if err != nil {
        panic("Invalid log file")
    }
    return log.New(logfile, "[sunny-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
