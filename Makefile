PACKAGE := github.com/c0va23/redirector
ERRCHECK_EXCLUDE_PATTERN := '(cmd/redirector-server|restapi/operations)'
GOSWAGGER_VERSION := 0.15.0
BUILD_ARCH := amd64

dev-deps:
	go get -u github.com/golang/dep/cmd/dep
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u honnef.co/go/tools/cmd/staticcheck
	go get -u github.com/phogolabs/parcello/cmd/parcello

deps:
	dep ensure -vendor-only

gen-swagger: bin/swagger
	bin/swagger generate server \
		--model FieldValidationError \
		--model Locales \
		--model LocaleTranslations \
		--model ModelValidationError \
		--model ServerError \
		--model Translation \
		--model ValidationError \
		-f api.yml

lint:
	go vet ./...
	test $$(goimports -d $$(git ls-files *.go) | tee /dev/stderr | wc -l) -eq 0
	golint -set_exit_status $$(go list ./...)
	env CGO_ENABLED=0 errcheck $$(go list ./... | grep -v -E $(ERRCHECK_EXCLUDE_PATTERN))
	env CGO_ENABLED=0 staticcheck ./...

run-test:
	go test -cover ./...

gen-locales:
	go generate ./locales/build.go

bin/:
	mkdir bin/

bin/swagger: bin/
	curl -o bin/swagger -L https://github.com/go-swagger/go-swagger/releases/download/$(GOSWAGGER_VERSION)/swagger_linux_$(BUILD_ARCH)
	chmod +x bin/swagger

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
