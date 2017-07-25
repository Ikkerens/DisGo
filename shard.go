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
	readLock, writeLock sync.Mutex

	// Channels used to synchronize shutdown of this shards goroutines
	isShuttingDown    bool
	closeMainLoop     chan bool
	closeConfirmation chan bool
}

// Shard constructor, automatically connects with the given shardNum
func newShard(session *Session, shardNum int) (*shard, error) {
	s := &shard{
		session:           session,
		shard:             shardNum,
		closeMainLoop:     make(chan bool),
		closeConfirmation: make(chan bool),
	}

	// These will be unlocked once a connection is made.
	s.readLock.Lock()
	s.writeLock.Lock()

	if err := s.connect(); err != nil {
		return nil, err
	}

	return s, nil
}

// Builds a new connection with Discord, waits for the "hello" frame and then proceeds to identify itself to the Discord service
func (s *shard) connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(s.session.wsUrl, http.Header{})
	if err != nil {
		return err
	}

	s.webSocket = conn
	s.webSocket.SetCloseHandler(s.onClose)

	helloFrame, err := s.readFrame(true)
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

// identify takes care of the identification using the bot token on the discord server
func (s *shard) identify() error {
	if s.sessionID != "" {
		logger.Debugf("Resuming connection starting at sequence %d.", s.sequence)
		s.sendFrame(&gatewayFrame{opResume, resumePayload{
			Token:     s.session.token,
			SessionID: s.sessionID,
			Sequence:  s.sequence,
		}}, true)
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
		}}, true)
	}

	for {
		frame, err := s.readFrame(true)
		if err != nil {
			return err
		}

		switch frame.Op {
		case opDispatch:
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

				s.readLock.Unlock()
				s.writeLock.Unlock()

				return nil // Break out of the loop, we have what we want
			default:
				s.session.dispatchEvent(frame) // Nope, resume action, let's wait for more frames
			}
		case opInvalidSession:
			s.sessionID = "" // Invalidate session and retry
			return s.identify()
		default:
			return fmt.Errorf("Unexpected opCode received from Discord: %d", frame.Op)
		}
	}
}

// Main loop of the shard connection, also responsible for starting the read loop.
// This goroutine will pass around all the messages through the rest of the library.
func (s *shard) mainLoop() {
	logger.Debugf("Starting main loop for shard [%d/%d]", s.shard+1, cap(s.session.shards))
	defer logger.Debugf("Exiting main loop for shard [%d/%d]", s.shard+1, cap(s.session.shards))
	defer func() { s.closeConfirmation <- true }()

	heartbeat := time.NewTicker(time.Duration(s.heartbeat) * time.Millisecond)
	defer heartbeat.Stop()
	sentHeartBeat := false

	reader := make(chan *receivedFrame, 1)
	go s.readWebSocket(reader)

	for {
		select {
		case <-heartbeat.C:
			if !sentHeartBeat {
				go s.sendFrame(&gatewayFrame{opHeartbeat, s.sequence}, false)
				sentHeartBeat = true
			} else {
				go s.disconnect(websocket.ClosePolicyViolation, "Did not respond to previous heartbeat")
			}
		case <-s.closeMainLoop:
			return
		case frame := <-reader:
			switch frame.Op {
			case opHeartbeat:
				s.sendFrame(&gatewayFrame{Op: opHeartbeatAck}, false)
			case opHeartbeatAck:
				sentHeartBeat = false
			case opReconnect:
				go s.disconnect(websocket.CloseNormalClosure, "op Reconnect")
			case opInvalidSession:
				s.sessionID = ""
				s.identify()
			case opDispatch:
				s.session.dispatchEvent(frame)
			default:
				logger.Errorf("Unexpected opCode received: %d", frame.Op)
			}
		}
	}
}

// The secondary goroutine of each shard, responsible for reading frames and putting them in the channel.
func (s *shard) readWebSocket(reader chan *receivedFrame) {
	logger.Debugf("Starting read loop for shard [%d/%d]", s.shard+1, cap(s.session.shards))
	defer logger.Debugf("Exiting read loop for shard [%d/%d]", s.shard+1, cap(s.session.shards))
	defer func() { s.closeConfirmation <- true }()

	for {
		frame, err := s.readFrame(false)
		if err != nil {
			if !s.isShuttingDown {
				go s.onClose(websocket.CloseAbnormalClosure, err.Error())
			}
			return
		}

		reader <- frame
	}

}

// Reads 1 frame from the websocket.
func (s *shard) readFrame(isConnecting bool) (*receivedFrame, error) {
	if !isConnecting {
		s.readLock.Lock()
		defer s.readLock.Unlock()
	}

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
	err = json.NewDecoder(reader).Decode(&frame)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Received frame with opCode: %d", frame.Op)

	if frame.Sequence > s.sequence {
		logger.Tracef("Last sequence received set to %d.", frame.Sequence)
		s.sequence = frame.Sequence
	}
	return &frame, nil
}

// Sends 1 frame to the websocket
func (s *shard) sendFrame(frame *gatewayFrame, isConnecting bool) {
	if !isConnecting {
		s.writeLock.Lock()
		defer s.writeLock.Unlock()
	}

	logger.Debugf("Sending frame with opCode: %d", frame.Op)
	s.webSocket.WriteJSON(frame)
}

// Called when we have received a closing intention that we have not initiated (ws close message, recv error)
func (s *shard) onClose(code int, text string) error {
	logger.Warnf("Received Close Frame from Discord. Code: %d. Text: %s", code, text)

	s.isShuttingDown = true
	s.stopRoutines()
	s.readLock.Lock()
	s.writeLock.Lock()

	if err := s.webSocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, text)); err != nil {
		logger.Error("Could not confirm close message.")
		logger.ErrorE(err)
	}

	s.cleanupWebSocket()

	return nil
}

// Called by us to disconnect the websocket (shutdown, protocol errors, missed heartbeats)
func (s *shard) disconnect(code int, text string) {
	s.isShuttingDown = true
	s.writeLock.Lock()

	closeMessage := make(chan int)
	s.webSocket.SetCloseHandler(func(code int, _ string) error {
		closeMessage <- code
		return nil
	})

	if err := s.webSocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, text)); err != nil {
		logger.Error("Could not write close message to Discord, connection already dead?")
		logger.ErrorE(err)
	} else {
		select {
		case code := <-closeMessage:
			logger.Debugf("Discord connection closed with code %d", code)
		case <-time.After(1 * time.Second):
			logger.Warn("Discord did not reply to the close message, force-closing connection.")
		}
	}

	s.stopRoutines()
	s.readLock.Lock()
	s.cleanupWebSocket()
}

// Stops the two shard goroutines
func (s *shard) stopRoutines() {
	s.webSocket.SetReadDeadline(time.Now()) // Stop the read loop, in case it hasn't already due to an error
	<-s.closeConfirmation                   // Wait for the read loop
	s.closeMainLoop <- true                 // Stop the main loop
	<-s.closeConfirmation                   // Wait for the main loop
}

// Cleans up the current websocket and tries to reconnect it if needed
func (s *shard) cleanupWebSocket() {
	s.webSocket.Close()
	s.webSocket = nil

	for !s.session.isShuttingDown() {
		if err := s.connect(); err != nil {
			logger.Error("Could not reconnect to Discord.")
			logger.ErrorE(err)

			time.Sleep(500 * time.Millisecond)
		} else {
			break
		}
	}
}
