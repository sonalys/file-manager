img_name = "file-manager"
tag = "v0.0.1"
binary_name = "../../manager/scripts/compress2heic"
build:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ${binary_name} .

image:
	make build
	docker build -t ${img_name}/imagemagick:${tag} ./
	rm ${binary_name}