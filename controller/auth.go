package controller

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"github.com/haran/biophonie-api/controller/user"
)

// location of the files used for signing and verification
var (
	privKeyPath = os.Getenv("KEYS_FOLDER") + "/app.rsa"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = os.Getenv("KEYS_FOLDER") + "/app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

func (c *Controller) Authorize(ctx *gin.Context) {
	token, err := request.ParseFromRequest(ctx.Request, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		return c.verifyKey, nil
	}, request.WithClaims(&user.CustomClaims{}))

	// If the token is missing or invalid, return error
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err).SetType(gin.ErrorTypePublic)
		return
	}

	var userId int
	if err := c.Db.Get(&userId, "SELECT id FROM accounts WHERE name = $1", token.Claims.(*user.CustomClaims).Name); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeAny).SetMeta("-> could not get user for auth")
		ctx.Abort()
		return
	}

	ctx.Set("userId", userId)
	ctx.Set("admin", token.Claims.(*user.CustomClaims).Admin)
	ctx.Next()
}

func (c *Controller) AuthorizeAdmin(ctx *gin.Context) {
	admin := ctx.GetBool("admin")
	if !admin {
		ctx.AbortWithError(http.StatusUnauthorized, errors.New("restricted to admins")).SetType(gin.ErrorTypePublic)
		return
	}
	ctx.Next()
}

// read the key files before starting http handlers
func (c *Controller) readKeys() {
	signBytes, err := os.ReadFile(privKeyPath)
	fatal(err)

	c.signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := os.ReadFile(pubKeyPath)
	fatal(err)

	c.verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Controller) createToken(name string, admin bool) (string, error) {
	// create a signer for rsa 256
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	t.Claims = &user.CustomClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 365)),
		},
		UserInfo: user.UserInfo{Name: name, Admin: admin},
	}

	// create token string
	return t.SignedString(c.signKey)
}
