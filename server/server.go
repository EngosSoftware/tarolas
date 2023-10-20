// Package server implements the file storage service functionality and API.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	routeDirectoryRead   = "/directory/read"   // Reads directory content.
	routeDirectoryTree   = "/directory/tree"   // Reads directory tree.
	routeDirectoryList   = "/directory/list"   // Lists all directories in tree with full relative paths.
	routeDirectoryCreate = "/directory/create" // Creates new directory.
	routeDirectoryDelete = "/directory/delete" // Deletes existing directory.
	routeFileRead        = "/file/read"        // Reads file's content.
	routeFileWrite       = "/file/write"       // Writes to existing file or creates a new one and writes to it.
	routeFileAppend      = "/file/append"      // Appends an existing file or creates a new one and appends it.
	routeFileDelete      = "/file/delete"      // Deletes existing file.
	routeFileExists      = "/file/exists"      // Checks if file exists.
	routeFileChecksum    = "/file/checksum"    // Calculates file checksum.
	routeFileShared      = "/shared/"          // Shares the file content (accessible as link to file).
	HttpGET              = "GET"               // HTTP get method.
	HttpPOST             = "POST"              // HTTP post method.
	HttpPUT              = "PUT"               // HTTP put method.
	HttpDELETE           = "DELETE"            // HTTP delete method.
	HttpOPTIONS          = "OPTIONS"           // HTTP options method.
)

// Handler defines custom type for declaring request handlers.
type Handler func(w http.ResponseWriter, req *http.Request)

// RouteHandler defines custom type for declaring route handlers.
type RouteHandler func(cfg *Configuration, w http.ResponseWriter, req *http.Request)

// requiredNameParam searches for required parameter named 'name' and validates
// the value against file and directory naming rules.
func requiredNameParam(w http.ResponseWriter, req *http.Request) (string, bool) {
	if name, ok := requiredSingleParam(w, req, "name"); ok {
		if !strings.HasPrefix(name, "/") {
			writeResultError(w, errorDto(errNoSlashInFileOrDirectoryName, name))
			return "", false
		}
		return name, true
	}
	return "", false
}

func requiredIntParam(w http.ResponseWriter, req *http.Request, name string) (int64, bool) {
	if strValue, ok := requiredSingleParam(w, req, name); !ok {
		return 0, false
	} else {
		if value, err := strconv.ParseInt(strValue, 10, 64); err != nil {
			writeResultError(w, errorDto(errRequiredParameterIsNotAnInteger, name))
			return 0, false
		} else {
			return value, true
		}
	}
}

// requiredSingleParam searches for a parameter with specified name.
// Parameter with this name should be present and may not be given more than once.
func requiredSingleParam(w http.ResponseWriter, req *http.Request, name string) (string, bool) {
	keys, ok := req.URL.Query()[name]
	if !ok {
		writeResultError(w, errorDto(errRequiredParameterIsMissing, name))
		return "", false
	}
	if len(keys) > 1 {
		writeResultError(w, errorDto(errOnlyOneParameterAllowed, name))
		return "", false
	}
	value := strings.TrimSpace(keys[0])
	if value == "" {
		writeResultError(w, errorDto(errRequiredParameterIsEmpty, name))
		return "", false
	}
	return value, true
}

// optionalSingleParam searches for a parameter with specified name.
// Parameter with this name may be present and may not be given more than once.
func optionalSingleParam(w http.ResponseWriter, req *http.Request, name string, defaultValue string) (string, bool) {
	keys, ok := req.URL.Query()[name]
	if !ok {
		return defaultValue, true
	}
	if len(keys) > 1 {
		writeResultError(w, errorDto(errOnlyOneParameterAllowed, name))
		return "", false
	}
	value := strings.TrimSpace(keys[0])
	if value == "" {
		return defaultValue, true
	}
	return value, true
}

// writeResultFile is a helper method for returning file DTO to caller with status 200.
func writeResultFile(w http.ResponseWriter, file *File) {
	writeResultData(w, FileDto{Data: file})
}

