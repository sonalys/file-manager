package model

import "fmt"

// UploadData serves as generic script input data structure.
type UploadData struct {
	Filename    string                  `json:"filename"`  // Path for the source file.
	Extension   string                  `json:"extension"` // Path for the source file.
	Destination string                  `json:"destination"`
	Metadata    map[string]ScriptOutput `json:"metadata"` // Metadata created by the scripts
	Children    []string                `json:"children"` // Children created from the upload.
}

// GetFullName returns the filename and extension. Example: photo.jpg.
func (u UploadData) GetFullName() string {
	return fmt.Sprintf("%s.%s", u.Filename, u.Extension)
}

// GetFullPath returns the absolute path with filename and extension. Example: /media/photo.jpg.
func (u UploadData) GetFullPath() string {
	return fmt.Sprintf("%s/%s", u.Destination, u.GetFullName())
}
