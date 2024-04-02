DB_URL=mysql://root:314159@tcp(localhost:3306)/hygieia
migrateup:
	migrate -path db/migration -database "${DB_URL}" -verbose up
migrateup1:
	migrate -path db/migration -database "${DB_URL}" -verbose up 1
migratedown:
	migrate -path db/migration -database "${DB_URL}" -verbose down
migratedown1:
	migrate -path db/migration -database "${DB_URL}" -verbose down 1
new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)
sqlc:
	sqlc generate
proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative proto/*.proto
.PHONY:migrateup1 migrateup migratedown1 migratedown new_migration sqlc proto