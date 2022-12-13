package login

import (
	"errors"
	"math/rand"
	"time"
)

func NewTokenGenerator() TokenGenerator {
	rand.Seed(time.Now().Unix())
	return tokenGenerator{}
}

type TokenGenerator interface {
	NewToken(size int) Token
}
type tokenGenerator struct {
}

// NewToken will generate a random token using the characters defined in TokenCharacters and the Unix time as a seed for
// the random generator. TokenCharacters has a size of 65, so the Token can have 65^size possibilities and a 1/(65^size)
// probability of generating a repeated token. Tokens are supposed to be temporary, so it can be re-used later to represent
// a different user.
// Use 0 as the size to generate a Token using the DefaultTokenSize constant.
func (t tokenGenerator) NewToken(size int) Token {
	if size == 0 {
		size = DefaultTokenSize
	}

	b := make([]byte, size)
	for i := range b {
		b[i] = TokenCharacters[rand.Intn(len(TokenCharacters))]
	}
	return Token(b)

}

type Credential struct {
	Token
	TokenDetails
}
type Token string
type TokenDetails struct {
	User       UserName
	Expiration time.Time
}

type Tokens map[Token]TokenDetails

// TokenCharacters defines the characters that can be used to generate a token
const TokenCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_~"

// DefaultTokenSize is "The Answer to life, the Universe and Everything", it has space for around 1.39*10^76
// possibilities, so it should not repeat itself
const DefaultTokenSize = 42

// UseDefaultSize is used when generating a NewToken as an alternate way to ask for it to be of the DefaultTokenSize
const UseDefaultSize = 0

var tokenExpiredError = errors.New("token expired")
var tokenNotFoundError = errors.New("token not found")
