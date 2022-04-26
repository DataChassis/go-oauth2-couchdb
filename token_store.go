package couchdb

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"

	_ "github.com/go-kivik/couchdb/v3" // The CouchDB driver
	kivik "github.com/go-kivik/kivik/v3"
)

type newDocument struct {
	Payload oauth2.TokenInfo
}

type existingDocument struct {
	ID      string `json:"_id"`
	Rev     string `json:"_rev"`
	Payload models.Token
}

// TokenStore CouchDB storage for OAuth 2.0
type TokenStore struct {
	tcfg    *TokenConfig
	dbName  string
	session *kivik.Client
}

// TokenConfig token configuration parameters
type TokenConfig struct {
	// store token based data collection name
	BasicCName string
}

// NewDefaultTokenConfig create a default token configuration
func NewDefaultTokenConfig() *TokenConfig {
	return &TokenConfig{
		BasicCName: "dc-oauth2-tokens",
	}
}

// NewTokenStore create a token store instance based on mongodb
func NewTokenStore(cfg *Config, tcfgs ...*TokenConfig) (store *TokenStore) {
	url := strings.Replace(cfg.URL, "//", fmt.Sprintf("//%s:%s@", cfg.Username, cfg.Password), 1)
	session, err := kivik.New("couch", url)
	if err != nil {
		panic(err)
	}

	return NewTokenStoreWithSession(session, cfg.DB, tcfgs...)
}

// NewTokenStoreWithSession create a token store instance based on mongodb
func NewTokenStoreWithSession(session *kivik.Client, dbName string, tcfgs ...*TokenConfig) (store *TokenStore) {
	ts := &TokenStore{
		dbName:  dbName,
		session: session,
		tcfg:    NewDefaultTokenConfig(),
	}
	if len(tcfgs) > 0 {
		ts.tcfg = tcfgs[0]
	}

	/* Ensure the required indexes exist */

	store = ts
	return
}

// Close close the couchdb session
func (ts *TokenStore) Close() {
	ts.session.Close(context.TODO())
}

// Create create and store the new token information
func (ts *TokenStore) Create(info oauth2.TokenInfo) (err error) {
	db := ts.session.DB(context.TODO(), ts.dbName)
	doc := &newDocument{Payload: info}
	_, _, err = db.CreateDoc(context.TODO(), doc)
	return
}

// RemoveByCode use the authorization code to delete the token information
func (ts *TokenStore) RemoveByCode(code string) (err error) {
	db := ts.session.DB(context.TODO(), ts.dbName)
	id, cerr := getIdByView(db, code, "code")
	if cerr != nil {
		err = cerr
		return
	}
	var rev string
	if _, rev, cerr = db.GetMeta(context.TODO(), id); cerr != nil {
		err = cerr
		return
	}
	_, cerr = db.Delete(context.TODO(), id, rev)
	if cerr != nil {
		err = cerr
	}
	return
}

// RemoveByAccess use the access token to delete the token information
func (ts *TokenStore) RemoveByAccess(access string) (err error) {
	db := ts.session.DB(context.TODO(), ts.dbName)
	id, cerr := getIdByView(db, access, "access")
	if cerr != nil {
		err = cerr
		return
	}
	var rev string
	if _, rev, cerr = db.GetMeta(context.TODO(), id); cerr != nil {
		err = cerr
		return
	}
	_, cerr = db.Delete(context.TODO(), id, rev)
	if cerr != nil {
		err = cerr
	}
	return
}

// RemoveByRefresh use the refresh token to delete the token information
func (ts *TokenStore) RemoveByRefresh(refresh string) (err error) {
	db := ts.session.DB(context.TODO(), ts.dbName)
	id, cerr := getIdByView(db, refresh, "refresh")
	if cerr != nil {
		err = cerr
		return
	}
	var rev string
	if _, rev, cerr = db.GetMeta(context.TODO(), id); cerr != nil {
		err = cerr
		return
	}
	_, cerr = db.Delete(context.TODO(), id, rev)
	if cerr != nil {
		err = cerr
	}
	return
}

func getData(db *kivik.DB, id string) (ti oauth2.TokenInfo, err error) {
	doc := &existingDocument{}
	row := db.Get(context.TODO(), id)
	if cerr := row.ScanDoc(doc); cerr != nil {
		err = cerr
		return
	}
	ti = &doc.Payload
	return
}

func getIdByView(db *kivik.DB, token string, tokenType string) (id string, err error) {
	rows, err := db.Query(context.TODO(), "token_views", "by_"+tokenType, kivik.Options{"key": token})
	if err != nil {
		return
	}
	for rows.Next() {
		if id != "" {
			err = fmt.Errorf("multiple token documents exist with the code [%s]", token)
			return
		}
		id = rows.ID()
	}
	return
}

// GetByCode use the authorization code for token information data
func (ts *TokenStore) GetByCode(code string) (ti oauth2.TokenInfo, err error) {
	db := ts.session.DB(context.TODO(), ts.dbName)
	id, cerr := getIdByView(db, code, "code")
	if id == "" {
		return
	}
	if cerr != nil {
		err = cerr
		return
	}
	ti, err = getData(db, id)
	return
}

// GetByAccess use the access token for token information data
func (ts *TokenStore) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {
	db := ts.session.DB(context.TODO(), ts.dbName)
	id, cerr := getIdByView(db, access, "access")
	if id == "" {
		return
	}
	if cerr != nil {
		err = cerr
		return
	}
	ti, err = getData(db, id)
	return
}

// GetByRefresh use the refresh token for token information data
func (ts *TokenStore) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {
	db := ts.session.DB(context.TODO(), ts.dbName)
	id, cerr := getIdByView(db, refresh, "refresh")
	if id == "" {
		return
	}
	if cerr != nil {
		err = cerr
		return
	}
	ti, err = getData(db, id)
	return
}
