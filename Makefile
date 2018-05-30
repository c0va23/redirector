PACKAGE := github.com/c0va23/redirector
ERRCHECK_EXCLUDE_PATTERN := '(cmd/redirector-server|restapi/operations)'

dev-deps:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get -u github.com/golang/dep/cmd/dep
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u honnef.co/go/tools/cmd/staticcheck

deps:
	dep ensure -vendor-only

gen-swagger:
	swagger generate server -f api.yml

lint:
	go vet ./...
	bash -c "test $$(goimports -d $$(git ls-files *.go) | tee /dev/stderr | wc -l) -eq 0"
	golint -set_exit_status $$(go list ./...)
	errcheck $$(go list ./... | grep -v -E $(ERRCHECK_EXCLUDE_PATTERN))
	staticcheck ./...

run-test:
	go test -cover ./...

bin/:
	mkdir bin/

bin/redirector-server-linux-amd64: bin/
	GOOS=linux \
	GOARCH=amd64 \
	go build -o $@ $(PACKAGE)/cmd/redirector-server

bin/redirector-server-linux-arm: bin/
	GOOS=linux \
	GOARCH=arm \
	go build -o $@ $(PACKAGE)/cmd/redirector-server

bin/redirector-server-linux-arm64: bin/
	GOOS=linux \
	GOARCH=arm64 \
	go build -o $@ $(PACKAGE)/cmd/redirector-server

clean:
	git clean -f -d -X  -- **
	git clean -f -d -X  -- cmd/**

bump-minor-version:
	./scripts/bump_version.go VERSION minor

bump-patch-version:
	./scripts/bump_version.go VERSION patch

commit-version:
	git commit VERSION -m "Bump version $(shell cat VERSION)"

tag-version:
	git tag $(shell cat VERSION)

tag-push:
	git push origin $(shell cat VERSION)
