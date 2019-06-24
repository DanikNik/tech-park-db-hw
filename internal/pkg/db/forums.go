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
