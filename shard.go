package disgo

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slf4go/logger"
)

type Shard struct {
	session    *Session
	webSocket  *websocket.Conn
	sequence   int
	heartbeat  int
	stopListen chan bool
}

func connectShard(session *Session, shard int) (*Shard, error) {
	conn, _, err := websocket.DefaultDialer.Dial(session.wsUrl, http.Header{})
	if err != nil {
		return nil, err
	}

	s := &Shard{session: session, webSocket: conn, stopListen: make(chan bool)}

	helloFrame, err := s.readFrame()
	if err != nil {
		return nil, err
	} else if helloFrame.Op != opHello {
		return nil, errors.New("First frame sent from the Discord Gateway is not hello")
	}

	hello := helloPayload{}
	err = json.Unmarshal(helloFrame.Data, &hello)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Connected to Discord servers: %s", strings.Join(hello.Servers, ", "))
	logger.Debugf("Setting up a heartbeat interval of %d ms", hello.HeartbeatInterval)
	s.heartbeat = hello.HeartbeatInterval
	go s.listen()

	conn.WriteJSON(gatewayFrame{opIdentify, identifyPayload{
		Token:          session.token,
		Compress:       true,
		LargeThreshold: 250,
		Shard:          [2]int{shard, cap(session.shards)},
		Properties: propertiesPayload{
			OS:      runtime.GOOS,
			Browser: "DisGo",
			Device:  "DisGo",
		},
	}})

	return s, nil
}

func (s *Shard) disconnect() {
	if s.webSocket != nil {
		s.stopListen <- true
		s.webSocket = nil
	}
}

func (s *Shard) listen() {
	defer s.webSocket.Close()

	s.stopListen = make(chan bool)
	heartbeat := time.NewTicker(time.Duration(s.heartbeat) * time.Millisecond)
	reader := make(chan *receivedFrame)

	go s.readWebSocket(reader)

	var sentHeartBeat = false

listenLoop:
	for {
		select {
		case <-heartbeat.C:
			logger.Debug("Sending heartbeat")
			s.webSocket.WriteJSON(gatewayFrame{opHeartbeat, s.sequence})
		case <-s.stopListen:
			break listenLoop
		case message := <-reader:
			switch opCode := message.Op; opCode {
			case opHeartbeat:
				if !sentHeartBeat {
					logger.Debug("Sending heartbeat Ack")
					s.webSocket.WriteJSON(gatewayFrame{Op: opHeartbeatAck})
					sentHeartBeat = true
				} else {
					// TODO Disconnect and reconnect
				}
			case opHeartbeatAck:
				sentHeartBeat = false
			case opDispatch:
				s.session.dispatchEvent(message)
			default:
				logger.Errorf("Invalid opCode received: %d", opCode)
			}
		}
	}

	heartbeat.Stop()
}

func (s *Shard) readWebSocket(reader chan *receivedFrame) {
	for {
		frame, err := s.readFrame()
		if err != nil {
			s.stopListen <- true
			break
		}

		reader <- frame
	}
}

func (s *Shard) readFrame() (*receivedFrame, error) {
	msgType, msg, err := s.webSocket.ReadMessage()
	if err != nil {
		return nil, err
	}

	var reader io.Reader
	reader = bytes.NewBuffer(msg)

	if msgType == websocket.BinaryMessage {
		zReader, err := zlib.NewReader(reader)
		if err != nil {
			return nil, err
		}

		defer zReader.Close()

		reader = zReader
	}

	frame := receivedFrame{}
	json.NewDecoder(reader).Decode(&frame)

	logger.Debugf("Received frame: %+v", struct {
		Op        opCode
		Sequence  int
		EventName string
	}{frame.Op, frame.Sequence, frame.EventName})

	if frame.Sequence != 0 {
		s.sequence = frame.Sequence
	}
	return &frame, nil
}
