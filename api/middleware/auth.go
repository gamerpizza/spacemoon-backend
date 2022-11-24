package middleware

import (
	"context"
	"errors"
	"fmt"
	"moonspace/utils"
	"net/http"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v11"
	"github.com/gin-gonic/gin"
)

const decodeDeadline = time.Second * 5

var (
	errNoAuthHeader            = errors.New("No Authorization Header Provided.")
	errUnknownUser             = errors.New("Access Denied. Unknown User.")
	errExtractingId            = errors.New("Coudln't extract user id from token.")
	errSharedVariable          = errors.New("No UID in context. Auth middleware.")
	errInsufficientPermissions = errors.New("Insufficient permissions.")
)

const (
	RoleAdmin    = "admin"
	RoleCustomer = "customer"
	RoleAnon     = "anon"
)

func Authentication(keycloak gocloak.GoCloak, realm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.AbortWithError(http.StatusBadRequest, errNoAuthHeader)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), decodeDeadline)
		defer cancel()

		t, _, err := keycloak.DecodeAccessToken(ctx, token, realm)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Failed Authenticating User: Details: %w", err))
			return
		}

		if !t.Valid {
			c.AbortWithError(http.StatusUnauthorized, errUnknownUser)
			return
		}

		ex := utils.NewClaimsExtractorKeycloak(keycloak, realm)
		uid, err := ex.Extract(token, utils.OAuthClaimSubject)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errExtractingId)
		}
		c.Set("uid", string(uid))
		c.Next()
	}
}

type AdminCredentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Realm    string `yaml:"realm"`
}

type AuthorizationRoleCallback func(keycloak gocloak.GoCloak, realm, cid string) gin.HandlerFunc

func AuthorizationForRoles(ac *AdminCredentials, roles []string) AuthorizationRoleCallback {
	return func(keycloak gocloak.GoCloak, realm, cid string) gin.HandlerFunc {
		return func(c *gin.Context) {
			token := c.Request.Header.Get("Authorization")
			if token == "" {
				c.AbortWithError(http.StatusBadRequest, errNoAuthHeader)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), decodeDeadline)
			defer cancel()

			uid, exists := c.Get("uid")
			if !exists {
				c.AbortWithError(http.StatusInternalServerError, errSharedVariable)
			}

			t, err := keycloak.LoginAdmin(ctx, ac.Username, ac.Password, ac.Realm)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Couldn't login admin. Reason: %w", err))
				return
			}

			rs, err := keycloak.GetRealmRolesByUserID(ctx, t.AccessToken, realm, uid.(string))
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Couldn't get permissions for user. Reason: %w", err))
				return
			}

			rolematch := 0
			roleLen := len(roles)
			for _, r := range rs {
				for _, pr := range roles {
					if strings.Compare(pr, *r.Name) == 0 {
						rolematch += 1
						break
					}
				}
			}

			if roleLen != rolematch {
				c.AbortWithError(http.StatusUnauthorized, errInsufficientPermissions)
				return
			}
		}
	}
}
