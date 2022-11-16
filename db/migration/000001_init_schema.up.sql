CREATE TABLE `teacher` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `full_name` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `is_active` boolean NOT NULL DEFAULT true,
  `created_at` datetime NOT NULL DEFAULT (now()),
  `updated_at` datetime NOT NULL DEFAULT (now()),
  `deleted_at` datetime NOT NULL DEFAULT "0001-01-01 00:00:00"
);

CREATE TABLE `student` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `full_name` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `is_suspended` boolean NOT NULL DEFAULT false,
  `is_active` boolean NOT NULL DEFAULT true,
  `suspended_at` datetime NOT NULL DEFAULT "0001-01-01 00:00:00",
  `created_at` datetime NOT NULL DEFAULT (now()),
  `updated_at` datetime NOT NULL DEFAULT (now()),
  `deleted_at` datetime NOT NULL DEFAULT "0001-01-01 00:00:00"
);

CREATE TABLE `register` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `teacher_id` bigint NOT NULL,
  `student_id` bigint NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (now())
);

CREATE INDEX `teacher_index_0` ON `teacher` (`email`);

CREATE INDEX `student_index_1` ON `student` (`email`);

CREATE INDEX `register_index_2` ON `register` (`teacher_id`);

CREATE INDEX `register_index_3` ON `register` (`student_id`);

CREATE UNIQUE INDEX `register_index_4` ON `register` (`teacher_id`, `student_id`);

ALTER TABLE `register` ADD FOREIGN KEY (`teacher_id`) REFERENCES `teacher` (`id`);

ALTER TABLE `register` ADD FOREIGN KEY (`student_id`) REFERENCES `student` (`id`);
