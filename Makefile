export tag=v1.0.0

build:
	go build \
	-race \
	-o app

release:
	docker build . -t listenzz/gohttpserver:${tag}

run: release
	docker run -p 8000:8000 listenzz/gohttpserver:${tag} -- -address :8000

push: release
	docker push listenzz/gohttpserver:${tag}
