package generate

// This file contains some of the constants used in go code generation

var RegisteredTypes = [...]string{"User", "Guild", "Channel", "Message", "Role"}

func IsRegisteredType(typ string) bool {
	for _, compare := range RegisteredTypes {
		if compare == typ {
			return true
		}
	}

	return false
}
