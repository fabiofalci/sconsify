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

func TestQueueAddToLimit(t *testing.T) {
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

}
