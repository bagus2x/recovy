export DB_USERNAME=postgres
export DB_PASSWORD=admin123
export DB_PORT=5432
export DB_HOST=localhost
export DB_NAME=recovy
export APP_PORT=8080
export ACCESS_TOKEN_KEY=test
export REFRESH_TOKEN_KEY=test
export ACCESS_TOKEN_LIFETIME=900
export REFRESH_TOKEN_LIFETIME=604800

export POSTGRES_URI=postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable

test:
	echo migrate -source file://path/to/migration -database postgres://localhost:5432/database up 2

create_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrate:
	migrate -source file://db/migration -database ${POSTGRES_URI} $(t)

force:
	migrate -path db/migration -database ${POSTGRES_URI} force $(v)

dev:
	go run main.go