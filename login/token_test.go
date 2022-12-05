package login

import (
	"testing"
)

func TestMakeToken_DefaultSizedTokenShouldBeUniqueAmongAMillionUsers(t *testing.T) {
	var tokens = make(map[Token]interface{})
	generator := NewTokenGenerator()
	//Testing can be done for larger numbers, but it will take a lot of time
	for i := 0; i < 1000000; i++ {
		newToken := generator.NewToken(42)
		if _, exists := tokens[newToken]; exists {
			t.Fatalf("repeated token on : %s", newToken)
		}
		tokens[newToken] = nil
	}

}
