CREATE TABLE IF NOT EXISTS user_status
(
    id   INT         NOT NULL AUTO_INCREMENT PRIMARY KEY,
    code TINYINT     NOT NULL DEFAULT 0,
    text VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS user
(
    id          INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    image       TEXT,
    email       varchar(100) NOT NULL,
    password    TEXT         NOT NULL,
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    status      INT          NOT NULL DEFAULT 1,
    FOREIGN KEY (status) REFERENCES user_status (id),
    UNIQUE (email)
);