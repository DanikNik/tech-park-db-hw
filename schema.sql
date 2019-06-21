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
    id       SERIAL PRIMARY KEY NOT NULL,

    nickname CITEXT        NOT NULL,
    email    CITEXT        NOT NULL,

    about    TEXT DEFAULT NULL,
    fullname TEXT               NOT NULL
);

-- CREATE INDEX users_covering_index
--   ON users (nickname, email, about, fullname);
--
CREATE UNIQUE INDEX users_nickname_index
  ON tp_forum.users (nickname);
--
CREATE UNIQUE INDEX users_email_index
  ON tp_forum.users(email);
--
-- CREATE INDEX ON users (nickname, email);

--
-- FORUM
--
CREATE TABLE tp_forum.forum
(
    id      SERIAL PRIMARY KEY,
    slug    CITEXT                            NOT NULL,
    title   TEXT                                   NOT NULL,

    author  INTEGER references tp_forum.users (id) NOT NULL,
    threads INTEGER                                NOT NULL DEFAULT 0,
    posts   BIGINT                                 NOT NULL DEFAULT 0
);

-- CREATE UNIQUE INDEX forum_slug_index
--     ON forum (slug);
--
-- CREATE INDEX forum_slug_id_index
--     ON forum (slug, id);
--
-- CREATE INDEX on forum (slug, id, title, moderator, threads, posts);
--
-- THREAD
--
CREATE TABLE tp_forum.thread
(
    id      SERIAL PRIMARY KEY                     NOT NULL,
    slug    CITEXT DEFAULT NULL,

    title   TEXT                                   NOT NULL,
    message TEXT                                   NOT NULL,

    forum   INTEGER REFERENCES tp_forum.forum (id) NOT NULL,

    author  INTEGER REFERENCES tp_forum.users (id) NOT NULL,

    created TIMESTAMPTZ
);

-- CREATE FUNCTION thread_insert()
--     RETURNS TRIGGER AS
-- $BODY$
-- BEGIN
--     UPDATE forum
--     SET threads = forum.threads + 1
--     WHERE slug = NEW.forum_slug;
--     RETURN NULL;
-- END;
-- $BODY$
--     LANGUAGE plpgsql;

-- CREATE TRIGGER on_thread_insert
--     AFTER INSERT
--     ON thread
--     FOR EACH ROW
-- EXECUTE PROCEDURE thread_insert();

-- CREATE UNIQUE INDEX thread_slug_index
--     ON thread (slug);
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
-- CREATE UNIQUE INDEX thread_id_forum_slug_index
--     ON thread (id, forum_slug);
--
-- CREATE UNIQUE INDEX thread_slug_forum_slug_index
--     ON thread (slug, forum_slug);
--
-- CREATE UNIQUE INDEX thread_covering_index
--     ON thread (forum_id, created, id, slug, title, message, forum_slug, user_nick, created, votes_count);

--
-- POST
--
CREATE TABLE tp_forum.post
(
    id        SERIAL primary key,

    user_nick TEXT                                    NOT NULL,

    message   TEXT                                    NOT NULL,
    created   TIMESTAMPTZ,

    thread    INTEGER REFERENCES tp_forum.thread (id) NOT NULL,

    parent    INTEGER references tp_forum.post (id)            DEFAULT 0,
    is_edited BOOLEAN                                 NOT NULL DEFAULT FALSE
);



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
    id        SERIAL,

    user_id   INTEGER references tp_forum.users (id)  NOT NULL,
    thread_id INTEGER REFERENCES tp_forum.thread (id) NOT NULL,

    vote_val  INTEGER,
    prev_vote INTEGER DEFAULT 0
--     CONSTRAINT unique_user_and_thread UNIQUE (user_id, thread_id)
);


-- CREATE UNIQUE INDEX forum_users_forum_id_nickname_index2
--     ON forum_users (forumId, lower(nickname));
--
-- CREATE INDEX forum_users_covering_index2
--     ON forum_users (forumId, lower(nickname), nickname, email, about, fullname);

-- SELECT * FROM tp_forum.users where nickname ILIKE 'шФагф';

INSERT INTO tp_forum.users (nickname, email, about, fullname) VALUES ('qwerty', 'qwerty', 'qwerty', 'qwerty');

TRUNCATE TABLE  tp_forum.users, tp_forum.forum, tp_forum.thread, tp_forum.post, tp_forum.vote CASCADE;