quality ?= 55
img_name = "file-manager"
tag = "v0.0.1"
convert:
	docker run -it --rm -v $(CURDIR):/images ${file-manager}/imagemagick:${tag} convert -quality ${quality} /images/${from} /images/${to}

run:
	docker run -it --rm -v $(CURDIR):/images ${file-manager}/imagemagick:${tag} ${run}

build:
	docker build -t ${file-manager}/imagemagick:${VERSION} ./imagemagick