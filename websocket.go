package disgo

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slf4go/logger"
)

func (s *Session) Open() error {
	s.checkValid()
	logger.Trace("Open() called")

	conn, _, err := websocket.DefaultDialer.Dial(s.wsUrl, http.Header{})
	if err != nil {
		return err
	}

	s.webSocket = conn

	hello := helloPayload{}
	helloFrame := gatewayFrame{opHello, &hello}
	conn.ReadJSON(&helloFrame)
	logger.Debugf("Connected to Discord servers: %s", strings.Join(hello.Servers, ", "))
	logger.Debugf("Setting up a heartbeat interval of %d ms", hello.HeartbeatInterval)

	s.heartbeat = hello.HeartbeatInterval
	go s.listen()

	conn.WriteJSON(gatewayFrame{opIdentify, identifyPayload{
		Token:          s.token,
		Compress:       true,
		LargeThreshold: 250,
		Shard:          [2]int{0, 1},
		Properties: propertiesPayload{
			OS:      runtime.GOOS,
			Browser: "DisGo",
			Device:  "DisGo",
		},
	}})

	<-s.stopListen
	return nil
}

func (s *Session) Close() {
	s.checkValid()
	logger.Trace("Close() called")

	if s.webSocket != nil {
		s.stopListen <- true
		s.webSocket = nil
	}
}

func (s *Session) listen() {
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
				logger.Infof("Received event %s", message.EventName)
			default:
				logger.Errorf("Invalid opCode received: %d", opCode)
			}
		}
	}

	heartbeat.Stop()
}

func (s *Session) readWebSocket(reader chan *receivedFrame) {
	for {
		err := s.readMessage(reader)
		if err != nil {
			s.stopListen <- true
			break
		}
	}
}

func (s *Session) readMessage(frameDst chan *receivedFrame) error {
	msgType, msg, err := s.webSocket.ReadMessage()
	if err != nil {
		return err
	}

	var reader io.Reader
	reader = bytes.NewBuffer(msg)

	if msgType == websocket.BinaryMessage {
		zReader, err := zlib.NewReader(reader)
		if err != nil {
			return err
		}

		defer zReader.Close()

		reader = zReader
	}

	frame := receivedFrame{}
	json.NewDecoder(reader).Decode(&frame)
	frameDst <- &frame
	return nil
}
