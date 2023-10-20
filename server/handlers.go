package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

// handlerDirectoryRead processes requests that read single directory content.
func handlerDirectoryRead(cfg *Configuration, w http.ResponseWriter, req *http.Request) {
	if name, ok := requiredNameParam(w, req); ok {
		if directory, errorDto := directoryContent(cfg, name); errorDto == nil {
			writeResultDirectory(w, directory)
		} else {
			writeResultError(w, errorDto)
		}
	}
}

// handlerDirectoryTree processes requests that read directory tree content.
func handlerDirectoryTree(cfg *Configuration, w http.ResponseWriter, _ *http.Request) {
	// TODO add parameter
	directory, _ := directoryTree(cfg)
	data, _ := json.MarshalIndent(directory, "", "  ")
	_, _ = w.Write(data)
}

// handlerDirectoryList processes requests that list whole directory tree with full relative paths.
func handlerDirectoryList(cfg *Configuration, w http.ResponseWriter, req *http.Request) {
	if name, ok := requiredNameParam(w, req); ok {
		if dirList, errorDto := directoryList(cfg, name); errorDto == nil {
			writeResultData(w, DirectoryListDto{dirList})
		} else {
			writeResultError(w, errorDto)
		}
	}
}

// handlerDirectoryCreate processes requests that create new directory.
func handlerDirectoryCreate(cfg *Configuration, w http.ResponseWriter, req *http.Request) {
	var name, all string
	var ok bool
	if name, ok = requiredNameParam(w, req); !ok {
		return
	}
	if all, ok = optionalSingleParam(w, req, "all", "false"); !ok {
		return
	}
	if directory, errorDto := createDirectory(cfg, name, strings.ToLower(all) == "true"); errorDto == nil {
		writeResultDirectory(w, directory)
	} else {
		writeResultError(w, errorDto)
	}
}

// handlerDirectoryDelete processes requests that delete existing directory.
func handlerDirectoryDelete(cfg *Configuration, w http.ResponseWriter, req *http.Request) {
	var name, all string
	var ok bool
	// check if single and required directory name is given as parameter
	if name, ok = requiredNameParam(w, req); !ok {
		return
	}
	// get deep delete flag or set it to false if not present
	if all, ok = optionalSingleParam(w, req, "all", "false"); !ok {
		return
	}
	// delete directory and optionally its whole content
	if directory, errorDto := deleteDirectory(cfg, name, strings.ToLower(all) == "true"); errorDto == nil {
		writeResultDirectory(w, directory)
	} else {
		writeResultError(w, errorDto)
	}
}

// TODO add documentation
func handlerFileRead(cfg *Configuration, w http.ResponseWriter, req *http.Request) {
	if name, ok := requiredNameParam(w, req); ok {
		if offset, ok := requiredIntParam(w, req, "offset"); ok {
			if size, ok := requiredIntParam(w, req, "size"); ok {
				if errorDto := readFile(cfg, w, name, offset, size); errorDto != nil {
					writeResultError(w, errorDto)
				}
			}
		}
	}
}

// TODO add documentation
func handlerFileWrite(cfg *Configuration, w http.ResponseWriter, req *http.Request) {
	if name, ok := requiredNameParam(w, req); ok {
		if file, errorDto := fileWrite(cfg, req, name); errorDto == nil {
			writeResultFile(w, file)
		} else {
			writeResultError(w, errorDto)
		}
	}
}

// TODO add documentation
func handlerFileAppend(cfg *Configuration, w http.ResponseWriter, req *http.Request) {
	if name, ok := requiredNameParam(w, req); ok {
		if file, errorDto := fileAppend(cfg, req, name); errorDto == nil {
			writeResultFile(w, file)
		} else {
			writeResultError(w, errorDto)
		}
	}
}

// handlerFileDelete processes requests that delete specified file.
func handlerFileDelete(cfg *Configuration, w http.ResponseWriter, req *http.Request) {
	if name, ok := requiredNameParam(w, req); ok {
		if file, errorDto := fileDelete(cfg, name); errorDto == nil {
			writeResultFile(w, file)
		} else {
			writeResultError(w, errorDto)
		}
	}
}

// handlerFileExists processes requests that check if specified file exists.
func handlerFileExists(cfg *Configuration, w http.ResponseWriter, req *http.Request) {
	if name, ok := requiredNameParam(w, req); ok {
		if file, errorDto := fileExists(cfg, name); errorDto == nil {
			writeResultFile(w, file)
		} else {
			writeResultError(w, errorDto)
		}
	}
}

// handlerFileChecksum processes requests that calculate file checksum.
func handlerFileChecksum(cfg *Configuration, w http.ResponseWriter, req *http.Request) {
	if name, ok := requiredNameParam(w, req); ok {
		if file, errorDto := fileChecksum(cfg, name); errorDto == nil {
			writeResultFile(w, file)
		} else {
			writeResultError(w, errorDto)
		}
	}
}

// handlerFileShared processes requests that read file contents shared as link.
func handlerFileShared(cfg *Configuration, w http.ResponseWriter, req *http.Request) {
	uriPrefix := cfg.UrlPrefix + routeFileShared
	name := "/" + strings.TrimPrefix(req.RequestURI, uriPrefix)
	if errorDto := writeSharedFileContent(cfg, w, name); errorDto != nil {
		writeResultError(w, errorDto)
	}
}
