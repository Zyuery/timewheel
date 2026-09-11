package wheel

import (
	"container/list"
	"sync"
	"time"
)

type TimeWheel struct {
	sync.Once
	interval     time.Duration
	ticker       *time.Ticker
	stopCh       chan struct{}
	addTaskCh    chan *taskElement
	removeTaskCh chan string
	slots        []*list.List
	curSlot      int
	keyToETask   map[string]*list.Element
}

type taskElement struct {
	fn    func()
	pos   int //环形数组索引位置
	cycle int //curSlot还要扫几轮
	key   string
}

type Option func(*TimeWheel)

func WithInterval(interval time.Duration) Option {
	return func(tw *TimeWheel) {
		tw.interval = interval
	}
}

func WithSlotNum(slotNum int) Option {
	return func(tw *TimeWheel) {
		tw.slots = make([]*list.List, 0, slotNum)
	}
}

func NewTimeWheel(options ...Option) *TimeWheel {
	wheel := &TimeWheel{
		interval: time.Second,
		slots:    make([]*list.List, 0, 10),
	}
	for _, option := range options {
		option(wheel)
	}
	wheel.ticker = time.NewTicker(wheel.interval)
	wheel.stopCh = make(chan struct{})
	wheel.addTaskCh = make(chan *taskElement)
	wheel.removeTaskCh = make(chan string)
	wheel.keyToETask = make(map[string]*list.Element)
	for i := 0; i < len(wheel.slots); i++ {
		wheel.slots[i] = list.New()
	}
	go wheel.run()
	return wheel
}
