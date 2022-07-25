CREATE TABLE IF NOT EXISTS trips(
  id CHAR(36) NOT NULL PRIMARY KEY,
  year SMALLINT UNSIGNED NOT NULL,
  month TINYINT UNSIGNED NOT NULL,
  day TINYINT UNSIGNED NOT NULL,
  unit VARCHAR(255) NOT NULL,
  INDEX idx_trips_unit_year (unit, year),
  INDEX idx_trips_unit_year_month (unit, year, month),
  INDEX idx_trips_unit_year_month_day (unit, year, month, day),
  CONSTRAINT fk_trips_unit FOREIGN KEY (unit) REFERENCES users(unit)
);
