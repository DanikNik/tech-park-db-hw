package db

import (
	"github.com/jackc/pgx"
	"tech-park-db-hw/internal/pkg/models"
	"time"
)

const (
	getParentThreadQuery = `
	SELECT thread
	FROM tp_forum.post
	WHERE id = $1
`
)

func CreatePostsBulk(slugOrId string, posts []models.Post) (*[]models.Post, error) {
	threadData, err := GetThread(slugOrId)
	if err == ErrNotFound {
		return nil, ErrNotFound
	}

	resultPosts := []models.Post{}
	if posts == nil || len(posts) == 0 {
		return &resultPosts, nil
	}

	tx, _ := dbObj.Begin()
	creationTime := time.Now()

	for _, post := range posts {
		newPost := models.Post{}
		row := tx.QueryRow(
			CreatePostsQuery,
			threadData.Forum,
			threadData.Id,
			post.Author,
			creationTime,
			post.Message,
			post.Parent,
		)
		err := row.Scan(
			&newPost.Id,
			&newPost.Author,
			&newPost.Message,
			&newPost.Created,
			&newPost.Thread,
			&newPost.Forum,
			&newPost.Parent,
			&newPost.IsEdited,
		)

		if newPost.Parent != 0 {
			var parentThread int
			tx.QueryRow(getParentThreadQuery, newPost.Parent).Scan(&parentThread)
			if parentThread != newPost.Thread {
				tx.Rollback()
				return nil, ErrConflict
			}
		}

		if err != nil {
			tx.Rollback()

			if pqError, ok := err.(pgx.PgError); ok {
				switch pqError.Code {
				case foreignKeyError:
					if pqError.ConstraintName == "post_parent_fkey" {
						return nil, ErrConflict
					}
					if pqError.ConstraintName == "post_author_fkey" {
						return nil, ErrNotFound
					}
				}
			}
			return nil, err
		}
		resultPosts = append(resultPosts, newPost)
	}

	tx.Commit()
	increasePostCount(int32(len(resultPosts)))

	return &resultPosts, nil
}
