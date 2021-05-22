img_name = "docker.pkg.github.com/sonalys/file-manager/detectron2"
tag = "v0.0.1"
binary_name = "../../manager/scripts/detectron2"

build:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ${binary_name} .

image:
	docker build -t ${img_name}:${tag} ./

push:
	docker push ${img_name}:${tag}