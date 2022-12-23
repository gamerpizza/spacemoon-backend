package bucket

import (
	"bytes"
	"context"
	"io"
	"os"
	"spacemoon/network"
	"testing"
)

func TestName(t *testing.T) {
	f, err := os.Open("test.jpg")
	if err != nil {
		t.Fatal(err.Error())
	}
	const fileName = "filename"
	var p network.MediaFilePersistence
	p, err = New(context.Background())
	if err != nil {
		t.Fatal(err.Error())
	}

	files := make(map[string]io.Reader)
	files[f.Name()] = f
	urlsMap, _ := p.SaveFiles(files, "/test/")
	var urls []string
	for url := range urlsMap {
		urls = append(urls, url)
	}

	if len(urls) == 0 {
		t.Fatal("did not retrieve any uri")
	}
	var retrievedFile io.Reader
	retrievedFile, _ = p.GetFile(urls[0])
	if err != nil {
		t.Fatal(err.Error())
	}
	retrievedFileContents, err := io.ReadAll(retrievedFile)
	if err != nil {
		t.Fatal(err.Error())
	}
	expectedFile, err := os.Open("test.jpg")
	if err != nil {
		t.Fatal(err.Error())
	}
	expectedFileContents, err := io.ReadAll(expectedFile)
	if err != nil {
		t.Fatal(err.Error())
	}
	if bytes.Compare(retrievedFileContents, expectedFileContents) != 0 {
		t.Fatal("did not retrieve saved file")
	}

	err = p.Delete(urls[0])
	if err != nil {
		t.Fatal(err.Error())
	}
}
