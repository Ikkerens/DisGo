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
	Url    string `json:"Url"`
	Shards int    `json:"shards,omitempty"`
}

type gatewayFrame struct {
	Op   opCode      `json:"op"`
	Data interface{} `json:"d,omitempty"`
}

type receivedFrame struct {
	Op        opCode          `json:"op"`
	Data      json.RawMessage `json:"d,omitempty"`
	Sequence  uint64          `json:"s,omitempty"`
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

type resumePayload struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Sequence  uint64 `json:"seq"`
}

type statusPayload struct {
	Since  uint64 `json:"since"`
	Game   *Game  `json:"game"`
	Status Status `json:"status,string"`
	AFK    bool   `json:"afk"`
}
