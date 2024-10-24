DROP TABLE `messages`;
CREATE TABLE `messages`
(
    `id`          INT AUTO_INCREMENT PRIMARY KEY,
    `sender_id`   VARCHAR(10),
    `receiver_id` VARCHAR(10),
    `message`     TEXT,
    `created_at`  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
SELECT *
FROM (SELECT *
      FROM `messages`
      WHERE `sender_id` = 1 and `receiver_id` = 2
         or `sender_id` = 2 and `receiver_id` = 1
      ORDER BY `created_at` DESC
      LIMIT 1000) AS recent_messages
ORDER BY `created_at`;

SELECT *
FROM (SELECT * FROM `group_messages` WHERE `group_id` = 1 ORDER BY `created_at` DESC LIMIT 1000) AS recent_messages
ORDER BY `created_at`;

SELECT COUNT(`id`) AS message_num
FROM `messages`
WHERE DATE(`created_at`) = CURDATE();

SELECT DATE(created_at) AS date, COUNT(`id`) AS message_num
FROM `messages`
WHERE `created_at` >= NOW() - INTERVAL 7 DAY
GROUP BY DATE(`created_at`)
ORDER BY DATE(`created_at`) DESC;

DROP TABLE `friends`;
CREATE TABLE `friends`
(
    `id`         INT AUTO_INCREMENT PRIMARY KEY,
    `self_id`    varchar(10),
    `friend_id`  varchar(10),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

DROP TABLE `groups`;
CREATE TABLE `groups`
(
    `id`         int auto_increment primary key,
    `name`       varchar(10),
    `avatar`     varchar(20),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

INSERT INTO `groups` (`name`, `avatar`)
VALUES ('WorldGroup', '1.jpg');

DROP TABLE group_members;
CREATE TABLE `group_members`
(
    `id`       int auto_increment primary key,
    `group_id` varchar(10),
    `user_id`  varchar(10),
    `status`   varchar(2) -- 0 表示是普通成员，1 表示是群主
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


INSERT INTO `group_members` (`group_id`, `user_id`, `status`)
VALUES ('1', '', '');

DROP TABLE `group_messages`;
CREATE TABLE `group_messages`
(
    `id`         INT AUTO_INCREMENT PRIMARY KEY,
    `sender_id`  VARCHAR(10),
    `group_id`   VARCHAR(10),
    `message`    TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

SELECT COUNT(`id`) AS group_message_num
FROM `group_messages`
WHERE DATE(`created_at`) = CURDATE();

SELECT DATE(created_at) AS date, COUNT(`id`) AS message_num
FROM `group_messages`
WHERE `created_at` >= NOW() - INTERVAL 7 DAY
GROUP BY DATE(`created_at`)
ORDER BY DATE(`created_at`) DESC;


DROP TABLE `requests`;
CREATE TABLE `requests`
(
    `id`          INT AUTO_INCREMENT PRIMARY KEY,
    `user_id`     VARCHAR(10) DEFAULT '',
    `group_id`    VARCHAR(10) DEFAULT '',
    `receiver_id` VARCHAR(10) DEFAULT '',
    `remark`      VARCHAR(50) DEFAULT '',
    `type`        VARCHAR(10) DEFAULT '',  -- 0 表示加好友请求，1 表示进群请求
    `finish`      VARCHAR(10) DEFAULT '0', -- 0 表示未完成，1 表示已完成
    `result`      VARCHAR(10) DEFAULT '',  -- 0 表示拒接，1 表示接受
    `created_at`  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
