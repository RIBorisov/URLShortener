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
DSN=postgres://shortenerodmen:shortenerodmen@172.18.0.2:5432/urlshortener?sslmode=disable
.PHONY: migration
migration: #  example: make migration name=add-smth
	docker run --rm \
    -v $(realpath ./internal/storage/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        create \
        -dir /migrations \
        -ext .sql \
        -seq -digits 3 \
        $(name)

.PHONY: db-upgrade
db-upgrade:
	docker run --rm \
    -v $(realpath ./internal/storage/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database $(DSN) \
        up

.PHONY: db-downgrade
db-downgrade:
	docker run --rm \
    -v $(realpath ./internal/storage/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database $(DSN) \
        down

RAWFILE:=coverage.out
HTMLREPORT:=coverage.html

.PHONY: coverage
coverage:
	go test ./internal/handlers -coverprofile=$(RAWFILE) && \
 	go tool cover -html=$(RAWFILE) -o $(HTMLREPORT)

PACKAGES := $(shell go list ./... | grep -vE "mocks|models|logger|storage|service" | tr '\n' ' ')

.PHONY: tests
tests:
	#go list ./... | grep -vE "mocks|models|logger|cmd|storage|service|middleware"|xargs go test -v -coverpkg=$1 -coverprofile=profile.cov $1
	go list ./... | grep -vE "mocks"|xargs go test -v -coverpkg=$1 -coverprofile=profile.cov $1
	go tool cover -func profile.cov

#go tool cover -html=$(RAWFILE) -o $(HTMLREPORT)


.PHONY: prof
prof:
	go tool pprof -http=":9090" -seconds=30 http://localhost:8081/debug/pprof/profile


.PHONY: save-base-pprof
save-base-pprof:
	curl http://127.0.0.1:8081/debug/pprof/profile > ./profiles/base.pprof
	#go tool pprof -http=":9090" ./profiles/base.pprof


.PHONY: gen-mocks
gen-mocks:
	mockgen -source=internal/service/service.go -destination=internal/service/mocks/service_mock.gen.go -package=mocks

