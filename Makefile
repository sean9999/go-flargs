REPO=github.com/sean9999/go-flargs
SEMVER := $$(git tag --sort=-version:refname | head -n 1)
BRANCH := $$(git branch --show-current)
REF := $$(git describe --dirty --tags --always)

.PHONY: test

info:
	@printf "REPO:\t%s\nSEMVER:\t%s\nBRANCH:\t%s\nREF:\t%s\n" $(REPO) $(SEMVER) $(BRANCH) $(REF)

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
pkgsite:
	if [ -z "$$(command -v pkgsite)" ]; then go install golang.org/x/pkgsite/cmd/pkgsite@latest; fi

docs: pkgsite
	pkgsite -open .

publish:
	GOPROXY=https://goproxy.io,direct go list -m ${REPO}@${SEMVER}