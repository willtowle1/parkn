package errs

import "fmt"

func WrapError(message string, err error) error { return fmt.Errorf("%s: %w", message, err) }

type ApiErr struct {
	Status  int                    `json:"status"`
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Detail  map[string]interface{} `json:"detail"`
}

func NewApiError(status int, code, message string, detail ...interface{}) *ApiErr {

	detailMap := make(map[string]interface{})
	n := len(detail)
	for i := 0; i < n; i += 2 {
		key := detail[i]
		var value interface{} = "MISSING"
		if i+1 < n {
			value = detail[i+1]
		}
		detailMap[fmt.Sprintf("%v", key)] = value
	}

	return &ApiErr{
		Status:  status,
		Code:    code,
		Message: message,
		Detail:  detailMap,
	}
}
