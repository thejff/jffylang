package jerror

import "fmt"

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {

	if where != "" {
		where = fmt.Sprintf(" %s", where)
	}

	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
}
