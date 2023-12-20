package jerror

import (
	"fmt"
)

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {

	if where != "" {
		where = fmt.Sprintf(" %s", where)
	}

	fmt.Printf("Error%s: %s [line %d] \n", where, message, line)
}

func RuntimeError(line int, message string) {
	fmt.Printf("Error: %s [line %d]\n", message, line)
}
