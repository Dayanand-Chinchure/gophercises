package toydb

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

//Database name to be used and collection "Employee"
var Database = "toy_db"

//Bootstrap package which actually does the db connection using (s *mgo.Session)
func Bootstrap(s *mgo.Session) *mgo.Collection {

	c := s.DB(Database).C("Employee")
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		Background: true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		fmt.Println("Error connecting to collection !")
	}

	return c
}

//Createconnection to the mongo db using the mgo.Dial operation
func Createconnection() *mgo.Collection {
	//defer session.Close()

	session, err := mgo.Dial("mongodb://localhost:27017/" + Database)

	if err != nil {
		fmt.Println(err)
	}
	c := Bootstrap(session)
	return c
}
