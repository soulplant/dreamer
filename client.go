package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Client struct {
	db    *gorm.DB
	clock Clock
}

func NewClient(db *gorm.DB, clock Clock) *Client {
	return &Client{db, clock}
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) NewProject(name string) (uint, error) {
	p := Project{Name: name}
	if db := c.db.Create(&p); db.Error != nil {
		return 0, db.Error
	}
	return p.ID, nil
}

func (c *Client) NewCategory(projectID uint, name string) (uint, error) {
	cat := Category{Name: name}
	if db := c.db.Create(&cat); db.Error != nil {
		return 0, db.Error
	}
	return cat.ID, nil
}

func (c *Client) NewUser(email string) (uint, error) {
	user := User{
		Email: email,
	}
	if e := c.db.Create(&user); e.Error != nil {
		return 0, e.Error
	}
	return user.ID, nil
}

func (c *Client) CreateTask(categoryID uint, userID uint, description string) (uint, error) {
	task := Task{
		CategoryID:  categoryID,
		UserID:      userID,
		Description: description,
	}
	if db := c.db.Create(&task); db.Error != nil {
		return 0, db.Error
	}
	return task.ID, nil
}

func (c *Client) StartTimer(taskID uint) (uint, error) {
	// Stop any existing timer for the user.
	db := c.db.Begin()
	commit := false
	defer func() {
		if commit {
			db.Commit()
		} else {
			db.Rollback()
		}
	}()

	// Stop the currently running stretch if any.
	var oldStretch Stretch
	if e := db.First(&oldStretch, "end = 0"); e.Error == nil {
		if e := db.Model(&oldStretch).Update("end = ?", c.clock.Now()); e.Error != nil {
			return 0, e.Error
		}
	}

	var task Task
	if e := c.db.First(&task, taskID); e.Error != nil {
		return 0, e.Error
	}

	newStretch := Stretch{
		Start:  c.clock.Now(),
		TaskID: task.ID,
	}
	if e := db.Create(&newStretch); e.Error != nil {
		return 0, e.Error
	}
	commit = true
	return newStretch.ID, nil
}

// GetElapsedTime returns the total duration of Stretches for the given Task.
func (c *Client) GetElapsedTime(taskID uint) time.Duration {
	var stretches []Stretch
	c.db.Where("task_id = ?", taskID).Find(&stretches)
	var elapsed time.Duration
	for _, frame := range stretches {
		elapsed += c.Duration(frame)
	}
	return elapsed
}

// Duration returns the total duration of a frame. If Stretch.End is undefined it treats it as time.Now().
func (c *Client) Duration(frame Stretch) time.Duration {
	var t time.Time
	if frame.End == t {
		return c.clock.Now().Sub(frame.Start)
	}
	return frame.End.Sub(frame.Start)
}

func (c *Client) StopTimer(taskID uint) error {
	var stretch Stretch
	if e := c.db.First(&stretch, "task_id = ? AND end = ?", taskID, time.Time{}); e.Error != nil {
		return e.Error
	}
	if e := c.db.Model(&stretch).Update(Stretch{End: c.clock.Now()}); e.Error != nil {
		return e.Error
	}
	return nil
}

func (c *Client) GetStretches(taskID uint) ([]Stretch, error) {
	var stretches []Stretch
	c.db.Find(&stretches, "task_id = ?", taskID)
	return stretches, nil
}
