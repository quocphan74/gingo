package utils

type ApiError struct {
	Field string
	Msg   string
}

func MsgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "min":
		return "Minimum 3 lengths automatically."
	case "max":
		return "Exceeds the specified length of 100 characters."
	case "uniqueEmail":
		return "Email already exists."
	}

	return ""
}
