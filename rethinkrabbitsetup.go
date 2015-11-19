package main

import (
	"crypto/tls"
	"log"

	r "github.com/dancannon/gorethink"
)

// Config contains various config data populated from YAML

func setup(config Config) {
	conn, err := r.Connect(r.ConnectOpts{
		Address:  config.RethinkDBAddress,
		Database: config.RethinkDBDatabase,
		AuthKey:  config.RethinkDBAuthkey,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	r.DB(config.RethinkDBDatabase).TableDrop("urls").Run(conn)

	r.DB(config.RethinkDBDatabase).TableCreate("urls").Run(conn)

	conn.Close()
}
