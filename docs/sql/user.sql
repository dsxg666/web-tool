DROP TABLE users;

CREATE TABLE users
(
    `id`              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `username`        VARCHAR(10)                 NOT NULL,
    `password`        VARCHAR(60)                 NOT NULL,
    `email`           VARCHAR(20)                 NOT NULL UNIQUE,
    `avatar`          VARCHAR(20) DEFAULT '0.png' NOT NULL,
    `email_update_at` TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    `created_at`      TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      TIMESTAMP   DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;