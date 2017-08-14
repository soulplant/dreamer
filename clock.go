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

type FakeClock time.Time

func NewFakeClock(t time.Time) *FakeClock {
	c := FakeClock(t)
	return &c
}

func (c FakeClock) Now() time.Time {
	return time.Time(c)
}

func (c *FakeClock) SetTime(t time.Time) {
	*c = FakeClock(t)
}

func (c *FakeClock) Advance(d time.Duration) {
	*c = FakeClock(time.Time(*c).Add(d))
}
