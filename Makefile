.PHONY:
run:
	@docker-compose up -d --build

.PHONY:
server:
	@docker-compose up -d mysql-server

.PHONY:
i:
	@echo "WARNING: before run the tests, make sure to start the containers with:"
	@echo "    $ make up"
	@echo
	@echo ">> running all tests in the application container"
	@docker-compose build bookstore
	@docker-compose run --rm bookstore /app/test.sh

.PHONY:
lint: gopkgcache
	@echo ">> running golangci-lint"
	@golangci-lint run ./...

.PHONY:
gopkgcache:
# 'go list' needs to be executed before staticcheck to prepopulate the modules cache.
# Otherwise staticcheck might fail randomly for some reason not yet explained.
	@go list -e -compiled -test=true -export=false -deps=true -find=false -tags= -- ./... > /dev/null
