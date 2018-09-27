package profiler

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"log"
	"time"
)

type DBHelper struct {
	Server   string
	Database string
	Session  *mgo.Session
}

var db *mgo.Database

const (
	COLLECTION = "ldoe"
)

// Establish a connection to database
func (m *DBHelper) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	m.Session = session
	fmt.Println("ping db:", session.Ping() == nil)
	db = session.DB(m.Database)
}

// Find list of movies
func (m *DBHelper) FindAll() (iter *mgo.Iter) {
	return db.C(COLLECTION).Find(nil).Iter()
}

// Find a movie by its id
func (m *DBHelper) FindByDate(start time.Time, end time.Time) (iter *mgo.Iter) {
	query := bson.M{}
	if !start.IsZero() {
		query["start_ts"] = start
	}
	if !end.IsZero() {
		query["end_ts"] = end
	}
	return db.C(COLLECTION).Find(query).Iter()
}

//// Find a movie by its id
//func (m *DBHelper) FindById(id string) (Movie, error) {
//	var movie Movie
//	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&movie)
//	return movie, err
//}
//
//// Insert a movie into database
//func (m *DBHelper) Insert(movie Movie) error {
//	err := db.C(COLLECTION).Insert(&movie)
//	return err
//}
//
//// Delete an existing movie
//func (m *DBHelper) Delete(movie Movie) error {
//	err := db.C(COLLECTION).Remove(&movie)
//	return err
//}
//
//// Update an existing movie
//func (m *DBHelper) Update(movie Movie) error {
//	err := db.C(COLLECTION).UpdateId(movie.ID, &movie)
//	return err
//}
