package ui

import (
	"github.com/fabiofalci/sconsify/sconsify"
)

type Queue struct {
	queue []*sconsify.Track
}

const QUEUE_MAX_ELEMENTS = 100

func InitQueue() *Queue {
	return &Queue{queue: make([]*sconsify.Track, 0, 0)}
}

func (queue *Queue) Add(track *sconsify.Track) *sconsify.Track {
	n := len(queue.queue)
	if n >= QUEUE_MAX_ELEMENTS {
		return nil
	}
	queue.queue = append(queue.queue, track)

	return queue.queue[n]
}

func (queue *Queue) Insert(track *sconsify.Track) *sconsify.Track {
	n := len(queue.queue)
	if n >= QUEUE_MAX_ELEMENTS {
		queue.Remove(QUEUE_MAX_ELEMENTS - 1)
	}

	queue.queue = append(queue.queue, nil)

	copy(queue.queue[1:], queue.queue)
	queue.queue[0] = track

	return queue.queue[0]
}

func (queue *Queue) Pop() *sconsify.Track {
	if len(queue.queue) == 0 {
		return nil
	}
	track := queue.queue[0]
	queue.queue = queue.queue[1:len(queue.queue)]
	return track
}

func (queue *Queue) RemoveAll() {
	if len(queue.queue) == 0 {
		return
	}

	queue.queue = make([]*sconsify.Track, 0, 0)
}

func (queue *Queue) Remove(index int) *sconsify.Track {
	if len(queue.queue) == 0 || index < 0 || index >= len(queue.queue) {
		return nil
	}
	track := queue.queue[index]
	queue.queue = append(queue.queue[:index], queue.queue[index+1:]...)
	return track
}

func (queue *Queue) Contents() []*sconsify.Track {
	return queue.queue
}

func (queue *Queue) IsEmpty() bool {
	return len(queue.queue) == 0
}
