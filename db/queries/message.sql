-- name: InsertMessage :exec
INSERT INTO `message` (SEND_UID, RCV_UID, CREATED_AT, HAS_READ, TYPE, TEXT, SUBTYPE, ATTACHMENT) VALUES
(
 ?,?,?,?,?,?,?,?
);

-- name: ListUserReceivedMessages :many
SELECT * FROM `message` WHERE rcv_uid = ? LIMIT ? OFFSET ?;

-- name: ListUserSentMessages :many
SELECT * FROM `message` WHERE send_uid = ? LIMIT ? OFFSET ?;

-- name: GetMessage :many
SELECT * FROM `message` WHERE send_uid = ? AND rcv_uid = ? AND created_at = ? LIMIT 1;
