package event

import (
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type EventParser struct {
	// Cache of parsed event ABIs
	eventSignatures map[string]abi.Event
}

func NewEventParser() *EventParser {
	return &EventParser{
		eventSignatures: make(map[string]abi.Event),
	}
}

// ParseEvent parses event topics and data into named parameters
func (p *EventParser) ParseEvent(topics [][]byte, data []byte) (string, map[string]interface{}, map[string]string) {
	if len(topics) == 0 {
		return "", nil, nil
	}

	// First topic is event signature
	eventSig := topics[0]
	eventName := p.getEventName(eventSig)

	// Parse indexed parameters (from topics[1:])
	indexedParams := make(map[string]interface{})
	indexedTypes := make(map[string]string)

	// Parse non-indexed parameters (from data)
	nonIndexedParams := make(map[string]interface{})
	nonIndexedTypes := make(map[string]string)

	// This is a simplified parser - in production, you'd need the full ABI
	// For now, we'll handle common event types

	switch eventName {
	case "Transfer":
		if len(topics) >= 3 {
			// Transfer(address indexed from, address indexed to, uint256 value)
			indexedParams["from"] = common.BytesToAddress(topics[1]).Hex()
			indexedTypes["from"] = "address"
			
			indexedParams["to"] = common.BytesToAddress(topics[2]).Hex()
			indexedTypes["to"] = "address"
			
			if len(data) >= 32 {
				value := new(big.Int).SetBytes(data[:32])
				nonIndexedParams["value"] = value.String()
				nonIndexedTypes["value"] = "uint256"
			}
		}

	case "Approval":
		if len(topics) >= 3 {
			// Approval(address indexed owner, address indexed spender, uint256 value)
			indexedParams["owner"] = common.BytesToAddress(topics[1]).Hex()
			indexedTypes["owner"] = "address"
			
			indexedParams["spender"] = common.BytesToAddress(topics[2]).Hex()
			indexedTypes["spender"] = "address"
			
			if len(data) >= 32 {
				value := new(big.Int).SetBytes(data[:32])
				nonIndexedParams["value"] = value.String()
				nonIndexedTypes["value"] = "uint256"
			}
		}

	default:
		// Generic parsing
		for i, topic := range topics[1:] {
			paramName := "topic" + string(rune(i))
			indexedParams[paramName] = hex.EncodeToString(topic)
			indexedTypes[paramName] = "bytes32"
		}
		
		if len(data) > 0 {
			nonIndexedParams["data"] = hex.EncodeToString(data)
			nonIndexedTypes["data"] = "bytes"
		}
	}

	// Merge indexed and non-indexed parameters
	result := make(map[string]interface{})
	resultTypes := make(map[string]string)

	for k, v := range indexedParams {
		result[k] = v
		resultTypes[k] = indexedTypes[k]
	}
	for k, v := range nonIndexedParams {
		result[k] = v
		resultTypes[k] = nonIndexedTypes[k]
	}

	return eventName, result, resultTypes
}

// getEventName attempts to get event name from signature hash
func (p *EventParser) getEventName(signatureHash []byte) string {
	// Common event signatures
	commonEvents := map[string]string{
		crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")).Hex(): "Transfer",
		crypto.Keccak256Hash([]byte("Approval(address,address,uint256)")).Hex():  "Approval",
	}

	hashHex := "0x" + hex.EncodeToString(signatureHash)
	if name, ok := commonEvents[hashHex]; ok {
		return name
	}

	return "UnknownEvent"
}

// DecodeEventData decodes event data using ABI
func (p *EventParser) DecodeEventData(eventABI abi.Event, data []byte) (map[string]interface{}, error) {
	// This would use the actual ABI to decode parameters
	// For now, return empty map
	return make(map[string]interface{}), nil
}

// EncodeEventTopics encodes indexed parameters into topics
func (p *EventParser) EncodeEventTopics(eventABI abi.Event, params map[string]interface{}) ([][]byte, error) {
	// This would encode indexed parameters according to ABI
	// For now, return empty topics
	return [][]byte{}, nil
}