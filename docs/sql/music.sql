DROP TABLE `songs`;
CREATE TABLE `songs`
(
    `id`     INT AUTO_INCREMENT PRIMARY KEY,
    `title`  VARCHAR(50) NOT NULL,
    `artist` VARCHAR(50) NOT NULL,
    `time` varchar(5) NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

INSERT INTO `songs` (`title`, `artist`, `time`)
VALUES ('不要说话', '陈奕迅', '04:45'),
       ('天外来物', '薛之谦', '04:17'),
       ('Free Loop', 'Daniel Powter', '03:48'),
       ('Dear John', '比莉', '05:11'),
        ('我好像在哪见过你', '薛之谦', '04:39'),
        ('Dancing With Your Ghost', 'Sasha Alex Sloan', '03:17');

DROP TABLE `song_favorites`;
CREATE TABLE `song_favorites`
(
    `id`      INT AUTO_INCREMENT PRIMARY KEY,
    `user_id` INT NOT NULL,
    `song_id` INT NOT NULL
);