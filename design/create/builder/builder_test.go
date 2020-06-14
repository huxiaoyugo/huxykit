package builder

import "testing"

func TestBuilderPattern(t *testing.T) {
	CreateDirector(BuilderA{})
	CreateDirector(BuilderB{})
}
