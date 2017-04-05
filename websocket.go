package disgo

import "github.com/slf4go/logger"

func (s *Session) Open() {
	logger.Trace("Open() called")
}

func (s *Session) Close() {
	logger.Trace("Close() called")
}
