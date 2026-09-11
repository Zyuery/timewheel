package wheel

import (
	"errors"
	"time"
)

var ErrStopped = errors.New("timewheel stopped")

func (w *TimeWheel) AddTask(key string, fn func(), executeAt time.Time) error {
	select {
	case <-w.stopCh:
		return ErrStopped
	default:
	}

	pos, cycle := w.getPosAndCycle(executeAt)
	task := &taskElement{
		key:   key,
		fn:    fn,
		pos:   pos,
		cycle: cycle,
	}

	select {
	case <-w.stopCh:
		return ErrStopped
	case w.addTaskCh <- task:
		return nil
	}
}

func (w *TimeWheel) RemoveTask(key string) error {
	select {
	case <-w.stopCh:
		return ErrStopped
	default:
	}

	select {
	case <-w.stopCh:
		return ErrStopped
	case w.removeTaskCh <- key:
		return nil
	}
}
