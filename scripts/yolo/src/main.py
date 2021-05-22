import os
import json
import os.path
import sys
import signal
import darknet

# Load YOLO model:
configPath = "/app/data/yolov3-openimages.cfg"
weightPath = "/app/data/yolov3-openimages.weights"
metaPath = "/app/data/openimages.data"

network, class_names, class_colors = darknet.load_network(
    configPath,
    metaPath,
    weightPath,
    batch_size=1
)

def detect(filename, threshold):
    im = darknet.load_image(bytes(filename, "ascii"), 0, 0)
    r = darknet.detect_image(network, class_names, im, thresh=threshold)
    darknet.free_image(im)
    print(r)
    return [item[0] for item in r]

def main():
    args = sys.argv[1:]
    prediction = detect(args[0], 0.25)
    with open("/dump.json", "a") as dump:
        dump.write(json.dumps(prediction))


if __name__ == '__main__':
    main()
