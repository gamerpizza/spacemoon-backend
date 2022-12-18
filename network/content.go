package network

type Content struct {
	URLS PostContentURLS `json:"URLS"`
}

func (c Content) GetURLS() PostContentURLS {
	return c.URLS
}

type PostContentURLS map[PostContentURL]interface{}
