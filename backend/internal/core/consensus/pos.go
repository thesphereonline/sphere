package consensus

import "math/rand"

var validators = []string{
	"validator1",
	"validator2",
	"validator3",
}

// PickValidator selects a pseudo-random validator from the list
func PickValidator() string {
	return validators[rand.Intn(len(validators))]
}
