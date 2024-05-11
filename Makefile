REPO=github.com/sean9999/go-flargs
SEMVER := $$(git tag --sort=-version:refname | head -n 1)

.PHONY: test

info:
	echo REPO is ${REPO} and SEMVER is ${SEMVER}

bin/kat:
	go build -v -ldflags="-X 'main.Version=${SEMVER}' -s -w" -o ./bin/kat		./kat/cmd

bin/rot13:
	go build -v -ldflags="-X 'main.Version=${SEMVER}' -s -w" -o ./bin/rot13		./rot13/cmd

bin/proverbs:
	go build -v -ldflags="-X 'main.Version=${SEMVER}' -s -w" -o ./bin/proverbs	./proverbs/cmd

bin/noop:
	go build -v -ldflags="-X 'main.Version=${SEMVER}' -s -w" -o ./bin/noop		./noop/cmd

binaries: bin/kat bin/rot13 bin/proverbs bin/noop

tidy:
	go mod tidy

test:
	go test ./...

clean:
	go clean
	go clean -modcache
	rm -f bin/kat bin/noop bin/proverbs bin/rot13

docs:
	pkgsite -open .

publish:
	GOPROXY=https://goproxy.io,direct go list -m ${REPO}@${SEMVER}