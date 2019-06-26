package db

const (
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

const (
	CreateForumQuery = `
INSERT INTO tp_forum.forum (slug, title, author) 
VALUES ($1, $2, (SELECT nickname FROM tp_forum.users WHERE lower(nickname) = lower($3))) 
RETURNING  author, slug, title, posts, threads;`

	GetForumQuery = `
SELECT slug, title, author, posts, threads 
FROM tp_forum.forum 
WHERE lower(slug) = lower($1);`
)

const (
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

const (
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

const (
	CreateVoteQuery = `
	INSERT INTO tp_forum.vote (user_nickname, thread, vote_val)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_nickname, thread)
	DO UPDATE SET vote_val = EXCLUDED.vote_val;
`

	threadUpdateVotesCountQuery = `
	UPDATE tp_forum.thread t
	SET votes = (
		SELECT SUM(vote_val)
		FROM tp_forum.vote
		WHERE thread=$1
	)
	WHERE id=$2
	RETURNING votes`
)

const (
	CreatePostsQuery = `
	INSERT INTO tp_forum.post
	(forum, thread, author, created, message, parent)
	VALUES
	($1, $2, $3, $4, $5, $6)
	RETURNING
	id, author, message, created, thread, forum, parent, is_edited 
`
	UpdatePostMessageQuery = `
	UPDATE tp_forum.post
	SET message = $1
	WHERE id = $2	
	RETURNING id, author, created, is_edited, message, parent, thread, forum`
)
