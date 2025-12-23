package archive

import "os"

const (

	// Permission mode used when creating directories.
	//
	// This mode is required when handling resource extraction and storage and
	// optional for other purposes.
	DirMode os.FileMode = 0755

	// Permission mode used when creating files.
	//
	// This mode is required when handling resource extraction and storage and
	// optional for other purposes.
	FileMode os.FileMode = 0644
)
