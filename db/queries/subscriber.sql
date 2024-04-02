-- name: GetSubscribePair :one
SELECT * FROM `subscribers` WHERE uid = ? and sbs_id = ?;

-- name: InsertSubscribePair :exec
INSERT INTO `subscribers` (uid, sbs_id, subscribe_time, cancel_time) VALUES
(
 ?,?,?,?
);

-- name: UpdateSubscribePair :exec
UPDATE `subscribers` SET
subscribe_time = COALESCE(sqlc.narg(subscribe_time),subscribe_time),
cancel_time = COALESCE(sqlc.narg(cancel_time),cancel_time),
is_cancel = COALESCE(sqlc.narg(is_cancel),is_cancel)
WHERE id = ?;

-- name: ListSubscribePairByUid :many
SELECT * FROM `subscribers` WHERE uid = ? AND is_cancel = false LIMIT ? OFFSET ?;

-- name: ListSubscribePairBySuid :many
SELECT * FROM `subscribers` WHERE sbs_id = ? AND is_cancel = false LIMIT ? OFFSET ?;