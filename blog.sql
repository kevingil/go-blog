-- MySQL Setup

CREATE TABLE `users` (
	`id` int NOT NULL AUTO_INCREMENT,
	`name` varchar(64) NOT NULL,
	`email` varchar(320) NOT NULL,
	`password` varchar(255) NOT NULL,
	`about` varchar(64),
	`contact` text,
	PRIMARY KEY (`id`),
	UNIQUE KEY `email` (`email`)
) ENGINE InnoDB,
  CHARSET utf8mb4,
  COLLATE utf8mb4_0900_ai_ci;

  CREATE TABLE `projects` (
	`id` int NOT NULL AUTO_INCREMENT,
	`title` varchar(255) NOT NULL,
	`description` varchar(255) NOT NULL,
	`url` varchar(255) NOT NULL,
	`image` varchar(255),
	`classes` varchar(255),
	`author` int NOT NULL,
	PRIMARY KEY (`id`),
	KEY `author` (`author`)
) ENGINE InnoDB,
  CHARSET utf8mb4,
  COLLATE utf8mb4_0900_ai_ci;

  CREATE TABLE `articles` (
	`id` int NOT NULL AUTO_INCREMENT,
	`image` varchar(255),
	`slug` varchar(255) NOT NULL,
	`title` varchar(60) NOT NULL,
	`content` text NOT NULL,
	`author` int NOT NULL,
	`created_at` datetime NOT NULL,
	`is_draft` tinyint(1) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`),
	UNIQUE KEY `slug` (`slug`),
	KEY `author` (`author`)
) ENGINE InnoDB,
  CHARSET utf8mb4,
  COLLATE utf8mb4_0900_ai_ci;

  CREATE TABLE `tags` (
	`tag_id` int NOT NULL AUTO_INCREMENT,
	`tag_name` varchar(50),
	PRIMARY KEY (`tag_id`),
	UNIQUE KEY `tag_name` (`tag_name`)
) ENGINE InnoDB,
  CHARSET utf8mb4,
  COLLATE utf8mb4_0900_ai_ci;

  CREATE TABLE `article_tags` (
	`article_id` int NOT NULL,
	`tag_id` int NOT NULL,
	PRIMARY KEY (`article_id`, `tag_id`)
) ENGINE InnoDB,
  CHARSET utf8mb4,
  COLLATE utf8mb4_0900_ai_ci;
