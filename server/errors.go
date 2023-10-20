package server

import "log"

const (
	errMsgCheckServerLogForDetails = "check server log for details"
)

var (
	errRequestMethodNotSupported       = ErrorDto{"400", "10001", "request method not supported", ""}
	errRequiredParameterIsMissing      = ErrorDto{"400", "10837", "required parameter is missing", ""}
	errRequiredParameterIsEmpty        = ErrorDto{"400", "10767", "required parameter is empty", ""}
	errRequiredParameterIsNotAnInteger = ErrorDto{"400", "10767", "required parameter is not an integer", ""}
	errNoSlashInFileOrDirectoryName    = ErrorDto{"400", "10767", "file or directory name must begin with slash", ""}
	errOnlyOneParameterAllowed         = ErrorDto{"400", "10364", "only one parameter allowed", ""}
	errReadingDirectoryContentFailed   = ErrorDto{"400", "10147", "reading directory content failed", ""}
	errDirectoryAlreadyExists          = ErrorDto{"400", "10237", "directory already exists", ""}
	errCreatingDirectoryFailed         = ErrorDto{"400", "10152", "creating directory failed", ""}
	errCreatingDirectoriesFailed       = ErrorDto{"400", "10141", "creating directories failed", ""}
	errDeletingDirectoryFailed         = ErrorDto{"400", "10371", "deleting directory failed", ""}
	errDeletingFileFailed              = ErrorDto{"400", "10923", "deleting file failed", ""}
	errOpeningFileForAppendFailed      = ErrorDto{"400", "10347", "opening file for append failed", ""}
	errOpeningFileForWritingFailed     = ErrorDto{"400", "10747", "opening file for writing failed", ""}
	errOpeningFileForReadingFailed     = ErrorDto{"400", "10352", "opening file for reading failed", ""}
	errRetrievingFileInfoFailed        = ErrorDto{"400", "10489", "retrieving file info failed", ""}
	errAppendingFileFailed             = ErrorDto{"400", "10429", "appending file failed", ""}
	errWritingFileFailed               = ErrorDto{"400", "10451", "writing file failed", ""}
	errReadingFileFailed               = ErrorDto{"400", "10455", "reading file failed", ""}
	errSeekingFileFailed               = ErrorDto{"400", "10458", "seeking file failed", ""}
	errFileNotFound                    = ErrorDto{"400", "10938", "file not found", ""}
	errCalculatingChecksumFailed       = ErrorDto{"400", "10956", "calculating checksum failed", ""}
	errNotAFile                        = ErrorDto{"400", "10462", "not a file", ""}
	errInvalidParameterValue           = ErrorDto{"400", "10901", "invalid parameter value", ""}
	errWalkingDirectoryTreeFailed      = ErrorDto{"400", "10177", "walking directory tree failed", ""}
)

type ErrorDto struct {
	Status string `json:"status"  api:"The HTTP status code applicable to this problem."`
	Code   string `json:"code"    api:"An application-specific error code."`
	Title  string `json:"title"   api:"A short, human-readable summary of the problem that SHOULD NOT change from occurrence to occurrence of the problem."`
	Detail string `json:"detail"  api:"A human-readable explanation specific to this occurrence of the problem."`
}

type ErrorsDto struct {
	Errors []ErrorDto `json:"errors"  api:"List of encountered problems"`
}

// TODO add documentation
func errorDto(errorDto ErrorDto, detail string) *ErrorDto {
	return &ErrorDto{
		Status: errorDto.Status,
		Code:   errorDto.Code,
		Title:  errorDto.Title,
		Detail: detail}
}

// TODO add documentation
func errorsDto(errorDto *ErrorDto) *ErrorsDto {
	return &ErrorsDto{Errors: []ErrorDto{*errorDto}}
}

// TODO add documentation
func logError(err error) {
	_ = log.Output(2, err.Error())
}
