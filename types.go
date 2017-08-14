package main

import "time"

type Project struct {
	ID   uint `gorm:"primary_key"`
	Name string
}

type User struct {
	ID        uint `gorm:"primary_key"`
	FirstName string
	LastName  string
	Email     string `gorm:"unique_index"`
}

// A Task represents a single task that time is tracked against cumulatively. A Task is comprised of many Frames, each of which represents
// an uninterrupted stretch of time that was allocated to the Task.
type Task struct {
	ID uint `gorm:"primary_key"`

	// The user who authored this Task and is working on it.
	UserID uint `sql:"not null"`

	// The category this task belongs to.
	CategoryID uint `sql:"not null"`

	// Human readable description of this task.
	Description string
}

// A Stretch represents an uninterrupted block of time dedicated to a particular Task.
type Stretch struct {
	ID uint `gorm:"primary_key"`

	// The Task this Stretch belongs to.
	TaskID uint `sql:"not null"`

	// The time this stretch began.
	Start time.Time

	// The time this stretch ended, or zero if this stretch is in progress.
	End time.Time
}

type Category struct {
	ID   uint `gorm:"primary_key"`
	Name string
}

