CREATE TABLE IF NOT EXISTS notification_queue (
    id INT AUTO_INCREMENT NOT NULL,
    message_queue_id INT NOT NULL,
    chat_id INT NOT NULL,
    is_sent TINYINT(1) DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    scheduled_for TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY(id),
    INDEX scheduled_for_idx (scheduled_for),
    CONSTRAINT message_queue_id_fk_1 FOREIGN KEY (message_queue_id) REFERENCES message_queue(id)
);

CREATE INDEX queue_id_idx ON message_queue(queue_id);
CREATE INDEX username_idx ON user(username);
CREATE INDEX external_chat_id_idx ON chat(external_chat_id);

ALTER TABLE message_queue MODIFY `message` TEXT NOT NULL;