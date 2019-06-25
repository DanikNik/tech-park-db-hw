package db

import (
	"sync/atomic"
	"tech-park-db-hw/internal/pkg/models"
)

func Status(s *models.Status) error {

	s.Forum = atomic.LoadInt32(&forumCount)
	s.Thread = atomic.LoadInt32(&threadCount)
	s.Post = atomic.LoadInt32(&postCount)
	s.User = atomic.LoadInt32(&userCount)

	return nil
}
