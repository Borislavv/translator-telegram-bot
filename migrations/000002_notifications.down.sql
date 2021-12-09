ALTER TABLE notification_queue DROP FOREIGN KEY message_queue_id_fk_1;

DROP INDEX external_chat_id_idx ON chat;
DROP INDEX username_idx ON user;
DROP INDEX queue_id_idx ON message_queue;
DROP INDEX scheduled_for_idx ON notification_queue;

DROP TABLE IF EXISTS notification_queue;

ALTER TABLE message_queue MODIFY `message` VARCHAR(1024) NOT NULL;