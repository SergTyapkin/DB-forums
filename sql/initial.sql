CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE Users (
    nickname CITEXT PRIMARY KEY,
    name     TEXT NOT NULL,
    email    CITEXT UNIQUE,
    about    TEXT
);

CREATE UNLOGGED TABLE Forums (
    slug    CITEXT PRIMARY KEY,
    title   TEXT NOT NULL,
    author  CITEXT NOT null,
    threads INT DEFAULT 0,
    posts   BIGINT DEFAULT 0,

    FOREIGN KEY (author) REFERENCES Users (nickname)
);

CREATE UNLOGGED TABLE Threads (
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

CREATE UNLOGGED TABLE Posts (
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
    --FOREIGN KEY (parent) REFERENCES Posts (id)
);

CREATE UNLOGGED TABLE Votes (
    nickname CITEXT,
    result   BOOLEAN,
    thread   INT,
  
    FOREIGN KEY (nickname) REFERENCES Users (nickname),
    FOREIGN KEY (thread) REFERENCES Threads (id),
    UNIQUE(nickname, thread)
);

CREATE UNLOGGED TABLE forums_to_users (
    slug     CITEXT NOT NULL,
    nickname CITEXT NOT NULL,
    name     TEXT NOT NULL,
    about    TEXT,
    email    CITEXT,
    FOREIGN KEY (nickname) REFERENCES users (nickname),
    FOREIGN KEY (slug) REFERENCES forums (slug),
    UNIQUE (nickname, slug)
);

CREATE OR REPLACE FUNCTION update_paths() RETURNS TRIGGER AS $$
BEGIN
    IF (NEW.parent = 0) THEN
        NEW.paths := array_append(NEW.paths, NEW.id);
    ELSE
        NEW.paths := array_append((SELECT paths FROM posts WHERE id = NEW.parent), NEW.id);
    END IF;
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER trig_update_paths
    BEFORE INSERT ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE update_paths();


CREATE INDEX post_id_path1_index ON posts (id, (posts.paths[1]));
CREATE INDEX post_thread_id_path1_parent_index ON posts (thread, id, (posts.paths[1]), parent);
CREATE INDEX post_thread_path_id_index ON posts (thread, paths, id);
CREATE INDEX post_path1_index ON posts ((posts.paths[1]));
CREATE INDEX post_thread_id_index ON posts (thread, id);
CREATE INDEX post_thread_index ON posts (thread);

CREATE INDEX forum_slug_LOWER_index ON forums (LOWER(forums.Slug));

CREATE INDEX users_email_nickname_LOWER_index ON users (LOWER(users.email), LOWER(users.nickname));
CREATE INDEX users_nickname_index ON users (LOWER(users.nickname));

CREATE UNIQUE INDEX forum_users_unique ON forums_to_users (slug, nickname);
CLUSTER forums_to_users USING forum_users_unique;

CREATE INDEX thread_forum_LOWER_index ON threads (LOWER(forum));
CREATE INDEX thread_slug_index ON threads (LOWER(slug));
CREATE INDEX thread_slug_id_index ON threads (LOWER(forum), created);
CREATE INDEX thread_created_index ON threads (created);

CREATE INDEX vote_nickname ON votes (LOWER(nickname), thread);