// writeResultDirectory is a helper method for returning directory DTO to caller with status 200.
func writeResultDirectory(w http.ResponseWriter, directory *Directory) {
	writeResultData(w, DirectoryDto{Data: directory})
}

// writeResultData is a helper method for returning DTO to caller.
// DTO returns correct data being the answer to request with HTTP status 200.
func writeResultData(w http.ResponseWriter, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	dataLength := len(jsonData)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/vnd.api+json")
	count, err := w.Write(jsonData)
	if err != nil || count != dataLength {
		log.Fatal(err)
	}
}

// writeResultError is a helper function for writing error objects back to caller.
// Single error DTO object is encapsulated with error list DTO object, and then
//  converted to JSON body and returned to caller. HTTP status is set according
// status code defined in error.
func writeResultError(w http.ResponseWriter, errorDto *ErrorDto) {
	errorsDto := errorsDto(errorDto)
	jsonData, err := json.MarshalIndent(errorsDto, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	statusCode, err := strconv.Atoi(errorDto.Status)
	if err != nil {
		log.Fatal(err)
	}
	dataLength := len(jsonData)
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/vnd.api+json")
	count, err := w.Write(jsonData)
	if err != nil || count != dataLength {
		log.Fatal(err)
	}
}

// httpHandler creates handler that checks if current request has the same
// HTTP method as defined in parameter. If current request method differs
// from the one passed as an argument, then error is returned to caller.
// Only requests with configured HTTP method are further processed.
// Aditionally CORS is enabled and OPTIONS preflight is supported.
func httpHandler(cfg *Configuration, method string, handler RouteHandler) Handler {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if req.Method == HttpOPTIONS {
			return
		}
		if req.Method == method {
			handler(cfg, w, req)
		} else {
			writeResultError(w, errorDto(errRequestMethodNotSupported, req.Method))
		}
	}
}

// StartServer starts the file server.
func StartServer(cfg *Configuration) *http.Server {
	// configure all routes (with prefixes)
	prefix := cfg.UrlPrefix
	mux := http.NewServeMux()
	mux.HandleFunc(prefix+routeDirectoryRead, httpHandler(cfg, HttpGET, handlerDirectoryRead))
	mux.HandleFunc(prefix+routeDirectoryTree, httpHandler(cfg, HttpGET, handlerDirectoryTree))
	mux.HandleFunc(prefix+routeDirectoryList, httpHandler(cfg, HttpGET, handlerDirectoryList))
	mux.HandleFunc(prefix+routeDirectoryCreate, httpHandler(cfg, HttpPOST, handlerDirectoryCreate))
	mux.HandleFunc(prefix+routeDirectoryDelete, httpHandler(cfg, HttpDELETE, handlerDirectoryDelete))
	mux.HandleFunc(prefix+routeFileRead, httpHandler(cfg, HttpGET, handlerFileRead))
	mux.HandleFunc(prefix+routeFileWrite, httpHandler(cfg, HttpPOST, handlerFileWrite))
	mux.HandleFunc(prefix+routeFileAppend, httpHandler(cfg, HttpPUT, handlerFileAppend))
	mux.HandleFunc(prefix+routeFileDelete, httpHandler(cfg, HttpDELETE, handlerFileDelete))
	mux.HandleFunc(prefix+routeFileExists, httpHandler(cfg, HttpGET, handlerFileExists))
	mux.HandleFunc(prefix+routeFileChecksum, httpHandler(cfg, HttpGET, handlerFileChecksum))
	mux.HandleFunc(prefix+routeFileShared, httpHandler(cfg, HttpGET, handlerFileShared))
	// display configuration summary
	cfg.DisplaySummary()
	// start the server
	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", cfg.ServerPort), Handler: mux}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			errMsg := strings.ToLower(err.Error())
			if !strings.Contains(errMsg, "server") && !strings.Contains(errMsg, "closed") {
				log.Fatal(err)
			}
		}
	}()
	return httpServer
}

// StopServer gracefully stops the server.
func StopServer(httpServer *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
