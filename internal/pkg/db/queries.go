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

	threadGetVotesCountQuery = `
	SELECT votes FROM tp_forum.thread WHERE id = $1;	
`
)

const (
	CreatePostsQueryBase = `
	INSERT INTO tp_forum.post
	(forum, thread, author, created, message, parent)
	VALUES
	%s
	RETURNING
	id, author, message, created, thread, forum, parent, is_edited 
`
	UpdatePostMessageQuery = `
	UPDATE tp_forum.post
	SET message = $1
	WHERE id = $2	
	RETURNING id, author, created, is_edited, message, parent, thread, forum`

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

	selectPostWithUserThreadQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum,
	t.id, t.slug, t.author, t.created, t.forum, t.title, t.message, t.votes,
	u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.post p 
	JOIN tp_forum.thread t ON p.thread = t.id
	JOIN tp_forum.users u ON lower(p.author) = lower(u.nickname)
	WHERE p.id = $1`

	selectPostWithForumThreadQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum,
	f.author, f.slug, f.title, f.threads, f.posts,
	t.id, t.slug, t.author, t.created, t.forum, t.title, t.message, t.votes
	FROM tp_forum.post p 
	JOIN tp_forum.thread t ON p.thread = t.id
	JOIN tp_forum.forum f ON lower(p.forum) = lower(f.slug)
	WHERE p.id = $1`

	selectPostWithForumUserQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum,
	f.author, f.slug, f.title, f.threads, f.posts,
	u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.post p 
	JOIN tp_forum.users u ON lower(p.author) = lower(u.nickname)
	JOIN tp_forum.forum f ON lower(p.forum) = lower(f.slug)
	WHERE p.id = $1`

	selectPostWithThreadQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum,
	t.id, t.slug, t.author, t.created, t.forum, t.title, t.message, t.votes
	FROM tp_forum.post p 
	JOIN tp_forum.thread t ON p.thread = t.id
	WHERE p.id = $1`

	selectPostWithForumQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum,
	f.author, f.slug, f.title, f.threads, f.posts
	FROM tp_forum.post p 
	JOIN tp_forum.thread t ON p.thread = t.id
	JOIN tp_forum.forum f ON lower(p.forum) = lower(f.slug)
	WHERE p.id = $1`

	selectPostWithUserQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum,
	u.nickname, u.fullname, u.about, u.email
	FROM tp_forum.post p 
	JOIN tp_forum.users u ON lower(p.author) = lower(u.nickname)
	WHERE p.id = $1`

	selectPostQuery = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p 
	WHERE p.id = $1`
)

const (
	selectPostsFlatLimitByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1
	ORDER BY p.id
	LIMIT $2
`
	selectPostsFlatLimitDescByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1
	ORDER BY p.id DESC
	LIMIT $2
`
	selectPostsFlatLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1 and p.id > $2
	ORDER BY p.id
	LIMIT $3
`
	selectPostsFlatLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1 and p.id < $2
	ORDER BY p.id DESC
	LIMIT $3
`
	selectPostsTreeLimitByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1
	ORDER BY p.path
	LIMIT $2
`
	selectPostsTreeLimitDescByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1
	ORDER BY path DESC
	LIMIT $2
`
	selectPostsTreeLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1 and (p.path > (SELECT p2.path from tp_forum.post p2 where p2.id = $2))
	ORDER BY p.path
	LIMIT $3
`
	selectPostsTreeLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum
	FROM tp_forum.post p
	WHERE p.thread = $1 and (p.path < (SELECT p2.path from tp_forum.post p2 where p2.id = $2))
	ORDER BY p.path DESC
	LIMIT $3
`
	selectPostsParentTreeLimitByID = `
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
	selectPostsParentTreeLimitDescByID = `
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
	selectPostsParentTreeLimitSinceByID = `
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
	selectPostsParentTreeLimitSinceDescByID = `
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
)
