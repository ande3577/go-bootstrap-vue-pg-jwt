package app

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/auth"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"

	"bitbucket.org/liamstask/goose/lib/goose"
	"database/sql"
	"os"
	"path/filepath"
)

func Initialize(settings *ApplicationSettings) {
	cookieStoreAuthenticationKey := os.Getenv("GO_TEMPLATE_COOKIE_STORE_AUTHENTICATION_KEY")
	if len(cookieStoreAuthenticationKey) == 0 {
		panic("GO_TEMPLATE_COOKIE_STORE_AUTHENTICATION_KEY not defined")
	}

	cookieStoreEncryptionKey := os.Getenv("GO_TEMPLATE_COOKIE_STORE_ENCRYPTION_KEY")
	if len(cookieStoreEncryptionKey) == 0 {
		panic("GO_TEMPLATE_COOKIE_STORE_ENCRYPTION_KEY not defined")
	}

	jwtSigningKeyString := os.Getenv("GO_TEMPLATE_JWT_KEY")
	if len(jwtSigningKeyString) == 0 {
		panic("GO_TEMPLATE_JWT_KEY not defined")
	}

	auth.Initialize(&model.Session{}, cookieStoreAuthenticationKey, cookieStoreEncryptionKey, jwtSigningKeyString, "go_template_test")
}

func OpenDB(settings *ApplicationSettings) (db *sql.DB, err error) {
	dbConf, err := goose.NewDBConf(filepath.Join(settings.RootDirectory, "db"), settings.Environment, "")
	if err != nil {
		return nil, err
	}

	return goose.OpenDBFromDBConf(dbConf)
}
