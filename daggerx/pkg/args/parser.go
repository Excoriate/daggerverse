package args

import "strings"

// ParseArgsFromStrToSlice parses a string of arguments in the form "arg1, arg2, arg3"
func ParseArgsFromStrToSlice(argStr string) []string {
	if argStr == "" {
		return nil
	}

	var parsedArgs []string
	args := strings.Split(argStr, ",")
	for _, arg := range args {
		trimmedArg := strings.TrimSpace(arg)
		if trimmedArg != "" {
			parsedArgs = append(parsedArgs, trimmedArg)
		}
	}
	if len(parsedArgs) == 0 {
		return nil
	}
	return parsedArgs
}
