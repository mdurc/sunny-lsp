package analysis

import "sunny-lsp/lsp"

func (s *State) Completion(id int, uri string) lsp.CompletionResponse {
    keywordCompletions := []lsp.CompletionItem{
        {Label: "func", Detail: "Function declaration", Documentation: "Define a new function"},
        {Label: "mut", Detail: "Mutable declaration", Documentation: "Declare a mutable variable"},
        {Label: "if", Detail: "If statement", Documentation: "Conditional execution"},
        {Label: "else", Detail: "Else clause", Documentation: "Alternative branch for if statements"},
        {Label: "for", Detail: "For loop", Documentation: "Iteration construct"},
        {Label: "while", Detail: "While loop", Documentation: "Conditional loop"},
        {Label: "print", Detail: "Print function", Documentation: "Output to console"},
        {Label: "return", Detail: "Return statement", Documentation: "Exit function with value"},
        {Label: "returns", Detail: "Return type declaration", Documentation: "Specify function return type"},
        {Label: "break", Detail: "Break statement", Documentation: "Exit loop"},
        {Label: "continue", Detail: "Continue statement", Documentation: "Skip to next iteration"},
        {Label: "true", Detail: "Boolean true", Documentation: "Literal true value"},
        {Label: "false", Detail: "Boolean false", Documentation: "Literal false value"},
        {Label: "null", Detail: "Null value", Documentation: "Representation of no value"},
        {Label: "and", Detail: "Logical AND", Documentation: "Boolean AND operation"},
        {Label: "or", Detail: "Logical OR", Documentation: "Boolean OR operation"},
    }

    typeCompletions := []lsp.CompletionItem{
        {Label: "u8", Detail: "8-bit unsigned integer", Documentation: "Unsigned 8-bit integer type"},
        {Label: "u16", Detail: "16-bit unsigned integer", Documentation: "Unsigned 16-bit integer type"},
        {Label: "u32", Detail: "32-bit unsigned integer", Documentation: "Unsigned 32-bit integer type"},
        {Label: "u64", Detail: "64-bit unsigned integer", Documentation: "Unsigned 64-bit integer type"},
        {Label: "i8", Detail: "8-bit signed integer", Documentation: "Signed 8-bit integer type"},
        {Label: "i16", Detail: "16-bit signed integer", Documentation: "Signed 16-bit integer type"},
        {Label: "i32", Detail: "32-bit signed integer", Documentation: "Signed 32-bit integer type"},
        {Label: "i64", Detail: "64-bit signed integer", Documentation: "Signed 64-bit integer type"},
        {Label: "f64", Detail: "64-bit float", Documentation: "64-bit floating point number"},
        {Label: "bool", Detail: "Boolean type", Documentation: "True/false values"},
        {Label: "String", Detail: "String type", Documentation: "UTF-8 string type"},
        {Label: "u0", Detail: "Void type", Documentation: "Absence of type"},
    }

    return lsp.CompletionResponse{
        Response: lsp.Response{
            RPC: "2.0",
            ID:  &id,
        },
        Result: append(keywordCompletions, typeCompletions...),
    }
}
