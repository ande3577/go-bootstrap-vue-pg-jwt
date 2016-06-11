package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var sessionStore sessions.Store
var jwtSigningKey []byte
var cookieIssuerVar string

type TokenData struct {
	UserId            string
	TokenString       string
	XsrfToken         string
	FromHttp          bool
	SessionIdentifier string
}

func lookupTokenSigningKey(kid interface{}) []byte {
	return jwtSigningKey
}

func Initialize(cookieStoreAuthenticationKey string, cookieStoreEncryptionKey string, jwtSigningKeyString string, cookieIssuer string) {
	if len(cookieStoreAuthenticationKey) == 0 {
		panic("Cookie store authentication key not defined")
	}

	if len(cookieStoreEncryptionKey) == 0 {
		panic("Cookie store authentication key not defined")
	}

	if len(jwtSigningKeyString) == 0 {
		panic("jwt signing key not defined")
	}

	sessionStore = sessions.NewCookieStore([]byte(cookieStoreAuthenticationKey), []byte(cookieStoreEncryptionKey))

	jwtSigningKey = []byte(jwtSigningKeyString)
	cookieIssuerVar = cookieIssuer
}

func GenerateRandomString(length int) (string, error) {
	bytes_ := make([]byte, length)
	_, err := rand.Read(bytes_)
	return base64.URLEncoding.EncodeToString(bytes_), err
}

func generateXSRFHeader() (string, error) {
	return GenerateRandomString(32)
}

func GenerateHashFromPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func CompareHashAndPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func createSessionIdentifier(userId string) (string, error) {
	if len(userId) > 0 {
		if sessionIdentifier, err := GenerateRandomString(40); err != nil {
			return "", err
		} else {
			return sessionIdentifier, nil
		}
	} else {
		return "", nil
	}
}

func Login(userId string, fromHttp bool, developmentMode bool) (tokenData *TokenData, err error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["sub"] = userId

	issuedTime := time.Now()
	token.Claims["iss"] = cookieIssuerVar
	token.Claims["iat"] = issuedTime.Unix()
	token.Claims["exp"] = issuedTime.Add(time.Hour * 72).Unix()

	sessionIdentifier, err := createSessionIdentifier(userId)
	if err != nil {
		return tokenData, err
	}
	token.Claims["session_id"] = sessionIdentifier

	var xsrfToken string
	if xsrfToken, err = generateXSRFHeader(); err != nil {
		return tokenData, err
	}

	token.Claims["xsrftoken"] = xsrfToken
	token.Claims["from_http"] = fromHttp

	if developmentMode {
		token.Claims["development_mode"] = true
	}

	var tokenString string
	if tokenString, err = token.SignedString(jwtSigningKey); err != nil {
		return tokenData, err
	}

	if xsrfToken, err = caculateXSRFHash(xsrfToken); err != nil {
		return tokenData, err
	}

	return &TokenData{FromHttp: fromHttp,
			TokenString:       tokenString,
			XsrfToken:         xsrfToken,
			UserId:            userId,
			SessionIdentifier: sessionIdentifier},
		err
}

func Logout(session *sessions.Session, developmentMode bool) (originalTokenData *TokenData) {
	// create a new token with an unauthenticated user
	if tokenData, err := Login("", true, developmentMode); err != nil {
		return &TokenData{}
	} else {
		session.Values["token"] = tokenData.TokenString
		return tokenData
	}
}

func caculateXSRFHash(xsrfToken string) (string, error) {
	// if xsrfHeaderBytes, err := bcrypt.GenerateFromPassword([]byte(xsrfToken), bcrypt.DefaultCost); err != nil {
	// 	return "", err
	// } else {
	// 	xsrfToken = base64.URLEncoding.EncodeToString(xsrfHeaderBytes)
	// }
	return xsrfToken, nil

}

func getDevelopmentModeFromToken(token *jwt.Token) bool {
	developmentModeInterface := token.Claims["development_mode"]
	if developmentModeInterface == nil {
		return false
	}
	return developmentModeInterface.(bool)
}

