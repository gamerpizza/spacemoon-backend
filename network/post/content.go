package post

type Content struct {
	URLS ContentURLS `json:"URLS"`
}

func (c Content) GetURLS() ContentURLS {
	return c.URLS
}

type ContentURLS map[ContentURI]interface{}

func (u ContentURLS) Is(url ContentURI) Verifier {
	return Verifier{urls: u, url: url}
}

type ContentURI string

type Verifier struct {
	url  ContentURI
	urls ContentURLS
}

func (v Verifier) NotPresent() bool {
	_, exists := v.urls[v.url]
	return !exists
}
