DROP TABLE todolist;

CREATE TABLE todolist
(
    `id`          INT AUTO_INCREMENT PRIMARY KEY,
    `user_id`     INT          NOT NULL,
    `title`       VARCHAR(255) NOT NULL,
    `description` TEXT,
    `priority`    ENUM ('Low', 'Medium', 'High')            DEFAULT 'Low',
    `status`      ENUM ('Todo', 'In Progress', 'Completed') DEFAULT 'Todo',
    `due_date`    DATETIME,
    `created_at`  TIMESTAMP                                 DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  TIMESTAMP                                 DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

INSERT INTO todolist (user_id, title, description, priority, status, due_date)
VALUES (1, '完成项目报告', '撰写并提交项目报告', 'High', 'Todo', '2024-09-30 12:30:00');

SELECT COUNT(`id`) AS num FROM `todolist` WHERE DATE(`created_at`) = DATE_SUB(CURDATE(), INTERVAL 1 DAY);

SELECT DATE(created_at) AS date, COUNT(`id`) AS num FROM `todolist` WHERE `created_at` >= NOW() - INTERVAL 7 DAY GROUP BY DATE(`created_at`) ORDER BY DATE(`created_at`) DESC;