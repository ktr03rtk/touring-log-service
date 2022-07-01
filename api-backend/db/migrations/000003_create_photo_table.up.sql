CREATE TABLE IF NOT EXISTS photos(
  id CHAR(36) NOT NULL PRIMARY KEY,
  year SMALLINT UNSIGNED NOT NULL,
  month TINYINT UNSIGNED NOT NULL,
  day TINYINT UNSIGNED NOT NULL,
  lat DECIMAL(8, 6) NOT NULL,
  lon DECIMAL(9, 6) NOT NULL,
  timestamp TIME NOT NULL,
  s3_object_key VARCHAR(256) NOT NULL,
  user_id CHAR(36) NOT NULL,
  INDEX idx_photos_userid_year (user_id, year),
  INDEX idx_photos_userid_year_month (user_id, year, month),
  CONSTRAINT fk_photos_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);
