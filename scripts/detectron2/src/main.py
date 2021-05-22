import time
import os

from PIL.Image import Image

from misc.image import rotate_exif
from classification.model import ImageClassifier
from dotenv import load_dotenv

load_dotenv()

BASE_ADDRESS = "http://" + os.getenv("BASE_ADDRESS", 'localhost:4000')
MEDIA_ADDRESS = BASE_ADDRESS + '/media/thumb_'
LABEL_ADDRESS = BASE_ADDRESS + '/labels'
FETCH_INTERVAL = int(os.getenv("FETCH_INTERVAL", '5'))

def main():
    im_clf = ImageClassifier()
    time.sleep(FETCH_INTERVAL)
    im = Image.load("/buffer/1.jpg")
    im = rotate_exif(im)
    prediction = im_clf.predict(im)


if __name__ == "__main__":
    main()
