DB_URL=postgresql://root:secret@localhost:5433/simple_real_estate?sslmode=disable

server:
	go run main.go

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_real_estate

dropdb:
	docker exec -it postgres dropdb --username=root --owner=root simple_real_estate

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate
 
.PHONY: server createdb dropdb migrateup migratedown new_migration sqlc