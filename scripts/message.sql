CREATE TABLE IF NOT EXISTS message
(
    id           INT     NOT NULL AUTO_INCREMENT PRIMARY KEY,
    id_sender    INT     NOT NULL,
    id_addressee INT     NOT NULL,
    FOREIGN KEY (id_sender) REFERENCES user (id),
    FOREIGN KEY (id_addressee) REFERENCES user (id)
);

CREATE TABLE message_type_text
(
    id         TINYINT NOT NULL PRIMARY KEY,
    id_message INT     NOT NULL REFERENCES message (id),
    text       TEXT    NOT NULL
);