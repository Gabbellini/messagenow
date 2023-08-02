CREATE TABLE IF NOT EXISTS user_status
(
    id   TINYINT     NOT NULL PRIMARY KEY,
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
    status      TINYINT      NOT NULL DEFAULT 0,
    FOREIGN KEY (status) REFERENCES user_status (id),
    UNIQUE (email)
);