package db

import (
	"github.com/jackc/pgx"
	"tech-park-db-hw/internal/pkg/models"
)

const (
	checkThreadExistByID = `
	SELECT FROM tp_forum.thread WHERE id = $1
	`
	checkThreadExistAndGetIDBySlug = `
	SELECT id FROM tp_forum.thread WHERE lower(slug) = lower($1)
	`
)

func isThreadExist(id int) (bool, error) {
	err := dbObj.QueryRow(checkThreadExistByID, id).Scan()
	if err == pgx.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func ifThreadExistGetID(slug string) (int, bool, error) {
	id := -1
	err := dbObj.QueryRow(checkThreadExistAndGetIDBySlug, slug).Scan(&id)
	if err == pgx.ErrNoRows {
		return id, false, nil
	}
	return id, true, nil
}

func SelectAllPostsByThread(slugOrIDThread string, limit int, desc bool,
	since int, sort string, posts *[]*models.Post) error {

	isExist := false
	threadID := 0
	if id, isID := isIdOrSlug(slugOrIDThread); isID {
		threadID = id
		isExist, _ = isThreadExist(threadID)
	} else {
		threadID, isExist, _ = ifThreadExistGetID(slugOrIDThread)
	}

	if !isExist {
		return ErrNotFound
	}

	return selectAllPostsByThreadID(threadID, limit, desc, since, sort, posts)
}

const selectPostsFlatLimitByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1
	ORDER BY p.id
	LIMIT $2
`

const selectPostsFlatLimitDescByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1
	ORDER BY p.id DESC
	LIMIT $2
`

const selectPostsFlatLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1 and p.id > $2
	ORDER BY p.id
	LIMIT $3
`
const selectPostsFlatLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1 and p.id < $2
	ORDER BY p.id DESC
	LIMIT $3
`

const selectPostsTreeLimitByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1
	ORDER BY p.path
	LIMIT $2
`

const selectPostsTreeLimitDescByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1
	ORDER BY path DESC
	LIMIT $2
`

const selectPostsTreeLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1 and (p.path > (SELECT p2.path from tp_forum.post p2 where p2.id = $2))
	ORDER BY p.path
	LIMIT $3
`

const selectPostsTreeLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1 and (p.path < (SELECT p2.path from tp_forum.post p2 where p2.id = $2))
	ORDER BY p.path DESC
	LIMIT $3
`

const selectPostsParentTreeLimitByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM tp_forum.post p2
		WHERE p2.thread = $2 AND p2.parent = 0
		ORDER BY p2.id
		LIMIT $3
	)
	ORDER BY path
`

const selectPostsParentTreeLimitDescByID = `
SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
FROM tp_forum.post p
WHERE p.thread = $1 and p.path[1] IN (
    SELECT p2.path[1]
    FROM tp_forum.post p2
	WHERE p2.thread = $2 AND p2.parent = 0
	ORDER BY p2.id DESC
    LIMIT $3
)
ORDER BY p.path[1] DESC, p.path[2:]
`

const selectPostsParentTreeLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM tp_forum.post p2
		WHERE p2.thread = $2 AND p2.parent = 0 and p2.path[1] > (SELECT p3.path[1] from tp_forum.post p3 where p3.id = $3)
		ORDER BY p2.id
		LIMIT $4
	)
	ORDER BY p.path
`

const selectPostsParentTreeLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM tp_forum.post p2
		WHERE p2.thread = $2 AND p2.parent = 0 and p2.path[1] < (SELECT p3.path[1] from tp_forum.post p3 where p3.id = $3)
		ORDER BY p2.id DESC
		LIMIT $4
	)
	ORDER BY p.path[1] DESC, p.path[2:]
`

func selectAllPostsByThreadID(id int, limit int, desc bool,
	since int, sort string, posts *[]*models.Post) error {

	rows, _ := doQuery(id, limit, desc, since, sort)
	defer rows.Close()
	for rows.Next() {
		post := &models.Post{}
		scanPostRows(rows, post)
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
