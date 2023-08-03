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

INSERT INTO user_status (code, text)
VALUES (1, 'OK');

INSERT INTO user (name, image, email, password)
VALUES ('GABRIEL DE BRITO BELLINI', NULL, 'gabrielbritobellini@gmail.com',
        '$2a$10$GeqS9DRDEU.UXmEfRtm8vO5fgddOVvvugXbJL0G0kAnBQ/ehtHoFq');

INSERT INTO user (name, image, email, password)
VALUES ('GREGORI SABEL', NULL, 'seriosabel@gmail.com',
        '$2a$10$GeqS9DRDEU.UXmEfRtm8vO5fgddOVvvugXbJL0G0kAnBQ/ehtHoFq');
