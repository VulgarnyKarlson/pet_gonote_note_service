package stream

func (s *StreamImpl) InClose() {
	if s.isClosed {
		return
	}
	close(s.inChan)
}

func (s *StreamImpl) InProxyClose() {
	if s.isClosed {
		return
	}
	close(s.inProxy)
}

func (s *StreamImpl) OutClose() {
	if s.isClosed {
		return
	}
	close(s.outChan)
}

func (s *StreamImpl) ErrClose() {
	if s.isClosed {
		return
	}
	close(s.errChan)
}
