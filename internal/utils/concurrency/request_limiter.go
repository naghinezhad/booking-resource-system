package concurrency

type RequestLimiter struct {
	sem chan struct{}
}

func NewRequestLimiter(max int) *RequestLimiter {
	return &RequestLimiter{
		sem: make(chan struct{}, max),
	}
}

func (l *RequestLimiter) Acquire() bool {
	select {
	case l.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

func (l *RequestLimiter) Release() {
	<-l.sem
}
