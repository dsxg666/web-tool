DROP TABLE codes;

CREATE TABLE codes
(
    `id`         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `email`      VARCHAR(20) NOT NULL,
    `code`       VARCHAR(6)  NOT NULL,
    `type`       VARCHAR(1) DEFAULT '0' COMMENT '0 is verification code | 1 is register code',
    `created_at` TIMESTAMP  DEFAULT CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;