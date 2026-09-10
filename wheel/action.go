package wheel

func (w *TimeWheel) run() {
	for {
		select {
		case <-w.stopCh:
			return
		case <-w.ticker.C:
			w.handle()
		case task := <-w.addTaskCh:
			w.add(task)
		case key := <-w.removeTaskCh:
			w.remove(key)
		default:
		}
	}
}
func (w *TimeWheel) Stop() {
	w.stopCh <- struct{}{}
}

func (w *TimeWheel) handle() {
}

func (w *TimeWheel) add(task *taskElement) {
}

func (w *TimeWheel) remove(key string) {
}
