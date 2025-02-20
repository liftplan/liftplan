
coverage_file=coverage.out

.PHONY: run dev

vendor:
	@cd serve && ./scripts/get-deps.sh

build:
	docker build . -t liftplan

run: build
	docker run -p 9000:9000 liftplan
	
# this is just a quality of life setting for watching all files that could
# change and rebuilding the server. This is mostly used for deving html templates
dev:
	find ./strategy ./gear ./serve -print | entr -r make run

test:
	go test -race -v ./strategy/... ./gear/...

coverage:
	go test -race -coverprofile=$(coverage_file) -covermode=atomic ./strategy/... ./gear/... && go tool cover -html=$(coverage_file)
