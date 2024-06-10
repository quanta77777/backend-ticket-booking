
CREATE TABLE IF NOT EXISTS `user` (
  `user_id` INT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(255),
  `email` VARCHAR(255),
  `password` VARCHAR(255),
  `role` VARCHAR(255) DEFAULT 'user',
   `image_url` varchar(255),
  `image_id` varchar(255),
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `cinema` (
  `cinema_id` INT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(255),
     `image_url` varchar(255),
  `image_id` varchar(255),
   `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `branch` (
  `branch_id` INT PRIMARY KEY AUTO_INCREMENT,
  `cinema_id` INT,
  `name` VARCHAR(255),
  `address` VARCHAR(255),
      `image_url` varchar(255),
  `image_id` varchar(255),
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  FOREIGN KEY (`cinema_id`) REFERENCES `cinema` (`cinema_id`) ON DELETE CASCADE

);

CREATE TABLE IF NOT EXISTS `movie` (
    `movie_id` INT AUTO_INCREMENT PRIMARY KEY,
    `title` VARCHAR(255) NOT NULL,
     `director` VARCHAR(255),
     `genre` VARCHAR(255),
  `duration` INT,
  `image_url` varchar(255),
  `image_id` varchar(255),
    `description` TEXT,
    `release_date` DATE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS `showtime` (
  `showtime_id` INT PRIMARY KEY AUTO_INCREMENT,
  `cinema_id` INT,
  `branch_id` INT,
  `theater_id` INT,
  `movie_id` INT,
  `start_time` DATETIME,
  `end_time` DATETIME,
   `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`branch_id`) REFERENCES `branch` (`branch_id`) ON DELETE CASCADE,
  FOREIGN KEY (`movie_id`) REFERENCES `movie` (`movie_id`) ON DELETE CASCADE,
  FOREIGN KEY (`theater_id`) REFERENCES `theater` (`theater_id`) ON DELETE CASCADE,
  FOREIGN KEY (`cinema_id`) REFERENCES `cinema` (`cinema_id`) ON DELETE CASCADE,
);


CREATE TABLE IF NOT EXISTS `theater` (
  `theater_id` INT PRIMARY KEY AUTO_INCREMENT,
  `branch_id` INT,
  `name` VARCHAR(255),
  FOREIGN KEY (`branch_id`) REFERENCES `branch` (`branch_id`) ON DELETE CASCADE,
);


CREATE TABLE IF NOT EXISTS `seat` (
  `seat_id` INT PRIMARY KEY AUTO_INCREMENT,
  `theater_id` INT,
  `seat_number` INT,
  `seat_type` ENUM('regular','vip','couple'),
  FOREIGN KEY (`theater_id`) REFERENCES `theater` (`theater_id`) ON DELETE CASCADE,
);

CREATE TABLE IF NOT EXISTS `product` (
    `product_id` INT AUTO_INCREMENT PRIMARY KEY,
    `branch_id` INT NOT NULL PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT NOT NULL,
    `price` INT NOT NULL,
    `image_url` varchar(255),
    `image_id` varchar(255),
     FOREIGN KEY (`branch_id`) REFERENCES `branch` (`branch_id`) ON DELETE CASCADE,
);

CREATE TABLE IF NOT EXISTS `price` (

    `price_id` INT AUTO_INCREMENT PRIMARY KEY,
    `showtime_id` INT NOT NULL PRIMARY KEY,
     `seat_type` ENUM('regular','vip','couple'),
    `price` INT NOT NULL,
     FOREIGN KEY (`showtime_id`) REFERENCES `showtime` (`showtime_id`) ON DELETE CASCADE,
);

CREATE TABLE  IF NOT EXISTS `ticket` (
    `ticket_id` INT AUTO_INCREMENT PRIMARY KEY,
    `user_id` INT NOT NULL,
    `showtime_id` INT NOT NULL,
    `seat_id` INT NOT NULL,
    `movie_id` INT NOT NULL,
    `price` INT NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE,
    FOREIGN KEY (`showtime_id`) REFERENCES `showtime` (`showtime_id`) ON DELETE CASCADE,
    FOREIGN KEY (`seat_id`) REFERENCES `seat` (`seat_id`) ON DELETE CASCADE,
    FOREIGN KEY (`movie_id`) REFERENCES `movie` (`movie_id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `ticket_product` (
    `ticket_product_id` INT PRIMARY KEY AUTO_INCREMENT,
    `ticket_id` INT NOT NULL PRIMARY KEY,
    `product_id` INT NOT NULL PRIMARY KEY,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`ticket_id`) REFERENCES `ticket` (`ticket_id`) ON DELETE CASCADE,
    FOREIGN KEY (`product_id`) REFERENCES `product` (`product_id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `ticket_seat` (
    `ticket_seat_id` INT PRIMARY KEY AUTO_INCREMENT,
    `seat_id` INT NOT NULL PRIMARY KEY,
    FOREIGN KEY (`seat_id`) REFERENCES `seat` (`seat_id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `review` (
    `review_id` INT AUTO_INCREMENT PRIMARY KEY,
    `user_id` INT NOT NULL,
    `movie_id` INT NOT NULL,
    `rating` INT CHECK (`rating` >= 1 AND `rating` <= 10),
    `comment` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(`user_id`, `movie_id`),
    FOREIGN KEY (`user_id`) REFERENCES `user`(`user_id`) ON DELETE CASCADE,
    FOREIGN KEY (`movie_id`) REFERENCES `movie`(`movie_id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `seat_reservation` (
    `reservation_id` INT AUTO_INCREMENT PRIMARY KEY,
    `seat_id` INT NOT NULL,
    `showtime_id` INT NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`seat_id`) REFERENCES `seat`(`seat_id`) ON DELETE CASCADE,
    FOREIGN KEY (`showtime_id`) REFERENCES `showtime`(`showtime_id`) ON DELETE CASCADE
);




