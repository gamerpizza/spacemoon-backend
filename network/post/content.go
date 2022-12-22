package post

type Content struct {
	URLS ContentURIS `json:"URLS"`
}

func (c Content) GetURLS() ContentURIS {
	return c.URLS
}

type ContentURIS map[ContentURI]bool

func (u ContentURIS) Is(url ContentURI) Verifier {
	return Verifier{urls: u, url: url}
}

type ContentURI string

type Verifier struct {
	url  ContentURI
	urls ContentURIS
}

func (v Verifier) NotPresent() bool {
	_, exists := v.urls[v.url]
	return !exists
}
