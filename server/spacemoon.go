package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"spacemoon/login"
	"spacemoon/network/profile"
	"spacemoon/network/profile/handler"
	"spacemoon/server/cors"
	"spacemoon/server/persistence/firestore"
	"strings"
	"time"
)

func main() {
	log.Default().Print("starting spacemoon server 🚀")
	log.Default().Print("v1.0.3")
	setupHandlers()
	listenAndServe()
}

func setupHandlers() {
	log.Default().Print("registering server handlers...")

	loginPersistence := getLoginPersistence()
	corsEnabledLoginHandler := cors.EnableCors(login.NewHandler(loginPersistence, time.Hour), http.MethodGet, http.MethodPost)
	http.Handle("/login", corsEnabledLoginHandler)

	protector := login.NewProtector(loginPersistence)

	mediaFilePersistence, err := getMediaFilePersistence(context.Background())
	if err != nil {
		panic(err)
	}
	socialNetworkPersistence := getSocialNetworkPersistence()
	setupSocialNetworkHandler(socialNetworkPersistence, loginPersistence, mediaFilePersistence, protector)
	setupProductHandlers(loginPersistence, protector)
	setupCategoryHandler(protector)

	profilePersistence, err := getProfilePersistence(context.Background())
	if err != nil {
		panic(err)
	}
	profileHandler := handler.New(profilePersistence, loginPersistence)
	protectedProfileHandler := protector.Protect(&profileHandler)
	protectedProfileHandler.Unprotect(http.MethodGet)
	corsEnabledProtectedProfileHandler := cors.EnableCors(protectedProfileHandler, http.MethodGet, http.MethodPut)
	http.Handle("/profile", corsEnabledProtectedProfileHandler)

	log.Default().Print("handler registration done, ready for takeoff")
}

func getProfilePersistence(ctx context.Context) (profile.Persistence, error) {
	creds := os.Getenv(googleCredentials)
	if strings.TrimSpace(creds) == "" {
		return nil, errors.New("no google credentials file set")
	}
	persistence, err := firestore.GetPersistence(ctx)
	if err != nil {
		return nil, err
	}
	return persistence, nil
}

func listenAndServe() {
	log.Default().Print(getRandomSpaceQuote())
	log.Default().Printf("listening on port %d", port)
	certFile := os.Getenv("cert_file")
	keyFile := os.Getenv("key_file")
	if strings.TrimSpace(certFile) == "" || strings.TrimSpace(keyFile) == "" {
		log.Default().Print("serving without TLS\n")
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			log.Default().Fatalf("error while performing listen and serve: %s", err.Error())
		}
		return
	}
	log.Default().Print("serving using TLS\n")
	err := http.ListenAndServeTLS(fmt.Sprintf(":%d", port), certFile, keyFile, nil)
	if err != nil {
		log.Default().Fatalf("error while performing listen and serve: %s", err.Error())
	}
}

const port = 1234
