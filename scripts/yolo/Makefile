img_name = "docker.pkg.github.com/sonalys/file-manager/yolo"
tag = "v0.0.1-openimages"

image:
	docker build -t ${img_name}:${tag} \
		--build-arg weights_file="yolov3-openimages.weights" \
		--build-arg config_file="data/yolov3-openimages.cfg" \
		--build-arg meta_file="data/openimages.data" \
					.