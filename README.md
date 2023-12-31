# Personal Blog


>Blog and CMS written in Go


[![Production](https://github.com/kevingil/blog/actions/workflows/update-production.yaml/badge.svg?branch=main)](https://github.com/kevingil/blog/actions/workflows/update-production.yaml)


Resume blog for Go developers


### Features

Markdown article editor

Project showcase/alt feed

HTMX frontend


### Database

I'm using Planetscale, but you can also run MySQL with Docker

```sql
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

```



## Setup 

Edit /views/*.gohtml with your own resume / style


#### Download and install
- Go 
- Docker (optional)

#### Download dependencies

gorilla/mux, mysql drivers, gosimple/slug, uuid, etc

`go mod download`


.env example

```sh
PORT=80 
PROD_DSN=your_database_key #MySQL connection string
S3_TOKEN=your_token # S3 is a work in progress
S3_ACCESS_KEY_ID=your_id 
S3_SECRET_ACCESS_KEY=your_key
ACCOUNT_ID=your_id
```

#### First run
You can register at /register

Then login at /login


## Screenshots


### Dashboard

![main](https://cdn.kevingil.com/dashboard-main.png)


### Editor

![editor](https://cdn.kevingil.com/dashboard-editor.png)



