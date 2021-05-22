img_name = "docker.pkg.github.com/sonalys/file-manager"
tag = "v0.0.1"
binary_name = "../../manager/scripts/compress2heic"

build:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ${binary_name} .

image:
	make build
	docker build -t ${img_name}/imagemagick:${tag} ./

push:
	docker push ${img_name}/imagemagick:${tag}