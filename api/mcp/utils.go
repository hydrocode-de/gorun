package mcp

type PongResult struct {
	Result string `json:"result"`
}

type PongResponse struct {
	McpResponse
	Result PongResult `json:"result"`
}

func Pong(id int) PongResponse {
	response := PongResponse{
		McpResponse: McpResponse{
			Jsonrpc: JSONRPC_VERSION,
			Id:      id,
		},
		Result: PongResult{Result: "pong"},
	}

	return response
}
