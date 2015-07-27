package nalog

type stopper struct {
	stop chan struct{}
	done chan struct{}
}

func newStopper() *stopper {
	return &stopper{
		make(chan struct{}),
		make(chan struct{}),
	}
}

func (s *stopper) Stop() {
	close(s.stop)
	<-s.done
}

func (s *stopper) Stopped() <-chan struct{} {
	return s.stop
}

func (s *stopper) Done() {
	close(s.done)
}
