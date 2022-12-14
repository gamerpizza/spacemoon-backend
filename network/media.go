package network

import (
	"fmt"
	"io"
	"spacemoon/network/post"
)

type MediaFileContentAdder interface {
	SaveNewPostWithMedia(post.Post, map[string]io.Reader) error
}

// NewMediaContentManager creates a new instance of a MediaFileContentAdder
// with a given Persistence and MediaFilePersistence
func NewMediaContentManager(p Persistence, mp MediaFilePersistence) MediaFileContentAdder {
	return mediaFileContentManager{postPersistence: p, mediaFilePersistence: mp}
}

// MediaFilePersistence defines how this package saves media files, represented in the more generic io.Reader form
type MediaFilePersistence interface {
	SaveFiles(files map[string]io.Reader, prefix string) (post.ContentURIS, error)
	GetFile(uri string) (io.Reader, error)
	Delete(uri string) error
}

type mediaFileContentManager struct {
	postPersistence      Persistence
	mediaFilePersistence MediaFilePersistence
}

func (cm mediaFileContentManager) SaveNewPostWithMedia(p post.Post, f map[string]io.Reader) error {
	var urls, err = cm.mediaFilePersistence.SaveFiles(f, prefix+string(p.GetId())+"/")
	if err != nil {
		return fmt.Errorf("could not save files: %w", err)
	}
	p.URLS = urls
	err = cm.postPersistence.AddPost(p)
	if err != nil {
		return fmt.Errorf("could not save post with media URLs: %w", err)
	}
	return nil
}

const prefix = "media/"
