DROP TABLE IF EXISTS tp_forum.users, tp_forum.forum, tp_forum.thread, tp_forum.post, tp_forum.vote CASCADE;
DROP SCHEMA IF EXISTS tp_forum CASCADE;
DROP EXTENSION IF EXISTS CITEXT CASCADE;

DROP FUNCTION IF EXISTS thread_insert();
CREATE SCHEMA IF NOT EXISTS tp_forum;

CREATE EXTENSION IF NOT EXISTS CITEXT;
--
-- USERS
--
CREATE UNLOGGED TABLE tp_forum.users
(
    nickname CITEXT NOT NULL,
    email    CITEXT NOT NULL,

    about    TEXT DEFAULT NULL,
    fullname TEXT   NOT NULL
);

CREATE INDEX users_covering_index
    ON tp_forum.users (nickname, email, about, fullname);
--
CREATE UNIQUE INDEX users_nickname_index
    ON tp_forum.users (nickname);
--
CREATE UNIQUE INDEX users_email_index
    ON tp_forum.users (email);
--
CREATE INDEX ON tp_forum.users (nickname, email);

--
-- FORUM
--
CREATE UNLOGGED TABLE tp_forum.forum
(
    id      BIGSERIAL PRIMARY KEY,
    slug    CITEXT                                      NOT NULL,
    title   TEXT                                        NOT NULL,

    author  citext references tp_forum.users (nickname) NOT NULL,
    threads INTEGER                                     NOT NULL DEFAULT 0,
    posts   INTEGER                                     NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX forum_slug_index
    ON tp_forum.forum (slug);
--
CREATE INDEX forum_slug_id_index
    ON forum (slug, id);
--
CREATE INDEX on forum (slug, id, title, author, threads, posts);

CREATE UNLOGGED TABLE tp_forum.forum_user
(
    nickname citext COLLATE "POSIX" references tp_forum.users (nickname),
    forum    citext references tp_forum.forum (slug),
    CONSTRAINT unique_forum_user UNIQUE (forum, nickname)
);

--
-- THREAD
--


CREATE UNLOGGED TABLE tp_forum.thread
(
    id      BIGSERIAL PRIMARY KEY                       NOT NULL,
    slug    CITEXT                                               DEFAULT NULL,

    title   TEXT                                        NOT NULL,
    message TEXT                                        NOT NULL,

    forum   CITEXT REFERENCES tp_forum.forum (slug)     NOT NULL,

    author  CITEXT REFERENCES tp_forum.users (nickname) NOT NULL,

    created TIMESTAMPTZ,

    votes   BIGINT                                      NOT NULL DEFAULT 0
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


CREATE OR REPLACE FUNCTION create_forum_user_on_thread_insert() returns trigger as
$$
BEGIN
    INSERT INTO tp_forum.forum_user
        (nickname, forum)
    VALUES (NEW.author, NEW.forum)
    ON CONFLICT DO NOTHING;
    RETURN NULL;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER on_thread_create_forum_user
    AFTER INSERT
    ON tp_forum.thread
    FOR EACH ROW
EXECUTE PROCEDURE create_forum_user_on_thread_insert();


CREATE UNIQUE INDEX thread_slug_index
    ON tp_forum.thread (slug);
--
CREATE INDEX thread_slug_id_index
    ON tp_forum.thread (slug, id);
--
CREATE INDEX thread_forum_id_created_index
    ON tp_forum.thread (forum, created);
--
CREATE INDEX thread_forum_id_created_index2
    ON tp_forum.thread (forum, created DESC);
--
CREATE UNIQUE INDEX thread_id_forum_slug_index
    ON tp_forum.thread (id, forum);
--
CREATE UNIQUE INDEX thread_slug_forum_slug_index
    ON tp_forum.thread (slug, forum);
--
-- CREATE UNIQUE INDEX thread_covering_index
--     ON thread (forum, created, id, slug, title, message, forum, author, created, votes);

--
-- POST
--
CREATE UNLOGGED TABLE tp_forum.post
(
    id        BIGSERIAL primary key,

    author    citext REFERENCES tp_forum.users (nickname),

    message   TEXT,
    created   TIMESTAMPTZ,

    thread    BIGINT REFERENCES tp_forum.thread (id),
    forum     CITEXT REFERENCES tp_forum.forum (slug),

    parent    BIGINT references tp_forum.post (id) DEFAULT 0,
    is_edited BOOLEAN                               DEFAULT FALSE,

    path      BIGINT[]
);

CREATE OR REPLACE FUNCTION change_edited_post() RETURNS trigger as
$change_edited_post$
BEGIN
    IF NEW.message <> OLD.message THEN
        NEW.is_edited = true;
    END IF;

    return NEW;
END;
$change_edited_post$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS change_edited_post ON tp_forum.post;

CREATE TRIGGER change_edited_post
    BEFORE UPDATE
    ON tp_forum.post
    FOR EACH ROW
EXECUTE PROCEDURE change_edited_post();

CREATE OR REPLACE FUNCTION create_path() RETURNS trigger as
$create_path$
BEGIN
    IF NEW.parent = 0 THEN
        NEW.path := (ARRAY [NEW.id]);
        return NEW;
    end if;

    NEW.path := (SELECT array_append(p.path, NEW.id::bigint)
                 from tp_forum.post p
                 where p.id = NEW.parent);
    RETURN NEW;
END;
$create_path$ LANGUAGE plpgsql;

CREATE TRIGGER create_path
    BEFORE INSERT
    ON tp_forum.post
    FOR EACH ROW
EXECUTE PROCEDURE create_path();

CREATE OR REPLACE FUNCTION post_insert()
    RETURNS TRIGGER AS
$BODY$
BEGIN
    UPDATE tp_forum.forum
    SET posts = posts + 1
    WHERE lower(slug) = lower((SELECT forum FROM tp_forum.thread WHERE id = NEW.thread));

    IF NEW.id = 0 THEN
        return null;
    end if;

    INSERT INTO tp_forum.forum_user
        (nickname, forum)
    VALUES (NEW.author,
            (SELECT forum FROM tp_forum.thread WHERE id = NEW.thread))
    ON CONFLICT DO NOTHING;
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
--     ON tp_forum.post (thread, id);

-- CREATE INDEX posts_thread_id_index2
--     ON tp_forum.post (thread);

-- CREATE INDEX posts_thread_path_index
--     ON tp_forum.post (thread, path);
--
-- CREATE INDEX posts_thread_parent_id_index
--     ON tp_forum.post (thread, parent, id);

CREATE INDEX parent_tree_index
    ON tp_forum.post ((path[1]), path DESC, id);

CREATE INDEX parent_tree_index2
    ON tp_forum.post (id, (path[1]));


-- VOTE

CREATE UNLOGGED TABLE tp_forum.vote
(
    user_nickname citext references tp_forum.users (nickname) NOT NULL,
    thread        BIGINT REFERENCES tp_forum.thread (id)     NOT NULL,

    vote_val      INTEGER
);

CREATE UNIQUE INDEX user_thread_unique_index
    ON tp_forum.vote (user_nickname, thread);


TRUNCATE TABLE tp_forum.users, tp_forum.forum, tp_forum.thread, tp_forum.post, tp_forum.vote, tp_forum.forum_user CASCADE;

insert into tp_forum.post (id)
values (0);