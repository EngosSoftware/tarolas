package server

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// File stores single file attributes like name and size.
type File struct {
	Name     *string `json:"name,omitempty"      api:"File name without parent path."`
	Size     *int64  `json:"size,omitempty"      api:"File size in bytes."`
	Checksum *string `json:"checksum,omitempty"  api:"File checksum (SHA256)."`
	Exists   *bool   `json:"exists,omitempty"    api:"Flag indicating if file exists."`
}

// FileDto is an implementation of DTO for file.
type FileDto struct {
	Data *File `json:"data,omitempty"  api:"File details."`
}

// TODO add documentation
func fileWrite(cfg *Configuration, req *http.Request, name string) (*File, *ErrorDto) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			logError(err)
		}
	}()
	fullName := prepareAbsolutePath(cfg, name)
	if file, err := os.OpenFile(fullName, os.O_CREATE|os.O_WRONLY, 0755); err == nil {
		defer func() {
			if err := file.Close(); err != nil {
				logError(err)
			}
		}()
		decoder := base64.NewDecoder(base64.StdEncoding, req.Body)
		if _, err := io.Copy(file, decoder); err == nil {
			fileInfo, _ := file.Stat()
			size := fileInfo.Size()
			return &File{Name: &name, Size: &size}, nil
		} else {
			return nil, errorDto(errWritingFileFailed, name)
		}
	} else {
		return nil, errorDto(errOpeningFileForWritingFailed, name)
	}
}

// readFile reads file part defined by offset and size. The returned content is base64 encoded.
func readFile(cfg *Configuration, w http.ResponseWriter, name string, offset, size int64) *ErrorDto {
	fullName := prepareAbsolutePath(cfg, name)
	// retrieve file info
	fileInfo, err := os.Stat(fullName)
	if err != nil {
		if os.IsNotExist(err) {
			return errorDto(errFileNotFound, name)
		}
		return errorDto(errRetrievingFileInfoFailed, name)
	}
	if fileInfo.IsDir() {
		return errorDto(errNotAFile, name)
	}
	// check the given offset
	if offset < 0 || offset > fileInfo.Size() {
		return errorDto(errInvalidParameterValue, "offset ("+strconv.FormatInt(offset, 10)+")")
	}
	// check the given size
	if size <= 0 || size > fileInfo.Size() {
		return errorDto(errInvalidParameterValue, "size ("+strconv.FormatInt(size, 10)+")")
	}
	// truncate the requested size if exceeds file size
	if offset+size > fileInfo.Size() {
		size = fileInfo.Size() - offset
	}
	// open file for reading
	if file, err := os.Open(fullName); err == nil {
		// make sure the file will be properly closed
		defer func() {
			if err := file.Close(); err != nil {
				logError(err)
			}
		}()
		if _, err := file.Seek(offset, 0); err != nil {
			return errorDto(errSeekingFileFailed, name)
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		encoder := base64.NewEncoder(base64.StdEncoding, w)
		if _, err := io.CopyN(encoder, file, size); err == nil {
			return nil
		} else {
			return errorDto(errReadingFileFailed, name)
		}
	} else {
		return errorDto(errRetrievingFileInfoFailed, name)
	}
}

// TODO add documentation
func fileAppend(cfg *Configuration, req *http.Request, name string) (*File, *ErrorDto) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			logError(err)
		}
	}()
	fullName := prepareAbsolutePath(cfg, name)
	if file, err := os.OpenFile(fullName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755); err == nil {
		defer func() {
			if err := file.Close(); err != nil {
				logError(err)
			}
		}()
		decoder := base64.NewDecoder(base64.StdEncoding, req.Body)
		if _, err := io.Copy(file, decoder); err == nil {
			fileInfo, _ := file.Stat()
			size := fileInfo.Size()
			return &File{Name: &name, Size: &size}, nil
		} else {
			logError(err)
			return nil, errorDto(errAppendingFileFailed, name)
		}
	} else {
		logError(err)
		return nil, errorDto(errOpeningFileForAppendFailed, name)
	}
}

