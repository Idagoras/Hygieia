-- name: InsertEEGData :exec
INSERT INTO `eeg_data`
(session_id, offset, user_id, collected_at, attention,
 meditation, blink_strength, alpha1, alpha2, beta1, beta2, gamma1, gamma2, delta, theta, raw)
VALUES(
     ?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?
);

-- name: GetEEGDataByEEGSessionId :many
SELECT * FROM `eeg_data`
WHERE session_id = ? LIMIT ?;