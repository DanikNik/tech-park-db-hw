package db

var (
	CreateUserQuery = `
INSERT INTO tp_forum.users (nickname, email, fullname, about) 
VALUES ($1, $2, $3, $4);`

	SelectUsersWithNickOrEmail = `
SELECT nickname, fullname, about, email 
FROM tp_forum.users 
WHERE lower(nickname) = lower($1) OR lower(email) = lower($2)`

	GetUserQuery = `
SELECT nickname, email, fullname, about 
FROM tp_forum.users 
WHERE lower(nickname) = lower($1);`

	UpdateUserQuery = `
UPDATE tp_forum.users 
SET fullname=$1, email=$2, about=$3 
WHERE lower(nickname) = lower($4) 
RETURNING nickname, fullname, about, email;`
)

var (
	CreateForumQuery = `
INSERT INTO tp_forum.forum (slug, title, author) 
VALUES ($1, $2, (SELECT nickname FROM tp_forum.users WHERE lower(nickname) = lower($3))) 
RETURNING  author, slug, title, posts, threads;`

	GetForumQuery = `
SELECT slug, title, author, posts, threads 
FROM tp_forum.forum 
WHERE lower(slug) = lower($1);`
)

var (
	CreateThreadQuery = `
INSERT INTO tp_forum.thread (slug, author, created, forum, title, message)
VALUES (
$1,
(SELECT nickname FROM tp_forum.users WHERE lower(nickname) = lower($2)),
$3,
(SELECT slug FROM tp_forum.forum WHERE lower(slug) = lower($4)),
$5,
$6
)
RETURNING
id, slug, author, created, forum, title, message, votes
`

	GetThreadByIdQuery = `
SELECT id, slug, author, created, forum, title, message, votes
FROM tp_forum.thread
WHERE id=$1
`
	GetThreadBySlugQuery = `
SELECT id, slug, author, created, forum, title, message, votes
FROM tp_forum.thread
WHERE lower(slug)=lower($1)
`

	UpdateThreadByIdQuery = `
UPDATE tp_forum.thread
SET message=$1, title=$2
WHERE id=$3
RETURNING
id, slug, author, created, forum, title, message, votes
`
	UpdateThreadBySlugQuery = `
UPDATE tp_forum.thread
SET message=$1, title=$2
WHERE lower(slug)=lower($3)
RETURNING
id, slug, author, created, forum, title, message, votes
`
)

var (
	SelectAllThreadsQuery = `
	SELECT id, slug, author, created, forum, title, message, votes
	FROM tp_forum.thread
	WHERE lower(forum) = lower($1)
	ORDER BY created
	`

	SelectAllThreadsDescQuery = `
	SELECT id, slug, author, created, forum, title, message, votes
	FROM tp_forum.thread
	WHERE lower(forum) = lower($1)
	ORDER BY created DESC
	`

	SelectAllThreadsLimitQuery = `
	SELECT id, slug, author, created, forum, title, message, votes
	FROM tp_forum.thread
	WHERE lower(forum) = lower($1)
	ORDER BY created
	LIMIT $2
	`

	SelectAllThreadsLimitDescQuery = `
	SELECT id, slug, author, created, forum, title, message, votes
	FROM tp_forum.thread
	WHERE lower(forum) = lower($1)
	ORDER BY created DESC
	LIMIT $2
	`

	SelectAllThreadsSinceQuery = `
	SELECT id, slug, author, created, forum, title, message, votes
	FROM tp_forum.thread
	WHERE lower(forum) = lower($1) AND created >= $2
	ORDER BY created
	`

	SelectAllThreadsSinceDescQuery = `
	SELECT id, slug, author, created, forum, title, message, votes
	FROM tp_forum.thread
	WHERE lower(forum) = lower($1) AND created <= $2
	ORDER BY created DESC
	`

	SelectAllThreadsSinceLimitQuery = `
	SELECT id, slug, author, created, forum, title, message, votes
	FROM tp_forum.thread
	WHERE lower(forum) = lower($1) AND created >= $2
	ORDER BY created
	LIMIT $3
	`

	SelectAllThreadsSinceLimitDescQuery = `
	SELECT id, slug, author, created, forum, title, message, votes
	FROM tp_forum.thread
	WHERE lower(forum) = lower($1) AND created <= $2
	ORDER BY created DESC
	LIMIT $3
	`
)