// fileDelete deletes file with specified name.
func fileDelete(cfg *Configuration, name string) (*File, *ErrorDto) {
	fullName := prepareAbsolutePath(cfg, name)
	if fileInfo, err := os.Stat(fullName); err == nil {
		if fileInfo.IsDir() {
			return nil, errorDto(errNotAFile, name)
		} else {
			size := fileInfo.Size()
			if err := os.Remove(fullName); err == nil {
				return &File{Name: &name, Size: &size}, nil
			} else {
				logError(err)
				return nil, errorDto(errDeletingFileFailed, name)
			}
		}
	} else {
		if os.IsNotExist(err) {
			return nil, errorDto(errFileNotFound, name)
		}
		logError(err)
		return nil, errorDto(errRetrievingFileInfoFailed, name)
	}
}

// fileExists checks if the file with specified name exists.
// When file with given name was found, then 'true' flag is returned.
func fileExists(cfg *Configuration, name string) (*File, *ErrorDto) {
	fullName := prepareAbsolutePath(cfg, name)
	if fileInfo, err := os.Stat(fullName); err == nil {
		if fileInfo.IsDir() {
			return nil, errorDto(errNotAFile, name)
		} else {
			size := fileInfo.Size()
			return &File{Name: &name, Size: &size, Exists: &FlagTrue}, nil
		}
	} else {
		if os.IsNotExist(err) {
			return &File{Name: &name, Exists: &FlagFalse}, nil
		}
		logError(err)
		return nil, errorDto(errRetrievingFileInfoFailed, name)
	}
}

// TODO add documentation
func fileChecksum(cfg *Configuration, name string) (*File, *ErrorDto) {
	fullName := prepareAbsolutePath(cfg, name)
	if file, err := os.Open(fullName); err == nil {
		defer func() {
			if err := file.Close(); err != nil {
				logError(err)
			}
		}()
		hash := sha256.New()
		if _, err := io.Copy(hash, file); err == nil {
			fileInfo, _ := file.Stat()
			size := fileInfo.Size()
			checksum := fmt.Sprintf("%x", hash.Sum(nil))
			return &File{Name: &name, Size: &size, Checksum: &checksum}, nil
		} else {
			logError(err)
			return nil, errorDto(errCalculatingChecksumFailed, name)
		}
	} else {
		if os.IsNotExist(err) {
			return nil, errorDto(errFileNotFound, name)
		}
		logError(err)
		return nil, errorDto(errOpeningFileForReadingFailed, name)
	}
}

func writeSharedFileContent(cfg *Configuration, w http.ResponseWriter, name string) *ErrorDto {
	fullName := prepareAbsolutePath(cfg, name)
	// open shared file for reading
	if file, err := os.Open(fullName); err == nil {
		defer func() {
			if err := file.Close(); err != nil {
				logError(err)
			}
		}()
		// stat file
		if fileInfo, err := file.Stat(); err == nil {
			if fileInfo.IsDir() {
				return errorDto(errNotAFile, name)
			}
			// detect file content
			buffer := make([]byte, 512)
			if _, err = file.Read(buffer); err != nil {
				return errorDto(errReadingFileFailed, name)
			}
			contentType := http.DetectContentType(buffer)
			if _, err := file.Seek(0, 0); err != nil {
				return errorDto(errSeekingFileFailed, name)
			}
			w.WriteHeader(200)
			w.Header().Set("Content-Type", contentType)
			if _, err = io.Copy(w, file); err == nil {
				return nil
			} else {
				return errorDto(errReadingFileFailed, name)
			}
		} else {
			return errorDto(errRetrievingFileInfoFailed, name)
		}
	} else {
		if os.IsNotExist(err) {
			return errorDto(errFileNotFound, name)
		}
		return errorDto(errOpeningFileForReadingFailed, name)
	}
}
