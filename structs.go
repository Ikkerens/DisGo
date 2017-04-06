package disgo

import "encoding/json"

type opCode int

const (
	opDispatch opCode = iota
	opHeartbeat
	opIdentify
	opStatusUpdate
	opVoiceStateUpdate
	opVoiceServerPing
	opResume
	opReconnect
	opRequestGuildMembers
	opInvalidSession
	opHello
	opHeartbeatAck
)

type gatewayGetResponse struct {
	Url    string `json:"url"`
	Shards int    `json:"shards,omitempty"`
}

type gatewayFrame struct {
	Op   opCode      `json:"op"`
	Data interface{} `json:"d,omitempty"`
}

type receivedFrame struct {
	Op        opCode          `json:"op"`
	Data      json.RawMessage `json:"d,omitempty"`
	Sequence  int             `json:"s,omitempty"`
	EventName string          `json:"t,omitempty"`
}

type helloPayload struct {
	HeartbeatInterval int      `json:"heartbeat_interval"`
	Servers           []string `json:"_trace"`
}

type identifyPayload struct {
	Token          string            `json:"token"`
	Properties     propertiesPayload `json:"properties"`
	Compress       bool              `json:"compress"`
	LargeThreshold int               `json:"large_threshold"`
	Shard          [2]int            `json:"shard"`
}

type propertiesPayload struct {
	OS              string `json:"$os"`
	Browser         string `json:"$browser"`
	Device          string `json:"$device"`
	Referrer        string `json:"$referrer"`
	ReferringDomain string `json:"$referring_domain"`
}
