---------------- Types
CREATE EXTENSION IF NOT EXISTS citext;

---------------- Tables
CREATE UNLOGGED TABLE IF NOT EXISTS Users (
    nickname CITEXT PRIMARY KEY,
    name     TEXT NOT NULL,
    email    CITEXT UNIQUE,
    about    TEXT
);

CREATE UNLOGGED TABLE IF NOT EXISTS Forums (
    slug    CITEXT PRIMARY KEY,
    title   TEXT NOT NULL,
    author  CITEXT NOT null,
    threads INT DEFAULT 0,
    posts   BIGINT DEFAULT 0,

    FOREIGN KEY (author) REFERENCES Users (nickname)
);

CREATE UNLOGGED TABLE IF NOT EXISTS Threads (
    id      SERIAL PRIMARY KEY,
    forum   CITEXT NOT NULL,
    author  CITEXT NOT NULL,
    created TIMESTAMP with time zone DEFAULT now(),
    message TEXT NOT NULL,
    title   TEXT NOT NULL,
    votes   INT DEFAULT 0,
    slug    CITEXT UNIQUE,

    FOREIGN KEY (author) REFERENCES Users (nickname),
    FOREIGN KEY (forum) REFERENCES Forums (slug)
);

CREATE UNLOGGED TABLE IF NOT EXISTS Posts (
    id      BIGSERIAL PRIMARY KEY,
    author  CITEXT NOT NULL,
    created TIMESTAMP with time zone default now(),
    forum   CITEXT NOT NULL,
    thread  INT,
    edited  BOOLEAN DEFAULT false,
    message TEXT NOT NULL,
    parent  BIGINT DEFAULT 0,
    paths   BIGINT[] DEFAULT ARRAY []::BIGINT[], --информация о дереве "над" постом

    FOREIGN KEY (author) REFERENCES Users (nickname),
    FOREIGN KEY (forum) REFERENCES Forums (slug),
    FOREIGN KEY (thread) REFERENCES Threads (id)
    --FOREIGN KEY (parent) REFERENCES Posts (id) --parent can be 0
);

CREATE UNLOGGED TABLE IF NOT EXISTS Votes (
    nickname CITEXT,
    result   BOOLEAN,
    thread   INT,

    FOREIGN KEY (nickname) REFERENCES Users (nickname),
    FOREIGN KEY (thread) REFERENCES Threads (id),
    UNIQUE(nickname, thread)
);

CREATE UNLOGGED TABLE IF NOT EXISTS forums_to_users (
    slug     CITEXT NOT NULL,
    nickname CITEXT NOT NULL,
    name     TEXT NOT NULL,
    about    TEXT,
    email    CITEXT,
    FOREIGN KEY (nickname) REFERENCES users (nickname),
    FOREIGN KEY (slug) REFERENCES forums (slug),
    UNIQUE (nickname, slug)
);

---------------- Procedures
CREATE OR REPLACE FUNCTION update_paths_in_post() RETURNS TRIGGER AS $$
DECLARE
    parent_thread   INT;
    parent_paths    BIGINT[];
BEGIN
    -- Update paths in `Posts`
    IF (NEW.parent = 0) THEN
        NEW.paths := array_append(NEW.paths, NEW.id);
    ELSE
        SELECT thread, paths FROM Posts WHERE id = NEW.parent INTO parent_thread, parent_paths;
        IF (NOT FOUND) OR parent_thread <> NEW.thread THEN
            RAISE EXCEPTION 'Parent post in another thread' USING ERRCODE = '00228';
        END IF;

        NEW.paths := array_append(parent_paths, NEW.id);
    END IF;

    RETURN NEW;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_before_insert_posts ON posts;
CREATE TRIGGER trig_before_insert_posts
    BEFORE INSERT ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE update_paths_in_post();

---------------
CREATE OR REPLACE FUNCTION update_forumsToUsers() RETURNS TRIGGER AS $$
DECLARE
    author_nickname CITEXT;
    author_name     TEXT;
    author_about    TEXT;
    author_email    CITEXT;
BEGIN
    -- Update forums_to_users
    SELECT nickname, name, about, email
    FROM Users
    WHERE nickname = NEW.author
    INTO author_nickname, author_name, author_about, author_email;

    INSERT INTO forums_to_users (slug, nickname, name, about, email)
    VALUES (NEW.forum, author_nickname, author_name, author_about, author_email)
    ON CONFLICT DO NOTHING;

    RETURN NEW;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_after_insert_posts ON posts;
CREATE TRIGGER trig_after_insert_posts
    AFTER INSERT ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_forumsToUsers();

----------------
CREATE OR REPLACE FUNCTION update_threads_count_in_forum_and_forumsToUsers() RETURNS TRIGGER AS $$
DECLARE
    author_nickname CITEXT;
    author_name     TEXT;
    author_about    TEXT;
    author_email    CITEXT;
BEGIN
    -- Update threads count in `Forums`
    -- UPDATE Forums SET threads = (threads + 1) WHERE LOWER(slug)=LOWER(NEW.forum);

    -- Update forums_to_users
    SELECT nickname, name, about, email
    FROM Users
    WHERE nickname = NEW.author
    INTO author_nickname, author_name, author_about, author_email;

    INSERT INTO forums_to_users (slug, nickname, name, about, email)
    VALUES (NEW.forum, author_nickname, author_name, author_about, author_email)
    ON CONFLICT DO NOTHING;

    RETURN NEW;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_before_insert_threads ON threads;
CREATE TRIGGER trig_before_insert_threads
    BEFORE INSERT ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_threads_count_in_forum_and_forumsToUsers();

---------------- Indexes
CREATE INDEX IF NOT EXISTS post_path1_paths_index ON posts ((paths[1]), paths);
CREATE INDEX IF NOT EXISTS post_id_index ON posts (id);
CREATE INDEX IF NOT EXISTS post_thread_parent_id_index ON posts (thread, parent, id);
CREATE INDEX IF NOT EXISTS post_thread_path1_parent_id_index ON posts (thread, (paths[1]), parent, id);
CREATE INDEX IF NOT EXISTS post_thread_paths_index ON posts (thread, paths);
CREATE INDEX IF NOT EXISTS post_thread_id_created_index ON posts (thread, id, created);

CREATE INDEX IF NOT EXISTS forum_slug_index ON forums (LOWER(slug));

CREATE INDEX IF NOT EXISTS users_email_index ON users (LOWER(email));
CREATE INDEX IF NOT EXISTS users_nickname_index ON users (nickname);

CREATE UNIQUE INDEX IF NOT EXISTS forum_to_users_unique_index ON forums_to_users (LOWER(slug), nickname);
CLUSTER forums_to_users USING forum_to_users_unique_index;

CREATE INDEX IF NOT EXISTS thread_id_index ON threads (id);
CREATE INDEX IF NOT EXISTS thread_slug_index ON threads (LOWER(slug));
CREATE INDEX IF NOT EXISTS thread_forum_created_index ON threads (LOWER(forum), created);

CREATE INDEX IF NOT EXISTS vote_nickname_thread_index ON votes (nickname, thread);