-- name: CreateEEGSession :exec
INSERT INTO `eeg_session` (
     uid,begin,end,data_count,expired_at
) VALUES (
    ?,?,?,?,?
);

-- name: GetEEGSessionByUID :many
SELECT * FROM `eeg_session` WHERE uid = ? ;

-- name: GetUnfinishedEEGSessionByUID :many
SELECT * FROM `eeg_session` WHERE uid = ? AND finish = FALSE ;

-- name: GetEEGSessionByID :one
SELECT * FROM `eeg_session` WHERE id = ? LIMIT 1;

-- name: GetUnfinishedEEGSession :many
SELECT * FROM `eeg_session` WHERE finish = false;

-- name: UpdateEEGSession :exec
UPDATE `eeg_session`
SET
end = COALESCE(sqlc.narg(end),end),
expired_at = COALESCE(sqlc.narg(expired_at),expired_at),
finish = COALESCE(sqlc.narg(finish),finish),
data_count = data_count +  COALESCE(sqlc.narg(data_count),0),
data_bits_map = COALESCE(sqlc.narg(data_bits_map),data_bits_map)
WHERE id = sqlc.arg(id);

-- name: GetLatestEEGSession :one
SELECT * FROM `eeg_session` WHERE id = LAST_INSERT_ID();

