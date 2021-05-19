# FILE MANAGER ( BOILERPLATE NAME )

This software handles file storage, running multiple pipelines that are configured by the user.
The files can be retrieved later with all the generated metadata. The user can dynamically add new algorithms to a file pipeline, for generating new metadata for existent and new files.

## Requirements

File Manager requires docker, since it needs to load all the available scripts, ready to run, and docker eliminates the need of configuration for multiple environments.

## Usages

The manager can be used to:

* Compress images to HEIC file format.
* Compress or rename files.
* Categorize images based on content, using computer vision.
* Resize and remove opacity from images.
* Extract text from images using OCR.
* Extract text from audio files.
* Segmentate files depending on metadata

## Configuration

The configuration file uses the following schema:

You can match using the filename and the save_location, then the service will proceed to run the specified scripts on pipeline, in order.

All the metadata will be saved under the same path as `$filename.metadata`.

Example:
  `photo1.jpeg` will have `photo1.jpeg.metadata`

```json
{
  "images": [
    "github.com/file-manager/registry/imagemagick:v0.0.1"
  ],
  "rules": [
    {
      "match": {
        "filename": "*.(png|bmp|jpg|jpeg)",
        "save_location": "/images/*"
      },
      "pipeline": [
        "object_detection",
        "face_recognition",
        "heic_converter"
      ],
    }
  ]
}
```

## Metadata

The generated metadata is stored using the following schema:

```json
{
  "object_detection": {
    "last_run": 213216565465,
    "script_version": 0,
    "detected_objects": [
      {
        "type": "notebook",
        "bounding_box": [120 314 400 563]
      }
    ]
  },
  "face_recognition": {
    "last_run": 213216565789,
    "script_version": 5,
    "faces": [
      {
        "geometry": [0.2145, 5.2157, 8.2546, 5.8842, 9.1547],
        "bounding_box": [542 402 860 756],
      }
    ]
  },
}
```

The fields `last_run` and `script_version` are common on all metadata files to ensure newer versions can run again.

## Run

To run the service, simply run `./file-manager -c config.json`
