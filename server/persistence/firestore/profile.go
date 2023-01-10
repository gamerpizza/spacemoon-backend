package firestore

import (
	"fmt"
	"spacemoon/network/profile"
	"strings"
)

func (p *fireStorePersistence) GetProfile(id profile.Id) (profile.Profile, error) {
	collection := p.storage.Collection(profilesCollection)
	doc, err := collection.Doc(string(id)).Get(p.ctx)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return profile.Profile{}, NotFoundError
		}
		return profile.Profile{}, fmt.Errorf("could not read from firestore: %w", err)
	}
	pr := profile.Profile{}
	err = doc.DataTo(&pr)
	if err != nil {
		return profile.Profile{}, fmt.Errorf("could not parse data from persistence into profile: %w", err)
	}
	return pr, nil
}

func (p *fireStorePersistence) SaveProfile(pr profile.Profile) error {
	collection := p.storage.Collection(profilesCollection)
	_, err := collection.Doc(string(pr.Id)).Set(p.ctx, pr)
	if err != nil {
		return fmt.Errorf("could not save to firestore: %w", err)
	}
	return nil
}
