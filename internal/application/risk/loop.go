package risk

func (e *engine) run() {
	for {
		select {
		case <-e.ctx.Done():
			return
		case signal := <-e.signalCh:
			e.handleSignal(signal)
		}
	}
}
