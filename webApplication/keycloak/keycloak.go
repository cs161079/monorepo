package keycloak

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
)

func NewKeycloakService() *oidc.IDTokenVerifier {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, os.Getenv("keycloak.url"))
	if err != nil {
		panic(err)
	}
	verifier := provider.Verifier(&oidc.Config{
		ClientID:          os.Getenv("keycloak.client_id"),
		SkipClientIDCheck: true,
	})

	return verifier
}

func AuthMiddleware(verifier *oidc.IDTokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		idToken, err := verifier.Verify(ctx, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token verification failed"})
			return
		}

		var claims struct {
			Email       string `json:"email"`
			RealmAccess struct {
				Roles []string `json:"roles"`
			} `json:"realm_access"`
		}

		if err := idToken.Claims(&claims); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Pass info to next handlers
		c.Set("email", claims.Email)
		c.Set("roles", claims.RealmAccess.Roles)
		c.Next()
	}
}

// RequireRole ensures the user has the required role
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesIface, exists := c.Get("roles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No roles found"})
			return
		}

		roles := rolesIface.([]string)
		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: missing role " + role})
	}
}
