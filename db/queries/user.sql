-- name: GetUserByMobile :one
SELECT * FROM `users` WHERE mobile = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM `users` WHERE email = ? LIMIT 1;

-- name: GetUserByUid :one
SELECT * FROM `users` WHERE id = ? LIMIT 1;

-- name: InsertUser :exec
INSERT INTO `users` (USERNAME, MOBILE, EMAIL,avatar,
                     HASHED_PASSWORD, IS_EMAIL_VERIFIED, PASSWORD_CHANGED_AT, CREATED_AT,
                     UPDATED_USERNAME_AT, UPDATED_AVATAR_AT, UPDATED_MOBILE_AT, UPDATED_EMAIL_AT)
VALUES (?,?,?,?,?,?,?,?,?,?,?,?);

-- name: UpdateUser :exec
UPDATE `users` SET
username = COALESCE(sqlc.narg(username),username),
avatar = COALESCE(sqlc.narg(avatar),avatar),
mobile = COALESCE(sqlc.narg(mobile),mobile),
email = COALESCE(sqlc.narg(email),email),
hashed_password = COALESCE(sqlc.narg(hashed_password),hashed_password),
is_email_verified = COALESCE(sqlc.narg(is_email_verified),is_email_verified),
password_changed_at = COALESCE(sqlc.narg(password_changed_at),password_changed_at),
updated_email_at = COALESCE(sqlc.narg(updated_email_at),updated_email_at),
updated_mobile_at = COALESCE(sqlc.narg(updated_mobile_at),updated_mobile_at),
updated_username_at = COALESCE(sqlc.narg(updated_username_at),updated_username_at),
updated_avatar_at = COALESCE(sqlc.narg(updated_avatar_at),updated_avatar_at)
WHERE id = sqlc.arg(id);

-- name: GetLatestUser :one
SELECT  * FROM `users` WHERE id = last_insert_id() ;