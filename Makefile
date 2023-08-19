devrun: dbrun migration-up run

devstop: migration-up-stop dbstop

run:
	export `grep -v '#' .env | xargs` && DB_HOST=localhost && go run cmd/main.go

dbrun: 
	docker compose -f postgres-db.yml up -d

dbstop: 
	docker compose -f postgres-db.yml down

migration-up:
	docker compose -f migrate-up.yml up -d

migration-up-stop:
	docker compose -f migrate-up.yml down

migration-down:
	docker compose -f migrate-down.yml up -d

build:
	docker compose -f production.yml up -d --build

down:
	docker compose -f production.yml down
