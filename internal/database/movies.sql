-- name: ListMovies :many
SELECT * FROM movie
ORDER BY title ASC;

-- name: ListMovie :one
SELECT * FROM movie
WHERE id = ?;

-- name: ListWatchlist :many
SELECT * FROM movie
WHERE status = 0
ORDER BY 
    CASE WHEN streaming_services = '' OR streaming_services IS NULL THEN 1 ELSE 0 END,
    title;

-- name: ListSettings :many
SELECT * FROM settings;

-- name: ListSetting :one
SELECT * FROM settings
WHERE id = ?;

-- name: ListSreamingServices :many
SELECT * FROM streaming_services;

-- name: CreateMovie :one
INSERT INTO movie (
    title, original_title, imdbid, tmdbid, length, genre, streaming_services, director, year, watchcount, rating, status, owned, owned_type, ripped, review, overview
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: AddSreamingService :exec
INSERT INTO streaming_services (title) VALUES ( ? );

-- name: DeleteSreamingService :exec
DELETE FROM streaming_services
WHERE id = ?;

-- name: DeleteMovie :exec
DELETE FROM movie
WHERE id = ?;

-- name: ChangeMovieStatus :exec
UPDATE movie
SET status = ?
WHERE id = ?;

-- name: ChangeMovieRating :exec
UPDATE movie
SET rating = ?
WHERE id = ?;

-- name: ChangeMovieOwned :exec
UPDATE movie
SET owned = ?, owned_type = ?
WHERE id = ?;

-- name: ChangeMovieStreamingServices :exec
UPDATE movie
SET streaming_services = ?
WHERE id = ?;

-- name: ChangeMovieRipped :exec
UPDATE movie
SET ripped = ?
WHERE id = ?;

-- name: ChangeSettingValue :exec
UPDATE settings
SET value = ?
WHERE id = ?;
