CREATE TABLE `users`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(50) NOT NULL,
    `coin` BIGINT NOT NULL
);
CREATE TABLE `user_items`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `item_id` BIGINT UNSIGNED NOT NULL,
    `count` INT NOT NULL
);
CREATE TABLE `items`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `rare` CHAR(2) NOT NULL,
    `weight` INT NOT NULL
);
ALTER TABLE
    `user_items` ADD CONSTRAINT `user_items_item_id_foreign` FOREIGN KEY(`item_id`) REFERENCES `items`(`id`);
ALTER TABLE
    `user_items` ADD CONSTRAINT `user_items_user_id_foreign` FOREIGN KEY(`user_id`) REFERENCES `users`(`id`);

INSERT INTO users (`id`, `name`, `coin`)
VALUES(NULL, "user1", 100);

INSERT INTO items (`id`, `name`, `rare`, `weight`) VALUES 
(NULL, "item1", "N", 15),
(NULL, "item2", "N", 15),
(NULL, "item3", "N", 15),
(NULL, "item4", "N", 15),
(NULL, "item5", "N", 15),
(NULL, "item6", "R", 6),
(NULL, "item7", "R", 6),
(NULL, "item8", "R", 6),
(NULL, "item9", "R", 6),
(NULL, "item10", "SR", 1);