package deleter

import "sync"

// Счетчик воркера для отслеживания количества URL в работе
type Counter struct {
	num int
	sync.Mutex
}

// Inc инкремент счетчика
func (c *Counter) Inc(count int) {
	c.Lock()
	defer c.Unlock()

	c.num += count
}

// Dec декремент счетчика
func (c *Counter) Dec(count int) {
	c.Lock()
	defer c.Unlock()

	c.num -= count
}

// Value получение значения счетчика
func (c *Counter) Value() int {
	return c.num
}

// Clean обнуление счетчика
func (c *Counter) Clean() {
	c.Lock()
	defer c.Unlock()

	c.num = 0
}