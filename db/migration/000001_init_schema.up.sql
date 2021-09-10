CREATE TABLE App_User (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    picture VARCHAR(512) NOT NULL DEFAULT '',
    password VARCHAR(512) NOT NULL,
    created_at INT NOT NULL,
    updated_at INT NOT NULL
);

CREATE TABLE Podcast (
    id SERIAL PRIMARY KEY,
    author_id INT NOT NULL REFERENCES App_User(id) ON DELETE CASCADE,
    picture VARCHAR(512) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL,
    description VARCHAR(512) NOT NULL DEFAULT '',
    file VARCHAR(512) NOT NULL DEFAULT '',
    created_at INT NOT NULL,
    updated_at INT NOT NULL
);

CREATE TABLE Starred_Podcast (
    id SERIAL PRIMARY KEY,
    podcast_id INT NOT NULL REFERENCES Podcast(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES App_User(id) ON DELETE CASCADE,
    created_at INT NOT NULL,
    UNIQUE(podcast_id, user_id)
);