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

func GetPost(id int) (*models.PostFull, error) {
	return nil, nil
}

func SelectPostFull(related []string, pf *models.PostFull) error {
	isIncludeUser, isIncludeForum, isIncludeThread := false, false, false
	for _, rel := range related {
		switch rel {
		case "user":
			pf.Author = &models.User{}
			isIncludeUser = true
		case "forum":
			pf.Forum = &models.Forum{}
			isIncludeForum = true
		case "thread":
			pf.Thread = &models.Thread{}
			isIncludeThread = true
		}
	}

	var err error
	if isIncludeForum && isIncludeUser && isIncludeThread {
		err = selectPostWithForumUserThread(pf)
	} else if !isIncludeForum && isIncludeUser && isIncludeThread {
		err = selectPostWithUserThread(pf)
	} else if isIncludeForum && !isIncludeUser && isIncludeThread {
		err = selectPostWithForumThread(pf)
	} else if isIncludeForum && isIncludeUser && !isIncludeThread {
		err = selectPostWithForumUser(pf)
	} else if !isIncludeForum && !isIncludeUser && isIncludeThread {
		err = selectPostWithThread(pf)
	} else if !isIncludeForum && isIncludeUser && !isIncludeThread {
		err = selectPostWithUser(pf)
	} else if isIncludeForum && !isIncludeUser && !isIncludeThread {
		err = selectPostWithForum(pf)
	} else if !isIncludeForum && !isIncludeUser && !isIncludeThread {
		err = selectPost(pf.Post)
	}

	if err == pgx.ErrNoRows {
		return ErrNotFound
	}

	return nil
}

const (
	selectPostWithForumUserThreadQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum,
	f.author, f.slug, f.title, f.threads, f.posts,
	t.id, t.slug, t.author, t.created, t.forum, t.title, t.message, t.votes,
	u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.post p 
	JOIN tp_forum.thread t ON p.thread = t.id
	JOIN tp_forum.users u ON lower(p.author) = lower(u.nickname)
	JOIN tp_forum.forum f ON p.forum = f.slug
	WHERE p.id = $1`
)

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

const (
	selectPostWithUserThreadQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum,
	t.id, t.slug, t.author, t.created, t.forum, t.title, t.message, t.votes,
	u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.post p 
	JOIN tp_forum.thread t ON p.thread = t.id
	JOIN tp_forum.users u ON lower(p.author) = lower(u.nickname)
	WHERE p.id = $1`
)

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

const (
	selectPostWithForumThreadQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum,
	f.author, f.slug, f.title, f.threads, f.posts,
	t.id, t.slug, t.author, t.created, t.forum, t.title, t.message, t.votes
	FROM tp_forum.post p 
	JOIN tp_forum.thread t ON p.thread = t.id
	JOIN tp_forum.forum f ON lower(p.forum) = lower(f.slug)
	WHERE p.id = $1`
)

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

const (
	selectPostWithForumUserQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum,
	f.author, f.slug, f.title, f.threads, f.posts,
	u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.post p 
	JOIN tp_forum.users u ON lower(p.author) = lower(u.nickname)
	JOIN tp_forum.forum f ON lower(p.forum) = lower(f.slug)
	WHERE p.id = $1`
)

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

const (
	selectPostWithThreadQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum,
	t.id, t.slug, t.author, t.created, t.forum, t.title, t.message, t.votes
	FROM tp_forum.post p 
	JOIN tp_forum.thread t ON p.thread = t.id
	WHERE p.id = $1`
)

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

const (
	selectPostWithForumQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum,
	f.author, f.slug, f.title, f.threads, f.posts
	FROM tp_forum.post p 
	JOIN tp_forum.thread t ON p.thread = t.id
	JOIN tp_forum.forum f ON lower(p.forum) = lower(f.slug)
	WHERE p.id = $1`
)

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

const (
	selectPostWithUserQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum,
	u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.post p 
	JOIN tp_forum.users u ON lower(p.author) = lower(u.nickname)
	WHERE p.id = $1`
)

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

const (
	selectPostQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p 
	WHERE p.id = $1`
)

func selectPost(pf *models.Post) error {
	return scanPost(dbObj.QueryRow(selectPostQuery, pf.Id), pf)
}

func scanPostRows(r *pgx.Rows, post *models.Post) error {
	err := r.Scan(&post.Id, &post.Author, &post.Created, &post.IsEdited,
		&post.Message, &post.Parent, &post.Thread, &post.Forum)
	return err
}

func scanPost(r *pgx.Row, post *models.Post) error {
	err := r.Scan(&post.Id, &post.Author, &post.Created, &post.IsEdited,
		&post.Message, &post.Parent, &post.Thread, &post.Forum)
	return err
}

func UpdatePost(post *models.Post, pu *models.PostUpdate) error {
	var err error
	if pu.Message == "" {
		err = selectPost(post)
	} else {
		err = scanPost(dbObj.QueryRow(UpdatePostMessageQuery, pu.Message, post.Id), post)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	return nil
}
