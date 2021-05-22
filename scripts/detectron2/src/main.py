import time
import os
import sys
import json

from PIL import Image

from misc.image import rotate_exif
from classification.model import ImageClassifier
from dotenv import load_dotenv

load_dotenv()

def main():
    args = sys.argv[1:]
    im_clf = ImageClassifier()
    time.sleep(5)
    with Image.open(args[0]) as im:
        im = rotate_exif(im)
        prediction = im_clf.predict(im)
        with open("/dump.json", "a") as dump:
            dump.write(json.dumps(prediction))


if __name__ == "__main__":
    main()
