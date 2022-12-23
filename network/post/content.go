package post

type Content struct {
	URLS ContentURIS `json:"URLS"`
}

func (c Content) GetURLS() ContentURIS {
	return c.URLS
}

type ContentURIS map[string]bool

func (u ContentURIS) Is(url string) Verifier {
	return Verifier{urls: u, url: url}
}

type Verifier struct {
	url  string
	urls ContentURIS
}

func (v Verifier) NotPresent() bool {
	_, exists := v.urls[v.url]
	return !exists
}
