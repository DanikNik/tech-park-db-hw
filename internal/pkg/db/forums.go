package db

import (
	"database/sql"
	"github.com/jackc/pgx"
	"tech-park-db-hw/internal/pkg/models"
)

func CreateForum(forum *models.Forum) error {
	row, _ := QueryRow(CreateForumQuery, forum.Slug, forum.Title, forum.User)
	err := row.Scan(
		&forum.User,
		&forum.Slug,
		&forum.Title,
		&forum.Posts,
		&forum.Threads,
	)
	if err != nil {
		if pqError, ok := err.(pgx.PgError); ok {
			switch pqError.Code {
			case uniqueIntegrityError:
				return ErrConflict
			case notNullError:
				return ErrNotFound
			}
		}
	}
	increaseForumCount()
	return err
}

func GetForum(slug string) (*models.Forum, error) {
	row, err := QueryRow(GetForumQuery, slug)
	if err != nil {
		return nil, err
	}
	forumData := &models.Forum{}
	err = row.Scan(
		&forumData.Slug,
		&forumData.Title,
		&forumData.User,
		&forumData.Posts,
		&forumData.Threads,
	)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	return forumData, nil
}

const (
	checkForumExistQuery = `SELECT FROM tp_forum.forum WHERE lower(slug) = lower($1)`
)

func checkForumExist(slug string) (bool, error) {
	err := dbObj.QueryRow(checkForumExistQuery, slug).Scan()
	if err == pgx.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func GetThreadsByForum(forumSlug string, limit int, desc bool, since string) (*[]models.Thread, error) {
	if isExist, _ := checkForumExist(forumSlug); !isExist {
		return nil, ErrNotFound
	}

	var rows *pgx.Rows
	if desc == true {
		if limit > 0 && since != "" {
			rows, _ = Query(SelectAllThreadsSinceLimitDescQuery, forumSlug, since, limit)
		} else if limit > 0 {
			rows, _ = Query(SelectAllThreadsLimitDescQuery, forumSlug, limit)
		} else if since != "" {
			rows, _ = Query(SelectAllThreadsSinceDescQuery, forumSlug, since)
		} else {
			rows, _ = Query(SelectAllThreadsDescQuery, forumSlug)
		}
	} else {
		if limit > 0 && since != "" {
			rows, _ = Query(SelectAllThreadsSinceLimitQuery, forumSlug, since, limit)
		} else if limit > 0 {
			rows, _ = Query(SelectAllThreadsLimitQuery, forumSlug, limit)
		} else if since != "" {
			rows, _ = Query(SelectAllThreadsSinceQuery, forumSlug, since)
		} else {
			rows, _ = Query(SelectAllThreadsQuery, forumSlug)
		}
	}

	defer rows.Close()

	ts := []models.Thread{}
	for rows.Next() {
		threadData := models.Thread{}
		slug := sql.NullString{}
		err := rows.Scan(&threadData.Id, &slug, &threadData.Author, &threadData.Created, &threadData.Forum, &threadData.Title, &threadData.Message, &threadData.Votes)
		if err != nil {
			return nil, err
		}
		if slug.Valid {
			threadData.Slug = slug.String
		}

		ts = append(ts, threadData)
	}
	return &ts, nil

}

const (
	selectAllUsersByForum = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.forum_user fu
	JOIN tp_forum.users u ON lower(fu.nickname) = lower(u.nickname)
	WHERE lower(fu.forum) = lower($1)
	ORDER BY lower(fu.nickname)`

	selectAllUsersByForumDesc = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.forum_user fu
	JOIN tp_forum.users u ON lower(fu.nickname) = lower(u.nickname)
	WHERE lower(fu.forum) = lower($1)
	ORDER BY lower(fu.nickname)`

	selectAllUsersByForumLimit = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.forum_user fu
	JOIN tp_forum.users u ON lower(fu.nickname) = lower(u.nickname)
	WHERE lower(fu.forum) = lower($1)
	ORDER BY lower(fu.nickname)
	LIMIT $2`

	selectAllUsersByForumLimitDesc = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.forum_user fu
	JOIN tp_forum.users u ON lower(fu.nickname) = lower(u.nickname)
	WHERE lower(fu.forum) = lower($1)
	ORDER BY lower(fu.nickname) DESC
	LIMIT $2`

	selectAllUsersByForumLimitSince = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.forum_user fu
	JOIN tp_forum.users u ON lower(fu.nickname) = lower(u.nickname)
	WHERE lower(fu.forum) = lower($1) AND lower(fu.nickname) > lower($2)
	ORDER BY lower(fu.nickname)
	LIMIT $3`

	selectAllUsersByForumLimitSinceDesc = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.forum_user fu
	JOIN tp_forum.users u ON lower(fu.nickname) = lower(u.nickname)
	WHERE lower(fu.forum) = lower($1) AND lower(fu.nickname) < lower($2)
	ORDER BY lower(fu.nickname) DESC
	LIMIT $3`

	selectAllUsersByForumSince = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.forum_user fu
	JOIN tp_forum.users u ON lower(fu.nickname) = lower(u.nickname)
	WHERE lower(fu.forum) = lower($1) AND lower(fu.nickname) > lower($2)
	ORDER BY lower(fu.nickname)`

	selectAllUsersByForumSinceDesc = `
	SELECT u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.forum_user fu
	JOIN tp_forum.users u ON lower(fu.nickname) = lower(u.nickname)
	WHERE lower(fu.forum) = lower($1) AND lower(fu.nickname) < lower($2)
	ORDER BY lower(fu.nickname) DESC`
)

func GetUsersByForum(slug string, limit int, desc bool, since string) (*[]models.User, error) {

	if isExist, _ := checkForumExist(slug); !isExist {
		return nil, ErrNotFound
	}

	var rows *pgx.Rows
	if desc == true {
		if since != "" && limit > 0 {
			rows, _ = Query(selectAllUsersByForumLimitSinceDesc, slug, since, limit)
		} else if since != "" {
			rows, _ = Query(selectAllUsersByForumSinceDesc, slug, since)
		} else if limit > 0 {
			rows, _ = Query(selectAllUsersByForumLimitDesc, slug, limit)
		} else {
			rows, _ = Query(selectAllUsersByForumDesc, slug)
		}
	} else {
		if since != "" && limit > 0 {
			rows, _ = Query(selectAllUsersByForumLimitSince, slug, since, limit)
		} else if since != "" {
			rows, _ = Query(selectAllUsersByForumSince, slug, since)
		} else if limit > 0 {
			rows, _ = Query(selectAllUsersByForumLimit, slug, limit)
		} else {
			rows, _ = Query(selectAllUsersByForum, slug)
		}
	}
	defer rows.Close()

	userList := []models.User{}
	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			panic(err)
		}
		userList = append(userList, user)
	}
	return &userList, nil
}
