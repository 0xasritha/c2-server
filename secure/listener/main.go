package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

/* https://www.youtube.com/watch?v=7K0gAYmWWho&list=PLShDm2AZYnK3cWZpOjV7nOpL7plH2Ztz0&index=2&ab_channel=kubucation
create self-signed SSL/TLS X.509 certificate: `openssl req -x509 -newkey rsa:4096 -sha256 -nodes -keyout tls.key -out tls.crt -subj "/CN=localhost" -days 365`
*/

var (
	CertFilePath = "./tls.crt"
	KeyFilePath  = "./tls.key"
)

type Implant struct {
	ID uint32
	// uuid ?
	LastSeen     time.Time
	PublicKey    string
	IP           string
	Status       string // enum?
	NonceCounter int
}

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hello": "world",
		})
	})

	/*
		func LoginHandler(w http.ResponseWriter, r *http.Request) {
		  if u.Username == "Chek" && u.Password == "123456" {
		    tokenString, err := CreateToken(u.Username)
		    if err != nil {
		       w.WriteHeader(http.StatusInternalServerError)
		       fmt.Errorf("No username found")
		     }
		    w.WriteHeader(http.StatusOK)
		    fmt.Fprint(w, tokenString)
		    return
		}
	*/

	// TODO: get better name: /auth?
	// getToken endpoint where implants can authenticate (provide their ID) + obtain a JWT token
	router.POST("/getToken", func(c *gin.Context) {
		// #TODO: make it so that only already existing implants get same JWT? have to already exist to get JWT?
		var implant Implant
		if err := c.ShouldBindJSON(&implant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // invalid implant?
		}
		if doesIDExist(implant.ID) {
			token, err := createToken(implant.ID)
			if err != nil {
				c.Header("Authorization", "Bearer "+token) // - i think sending server sending in JSON body to client is the right way to do it
				c.JSON(http.StatusOK, gin.H{
					//		"token": token,
				})
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized, invalid or no ID"})
		// invalid credentials
	})

	// protected-beacon will require a valid JWT token for access
	router.GET("/protected-beacon", func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		values := strings.Split(authHeader, " ")
		if len(values) == 2 {
			if values[0] == "Bearer" {
				token := values[1]
				err := verifyToken(token)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
				}
				c.JSON(http.StatusOK, gin.H{
					"PROTECTED-TASKS": "TASKJDLFAKSJFDLJAS",
				})
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT must be specified"})
	})

	router.RunTLS(":443", CertFilePath, KeyFilePath)
}

var secretKey = []byte("secret")

func doesIDExist(id uint32) bool {
	return true
}

// createToken will generate a new JWT token
func createToken(id uint32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // signing method is HS256: HMAC using SHA-256. symmetric signing algo used for JWTs (Hash-based Message Authentication Code) - uses a shared secret key + SHA-256 to generate a cryptographic SIGNATURE
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// verifyToken will verify its authenticity before granting access to protected resources
func verifyToken(tokenString string) error {
	// jwt.Parse parses and verifies the token (using the secret key)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { // callback function to retrieve the secret key used for signing the token
		return secretKey, nil
	})

	if err != nil { // token is invalid
		return fmt.Errorf("Token is invalid! %v", err)
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}
