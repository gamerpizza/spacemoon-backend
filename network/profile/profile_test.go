package profile

import "testing"

func TestProfile(t *testing.T) {
	const id Id = "user-id"
	const username UserName = "user-name"
	const motto Motto = "test everything"
	const avatar AvatarUrl = "avatar.jpg"
	p := New(id, username, motto, avatar)
	if p.UserName != username {
		t.Fatal("bad user name")
	}
	if p.Motto != motto {
		t.Fatal("bad motto")
	}
	if p.Avatar.Url != avatar {
		t.Fatal("bad avatar url")
	}
	if p.Id != id {
		t.Fatal("bad id")
	}
}

func TestNewReturnsAProfile(t *testing.T) {
	var _ Profile = New("", "", "", "")
}
