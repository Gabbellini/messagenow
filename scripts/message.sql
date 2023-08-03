CREATE TABLE IF NOT EXISTS message
(
    id           INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    id_user      INT NOT NULL,
    id_sender    INT NOT NULL,
    id_addressee INT NOT NULL,
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (id_user) REFERENCES user (id),
    FOREIGN KEY (id_sender) REFERENCES user (id),
    FOREIGN KEY (id_addressee) REFERENCES user (id)
);

CREATE TABLE message_text
(
    id         INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    id_message INT     NOT NULL,
    text       TEXT    NOT NULL,
    FOREIGN KEY (id_message) REFERENCES message(id)
);