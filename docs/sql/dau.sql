DROP TABLE `dau`;

CREATE TABLE `dau`
(
    `id`          INT AUTO_INCREMENT PRIMARY KEY,
    `user_id`   VARCHAR(10),
    `user_ip` VARCHAR(100),
    `created_at`  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

INSERT INTO `dau` (`user_id`,`user_ip`) VALUES ('1', '1');

SELECT COUNT(DISTINCT user_id) AS daily_active_users
FROM dau
WHERE DATE(created_at) = CURDATE();

SELECT COUNT(DISTINCT `user_id`) AS dau_num FROM dau WHERE DATE(`created_at`) = DATE_SUB(CURDATE(), INTERVAL 1 DAY);

SELECT DATE(created_at) AS date, COUNT(DISTINCT `user_id`) AS active_users FROM `dau` WHERE `created_at` >= NOW() - INTERVAL 7 DAY GROUP BY DATE(`created_at`) ORDER BY DATE(`created_at`) DESC;