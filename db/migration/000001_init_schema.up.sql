CREATE TABLE `users`(
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY ,
    `username` varchar(20) NOT NULL ,
    `mobile` varchar(20) UNIQUE NOT NULL ,
    `email` varchar(20) UNIQUE NOT NULL ,
    `avatar` varchar(20) NOT NULL ,
    `hashed_password` varchar(20) NOT NULL ,
    `is_email_verified` bool NOT NULL DEFAULT false,
    `password_changed_at` TIMESTAMP NOT NULL ,
    `created_at` TIMESTAMP NOT NULL ,
    `updated_username_at` TIMESTAMP NOT NULL ,
    `updated_avatar_at` TIMESTAMP NOT NULL ,
    `updated_mobile_at` TIMESTAMP NOT NULL ,
    `updated_email_at` TIMESTAMP NOT NULL
) CHARACTER SET = utf8mb4;

CREATE INDEX idx_mobile ON users(mobile);
CREATE INDEX idx_email ON users(email);
CREATE INDEX idx_username ON users(username);

CREATE TABLE `eeg_session`(
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY ,
    `uid` BIGINT UNSIGNED NOT NULL ,
    `begin` TIMESTAMP NOT NULL ,
    `end` TIMESTAMP NOT NULL ,
    `data_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `expired_at` TIMESTAMP NOT NULL ,
    `finish` BOOL DEFAULT FALSE,
    `data_bits_map` VARCHAR(1024) NOT NULL ,
    FOREIGN KEY (uid) REFERENCES users(id)
) CHARACTER SET = utf8mb4;

CREATE INDEX idx_uid ON eeg_session(uid);

CREATE TABLE `eeg_data`(
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY ,
    `session_id` BIGINT UNSIGNED NOT NULL ,
    `offset` INT UNSIGNED NOT NULL ,
    `user_id` BIGINT UNSIGNED NOT NULL ,
    `collected_at` TIMESTAMP NOT NULL ,
    `attention` TINYINT UNSIGNED NOT NULL ,
    `meditation` TINYINT UNSIGNED NOT NULL ,
    `blink_strength` TINYINT UNSIGNED NOT NULL ,
    `alpha1` TINYINT UNSIGNED NOT NULL ,
    `alpha2` TINYINT UNSIGNED NOT NULL ,
    `beta1` TINYINT UNSIGNED NOT NULL ,
    `beta2` TINYINT UNSIGNED NOT NULL ,
    `gamma1` TINYINT UNSIGNED NOT NULL ,
    `gamma2` TINYINT UNSIGNED NOT NULL ,
    `delta` TINYINT UNSIGNED NOT NULL ,
    `theta` TINYINT UNSIGNED NOT NULL ,
    `raw` JSON NOT NULL ,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (session_id) REFERENCES eeg_session(id)
) CHARACTER SET = utf8mb4;

CREATE INDEX idx_session_id ON eeg_data(session_id);

CREATE TABLE `message`(
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY ,
    `send_uid` BIGINT UNSIGNED NOT NULL ,
    `rcv_uid` BIGINT UNSIGNED NOT NULL ,
    `created_at` TIMESTAMP NOT NULL ,
    `has_read` BOOL NOT NULL DEFAULT FALSE,
    `type` TINYINT UNSIGNED NOT NULL ,
    `text` varchar(255) NOT NULL ,
    `subtype` TINYINT UNSIGNED NOT NULL ,
    `attachment` varchar(255) NOT NULL ,
    FOREIGN KEY (send_uid) REFERENCES users(id),
    FOREIGN KEY (rcv_uid) REFERENCES users(id)
)CHARACTER SET = utf8mb4;

CREATE INDEX  idx_send_uid_rcv_id ON message(send_uid,rcv_uid);

CREATE TABLE `subscribers`(
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY ,
    `uid` BIGINT UNSIGNED NOT NULL ,
    `sbs_id` BIGINT UNSIGNED NOT NULL ,
    `subscribe_time` TIMESTAMP NOT NULL ,
    `is_cancel` BOOL DEFAULT FALSE NOT NULL ,
    `cancel_time` TIMESTAMP NOT NULL,
    FOREIGN KEY (uid) REFERENCES users(id),
    FOREIGN KEY (sbs_id) REFERENCES users(id)
) CHARACTER SET = utf8mb4;

CREATE INDEX idx_uid ON subscribers(uid);
CREATE INDEX idx_sbs_id ON subscribers(sbs_id);