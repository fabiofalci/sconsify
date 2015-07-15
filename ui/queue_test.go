package ui

import (
	"github.com/fabiofalci/sconsify/sconsify"
	"testing"
)

func TestQueueAddPopEmptyAndContents(t *testing.T) {
	queue := InitQueue()

	track0 := &sconsify.Track{}
	queue.Add(track0)

	trackPop0 := queue.Pop()
	if track0 != trackPop0 {
		t.Error("Queue is not returning right element")
	}

	track1 := &sconsify.Track{}

	queue.Add(track0)
	queue.Add(track1)

	if queue.isEmpty() {
		t.Error("Queue is not adding elements")
	}

	contents := queue.Contents()
	if contents[0] != track0 || contents[1] != track1 {
		t.Error("Queue content is not correct")
	}

	trackPop0 = queue.Pop()
	if track0 != trackPop0 {
		t.Error("Queue is not returning right element")
	}

	trackPop1 := queue.Pop()
	if track1 != trackPop1 {
		t.Error("Queue is not returning right element")
	}

	track := queue.Pop()
	if track != nil {
		t.Error("Queue should return nil but it isn't")
	}
}

func TestQueueInsertAndContents(t *testing.T) {
	queue := InitQueue()

	track0 := &sconsify.Track{}
	track1 := &sconsify.Track{}
	queue.Add(track0)
	queue.Add(track1)

	track2 := &sconsify.Track{}
	queue.Insert(track2)

	contents := queue.Contents()
	if contents[0] != track2 || contents[1] != track0 || contents[2] != track1 {
		t.Error("Queue content is not correct")
	}

	trackPop1 := queue.Pop()
	if track2 != trackPop1 {
		t.Error("Queue is not returning right element")
	}

	trackPop2 := queue.Pop()
	if track0 != trackPop2 {
		t.Error("Queue is not returning right element")
	}

	trackPop3 := queue.Pop()
	if track1 != trackPop3 {
		t.Error("Queue is not returning right element")
	}

	track := queue.Pop()
	if track != nil {
		t.Error("Queue should return nil but it isn't")
	}
}

func TestQueueRemove(t *testing.T) {
	queue := InitQueue()

	track0 := &sconsify.Track{}
	queue.Add(track0)

	track1 := &sconsify.Track{}
	queue.Add(track1)

	track2 := &sconsify.Track{}
	queue.Add(track2)

	trackRemoved := queue.Remove(1)

	if trackRemoved != track1 {
		t.Error("Queue is not removing correctly")
	}

	contents := queue.Contents()
	if contents[0] != track0 || contents[1] != track2 {
		t.Error("Queue content is not correct")
	}
}

func TestQueueRemoveOutOfBounds(t *testing.T) {
	queue := InitQueue()

	queue.Add(&sconsify.Track{})
	queue.Add(&sconsify.Track{})
	queue.Add(&sconsify.Track{})

	if queue.Remove(-1) != nil {
		t.Error("Index -1 is not valid for removal")
	}
	if queue.Remove(3) != nil {
		t.Error("Index 3 is not valid for removal because Size is 3")
	}
	if queue.Remove(2) == nil {
		t.Error("Index 2 is valid for removal because Size is 3")
	}
}

func TestQueueRemoveAll(t *testing.T) {
	queue := InitQueue()

	track0 := &sconsify.Track{}
	queue.Add(track0)

	track1 := &sconsify.Track{}
	queue.Add(track1)

	if queue.isEmpty() {
		t.Error("Queue is not empty")
	}

	queue.RemoveAll()

	if !queue.isEmpty() {
		t.Error("Queue is empty")
	}
}

func TestQueueEmpty(t *testing.T) {
	queue := InitQueue()
	if !queue.isEmpty() {
		t.Error("Queue should be empty after init")
	}

	track := queue.Pop()
	if track != nil {
		t.Error("Queue should return nil but it isn't")
	}
}

func TestQueueAddAndInsertToLimit(t *testing.T) {
	queue := InitQueue()

	for i := 0; i < QUEUE_MAX_ELEMENTS; i++ {
		track := &sconsify.Track{}
		trackAdded := queue.Add(track)
		if track != trackAdded {
			t.Error("Queue add should return the very same element")
		}
	}

	track := &sconsify.Track{}
	trackAdded := queue.Add(track)
	if trackAdded != nil {
		t.Error("Queue reached its limit, it should not add anymore")
	}

	track = &sconsify.Track{}
	trackAdded = queue.Insert(track)
	if trackAdded != trackAdded {
		t.Error("Queue insert should always insert and discard the last one if that's the case")
	}

	queue.Remove(99)
	track = &sconsify.Track{}
	trackAdded = queue.Add(track)
	if track != trackAdded {
		t.Error("Queue add should return the very same element")
	}

	track = &sconsify.Track{}
	trackAdded = queue.Add(track)
	if trackAdded != nil {
		t.Error("Queue reached its limit, it should not add anymore")
	}
}

func TestQueueAddAndPopToLimit(t *testing.T) {
	queue := InitQueue()

	for i := 0; i < QUEUE_MAX_ELEMENTS; i++ {
		track := &sconsify.Track{}
		trackAdded := queue.Add(track)
		if track != trackAdded {
			t.Error("Queue add should return the very same element")
		}
	}

	track := &sconsify.Track{}
	trackAdded := queue.Add(track)
	if trackAdded != nil {
		t.Error("Queue reached its limit, it should not add anymore")
	}

	track = &sconsify.Track{}
	trackAdded = queue.Insert(track)
	if trackAdded != trackAdded {
		t.Error("Queue insert should always insert and discard the last one if that's the case")
	}

	queue.Pop()

	track = &sconsify.Track{}
	trackAdded = queue.Add(track)
	if track != trackAdded {
		t.Error("Queue add should return the very same element")
	}

	track = &sconsify.Track{}
	trackAdded = queue.Add(track)
	if trackAdded != nil {
		t.Error("Queue reached its limit, it should not add anymore")
	}
}
