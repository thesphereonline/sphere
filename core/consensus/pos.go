package consensus

var validators = []string{
	"0xValidator1", "0xValidator2", "0xValidator3",
}

var index = 0

func SelectValidator() string {
	validator := validators[index%len(validators)]
	index++
	return validator
}
