dev-deps:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get -u github.com/golang/dep/cmd/dep
	go get -u golang.org/x/tools/cmd/goimports

deps:
	dep ensure -vendor-only

gen-swagger:
	swagger generate server -f api.yml

lint:
	bash -c "test $$(goimports -d $$(git ls-files *.go) | tee /dev/stderr | wc -l) -eq 0"

run-test:
	go test ./...

clean:
	git clean -f -d -X  -- **
	git clean -f -d -X  -- cmd/**
