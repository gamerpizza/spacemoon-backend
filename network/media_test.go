package network

import (
	"bytes"
	"io"
	"os"
	"spacemoon/login"
	"spacemoon/network/post"
	"strings"
	"testing"
)

func TestCreatePostWithMedia(t *testing.T) {
	p := &fakePersistence{}
	mp := &fakeMediaFilePersistence{}
	cm := NewMediaContentManager(p, mp)
	pst := post.Post{
		Caption: "this post is a test",
		Author:  "mr Test PostManager",
		Id:      postId,
	}
	f1, f2, err := getTestImageFiles(t)
	if err != nil {
		t.Fatal(err.Error())
	}
	media := make(map[string]io.Reader)
	media[f1.Name()] = f1
	media[f2.Name()] = f2
	err = cm.SaveNewPostWithMedia(pst, media)
	if err != nil {
		t.Fatal(err.Error())
	}
	posts, err := p.GetAllPosts()
	if err != nil {
		t.Fatal(err.Error())
	}
	retrievedPost, exists := posts[postId]
	if !exists {
		t.Fatal("did not retrieve expected post")
	}
	for url, _ := range retrievedPost.Content().GetURLS() {
		if !strings.Contains(string(url), prefix) {
			t.Fatal("prefix not found on URL")
		}
		file, err := mp.GetFile(url)
		if err != nil {
			t.Fatal(err.Error())
		}
		fileContents, _ := io.ReadAll(file)
		f1Contents, _ := io.ReadAll(f1)
		f2Contents, _ := io.ReadAll(f2)
		if bytes.Compare(fileContents, f1Contents) != 0 && bytes.Compare(fileContents, f2Contents) != 0 {
			t.Fatalf("invalid urls: %+v", url)
		}
	}
}

func getTestImageFiles(t *testing.T) (*os.File, *os.File, error) {
	f, err := os.Open("test.jpg")
	if err != nil {
		t.Fatal(err.Error())
	}
	f2, err := os.Open("test2.jpg")
	if err != nil {
		t.Fatal(err.Error())
	}
	return f, f2, err
}

type fakeMediaFilePersistence struct {
	files map[string]io.Reader
}

func (s *fakeMediaFilePersistence) GetFile(url string) (io.Reader, error) {
	return s.files[url], nil
}

func (s *fakeMediaFilePersistence) Delete(url string) error {
	//TODO implement me
	panic("implement me")
	return nil
}

func (s *fakeMediaFilePersistence) SaveFiles(files map[string]io.Reader, prefix string) (post.ContentURIS, error) {
	urls := post.ContentURIS{}
	generator := login.NewTokenGenerator()
	if s.files == nil {
		s.files = make(map[string]io.Reader)
	}
	for name, file := range files {
		fileName := strings.Split(name, ".")
		fileExt := fileExtension(strings.ToLower(fileName[len(fileName)-1]))
		if _, isAccepted := acceptedFileTypes[fileExt]; !isAccepted {
			continue
		}
		url := prefix + string(generator.NewToken(8)) + string(fileExt)
		s.files[url] = file
		urls[url] = true
	}
	return urls, nil
}

type fakePersistence struct {
	posts post.Posts
}

func (f *fakePersistence) CheckIfPostExists(id post.Id) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (f *fakePersistence) DeletePost(post.Id) error {
	//TODO implement me
	panic("implement me")
	return nil
}

func (f *fakePersistence) AddPost(p post.Post) error {
	if f.posts == nil {
		f.posts = map[post.Id]post.Post{}
	}
	f.posts[p.GetId()] = p
	return nil
}

func (f *fakePersistence) GetAllPosts() (post.Posts, error) {
	if f.posts == nil {
		f.posts = post.Posts{postId: post.Post{
			Caption: "",
			Author:  "",
			URLS:    nil,
			Id:      postId,
		}}
	}

	return f.posts, nil
}

const postId = "post id"

type fileExtension string

var acceptedFileTypes = make(map[fileExtension]interface{})

func init() {
	acceptedFileTypes["jpg"] = nil
}