func ParseToken(tokenString string, fromHttp bool, developmentMode bool) (tokenData *TokenData, err error) {
	tokenData = &TokenData{TokenString: tokenString}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return lookupTokenSigningKey(token.Header["kid"]), nil
	})

	if err != nil {
		return &TokenData{}, err
	}

	if !developmentMode && getDevelopmentModeFromToken(token) {
		return nil, fmt.Errorf("access denied - cookie mismatch")
	}

	if err == nil && token.Valid {
		if userIdInterface := token.Claims["sub"]; userIdInterface != nil {
			tokenData.UserId = userIdInterface.(string)
		}
		if xsrfInterface := token.Claims["xsrftoken"]; xsrfInterface != nil {
			if tokenData.XsrfToken, err = caculateXSRFHash(xsrfInterface.(string)); err != nil {
				return &TokenData{}, err
			}
		}
		if fromHttpInterface := token.Claims["from_http"]; fromHttpInterface != nil {
			tokenData.FromHttp = fromHttpInterface.(bool)
			if tokenData.FromHttp != fromHttp {
				return &TokenData{}, errors.New("access denied - cookie http mismatch")
			}
		}
		if sessionIdentifierInterface := token.Claims["session_id"]; sessionIdentifierInterface != nil {
			tokenData.SessionIdentifier = sessionIdentifierInterface.(string)
		}
	} else {
		return &TokenData{}, err
	}

	return tokenData, err
}

func validateCSRFToken(receiveToken, jwtToken string) bool {
	if len(receiveToken) == 0 {
		return false
	} else if receiveToken != jwtToken {
		return false
	} else {
		return true
	}
}

func Authorize(r *http.Request, developmentMode bool) (session *sessions.Session, tokenData *TokenData, err error) {
	if session, err = sessionStore.Get(r, cookieIssuerVar); err != nil {
		session, _ = sessionStore.New(r, cookieIssuerVar)
	}

	// set the session to httponly
	session.Options.HttpOnly = true

	var tokenString string
	tokenStringInterface := session.Values["token"]
	if tokenStringInterface != nil {
		tokenString = tokenStringInterface.(string)
	} else {
		// if there is no token present in the cookie, create a new token with an unauthenticated user
		if tokenData, err = Login("", true, developmentMode); err != nil {
			return session, &TokenData{}, err
		}
		session.Values["token"] = tokenString
	}

	if tokenData, err = ParseToken(tokenString, true, developmentMode); err != nil {
		tokenData := Logout(session, developmentMode)
		return session, tokenData, nil
	}

	// any non-get operation requires checking XSRF protection
	if r.Method != "GET" {
		// ensure the CSRF protection token is included with the request and matches the value pulled
		// from the jwt.  only check for CSRF if user is logged in
		if !validateCSRFToken(r.FormValue("xsrf-token"), tokenData.XsrfToken) {
			return session, &TokenData{}, errors.New("access denied - csrf.")
		}
	}

	return session, tokenData, err
}

func AuthorizeJSON(r *http.Request, developmentMode bool) (tokenData *TokenData, err error) {

	fromSession := false

	// first look for a token in the secure session store
	session, err := sessionStore.Get(r, cookieIssuerVar)
	var tokenString string
	if err == nil {
		tokenStringInterface := session.Values["token"]
		if tokenStringInterface != nil {
			tokenString = tokenStringInterface.(string)
			fromSession = true
		}
	}

	// if its not there, look inside the url
	if len(tokenString) == 0 {
		tokenString = r.FormValue("auth-token")
	}

	tokenData, err = ParseToken(tokenString, fromSession, developmentMode)

	// xsrf protection is only required if the token was obtained from the browser session
	// if the user allows their jwt token to be stolen otherwise, they are on their own
	// ensure that the xsrf token in included with the request header
	if fromSession && !validateCSRFToken(r.Header.Get("X-XSRF-TOKEN"), tokenData.XsrfToken) {
		return &TokenData{}, errors.New("access denied - csrf.")
	}

	return tokenData, err
}
