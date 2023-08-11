CREATE TABLE IF NOT EXISTS room
(
    id          INT      NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_room
(
    id_room     INT      NOT NULL,
    id_user     INT      NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (id_user) REFERENCES user (id),
    FOREIGN KEY (id_room) REFERENCES room (id),
    UNIQUE (id_room, id_user)
);

CREATE TABLE IF NOT EXISTS message
(
    id          INT      NOT NULL AUTO_INCREMENT PRIMARY KEY,
    id_room     INT      NOT NULL,
    id_user     INT      NOT NULL,
    text        TEXT     NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (id_room) REFERENCES room (id),
    FOREIGN KEY (id_user) REFERENCES user (id)
);

INSERT INTO room (id)
VALUES (id);
INSERT INTO user_room (id_room, id_user)
VALUES (1, 1);
INSERT INTO user_room (id_room, id_user)
VALUES (1, 2);

INSERT INTO message (id_room, id_user, text)
VALUES (1, 1, 'Eai greg');

INSERT INTO message (id_room, id_user, text)
VALUES (1, 2, 'Eai bellini');