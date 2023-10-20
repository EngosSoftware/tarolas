package server

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	RootSymbol = "/"
)

// Directory stores directory attributes like name and contained files and directories.
type Directory struct {
	Name        string       `json:"name,omitempty"         api:"Name of the directory without parent path."`
	Directories []*Directory `json:"directories,omitempty"  api:"Child directories in this directory."`
	Files       []*File      `json:"files,omitempty"        api:"Files contained in this directory."`
}

// DirectoryDto is the implementation of DTO for directory.
type DirectoryDto struct {
	Data *Directory `json:"data"  api:"Directory details."`
}

// DirectoryListDto is the implementation of DTO for directory list.
type DirectoryListDto struct {
	Data []string `json:"data"  api:"List of directories."`
}

// AddDirectory adds new subdirectory with specified name.
func (d *Directory) AddDirectory(name string) *Directory {
	directory := Directory{Name: name}
	d.Directories = append(d.Directories, &directory)
	return &directory
}

// AddFile adds new file with specified name and size.
func (d *Directory) AddFile(name string, size int64) *File {
	file := File{Name: &name, Size: &size}
	d.Files = append(d.Files, &file)
	return &file
}

// directoryTree returns whole directory tree including files.
func directoryTree(cfg *Configuration) (*Directory, error) {
	rootDirectory := cfg.RootDirectory
	directory := Directory{Name: filepath.Base(rootDirectory)}
	lookup := make(map[string]*Directory)
	lookup[rootDirectory] = &directory
	err := filepath.Walk(rootDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path != rootDirectory {
			parentDir := filepath.Dir(path)
			if info.IsDir() {
				if dir, ok := lookup[parentDir]; ok {
					lookup[path] = dir.AddDirectory(info.Name())
				}
			} else {
				if dir, ok := lookup[parentDir]; ok {
					dir.AddFile(info.Name(), info.Size())
					lookup[parentDir] = dir
				}
			}
		}
		return nil
	})
	return &directory, err
}

// directoryList returns whole directory tree excluding files with relative paths.
// The result is a list of directory names (with relative paths from root directory)
// in the same order they are walked through by function filepath.Walk.
// The first element in directory list is the root directory itself (/).
func directoryList(cfg *Configuration, name string) ([]string, *ErrorDto) {
	fullName := prepareAbsolutePath(cfg, name)
	dirList := make([]string, 0)
	dirList = append(dirList, "/")
	err := filepath.Walk(fullName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path != fullName {
			if info.IsDir() {
				dirList = append(dirList, strings.TrimPrefix(path, fullName))
			}
		}
		return nil
	})
	if err != nil {
		logError(err)
		return nil, errorDto(errWalkingDirectoryTreeFailed, errMsgCheckServerLogForDetails)
	}
	return dirList, nil
}

// directoryContent lists the content of specified directory without subdirectories.
func directoryContent(cfg *Configuration, name string) (*Directory, *ErrorDto) {
	fullName := prepareAbsolutePath(cfg, name)
	rootName := filepath.Base(fullName)
	if fullName == cfg.RootDirectory {
		rootName = RootSymbol
	}
	directory := Directory{Name: rootName}
	entries, err := os.ReadDir(fullName)
	if err != nil {
		return nil, errorDto(errReadingDirectoryContentFailed, name)
	}
	for _, fileEntry := range entries {
		fileInfo, err := fileEntry.Info()
		if err != nil {
			return nil, errorDto(errReadingDirectoryContentFailed, name)
		}
		if fileInfo.IsDir() {
			directory.AddDirectory(fileInfo.Name())
		} else {
			directory.AddFile(fileInfo.Name(), fileInfo.Size())
		}
	}
	return &directory, nil
}

// TODO add documentation
func createDirectory(cfg *Configuration, name string, all bool) (*Directory, *ErrorDto) {
	const dirMode = 0755
	fullName := prepareAbsolutePath(cfg, name)
	if all {
		err := os.MkdirAll(fullName, dirMode)
		if err != nil {
			return nil, errorDto(errCreatingDirectoriesFailed, name)
		}
		return &Directory{Name: filepath.Base(name)}, nil
	} else {
		err := os.Mkdir(fullName, dirMode)
		if err != nil {
			if os.IsExist(err) {
				return nil, errorDto(errDirectoryAlreadyExists, name)
			}
			return nil, errorDto(errCreatingDirectoryFailed, name)
		}
		return &Directory{Name: name}, nil
	}
}

// deleteDirectory deletes directory with specified name.
// If 'all' flag is 'false' the directory must be empty. Deleting root directory has no effect.
// If 'all' flag is 'true' the directory may contain other directories or files.
// Deleting root directory removes all its content omitting root directory itself.
func deleteDirectory(cfg *Configuration, name string, all bool) (*Directory, *ErrorDto) {
	fullName := prepareAbsolutePath(cfg, name)
	if all {
		if fileInfos, err := ioutil.ReadDir(fullName); err == nil {
			for _, fileInfo := range fileInfos {
				fileName := prepareAbsolutePath(cfg, fileInfo.Name())
				if fileInfo.IsDir() {
					if err = os.RemoveAll(fileName); err != nil {
						return nil, errorDto(errDeletingDirectoryFailed, fileName)
					}
				} else {
					if err = os.Remove(fileName); err != nil {
						return nil, errorDto(errDeletingFileFailed, fileName)
					}
				}
			}
		} else {
			return nil, errorDto(errReadingDirectoryContentFailed, name)
		}
	}
	if fullName == cfg.RootDirectory {
		return &Directory{Name: RootSymbol}, nil
	}
	if err := os.Remove(fullName); err == nil {
		return &Directory{Name: filepath.Base(fullName)}, nil
	} else {
		return nil, errorDto(errDeletingDirectoryFailed, name)
	}
}
