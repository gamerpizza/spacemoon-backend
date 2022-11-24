package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"moonspace/model"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v11"
	"github.com/ajg/form"
	"github.com/gin-gonic/gin"
)

func DecodeRequestBody(body io.ReadCloser, result any) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, result)
	if err != nil {
		return err
	}

	return nil
}

func DecodeFormRequest(body io.ReadCloser, result any) error {
	d := form.NewDecoder(body)
	if err := d.Decode(result); err != nil {
		return err
	}

	return nil
}

const (
	timeout           = time.Second * 15
	OAuthClaimSubject = "sub"
	TokenHeader       = "Authorization"
)

var (
	errInvalidToken = errors.New("invalid token")
	errNoClaims     = errors.New("empty claims while decoding token")
	errClaimEmpty   = errors.New("claim doesn't exist")
)

type ClaimsExtractor interface {
	Extract(token, claim string) (model.UserID, error)
}

type claimsExtractorKeycloak struct {
	cli   gocloak.GoCloak
	realm string
}

func NewClaimsExtractorKeycloak(cli gocloak.GoCloak, realm string) claimsExtractorKeycloak {
	return claimsExtractorKeycloak{
		cli,
		realm,
	}
}

func (cek *claimsExtractorKeycloak) Extract(token, claim string) (model.UserID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	t, mc, err := cek.cli.DecodeAccessToken(ctx, token, cek.realm)
	if err != nil {
		return model.UserID("0"), err
	}

	if !t.Valid {
		return model.UserID("0"), errInvalidToken
	}

	if mc == nil {
		return model.UserID("0"), errNoClaims
	}

	x, ok := (*mc)[claim]
	if !ok {
		return model.UserID("0"), errClaimEmpty
	}

	return model.UserID(x.(string)), nil
}

type GrantType struct {
	GrantType string `json:"grant_type"`
}

func BasicAuthorization(id, secret, authUrl string, submitForm bool, tokenData any) error {
	str := id + ":" + secret
	basicAuth := base64.StdEncoding.EncodeToString([]byte(str))
	var body io.Reader
	if submitForm {
		form := url.Values{}
		form.Add("grant_type", "client_credentials")
		body = strings.NewReader(form.Encode())
	} else {
		gt := GrantType{
			GrantType: "client_credentials",
		}
		bb, err := json.Marshal(gt)
		if err != nil {
			return err
		}

		body = bytes.NewReader(bb)
	}

	request, err := http.NewRequest(http.MethodPost, authUrl, body)
	if err != nil {
		return err
	}

	request.Header.Add(TokenHeader, "Basic "+string(basicAuth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.DefaultClient
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return fmt.Errorf("Server responded with an error: %s", resp.Status)
	}

	return DecodeRequestBody(resp.Body, tokenData)
}

const ImagesPath = "images/"

func SaveImage(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
	_, err := os.Stat("images")
	if os.IsNotExist(err) {
		os.Mkdir("images", os.ModePerm)
	}

	dst := ImagesPath + filepath.Base(file.Filename)
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}

	return ImagesPath + file.Filename, nil
}
