package disgo

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slf4go/logger"
)

type shard struct {
	// Parent session
	session *Session

	// Websocket information used for identification and reconnection
	webSocket *websocket.Conn
	shard     int
	sessionID string
	sequence  int
	heartbeat int

	// Mutex locks, reconnect to make sure there is only 1 process reconnecting and concurrent read/write accesses on the socket
	reconnectLock       sync.Mutex
	readLock, writeLock sync.Mutex
	isReconnecting      bool

	// Channels to pass around messages
	closeMessage chan int
	stopListen   chan bool
	stopRead     chan bool
}

func connectShard(session *Session, shardNum int) (*shard, error) {
	conn, _, err := websocket.DefaultDialer.Dial(session.wsUrl, http.Header{})
	if err != nil {
		return nil, err
	}

	s := &shard{
		session:      session,
		webSocket:    conn,
		shard:        shardNum,
		closeMessage: make(chan int, 1),
		stopListen:   make(chan bool),
		stopRead:     make(chan bool),
	}

	if err = s.handshake(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *shard) handshake() error {
	s.webSocket.SetCloseHandler(s.onClose)

	helloFrame, err := s.readFrame()
	if err != nil {
		return err
	} else if helloFrame.Op != opHello {
		return errors.New("First frame sent from the Discord Gateway is not hello")
	}

	hello := helloPayload{}
	if err = json.Unmarshal(helloFrame.Data, &hello); err != nil {
		return err
	}

	logger.Debugf("Connected to Discord servers: %s", strings.Join(hello.Servers, ", "))
	logger.Debugf("Setting up a heartbeat interval of %d ms", hello.HeartbeatInterval)
	s.heartbeat = hello.HeartbeatInterval

	if err = s.identify(); err != nil {
		return err
	}

	go s.mainLoop()
	return nil
}

func (s *shard) identify() error {
	if s.sessionID != "" {
		logger.Debugf("Resuming connection starting at sequence %d.", s.sequence)
		s.sendFrame(&gatewayFrame{opResume, resumePayload{
			Token:     s.session.token,
			SessionID: s.sessionID,
			Sequence:  s.sequence,
		}})
	} else {
		logger.Debugf("Identifying to websocket")
		s.sendFrame(&gatewayFrame{opIdentify, identifyPayload{
			Token:          s.session.token,
			Compress:       true,
			LargeThreshold: 250,
			Shard:          [2]int{s.shard, cap(s.session.shards)},
			Properties: propertiesPayload{
				OS:      runtime.GOOS,
				Browser: "DisGo",
				Device:  "DisGo",
			},
		}})
	}

	for {
		frame, err := s.readFrame()
		if err != nil {
			return err
		}

		if frame.Op == opDispatch {
			switch frame.EventName {
			case "READY":
				ready := ReadyEvent{}
				if err := json.Unmarshal(frame.Data, &ready); err != nil {
					return err
				}

				s.sessionID = ready.SessionID

				fallthrough
			case "RESUMED":
				s.session.dispatchEvent(frame)
				return nil // Break out of the loop, we have what we want
			default:
				s.session.dispatchEvent(frame) // Nope, resume action, let's wait for more frames
			}
		} else if frame.Op == opInvalidSession {
			s.sessionID = "" // Invalidate session and retry
			s.identify()
			return nil
		} else {
			return fmt.Errorf("Unexpected opCode received from Discord: %d", frame.Op)
		}
	}
}

func (s *shard) mainLoop() {
	logger.Debugf("Starting main loop for shard [%d/%d]", s.shard+1, cap(s.session.shards))
	defer logger.Debugf("Exiting main loop for shard [%d/%d]", s.shard+1, cap(s.session.shards))

	heartbeat := time.NewTicker(time.Duration(s.heartbeat) * time.Millisecond)
	defer heartbeat.Stop()
	sentHeartBeat := false

	reader := make(chan *receivedFrame, 1)
	go s.readWebSocket(reader)

	for {
		select {
		case <-heartbeat.C:
			if !sentHeartBeat {
				go s.sendFrame(&gatewayFrame{opHeartbeat, s.sequence})
				sentHeartBeat = true
			} else {
				go s.disconnect(websocket.CloseAbnormalClosure, "Did not respond to previous heartbeat")
			}
		case <-s.stopListen:
			return
		case frame := <-reader:
			switch frame.Op {
			case opHeartbeat:
				s.sendFrame(&gatewayFrame{Op: opHeartbeatAck})
			case opHeartbeatAck:
				sentHeartBeat = false
			case opReconnect:
				s.disconnect(websocket.CloseNormalClosure, "op Reconnect")
			case opDispatch:
				s.session.dispatchEvent(frame)
			default:
				logger.Errorf("Unexpected opCode received: %d", frame.Op)
			}
		}
	}
}

func (s *shard) sendFrame(frame *gatewayFrame) {
	// Resume will only be sent within the reconnect lock, we don't relock to prevent a deadlock
	if frame.Op != opResume {
		s.reconnectLock.Lock()
		defer s.reconnectLock.Unlock()
	}

	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	logger.Debugf("Sending frame with opCode: %d", frame.Op)
	s.webSocket.WriteJSON(frame)
}

func (s *shard) readWebSocket(reader chan *receivedFrame) {
	logger.Debugf("Starting read loop for shard [%d/%d]", s.shard+1, cap(s.session.shards))
	defer logger.Debugf("Exiting read loop for shard [%d/%d]", s.shard+1, cap(s.session.shards))

	for {
		select {
		case <-s.stopRead:
			return
		default:
			frame, err := s.readFrame()
			if err != nil {
				if !s.session.isShuttingDown() {
					if closeErr, isClose := err.(*websocket.CloseError); isClose {
						s.closeMessage <- closeErr.Code
					}
					logger.ErrorE(err)
					go s.disconnect(websocket.CloseAbnormalClosure, err.Error())
				}
				<-s.stopRead // Wait for the connection to be closed
				return
			}

			reader <- frame
		}
	}

}

func (s *shard) readFrame() (*receivedFrame, error) {
	s.readLock.Lock()
	defer s.readLock.Unlock()

	logger.Tracef("Shard.readFrame() called")
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

	frame := receivedFrame{Sequence: -1}
	json.NewDecoder(reader).Decode(&frame)

	logger.Debugf("Received frame with opCode: %d", frame.Op)

	if frame.Sequence != -1 {
		logger.Tracef("Last sequence received set to %d.", frame.Sequence)
		s.sequence = frame.Sequence
	}
	return &frame, nil
}

func (s *shard) startReconnect() {
	if !s.isReconnecting {
		s.reconnectLock.Lock()
		s.isReconnecting = true
		go s.reconnect()
	}
}

func (s *shard) reconnect() {
	if !s.session.isShuttingDown() {
		logger.Noticef("Reconnecting shard [%d/%d]", s.shard+1, cap(s.session.shards))
		conn, _, err := websocket.DefaultDialer.Dial(s.session.wsUrl, http.Header{})

		if err == nil {
			s.webSocket = conn
			err = s.handshake()

			if err != nil {
				conn.Close()
			}
		}

		if err != nil {
			logger.ErrorE(err)
			time.Sleep(1 * time.Second)
			go s.reconnect()
		} else {
			s.isReconnecting = false
			s.reconnectLock.Unlock()
		}
	}
}

func (s *shard) onClose(code int, text string) error {
	logger.Debugf("Received Close Frame from Discord. Code: %d. Text: %s", code, text)
	s.closeMessage <- code
	s.startReconnect()
	return nil
}

func (s *shard) disconnect(code int, text string) {
	logger.Tracef("Shard.disconnect() called")
	s.stopListen <- true

	s.writeLock.Lock()
	err := s.webSocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, text))
	s.writeLock.Unlock()
	if err != nil {
		logger.ErrorE(err)
	}

	select {
	case code := <-s.closeMessage:
		logger.Debugf("Discord connection closed with code %d", code)
	case <-time.After(1 * time.Second):
		logger.Warn("Discord did not reply to the close message, force-closing connection.")
	}

	s.webSocket.Close()
	s.stopRead <- true

	s.startReconnect()
}
