package profile

func New(i Id, u UserName, m Motto, a AvatarUrl) Profile {
	return Profile{Id: i, UserName: u, Motto: m, Avatar: Avatar{Url: a}}
}

type Profile struct {
	Id       Id       `json:"id"`
	UserName UserName `json:"user_name"`
	Motto    Motto    `json:"motto"`
	Avatar   Avatar   `json:"avatar"`
}

type Avatar struct {
	Url AvatarUrl `json:"url"`
}

type Id string

type UserName string

func (n UserName) String() string {
	return string(n)
}

type Motto string

func (m Motto) String() string {
	return string(m)
}

type AvatarUrl string

func (u AvatarUrl) String() string {
	return string(u)
}
