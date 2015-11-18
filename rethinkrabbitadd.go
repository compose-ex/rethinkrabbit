package main

import (
	"crypto/tls"
	"log"

	r "github.com/dancannon/gorethink"
)

// Config contains various config data populated from YAML

func add(config Config) {
	conn, err := r.Connect(r.ConnectOpts{
		Address:  config.RethinkDBAddress,
		Database: config.RethinkDBDatabase,
		AuthKey:  config.RethinkDBAuthkey,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
	if err != nil {
		log.Fatal(conn, err)
	}

	type URLentry struct {
		ID  string `gorethink:"id,omitempty"`
		URL string `gorethink:"url"`
	}

	r.DB(config.RethinkDBDatabase).Table("urls").Insert(URLentry{URL: *urlarg}).RunWrite(conn)
	if err != nil {
		log.Fatal(err)
	}

	conn.Close()
}
