DROP TABLE IF EXISTS tp_forum.users, tp_forum.forum, tp_forum.thread, tp_forum.post, tp_forum.vote CASCADE;
DROP SCHEMA IF EXISTS tp_forum CASCADE;
DROP EXTENSION IF EXISTS CITEXT CASCADE;

DROP FUNCTION IF EXISTS thread_insert();
CREATE SCHEMA IF NOT EXISTS tp_forum;

CREATE EXTENSION IF NOT EXISTS CITEXT;
--
-- USERS
--
CREATE TABLE tp_forum.users
(
    nickname CITEXT NOT NULL,
    email    CITEXT NOT NULL,

    about    TEXT DEFAULT NULL,
    fullname TEXT   NOT NULL
);

-- CREATE INDEX users_covering_index
--   ON users (nickname, email, about, fullname);
--
CREATE UNIQUE INDEX users_nickname_index
    ON tp_forum.users (nickname);
--
CREATE UNIQUE INDEX users_email_index
    ON tp_forum.users (email);
--
-- CREATE INDEX ON users (nickname, email);

--
-- FORUM
--
CREATE TABLE tp_forum.forum
(
    id      SERIAL PRIMARY KEY,
    slug    CITEXT                                      NOT NULL,
    title   TEXT                                        NOT NULL,

    author  citext references tp_forum.users (nickname) NOT NULL,
    threads INTEGER                                     NOT NULL DEFAULT 0,
    posts   INTEGER                                     NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX forum_slug_index
    ON tp_forum.forum (slug);
--
-- CREATE INDEX forum_slug_id_index
--     ON forum (slug, id);
--
-- CREATE INDEX on forum (slug, id, title, author, threads, posts);

CREATE TABLE tp_forum.forum_user
(
    nickname   citext COLLATE "POSIX" references tp_forum.users (nickname),
    forum_slug citext references tp_forum.forum (slug),
    CONSTRAINT unique_forum_user UNIQUE (forum_slug, nickname)
);

-- CREATE UNIQUE INDEX forum_users_forum_id_nickname_index2
--     ON forum_users (forumId, lower(nickname));
--
-- CREATE INDEX forum_users_covering_index2
--     ON forum_users (forumId, lower(nickname), nickname, email, about, fullname);



--
-- THREAD
--


CREATE TABLE tp_forum.thread
(
    id      SERIAL PRIMARY KEY                          NOT NULL,
    slug    CITEXT                                               DEFAULT NULL,

    title   TEXT                                        NOT NULL,
    message TEXT                                        NOT NULL,

    forum   CITEXT REFERENCES tp_forum.forum (slug)     NOT NULL,

    author  CITEXT REFERENCES tp_forum.users (nickname) NOT NULL,

    created TIMESTAMPTZ,

    votes   INTEGER                                     NOT NULL DEFAULT 0
);

CREATE FUNCTION thread_insert()
    RETURNS TRIGGER AS
$BODY$
BEGIN
    UPDATE tp_forum.forum
    SET threads = threads + 1
    WHERE lower(slug) = lower(NEW.forum);
    RETURN NULL;
END;
$BODY$
    LANGUAGE plpgsql;
--
CREATE TRIGGER on_thread_insert
    AFTER INSERT
    ON tp_forum.thread
    FOR EACH ROW
EXECUTE PROCEDURE thread_insert();

CREATE UNIQUE INDEX thread_slug_index
    ON tp_forum.thread (slug);
--
-- CREATE INDEX thread_slug_id_index
--     ON thread (slug, id);
--
-- CREATE INDEX thread_forum_id_created_index
--     ON thread (forum_id, created);
--
-- CREATE INDEX thread_forum_id_created_index2
--     ON thread (forum_id, created DESC);
--
CREATE UNIQUE INDEX thread_id_forum_slug_index
    ON tp_forum.thread (id, forum);
--
CREATE UNIQUE INDEX thread_slug_forum_slug_index
    ON tp_forum.thread (slug, forum);
--
-- CREATE UNIQUE INDEX thread_covering_index
--     ON thread (forum_id, created, id, slug, title, message, forum_slug, user_nick, created, votes_count);

--
-- POST
--
CREATE TABLE tp_forum.post
(
    id        SERIAL primary key,

    author    citext REFERENCES tp_forum.users (nickname) NOT NULL,

    message   TEXT                                        NOT NULL,
    created   TIMESTAMPTZ,

    thread    INTEGER REFERENCES tp_forum.thread (id)     NOT NULL,

    parent    INTEGER references tp_forum.post (id)                DEFAULT 0,
    is_edited BOOLEAN                                     NOT NULL DEFAULT FALSE
);

CREATE FUNCTION post_insert()
    RETURNS TRIGGER AS
$BODY$
BEGIN
    UPDATE tp_forum.forum
    SET posts = posts + 1
    WHERE lower(slug) = lower((SELECT forum FROM tp_forum.thread WHERE id = NEW.thread));
    RETURN NULL;
END;
$BODY$
    LANGUAGE plpgsql;
--
CREATE TRIGGER on_post_insert
    AFTER INSERT
    ON tp_forum.post
    FOR EACH ROW
EXECUTE PROCEDURE post_insert();

-- CREATE INDEX posts_thread_id_index
--     ON post (thread_id, id);
--
-- CREATE INDEX posts_thread_id_index2
--     ON post (thread_id);
--
-- CREATE INDEX posts_thread_id_parents_index
--     ON post (thread_id, parents);
--
-- CREATE INDEX ON post (thread_id, id, parent, main_parent)
--     WHERE parent = 0;
--
-- CREATE INDEX parent_tree_3_1
--     ON post (main_parent, parents DESC, id);
--
-- CREATE INDEX parent_tree_4
--     ON post (id, main_parent);
--
--
-- VOTE
--
CREATE TABLE tp_forum.vote
(
    user_nickname citext references tp_forum.users (nickname) NOT NULL,
    thread        INTEGER REFERENCES tp_forum.thread (id)     NOT NULL,

    vote_val      INTEGER,
    prev_vote     INTEGER DEFAULT 0
);

CREATE UNIQUE INDEX user_thread_unique_index
    ON tp_forum.vote (user_nickname, thread);


TRUNCATE TABLE tp_forum.users, tp_forum.forum, tp_forum.thread, tp_forum.post, tp_forum.vote CASCADE;

