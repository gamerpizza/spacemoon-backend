package handler

import (
	"errors"
	"io"
	"spacemoon/login"
	"spacemoon/network/post"
	"time"
)

type failNetworkPersistence struct {
}

func (f failNetworkPersistence) DeletePost(post.Id) error {
	//TODO implement me
	panic("implement me")
	return nil
}

func (f failNetworkPersistence) AddPost(_ post.Post) error {
	return errors.New("some fake error")
}

func (f failNetworkPersistence) GetAllPosts() (post.Posts, error) {
	return nil, errors.New("some fake error")
}

type failLoginPersistence struct {
}

func (f failLoginPersistence) SetUserToken(_ login.UserName, _ login.Token, _ time.Duration) error {
	return fakeError
}

func (f failLoginPersistence) GetUser(_ login.Token) (login.UserName, error) {
	return "", fakeError
}

func (f failLoginPersistence) SignUpUser(u login.UserName, p login.Password) error {
	return fakeError
}

func (f failLoginPersistence) ValidateCredentials(u login.UserName, p login.Password) bool {
	//TODO implement me
	panic("implement me")
}

func (f failLoginPersistence) DeleteUser(name login.UserName) error {
	return fakeError
}

func (f failLoginPersistence) Check(name login.UserName) (bool, error) {
	return false, fakeError
}

type stubPersistence struct {
}

func (s stubPersistence) DeletePost(post.Id) error {

	return nil
}

func (s stubPersistence) AddPost(_ post.Post) error {
	return nil
}

func (s stubPersistence) GetAllPosts() (post.Posts, error) {
	return expectedPosts, nil
}

type mockNetworkPersistence struct {
	posts post.Posts
}

func (m *mockNetworkPersistence) DeletePost(p post.Id) error {
	delete(m.posts, p)
	return nil
}

func (m *mockNetworkPersistence) AddPost(p post.Post) error {
	if m.posts == nil {
		m.posts = make(post.Posts)
	}
	m.posts[p.GetId()] = p
	return nil
}

func (m *mockNetworkPersistence) GetAllPosts() (post.Posts, error) {
	return m.posts, nil
}

type stubLoginPersistence struct {
}

func (f stubLoginPersistence) Check(_ login.UserName) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (f stubLoginPersistence) SetUserToken(_ login.UserName, _ login.Token, _ time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (f stubLoginPersistence) GetUser(_ login.Token) (login.UserName, error) {
	return testAuthor, nil
}

func (f stubLoginPersistence) SignUpUser(u login.UserName, p login.Password) error {
	//TODO implement me
	panic("implement me")
}

func (f stubLoginPersistence) ValidateCredentials(u login.UserName, p login.Password) bool {
	//TODO implement me
	panic("implement me")
}

func (f stubLoginPersistence) DeleteUser(name login.UserName) error {
	//TODO implement me
	panic("implement me")
}

type stubMediaFilePersistence struct {
}

func (s stubMediaFilePersistence) SaveFiles(_ map[string]io.Reader, _ string) (post.ContentURIS, error) {
	return nil, nil
}

func (s stubMediaFilePersistence) GetFile(uri string) (io.Reader, error) {
	//TODO implement me
	panic("implement me")
}

func (s stubMediaFilePersistence) Delete(uri string) error {
	//TODO implement me
	panic("implement me")
}
