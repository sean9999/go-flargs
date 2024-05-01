REPO=github.com/sean9999/go-flargs
SEMVER := $$(git tag --sort=-version:refname | head -n 1)

.PHONY: test

info:
	echo REPO is ${REPO} and SEMVER is ${SEMVER}

build:
	go build -v -ldflags="-X 'main.Version=${SEMVER}' -s -w" -o ./bin/kat		./kat/cmd
	go build -v -ldflags="-X 'main.Version=${SEMVER}' -s -w" -o ./bin/rot13		./rot13/cmd
	go build -v -ldflags="-X 'main.Version=${SEMVER}' -s -w" -o ./bin/proverbs	./proverbs/cmd
	go build -v -ldflags="-X 'main.Version=${SEMVER}' -s -w" -o ./bin/noop		./noop/cmd

tidy:
	go mod tidy

test:
	go test ./...

clean:
	go clean

docs:
	pkgsite -open .

publish:
	GOPROXY=https://goproxy.io,direct go list -m ${REPO}@${SEMVER}