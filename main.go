package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"time"
)

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

func OpenTestDb() *gorm.DB {
	os.Remove("test.db")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&Project{}, &User{}, &Task{}, &Stretch{}, &Category{})
	return db
}

func main() {
	db := OpenTestDb()
	test(db)
}

func test(db *gorm.DB) {
	// Read
	var project Project
	if e := db.First(&project, 1000); e.Error != nil {
		fmt.Println("Couldn't find 1000")
	}
	if e := db.First(&project, 1); e.Error != nil {
		fmt.Println("Couldn't find 1")
	}
	db.First(&project, "Name = ?", "Dreamer")

	// Delete - delete project
	// db.Delete(&project)
}
