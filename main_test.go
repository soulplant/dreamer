package main

import (
	"log"
	"testing"
	"time"
)

func mkTime(str string) time.Time {
	res, err := time.Parse(time.Kitchen, str)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func Test_NewFakeClock(t *testing.T) {
	fc1 := NewFakeClock(mkTime("3:00PM"))
	fc2 := NewFakeClock(mkTime("4:00PM"))
	if !fc1.Now().Before(fc2.Now()) {
		t.Error("time not consistent")
	}
}
