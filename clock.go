package main

import "time"

// Clock represents a device that tells the time. Useful for ensuring repeatability in tests.
type Clock interface {
	Now() time.Time
}

type RealClock int

func (m RealClock) Now() time.Time {
	return time.Now()
}

type FakeClock struct {
	t time.Time
}

func NewFakeClock(t time.Time) *FakeClock {
	return &FakeClock{t}
}

func (c FakeClock) Now() time.Time {
	return c.t
}

func (c *FakeClock) SetTime(t time.Time) {
	c.t = t
}

func (c *FakeClock) Advance(d time.Duration) {
	c.t = c.t.Add(d)
}
