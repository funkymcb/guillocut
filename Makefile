run:
	go run main.go | jq '.'

local-build:
	docker build -t guillo-mongo cdk/mongodb/local
	docker run -p 27017:27017 guillo-mongo
