package utils

type ErrorMsg struct {
	Field   string
	Message string
}

func MsgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return ""
}
