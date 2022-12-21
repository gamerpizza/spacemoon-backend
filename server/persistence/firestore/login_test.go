package firestore

import (
	"context"
	"spacemoon/login"
	"testing"
)

func TestLoginPersistence(t *testing.T) {
	ctx := context.TODO()
	var lp login.Persistence
	lp, err := GetPersistence(ctx)
	defer func(persistence *fireStorePersistence) {
		err := persistence.Close()
		if err != nil {
			t.Fatal(err.Error())
		}
	}(lp.(*fireStorePersistence))
	if err != nil {
		t.Fatal(err.Error())
	}

	const testUser = "test-user"
	const testPassword = "test-pass"
	err = lp.SignUpUser(testUser, testPassword)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !lp.ValidateCredentials(testUser, testPassword) {
		t.Fatal("expected user not found")
	}

	const testToken = "test-token"
	err = lp.SetUserToken(testUser, testToken, login.DefaultTokenDuration)
	if err != nil {
		t.Fatal(err.Error())
	}
	username, err := lp.GetUser(testToken)
	if err != nil {
		t.Fatal(err.Error())
	}
	if username != testUser {
		t.Fatal("retrieved user does not match expected user")
	}

	err = lp.DeleteUser(testUser)
	if err != nil {
		t.Fatal(err.Error())
	}
	if lp.ValidateCredentials(testUser, testPassword) {
		t.Fatal("unexpected user found")
	}
	err = lp.(*fireStorePersistence).DeleteToken(testToken)
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = lp.GetUser(testToken)
	if err == nil {
		t.Fatal("did not erase token")
	}
}
