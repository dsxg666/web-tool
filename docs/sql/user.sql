DROP TABLE users;

CREATE TABLE users
(
    `id`              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `username`        VARCHAR(10)                 NOT NULL,
    `password`        VARCHAR(60)                 NOT NULL,
    `email`           VARCHAR(20)                 NOT NULL UNIQUE,
    `avatar`          VARCHAR(20) DEFAULT '0.png' NOT NULL,
    `path`            VARCHAR(20)                 NOT NULL,
    `email_update_at` TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    `created_at`      TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      TIMESTAMP   DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

INSERT INTO tool.users (id, username, password, email, avatar, path, email_update_at, created_at, updated_at) VALUES (1, 'weiguo', '$2a$10$/jKs8SpRQwz25udtS9XbruasUQQptKojRClN3dEudcrgf/qea1Zea', '2637046983@qq.com', '1.png', 'dsxg666', '2024-09-08 22:41:36', '2024-09-08 22:41:36', '2024-11-25 14:47:28');
INSERT INTO tool.users (id, username, password, email, avatar, path, email_update_at, created_at, updated_at) VALUES (2, '风飘', '$2a$10$AJZHVIH7aAP.oO7Mzc9qWuNudcFChGinimnmG1m1mi8AnG2ubPZXi', '1901728868@qq.com', '2.png', '123', '2024-09-09 09:51:33', '2024-09-09 09:51:33', '2024-09-12 19:48:15');
