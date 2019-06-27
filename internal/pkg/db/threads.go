package db

import (
	"database/sql"
	"github.com/jackc/pgx"
	"strconv"
	"tech-park-db-hw/internal/pkg/models"
)

func isIdOrSlug(slugOrId string) (int, bool) {
	if value, err := strconv.Atoi(slugOrId); err != nil {
		return -1, false
	} else {
		return value, true
	}
}

func slugToNullable(slug string) sql.NullString {
	nullable := sql.NullString{
		String: slug,
		Valid:  true,
	}
	if slug == "" {
		nullable.Valid = false
	}

	return nullable
}

func CreateThread(threadData *models.Thread) error {
	row, _ := QueryRow(CreateThreadQuery, slugToNullable(threadData.Slug), threadData.Author, threadData.Created, threadData.Forum, threadData.Title, threadData.Message)
	err := row.Scan(
		&threadData.Id,
		&threadData.Slug,
		&threadData.Author,
		&threadData.Created,
		&threadData.Forum,
		&threadData.Title,
		&threadData.Message,
		&threadData.Votes,
	)
	if err != nil {
		if err, ok := err.(pgx.PgError); ok {
			switch err.Code {
			case notNullError, foreignKeyError:
				return ErrNotFound
			case uniqueIntegrityError:
				return ErrConflict
			}
		}
		return err
	}
	increaseThreadCount()
	return nil
}

func GetThread(slugOrId string) (*models.Thread, error) {
	threadData := models.Thread{}
	potId, flag := isIdOrSlug(slugOrId)
	var row *pgx.Row
	if flag {
		row, _ = QueryRow(GetThreadByIdQuery, potId)
	} else {
		row, _ = QueryRow(GetThreadBySlugQuery, slugOrId)
	}
	var fSlug sql.NullString
	err := row.Scan(
		&threadData.Id,
		&fSlug,
		&threadData.Author,
		&threadData.Created,
		&threadData.Forum,
		&threadData.Title,
		&threadData.Message,
		&threadData.Votes,
	)
	threadData.Slug = fSlug.String

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &threadData, nil
}

func UpdateThread(slugOrId string, threadUpdateData *models.ThreadUpdate) (*models.Thread, error) {
	threadData := &models.Thread{}
	threadData, err := GetThread(slugOrId)
	if err != nil {
		return nil, err
	}
	if threadUpdateData.Message == "" {
		threadUpdateData.Message = threadData.Message
	}
	if threadUpdateData.Title == "" {
		threadUpdateData.Title = threadData.Title
	}

	potId, flag := isIdOrSlug(slugOrId)
	var row *pgx.Row
	if flag {
		row, _ = QueryRow(UpdateThreadByIdQuery, threadUpdateData.Message, threadUpdateData.Title, potId)
	} else {
		row, _ = QueryRow(UpdateThreadBySlugQuery, threadUpdateData.Message, threadUpdateData.Title, slugOrId)
	}
	err = row.Scan(
		&threadData.Id,
		&threadData.Slug,
		&threadData.Author,
		&threadData.Created,
		&threadData.Forum,
		&threadData.Title,
		&threadData.Message,
		&threadData.Votes,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}

	return threadData, nil
}
