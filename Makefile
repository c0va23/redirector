deps:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger

gen-swagger:
	swagger generate server -f api.yml

clean:
	git clean -X -f **/*
