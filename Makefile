.PHONY: lint
lint: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-golangci-lint run -c .golangci.yml > ./golangci-lint/report-unformatted.json

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: golangci-lint-clean
golangci-lint-clean:
	sudo rm -rf ./golangci-lint

# миграции
.PHONY: migration
migration: #  example: make migration name=add-smth
	docker run --rm \
    -v $(realpath ./internal/db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        create \
        -dir /migrations \
        -ext .sql \
        -seq -digits 3 \
        $(name)

.PHONY: db-upgrade
db-upgrade:
	docker run --rm \
    -v $(realpath ./internal/db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://shortenerodmen:shortenerodmen@172.21.0.2:5432/urlshortener?sslmode=disable \
        up

.PHONY: db-downgrade
db-downgrade:
	docker run --rm \
    -v $(realpath ./internal/db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://shortenerodmen:shortenerodmen@172.21.0.2:5432/urlshortener?sslmode=disable \
        down

.PHONY: db-upgrade-all
db-upgrade-all:
	docker run --rm \
    -v $(realpath ./internal/db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://shortenerodmen:shortenerodmen@172.21.0.2:5432/urlshortener?sslmode=disable \
        up -all

.PHONY: db-downgrade-all
db-downgrade-all:
	docker run --rm \
    -v $(realpath ./internal/db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://shortenerodmen:shortenerodmen@172.21.0.2:5432/urlshortener?sslmode=disable \
        down -all


.PHONY: db-upgrade-to # usage: make db-downgrade-to number=1
db-upgrade-to:
	docker run --rm \
    -v $(realpath ./internal/db/migrations-careful):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://shortenerodmen:shortenerodmen@172.21.0.2:5432/urlshortener?sslmode=disable \
        up $(number)
.PHONY: db-downgrade-to # usage: make db-downgrade-to number=1
db-downgrade-to:
	docker run --rm \
    -v $(realpath ./internal/db/migrations-careful):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://shortenerodmen:shortenerodmen@172.21.0.2:5432/urlshortener?sslmode=disable \
        down $(number)