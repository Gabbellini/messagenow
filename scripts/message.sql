CREATE TABLE IF NOT EXISTS message
(
    id           INT     NOT NULL AUTO_INCREMENT PRIMARY KEY,
    type         TINYINT NOT NULL REFERENCES message_type_text (id),
    id_sender    INT     NOT NULL REFERENCES user (id),
    id_addressee INT     NOT NULL REFERENCES user (id)
);

CREATE TABLE message_type_text
(
    id         TINYINT NOT NULL PRIMARY KEY,
    id_message INT     NOT NULL REFERENCES message (id),
    text       TEXT    NOT NULL
);