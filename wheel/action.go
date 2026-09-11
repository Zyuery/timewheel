package wheel

import (
	"container/list"
	"time"
)

func (w *TimeWheel) run() {
	defer func() {
		if r := recover(); r != nil {
			w.ticker.Stop()
		}
	}()
	for {
		select {
		case <-w.stopCh:
			return
		case <-w.ticker.C:
			w.tick()
		case task := <-w.addTaskCh:
			w.addTask(task)
		case key := <-w.removeTaskCh:
			w.removeTask(key)
		default:
		}
	}
}
func (w *TimeWheel) Stop() {
	w.Do(func() {
		w.ticker.Stop()
		close(w.stopCh)
	})
}

func (w *TimeWheel) tick() {
	w.handle()
}

func (w *TimeWheel) handle() {
	elist := w.slots[w.curSlot]
	defer w.cursorIncr()
	w.execute(elist)
}

func (w *TimeWheel) execute(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		task, _ := e.Value.(*taskElement)
		task.cycle--
		if task.cycle > 0 {
			continue
		}
		go func() {
			defer func() {
				if r := recover(); r != nil {
					w.removeTask(task.key)
				}
			}()
			task.fn()
		}()
	}
}

func (w *TimeWheel) cursorIncr() {
	w.curSlot = (w.curSlot + 1) % len(w.slots)
}

func (w *TimeWheel) RemoveTask(key string) {
	w.removeTaskCh <- key
}

func (w *TimeWheel) removeTask(key string) {
	element, ok := w.keyToETask[key]
	if !ok {
		return
	}
	task, _ := element.Value.(*taskElement)
	delete(w.keyToETask, key)
	w.slots[task.pos].Remove(element)
}

func (w *TimeWheel) AddTask(key string, task func(), executeAt time.Time) {
	pos, cycle := w.getPosAndCycle(executeAt)
	w.addTaskCh <- &taskElement{
		key:   key,
		fn:    task,
		pos:   pos,
		cycle: cycle,
	}
}

func (w *TimeWheel) getPosAndCycle(executeAt time.Time) (pos, cycle int) {
	delay := executeAt.Sub(w.now())
	cycle = int(delay / (time.Duration(len(w.slots)) * w.interval))
	pos = (w.curSlot + int(delay/w.interval)) % len(w.slots)
	return
}

func (w *TimeWheel) now() time.Time {
	return time.Now()
}

func (w *TimeWheel) addTask(task *taskElement) {
	if task == nil {
		return
	}
	if _, ok := w.keyToETask[task.key]; ok {
		w.removeTask(task.key)
	}
	w.keyToETask[task.key] = w.slots[task.pos].PushBack(task)
}
