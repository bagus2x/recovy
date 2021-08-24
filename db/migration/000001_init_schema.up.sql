CREATE TABLE App_User(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    picture VARCHAR(512) NULL,
    password VARCHAR(512) NOT NULL,
    created_at INT NOT NULL,
    updated_at INT NOT NULL,
    deleted_at INT NULL
);