-- Thêm chỉ mục cho các bảng
ALTER TABLE `user` ADD INDEX `idx_user_email` (`email`);
ALTER TABLE `branch` ADD INDEX `idx_branch_cinema_chain_id` (`cinema_chain_id`);
ALTER TABLE `movie_showtime` ADD INDEX `idx_movie_showtime_branch_id` (`branch_id`);
ALTER TABLE `movie_showtime` ADD INDEX `idx_movie_showtime_movie_id` (`movie_id`);
ALTER TABLE `theater` ADD INDEX `idx_theater_branch_id` (`branch_id`);
ALTER TABLE `seat` ADD INDEX `idx_seat_theater_id` (`theater_id`);
ALTER TABLE `booking` ADD INDEX `idx_booking_showtime_id` (`showtime_id`);
ALTER TABLE `booking` ADD INDEX `idx_booking_seat_id` (`seat_id`);
ALTER TABLE `booking` ADD INDEX `idx_booking_user_id`
