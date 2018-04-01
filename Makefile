dev-deps:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get -u github.com/golang/dep/cmd/dep

deps:
	dep ensure -vendor-only

gen-swagger:
	swagger generate server -f api.yml

clean:
	git clean -f -d -X  -- **
	git clean -f -d -X  -- cmd/**
