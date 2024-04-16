REPO=github.com/sean9999/go-flargs
SEMVER := $$(git tag --sort=-version:refname | head -n 1)

.PHONY: test

info:
	echo REPO is ${REPO} and SEMVER is ${SEMVER}

tidy:
	go mod tidy

test:
	go test ./...

clean:
	go clean

publish:
	GOPROXY=https://goproxy.io,direct go list -m ${REPO}@${SEMVER}