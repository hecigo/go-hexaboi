package core

type HPIResult struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	ErrorCode string      `json:"error_code"`
}

// core.HPIResult makes it compatible with the `error` interface.
func (r *HPIResult) Error() string {
	if r.Status >= 400 {
		return r.Message
	}
	return ""
}
