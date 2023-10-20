package server

import "fmt"

const (
	version = "0.0.8" // File server version.
)

// Configuration structure stores all configuration options available for the server.
// ServerPort defines the port number the file server will be listening for incoming API requests.
// RootDirectory defines the directory that is a parent for all other directories and files stored
// on server. All directory or file names used in API calls should be relative to root directory.
// UrlPrefix defines the prefix that will be prepended to all API endpoints.
type Configuration struct {
	ServerPort    int    `json:"serverPort"`    // Port number on which the server will be waiting for requests.
	RootDirectory string `json:"rootDirectory"` // Name of the root directory, where all content will be stored.
	UrlPrefix     string `json:"urlPrefix"`     // Prefix that will be prepended to all API endpoints.
}

// DisplaySummary prints current configuration setting to standard output.
func (c *Configuration) DisplaySummary() {
	fmt.Printf("Tarolas - the lightweight file server v%s\n", version)
	fmt.Printf("  configuration:\n")
	fmt.Printf("    - port           : %d\n", c.ServerPort)
	fmt.Printf("    - root directory : %s\n", c.RootDirectory)
	urlPrefix := c.UrlPrefix
	if urlPrefix == "" {
		urlPrefix = "(none)"
	}
	fmt.Printf("    - URL prefix     : %s\n", urlPrefix)
}
