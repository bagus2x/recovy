export DB_USERNAME=postgres
export DB_PASSWORD=admin123
export DB_PORT=5432
export DB_HOST=localhost
export DB_NAME=recovy
export PORT=8080
export ACCESS_TOKEN_KEY=test
export REFRESH_TOKEN_KEY=test
export ACCESS_TOKEN_LIFETIME=900
export REFRESH_TOKEN_LIFETIME=604800
export CACHE_SIZE=1

export POSTGRES_URI=postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable

test:
	echo migrate -source file://path/to/migration -database postgres://localhost:5432/database up 2

create_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrate:
	migrate -source file://db/migration -database postgres://amhdrsmvnfhpxf:ed200a2dfe9c307701c77ddfcdb5c354ceb4e19f31ce350004c5be0917b84a7c@ec2-34-196-238-94.compute-1.amazonaws.com:5432/d3veunub58s339 $(t)

force:
	migrate -path db/migration -database ${POSTGRES_URI} force $(v)

dev:
	go run main.go

redoc:
	npx redoc-cli bundle -o openapi.html openapi.json