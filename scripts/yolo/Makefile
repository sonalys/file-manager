binary_name = "../../manager/scripts/yolo/yolo"
build:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o ${binary_name} .