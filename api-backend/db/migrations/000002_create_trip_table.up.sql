CREATE TABLE IF NOT EXISTS trips(
  id CHAR(36) NOT NULL PRIMARY KEY,
  year TINYINT UNSIGNED NOT NULL,
  month TINYINT UNSIGNED NOT NULL,
  day TINYINT UNSIGNED NOT NULL,
  user_id CHAR(36) NOT NULL,
  INDEX idx_trips_userid_year (user_id, year),
  INDEX idx_trips_userid_year_month (user_id, year, month),
  CONSTRAINT fk_trips_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);
