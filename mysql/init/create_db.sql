DROP DATABASE IF EXISTS linq;
CREATE DATABASE linq;
USE linq;

CREATE TABLE IF NOT EXISTS `users` (
    `id` varchar(36) PRIMARY KEY NOT NULL,
    `name` varchar(36) NOT NULL UNIQUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `entrys` (
    `id` varchar(36) PRIMARY KEY NOT NULL,
    `url` text NOT NULL UNIQUE,
    `title` text NOT NULL,
    `caption` text,
    `thumbnail` text,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `bookmarks` (
    `id` int(11) AUTO_INCREMENT PRIMARY KEY NOT NULL,
    `user_id` varchar(36) NOT NULL,
    `entry_id` varchar(36) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY user_bookmark(`user_id`) REFERENCES `users`(`id`),
    FOREIGN KEY entry_bookmark(`entry_id`) REFERENCES `entrys`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tags` (
    `id` int(11) AUTO_INCREMENT PRIMARY KEY NOT NULL,
    `tag` varchar(32) NOT NULL,
    `entry_id` varchar(36) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY entry_tag(`entry_id`) REFERENCES `entrys`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `favorites` (
    `id` int(11) AUTO_INCREMENT PRIMARY KEY NOT NULL,
    `user_id` varchar(36) NOT NULL,
    `entry_id` varchar(36) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY user_favorite(`user_id`) REFERENCES `users`(`id`),
    FOREIGN KEY entry_favorite(`entry_id`) REFERENCES `entrys`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;