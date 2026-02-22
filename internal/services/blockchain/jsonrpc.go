package blockchain

import (
	"context"
	"encoding/json"

	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"google.golang.org/grpc"
)

// JsonRPCClient wraps the gRPC client for JSON-RPC operations
type JsonRPCClient struct {
	client lindapb.JsonRpcClient
	conn   *grpc.ClientConn
}

func NewJsonRPCClient(conn *grpc.ClientConn) *JsonRPCClient {
	return &JsonRPCClient{
		client: lindapb.NewJsonRpcClient(conn),
		conn:   conn,
	}
}

// Forward forwards a JSON-RPC request to the node
func (c *JsonRPCClient) Forward(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	// Convert map to JSON
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Create gRPC request
	grpcReq := &lindapb.JsonRpcRequest{
		Body: reqJSON,
	}

	// Forward
	resp, err := c.client.Forward(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Call executes a specific JSON-RPC method
func (c *JsonRPCClient) Call(ctx context.Context, method string, params []interface{}, id int) (map[string]interface{}, error) {
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      id,
	}
	return c.Forward(ctx, req)
}

// EthBlockNumber returns the latest block number
func (c *JsonRPCClient) EthBlockNumber(ctx context.Context) (string, error) {
	resp, err := c.Call(ctx, "eth_blockNumber", []interface{}{}, 1)
	if err != nil {
		return "", err
	}
	if result, ok := resp["result"].(string); ok {
		return result, nil
	}
	return "", nil
}

// EthGetBalance returns the balance of an address
func (c *JsonRPCClient) EthGetBalance(ctx context.Context, address, block string) (string, error) {
	resp, err := c.Call(ctx, "eth_getBalance", []interface{}{address, block}, 1)
	if err != nil {
		return "", err
	}
	if result, ok := resp["result"].(string); ok {
		return result, nil
	}
	return "", nil
}

// EthGetTransactionByHash returns a transaction by hash
func (c *JsonRPCClient) EthGetTransactionByHash(ctx context.Context, hash string) (map[string]interface{}, error) {
	resp, err := c.Call(ctx, "eth_getTransactionByHash", []interface{}{hash}, 1)
	if err != nil {
		return nil, err
	}
	if result, ok := resp["result"].(map[string]interface{}); ok {
		return result, nil
	}
	return nil, nil
}

// EthGetTransactionReceipt returns a transaction receipt
func (c *JsonRPCClient) EthGetTransactionReceipt(ctx context.Context, hash string) (map[string]interface{}, error) {
	resp, err := c.Call(ctx, "eth_getTransactionReceipt", []interface{}{hash}, 1)
	if err != nil {
		return nil, err
	}
	if result, ok := resp["result"].(map[string]interface{}); ok {
		return result, nil
	}
	return nil, nil
}

// EthGetLogs returns logs matching filter
func (c *JsonRPCClient) EthGetLogs(ctx context.Context, filter map[string]interface{}) ([]interface{}, error) {
	resp, err := c.Call(ctx, "eth_getLogs", []interface{}{filter}, 1)
	if err != nil {
		return nil, err
	}
	if result, ok := resp["result"].([]interface{}); ok {
		return result, nil
	}
	return nil, nil
}

// BuildTransaction creates a transaction (LINDA-specific JSON-RPC method)
func (c *JsonRPCClient) BuildTransaction(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	return c.Call(ctx, "buildTransaction", []interface{}{params}, 1)
}