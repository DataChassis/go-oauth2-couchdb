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

/* The data structure of the stored data */
type client struct {
	ID     string `json:"_id"`
	Secret string `json:"secret"`
	Domain string `json:"domain"`
	UserID string `json:"userid"`
}

// ClientConfig client configuration parameters
type ClientConfig struct {
	// store clients data collection name
	ClientsCName string
}

// ClientStore CouchDB storage for OAuth 2.0
type ClientStore struct {
	ccfg    *ClientConfig
	dbName  string
	session *kivik.Client
}

// NewDefaultClientConfig create a default client configuration
func NewDefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		ClientsCName: "dc-oauth2-clients",
	}
}

// NewClientStore create a client store instance based on mongodb
func NewClientStore(cfg *Config, ccfgs ...*ClientConfig) *ClientStore {
	url := strings.Replace(cfg.URL, "//", fmt.Sprintf("//%s:%s@", cfg.Username, cfg.Password), 1)
	session, err := kivik.New("couch", url)
	if err != nil {
		panic(err)
	}

	return NewClientStoreWithSession(session, cfg.DB, ccfgs...)
}

// NewClientStoreWithSession create a client store instance based on mongodb
func NewClientStoreWithSession(session *kivik.Client, dbName string, ccfgs ...*ClientConfig) *ClientStore {
	cs := &ClientStore{
		dbName:  dbName,
		session: session,
		ccfg:    NewDefaultClientConfig(),
	}
	if len(ccfgs) > 0 {
		cs.ccfg = ccfgs[0]
	}

	return cs
}

// Set set client information
func (cs *ClientStore) Set(info oauth2.ClientInfo) (err error) {
	entity := &client{
		ID:     info.GetID(),
		Secret: info.GetSecret(),
		Domain: info.GetDomain(),
		UserID: info.GetUserID(),
	}
	db := cs.session.DB(context.TODO(), cs.dbName)
	if _, cerr := db.Put(context.TODO(), info.GetID(), entity); cerr != nil {
		err = cerr
		return
	}

	return
}

// GetByID according to the ID for the client information
func (cs *ClientStore) GetByID(id string) (info oauth2.ClientInfo, err error) {
	db := cs.session.DB(context.TODO(), cs.dbName)
	entity := new(client)

	row := db.Get(context.TODO(), id)
	if cerr := row.ScanDoc(entity); cerr != nil {
		err = cerr
		return
	}

	info = &models.Client{
		ID:     entity.ID,
		Secret: entity.Secret,
		Domain: entity.Domain,
		UserID: entity.UserID,
	}

	return
}

// RemoveByID use the client id to delete the client information
func (cs *ClientStore) RemoveByID(id string) (err error) {
	db := cs.session.DB(context.TODO(), cs.dbName)
	var rev string
	var cerr error
	if _, rev, cerr = db.GetMeta(context.TODO(), id); cerr != nil {
		err = cerr
		return
	}
	if _, cerr := db.Delete(context.TODO(), id, rev); cerr != nil {
		err = cerr
		return
	}
	return
}
