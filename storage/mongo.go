package storage

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
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
type UserInfo struct {
	UserAccount string  `bson:"account"`
	Wallet      string  `bson:"walletAddress"`
	Fee         float64 `bson:"fee"`
	Phone       int64   `bson:"phone"`
	Email       string  `bson:"email"`
	Vip         int     `bson:"vip"`
}

type offLine struct {
	Wallet string `bson:"walletAddress"`
	Time   string `bson:"offlineTime"`
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

//GetUserInfo get user info by user account or wallet address or phonenumer or email
func GetUserInfo(key string, value string, vip int, cfg MongoConfig) []UserInfo {
	//TODO insert user info when user log in at first time
	connect(cfg)
	defer curSession.Close()

	infos := []UserInfo{}
	value = strings.ToLower(value)
	selector := bson.M{key: value, "vip": vip}

	db := curSession.DB("etc_pool")
	collection := db.C("user_info")
	err := collection.Find(selector).All(&infos)
	if err != nil && err != mgo.ErrNotFound {
		return nil
	}
	return infos
}

//GetOnlineStat get online or offline state by given walletaddress, return offline time and if online return true, otherwise false
func GetOnlineStat(wallet string, cfg MongoConfig) string {
	connect(cfg)
	defer curSession.Close()

	selector := bson.M{"walletAddress": wallet}

	offlineInfo := offLine{}
	db := curSession.DB("etc_pool")
	collection := db.C("off_line")
	err := collection.Find(selector).One(&offlineInfo)
	if err == nil {
		return offlineInfo.Time
	}
	return ""
}
