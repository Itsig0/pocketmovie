CREATE TABLE  IF NOT EXISTS movie (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title VARCHAR(255) NOT NULL,
		original_title VARCHAR(255) NOT NULL,
		imdbid VARCHAR(25) NOT NULL,
		tmdbid INTEGER NOT NULL,
		length INTEGER NOT NULL,
		genre VARCHAR(255) NOT NULL,
		streaming_services VARCHAR(255) NOT NULL,
		director VARCHAR(255) NOT NULL,
		year VARCHAR(10) NOT NULL,
		watchcount INTEGER NOT NULL,
		rating INTEGER NOT NULL,
		status INTEGER NOT NULL,
		owned INTEGER NOT NULL,
		owned_type VARCHAR(255) NOT NULL,
		ripped INTEGER NOT NULL,
		review TEXT NOT NULL,
		overview TEXT NOT NULL
		);

CREATE TABLE  IF NOT EXISTS genres (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title VARCHAR(30) NOT NULL
		);

CREATE TABLE  IF NOT EXISTS streaming_services (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title VARCHAR(30) NOT NULL
		);

CREATE TABLE  IF NOT EXISTS settings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(100) NOT NULL,
        value VARCHAR(100) NOT NULL,
        description VARCHAR(255) NOT NULL
		);

INSERT OR IGNORE INTO streaming_services
VALUES
    ("1", "Netflix"),
    ("2", "Disney Plus"),
    ("3", "Amazon Prime Video"),
    ("4", "Apple TV+");


INSERT OR IGNORE INTO settings
VALUES
    ("1", "HOME_GRID_VIEW", "false", "Grid or no grid on the Homepage"),
    ("2", "TMDB_API_KEY", "", "Your TMDB api key"),
    ("3", "REGION", "DE", "Your Region");
