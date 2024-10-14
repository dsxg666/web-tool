DROP TABLE `posts`;
CREATE TABLE `posts`
(
    `id`           INT AUTO_INCREMENT PRIMARY KEY,
    `user_id`      INT         NOT NULL,
    `title`        VARCHAR(20) NOT NULL,
    `content`      TEXT        NOT NULL,
    `category`     VARCHAR(20) NOT NULL,
    `view_count`   INT        DEFAULT 0,
    `is_public`    VARCHAR(1) DEFAULT '0',
    `published_at` TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
    `created_at`   TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
    `updated_at`   TIMESTAMP  DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

INSERT INTO `posts` (`user_id`, `title`, `content`, `category`)
VALUES ('1', '如何学好python', '啦啦啦', 'technology');

select *
from `posts`;

CREATE TABLE `comments`
(
    `id`         INT AUTO_INCREMENT PRIMARY KEY,
    `user_id`    INT  NOT NULL,
    `post_id`    INT  NOT NULL,
    `content`    TEXT NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

DROP TABLE `likes`;
CREATE TABLE `likes`
(
    `id`         INT AUTO_INCREMENT PRIMARY KEY,
    `user_id`    INT NOT NULL,
    `post_id`    INT NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

DROP TABLE `favorites`;
CREATE TABLE `favorites`
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id    INT NOT NULL,
    post_id    INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

SELECT COUNT(f.post_id) AS count FROM favorites f JOIN posts p ON f.post_id = p.id WHERE f.user_id = '1' AND p.is_public = '1';