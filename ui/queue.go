package ui

import (
	sp "github.com/op/go-libspotify/spotify"
)

type Queue struct {
	queue []*sp.Track
}

const QUEUE_MAX_ELEMENTS = 100

func InitQueue() *Queue {
	return &Queue{queue: make([]*sp.Track, 0, QUEUE_MAX_ELEMENTS)}
}

func (queue *Queue) Add(track *sp.Track) *sp.Track {
	n := len(queue.queue)
	if n+1 > cap(queue.queue) {
		return nil
	}
	queue.queue = queue.queue[0 : n+1]
	queue.queue[n] = track

	return queue.queue[n]
}

func (queue *Queue) Pop() *sp.Track {
	if len(queue.queue) == 0 {
		return nil
	}
	track := queue.queue[0]
	queue.queue = queue.queue[1:len(queue.queue)]
	return track
}

func (queue *Queue) Remove(index int) *sp.Track {
	if len(queue.queue) == 0 {
		return nil
	}
	track := queue.queue[index]
	queue.queue = append(queue.queue[:index], queue.queue[index+1:]...)
	return track
}

func (queue *Queue) Contents() []*sp.Track {
	return queue.queue
}

func (queue *Queue) isEmpty() bool {
	return len(queue.queue) == 0
}
