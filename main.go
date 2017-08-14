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
	UserID uint

	// The category this task belongs to.
	CategoryID uint

	// Human readable description of this task.
	Description string
}

// A Stretch represents an uninterrupted block of time dedicated to a particular Task.
type Stretch struct {
	ID uint `gorm:"primary_key"`

	// The Task this Stretch belongs to.
	TaskID uint

	// The time this stretch began.
	Start time.Time

	// The time this stretch ended, or zero if this stretch is in progress.
	End time.Time
}

type Category struct {
	ID   uint `gorm:"primary_key"`
	Name string
}

// Duration returns the total duration of a frame. If Stretch.End is undefined it treats it as time.Now().
func Duration(frame Stretch) time.Duration {
	var t time.Time
	if frame.End == t {
		return time.Now().Sub(frame.Start)
	}
	return frame.End.Sub(frame.Start)
}

// GetEntryElapsed returns the total duration of frames of the given Task.
func GetEntryElapsed(db *gorm.DB, taskID uint) time.Duration {
	var frames []Stretch
	db.Where("task_id = ?", taskID).Find(&frames)
	var elapsed time.Duration
	for _, frame := range frames {
		elapsed += Duration(frame)
	}
	return elapsed
}

func OpenTestDb() *gorm.DB {
	err := os.Remove("test.db")
	if err != nil {
		fmt.Println("Couldn't remove db", err)
	}
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
	defer db.Close()

	p := Project{Name: "Dreamer"}
	// Create
	db.Create(&p)

	t := Task{}
	db.Create(&t)
	fmt.Println("Created time entry", t.ID)

	start := time.Now().Add(-5 * time.Minute)
	f1 := Stretch{
		TaskID: t.ID,
		Start:  start,
		End:    start.Add(time.Minute),
	}
	f2 := Stretch{
		TaskID: t.ID,
		Start:  start.Add(5 * time.Minute),
		End:    start.Add(10 * time.Minute),
	}

	db.Create(&f1)
	db.Create(&f2)

	elapsed := GetEntryElapsed(db, t.ID)

	fmt.Println("elapsed", elapsed)

	// test(db)
}

func test(db *gorm.DB) {
	// Read
	var project Project
	if e := db.First(&project, 1000); e.Error != nil {
		fmt.Println("Couldn't find 1000")
	}
	if e := db.First(&project, 1); e.Error != nil {
		fmt.Println("Couldn't find 1")
	}                     // find project with id 1
	db.First(&project, "Name = ?", "Dreamer") // find project with code l1212

	// Update - update project's price to 2000
	db.Model(&project).Update("Price", 2000)

	// Delete - delete project
	// db.Delete(&project)
}
