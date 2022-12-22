package bucket

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"math/rand"
	"spacemoon/network"
	"spacemoon/network/post"
	"strings"
	"time"
)

func New(ctx context.Context) (network.MediaFilePersistence, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate cloud storage client: %w", err)
	}
	bucket := client.Bucket("spacemoon-posted-media")
	return persistence{client: client, bucket: bucket, ctx: ctx}, nil
}

type FileName string

type persistence struct {
	client *storage.Client
	bucket *storage.BucketHandle
	ctx    context.Context
}

func (p persistence) Delete(uri string) error {
	err := p.bucket.Object(string(uri)).Delete(p.ctx)
	if err != nil {
		return fmt.Errorf("could not delete from storage: %w", err)
	}
	return nil
}

func (p persistence) SaveFiles(files map[string]io.Reader, prefix string) (post.ContentURIS, error) {
	urls := post.ContentURIS{}
	generator := newUriGenerator()

	for name, file := range files {
		fileName := strings.Split(name, ".")
		fileExt := strings.ToLower(fileName[len(fileName)-1])
		u := prefix + string(generator.newToken()) + "." + fileExt

		w := p.bucket.Object(string(u)).NewWriter(p.ctx)

		_, err := io.Copy(w, file)
		if err != nil {
			return nil, fmt.Errorf("could not copy file to cloud: %w", err)
		}
		err = w.Close()
		if err != nil {
			return nil, fmt.Errorf("could not close file on cloud: %w", err)
		}
		urls[u] = true
	}
	return urls, nil

}

func (p persistence) GetFile(uri string) (io.Reader, error) {
	reader, err := p.bucket.Object(uri).NewReader(p.ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create object reader: %w", err)
	}
	return reader, nil
}

func newUriGenerator() urlGenerator {
	rand.Seed(time.Now().Unix())
	return urlGenerator{}
}

type urlGenerator struct {
}

type uri string

func (t urlGenerator) newToken() uri {
	size := 15

	b := make([]byte, size)
	for i := range b {
		b[i] = urlChars[rand.Intn(len(urlChars))]
	}
	return uri(b)
}

const urlChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-"
