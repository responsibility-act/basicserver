package basicserver

import (
	"log"

	"github.com/kataras/iris"

	mgo "github.com/globalsign/mgo"
)

const usersCollection = "users"
const statesCollection = "states"
const filesCollection = "files"

type collections struct {
	Users  *mgo.Collection
	States *mgo.Collection
	Files  *mgo.GridFS
}

// Settings values are used by BasicApp. At least `MongoString` and `ServerPort` are required.
//
//  Following values are possible:
//
//   `LogLevel` - available values are: "disable", "fatal", "error", "warn", "info", "debug"
//   `MongoString` - URI format described at http://docs.mongodb.org/manual/reference/connection-string/
//   `Secret` - secret value used by JWT parser
//   `ServerPort` - port on which the server should listen to
//
type Settings struct {
	LogLevel    string
	MongoString string
	Secret      []byte
	ServerPort  string
}

// BasicApp contains following fields:
//
//   `Coll.Users` - MongoDB "users" collection
//   `Coll.State` - MongoDB "states" collection
//   `Coll.File` - MongoDB "files" collection
//   `Db` - MongoDB named database
//   `Iris` - iris.Default() instance
//   `Settings` - Settings passed as an argument
//
type BasicApp struct {
	Coll     *collections
	Db       *mgo.Database
	Iris     *iris.Application
	Settings *Settings
}

// CreateApp returns BasicApp.
//
// `settings` argument should contain at least `MongoString` and `ServerPort` fields.
//
// BasicApp contains following fields:
//
//   `Coll.Users` - MongoDB "users" collection
//   `Coll.State` - MongoDB "states" collection
//   `Coll.File` - MongoDB "files" collection
//   `Db` - MongoDB named database
//   `Iris` - iris.Default() instance
//   `Settings` - Settings passed as an argument
//
func CreateApp(settings *Settings) *BasicApp {
	if settings.MongoString == "" {
		log.Fatal("MongoString cannot be empty!")
	}
	if settings.ServerPort == "" {
		log.Fatal("ServerPort cannot be empty!")
	}

	session, err := mgo.Dial(settings.MongoString)
	if err != nil {
		log.Fatal(err)
	}
	db := session.DB("")

	app := &BasicApp{
		Coll: &collections{
			Users:  db.C(usersCollection),
			States: db.C(statesCollection),
			Files:  db.GridFS(filesCollection),
		},
		Db:       db,
		Iris:     iris.Default(),
		Settings: settings,
	}

	app.Iris.Logger().SetLevel(settings.LogLevel)

	return app
}
