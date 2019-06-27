package db

import (
	"github.com/jackc/pgx"
	"tech-park-db-hw/internal/pkg/models"
)

const (
	checkThreadExistQuery = `
	SELECT FROM tp_forum.thread WHERE id = $1
	`
	threadGetIdQuery = `
	SELECT id FROM tp_forum.thread WHERE lower(slug) = lower($1)
	`
)

func checkThreadExists(id int) (bool, error) {
	err := dbObj.QueryRow(checkThreadExistQuery, id).Scan()
	if err == pgx.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func threadGetId(slug string) (int, bool, error) {
	id := -1
	err := dbObj.QueryRow(threadGetIdQuery, slug).Scan(&id)
	if err == pgx.ErrNoRows {
		return id, false, nil
	}
	return id, true, nil
}

func GetPostsByThread(slugOrId string, limit int, desc bool,
	since int, sort string, posts *models.Posts) error {

	flag := false
	threadID := 0
	if id, isID := isIdOrSlug(slugOrId); isID {
		threadID = id
		flag, _ = checkThreadExists(threadID)
	} else {
		threadID, flag, _ = threadGetId(slugOrId)
	}

	if !flag {
		return ErrNotFound
	}

	return getPostsByThreadId(threadID, limit, desc, since, sort, posts)
}

func getPostsByThreadId(id int, limit int, desc bool,
	since int, sort string, posts *models.Posts) error {

	rows, _ := doQuery(id, limit, desc, since, sort)
	defer rows.Close()
	for rows.Next() {
		post := &models.Post{}
		rows.Scan(&post.Id, &post.Author, &post.Created, &post.IsEdited,
			&post.Message, &post.Parent, &post.Thread, &post.Forum)
		*posts = append(*posts, post)
	}

	return nil
}

func doQuery(id int, limit int, desc bool,
	since int, sort string) (*pgx.Rows, error) {
	var rows *pgx.Rows
	switch sort {
	case "":
		fallthrough
	case "flat":
		if since > 0 {
			if desc {
				rows, _ = dbObj.Query(selectPostsFlatLimitSinceDescByID, id,
					since, limit)
			} else {
				rows, _ = dbObj.Query(selectPostsFlatLimitSinceByID, id,
					since, limit)
			}
		} else {
			if desc == true {
				rows, _ = dbObj.Query(selectPostsFlatLimitDescByID, id, limit)
			} else {
				rows, _ = dbObj.Query(selectPostsFlatLimitByID, id, limit)
			}
		}
	case "tree":
		if since > 0 {
			if desc {
				rows, _ = dbObj.Query(selectPostsTreeLimitSinceDescByID, id,
					since, limit)
			} else {
				rows, _ = dbObj.Query(selectPostsTreeLimitSinceByID, id,
					since, limit)
			}
		} else {
			if desc {
				rows, _ = dbObj.Query(selectPostsTreeLimitDescByID, id, limit)
			} else {
				rows, _ = dbObj.Query(selectPostsTreeLimitByID, id, limit)
			}
		}
	case "parent_tree":
		if since > 0 {
			if desc {
				rows, _ = dbObj.Query(selectPostsParentTreeLimitSinceDescByID, id, id,
					since, limit)
			} else {
				rows, _ = dbObj.Query(selectPostsParentTreeLimitSinceByID, id, id,
					since, limit)
			}
		} else {
			if desc {
				rows, _ = dbObj.Query(selectPostsParentTreeLimitDescByID, id, id,
					limit)
			} else {
				rows, _ = dbObj.Query(selectPostsParentTreeLimitByID, id, id,
					limit)
			}
		}
	}

	return rows, nil
}
