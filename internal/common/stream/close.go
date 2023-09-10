package stream

func (s *Impl) InClose() {
	if s.isClosed {
		return
	}
	close(s.inChan)
}

func (s *Impl) InProxyClose() {
	if s.isClosed {
		return
	}
	close(s.inProxy)
}

func (s *Impl) OutClose() {
	if s.isClosed {
		return
	}
	close(s.outChan)
}

func (s *Impl) ErrClose() {
	if s.isClosed {
		return
	}
	close(s.errChan)
}
