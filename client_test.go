package main

import (
	"testing"
	"time"
)

type testVars struct {
	projectID uint
	catID     uint
	userID    uint
	taskID    uint
}

func setupTest(t *testing.T, client *Client) testVars {
	projectID, err := client.NewProject("Dreamco")
	if err != nil {
		t.Fatal(err)
	}
	catID, err := client.NewCategory(projectID, "Dream Stack")
	if err != nil {
		t.Fatal(err)
	}
	userID, err := client.NewUser("james@dreamco.io")
	if err != nil {
		t.Fatal(err)
	}
	taskID, err := client.CreateTask(catID, userID, "Write some tests for dreamer")
	if err != nil {
		t.Fatal(err)
	}
	if taskID == 0 {
		t.Fatal("Got a task id of 0")
	}
	return testVars{
		projectID: projectID,
		catID:     catID,
		userID:    userID,
		taskID:    taskID,
	}
}

func TestClient_TimerDuration(t *testing.T) {
	db := OpenTestDb()
	startTime := mkTime("3:00PM")
	clock := NewFakeClock(startTime)
	client := NewClient(db, clock)
	defer client.Close()
	vars := setupTest(t, client)
	taskID := vars.taskID

	_, err := client.StartTimer(taskID)
	if err != nil {
		t.Fatal(err)
	}
	elapsed := client.GetElapsedTime(taskID)
	if elapsed != 0 {
		t.Error("Expected no time to have passed, actual", elapsed)
	}
	clock.SetTime(startTime.Add(time.Hour))
	expectElapsed(t, client, taskID, time.Hour)
	clock.SetTime(startTime.Add(2 * time.Hour))
	expectElapsed(t, client, taskID, 2*time.Hour)

	if err = client.StopTimer(taskID); err != nil {
		t.Error("Couldn't stop timer", err)
	}

	expectElapsed(t, client, taskID, 2 * time.Hour)
	clock.SetTime(startTime.Add(24 * time.Hour))
	expectElapsed(t, client, taskID, 2 * time.Hour)

	_, err = client.StartTimer(taskID)
	if err != nil {
		t.Error(err)
	}

	clock.SetTime(startTime.Add(48 * time.Hour));
	expectElapsed(t, client, taskID, 26 * time.Hour)

	stretches, err := client.GetStretches(taskID)
	if err != nil {
		t.Error(err)
	}
	if len(stretches) != 2 {
		t.Error("Expected 2 stretches, got ", len(stretches))
	}
}

func expectElapsed(t *testing.T, client *Client, taskID uint, expected time.Duration) {
	elapsed := client.GetElapsedTime(taskID)
	if elapsed != expected {
		t.Error("Expected", expected, "to have passed, but actually", elapsed)
	}
}
