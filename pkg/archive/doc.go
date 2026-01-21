// Package archive provides functions for creating and extracting zstd-compressed
// tar archive.
//
// Archives are compressed using zstd. Only regular files and directories are
// supported; symlinks and special files (devices, sockets, named pipes) are
// rejected with [ErrUnsupportedFileType].
//
// Example:
//
//	// Create an archive
//	err := archive.Create("mydir", "output.tar.zst")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Extract the archive
//	err = archive.Extract("output.tar.zst", "extracted")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Extract from an io.Reader
//	file, _ := os.Open("output.tar.zst")
//	defer file.Close()
//	err = archive.ExtractFromReader(file, "extracted")
//	if err != nil {
//		log.Fatal(err)
//	}
package archive
