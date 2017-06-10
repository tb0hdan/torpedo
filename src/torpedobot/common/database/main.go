package database

import (
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

type MongoDB struct {
	Server string
	Database string
}


func (mdb *MongoDB) GetSession() (session *mgo.Session, err error){
	session, err = mgo.Dial(mdb.Server)
	if err != nil {
		panic(err)
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return
}

func (mdb *MongoDB) GetCollection(collectionName string) (session *mgo.Session, collection *mgo.Collection, err error) {
	session, err = mdb.GetSession()
	if err != nil {
		return
	}
	collection = session.DB(mdb.Database).C(collectionName)
	return
}

func New(server, database string) (mongodb *MongoDB){
	if server == "" {
		server = "localhost"
	}
	if database == "" {
		database = "torpedobot"
	}
	mongodb = &MongoDB{Server:server,
		           Database:database}
	return
}

/*
func main() {
	session, err := mgo.Dial("server1.example.com,server2.example.com")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
}
*/