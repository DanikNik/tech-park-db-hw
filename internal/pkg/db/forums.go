package db

import (
	"tech-park-db-hw/internal/pkg/models"
)

func CreateForum(forum models.Forum) error {
	_, err := Exec(CreateForumQuery, forum.Slug, forum.Title, forum.User)
	return err
}
