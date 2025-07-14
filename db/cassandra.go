package db

import (
	"log"
	"time"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func InitCassandra() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "userks"
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = time.Second * 10

	var err error
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("unable to connect to Cassandra: %v", err)
	}
	log.Println("Cassandra connected")
}
