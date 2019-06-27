package db

import (
	"database/sql"
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
	if err != nil {
		if err == ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
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
					if pqError.ConstraintName == "post_author_fkey" || pqError.ConstraintName == "post_forum_fkey" {
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

func GetPostFullData(related []string, postFullData *models.PostFull) error {
	includeUser, includeForum, includeThread := false, false, false
	for _, rel := range related {
		switch rel {
		case "user":
			postFullData.Author = &models.User{}
			includeUser = true
		case "forum":
			postFullData.Forum = &models.Forum{}
			includeForum = true
		case "thread":
			postFullData.Thread = &models.Thread{}
			includeThread = true
		}
	}

	var err error
	if includeForum && includeUser && includeThread {
		err = selectPostWithForumUserThread(postFullData)
	} else if !includeForum && includeUser && includeThread {
		err = selectPostWithUserThread(postFullData)
	} else if includeForum && !includeUser && includeThread {
		err = selectPostWithForumThread(postFullData)
	} else if includeForum && includeUser && !includeThread {
		err = selectPostWithForumUser(postFullData)
	} else if !includeForum && !includeUser && includeThread {
		err = selectPostWithThread(postFullData)
	} else if !includeForum && includeUser && !includeThread {
		err = selectPostWithUser(postFullData)
	} else if includeForum && !includeUser && !includeThread {
		err = selectPostWithForum(postFullData)
	} else if !includeForum && !includeUser && !includeThread {
		err = getPost(postFullData.Post)
	}

	if err == pgx.ErrNoRows {
		return ErrNotFound
	}

	return nil
}

func selectPostWithForumUserThread(pf *models.PostFull) error {
	slugThread := sql.NullString{}
	err := dbObj.QueryRow(selectPostWithForumUserThreadQuery, pf.Post.Id).Scan(
		&pf.Post.Id,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&pf.Post.Parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Forum.User,
		&pf.Forum.Slug,
		&pf.Forum.Title,
		&pf.Forum.Threads,
		&pf.Forum.Posts,
		&pf.Thread.Id,
		&slugThread,
		&pf.Thread.Author,
		&pf.Thread.Created,
		&pf.Thread.Forum,
		&pf.Thread.Title,
		&pf.Thread.Message,
		&pf.Thread.Votes,
		&pf.Author.Nickname,
		&pf.Author.Fullname,
		&pf.Author.About,
		&pf.Author.Email,
	)
	if err != nil {
		return err
	}
	if slugThread.Valid {
		pf.Thread.Slug = slugThread.String
	} else {
		pf.Thread.Slug = ""
	}
	return nil
}

func selectPostWithUserThread(pf *models.PostFull) error {
	slugThread := sql.NullString{}
	err := dbObj.QueryRow(selectPostWithUserThreadQuery, pf.Post.Id).Scan(
		&pf.Post.Id,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&pf.Post.Parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Thread.Id,
		&slugThread,
		&pf.Thread.Author,
		&pf.Thread.Created,
		&pf.Thread.Forum,
		&pf.Thread.Title,
		&pf.Thread.Message,
		&pf.Thread.Votes,
		&pf.Author.Nickname,
		&pf.Author.Fullname,
		&pf.Author.About,
		&pf.Author.Email,
	)

	if err != nil {
		return err
	}

	if slugThread.Valid {
		pf.Thread.Slug = slugThread.String
	} else {
		pf.Thread.Slug = ""
	}
	return nil
}

func selectPostWithForumThread(pf *models.PostFull) error {
	slugThread := sql.NullString{}
	err := dbObj.QueryRow(selectPostWithForumThreadQuery, pf.Post.Id).Scan(
		&pf.Post.Id,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&pf.Post.Parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Forum.User,
		&pf.Forum.Slug,
		&pf.Forum.Title,
		&pf.Forum.Threads,
		&pf.Forum.Posts,
		&pf.Thread.Id,
		&slugThread,
		&pf.Thread.Author,
		&pf.Thread.Created,
		&pf.Thread.Forum,
		&pf.Thread.Title,
		&pf.Thread.Message,
		&pf.Thread.Votes,
	)

	if err != nil {
		return err
	}

	if slugThread.Valid {
		pf.Thread.Slug = slugThread.String
	} else {
		pf.Thread.Slug = ""
	}
	return nil
}

func selectPostWithForumUser(pf *models.PostFull) error {
	err := dbObj.QueryRow(selectPostWithForumUserQuery, pf.Post.Id).Scan(
		&pf.Post.Id,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&pf.Post.Parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Forum.User,
		&pf.Forum.Slug,
		&pf.Forum.Title,
		&pf.Forum.Threads,
		&pf.Forum.Posts,
		&pf.Author.Nickname,
		&pf.Author.Fullname,
		&pf.Author.About,
		&pf.Author.Email,
	)

	return err
}

func selectPostWithThread(pf *models.PostFull) error {
	slugThread := sql.NullString{}
	err := dbObj.QueryRow(selectPostWithThreadQuery, pf.Post.Id).Scan(
		&pf.Post.Id,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&pf.Post.Parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Thread.Id,
		&slugThread,
		&pf.Thread.Author,
		&pf.Thread.Created,
		&pf.Thread.Forum,
		&pf.Thread.Title,
		&pf.Thread.Message,
		&pf.Thread.Votes,
	)

	if err != nil {
		return err
	}
	if slugThread.Valid {
		pf.Thread.Slug = slugThread.String
	} else {
		pf.Thread.Slug = ""
	}
	return nil
}

func selectPostWithForum(pf *models.PostFull) error {
	err := dbObj.QueryRow(selectPostWithForumQuery, pf.Post.Id).Scan(
		&pf.Post.Id,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&pf.Post.Parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Forum.User,
		&pf.Forum.Slug,
		&pf.Forum.Title,
		&pf.Forum.Threads,
		&pf.Forum.Posts,
	)
	return err
}

func selectPostWithUser(pf *models.PostFull) error {
	err := dbObj.QueryRow(selectPostWithUserQuery, pf.Post.Id).Scan(
		&pf.Post.Id,
		&pf.Post.Author,
		&pf.Post.Created,
		&pf.Post.IsEdited,
		&pf.Post.Message,
		&pf.Post.Parent,
		&pf.Post.Thread,
		&pf.Post.Forum,
		&pf.Author.Nickname,
		&pf.Author.Fullname,
		&pf.Author.About,
		&pf.Author.Email,
	)
	return err
}

func getPost(postFullData *models.Post) error {
	return scanPostData(dbObj.QueryRow(selectPostQuery, postFullData.Id), postFullData)
}

func scanPostData(r *pgx.Row, post *models.Post) error {
	err := r.Scan(&post.Id, &post.Author, &post.Created, &post.IsEdited,
		&post.Message, &post.Parent, &post.Thread, &post.Forum)
	return err
}

func UpdatePost(post *models.Post, pu *models.PostUpdate) error {
	var err error
	if pu.Message == "" {
		err = getPost(post)
	} else {
		err = scanPostData(dbObj.QueryRow(UpdatePostMessageQuery, pu.Message, post.Id), post)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	return nil
}
