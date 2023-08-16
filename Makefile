run:
	go run main.go
unit-tests:
	go test ./... -v
build:
	docker build . -t adarshtw/passwordly_backend:${IMAGE_TAG}
push:
	docker push adarshtw/passwordly_backend:${IMAGE_TAG}