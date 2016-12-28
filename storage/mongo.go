package storage

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MongoConfig struct {
	Endpoint string `json:"endpoint"`
	Database string `json:"database"`
}

type admin struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

var curSession *mgo.Session

func connect(cfg MongoConfig) {
	var err error
	curSession, err = mgo.Dial(cfg.Endpoint)
	if err != nil {
		time.Sleep(time.Second * 10)
		connect(cfg)
	}
	curSession.SetMode(mgo.Monotonic, true)
}

//CheckUserAdmin check user is exist and user & pwd correct
func CheckUserAdmin(user string, pwd string, cfg MongoConfig) bool {
	connect(cfg)
	defer curSession.Close()

	info := admin{}
	selector := bson.M{"username": user}
	db := curSession.DB(cfg.Database)
	collection := db.C("etc_admin")

	err := collection.Find(selector).One(&info)
	if err != nil {
		return false
	}
	if info.Username == user && info.Password == pwd {
		return true
	}
	return false
}
