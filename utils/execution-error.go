package utils

const (
	UKNOWN_ERROR = 1
	HOMEDIR_DETECTION_ERROR = 2
)

type ExecutionError struct {
	Code int
	FriendlyMessage string
	OriginalError error
}

func (eerr ExecutionError) Error() string {
	if eerr.FriendlyMessage != "" {
		return eerr.FriendlyMessage
	}

	return eerr.OriginalError.Error()
}

func (eerr ExecutionError) Unwrap() error {
	return eerr.OriginalError
}
