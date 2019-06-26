package db

import (
	"sync/atomic"
	"tech-park-db-hw/internal/pkg/models"
)

func Truncate() error {
	if dbObj == nil {
		return NotInit
	}

	_, err := Exec("TRUNCATE TABLE tp_forum.users, tp_forum.forum, tp_forum.thread, tp_forum.post, tp_forum.vote CASCADE;")
	if err != nil {
		return err
	}
	Exec("INSERT INTO tp_forum.post (id) VALUES (0)")
	atomic.SwapInt32(&forumCount, 0)
	atomic.SwapInt32(&threadCount, 0)
	atomic.SwapInt32(&postCount, 0)
	atomic.SwapInt32(&userCount, 0)
	return nil
}

var forumCount int32
var threadCount int32
var postCount int32
var userCount int32

func increaseForumCount() {
	atomic.AddInt32(&forumCount, 1)
}

func increaseThreadCount() {
	atomic.AddInt32(&threadCount, 1)
}

func increaseUserCount() {
	atomic.AddInt32(&userCount, 1)
}

func increasePostCount(count int32) {
	atomic.AddInt32(&postCount, count)
}

func Status(s *models.Status) error {

	s.Forum = atomic.LoadInt32(&forumCount)
	s.Thread = atomic.LoadInt32(&threadCount)
	s.Post = atomic.LoadInt32(&postCount)
	s.User = atomic.LoadInt32(&userCount)

	return nil
}
