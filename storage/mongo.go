package storage

import (
	"errors"
	"github.com/cihub/seelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
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
	UserAccount string  `bson:"account" json:"account"`
	Wallet      string  `bson:"walletAddress" json:"walletAddress"`
	Fee         float64 `bson:"fee" json:"fee"`
	Phone       string  `bson:"phone" json:"phone"`
	Email       string  `bson:"email" json:"email"`
	Vip         int     `bson:"vip" json:"vip"`
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
func GetUserInfo(key string, value string, cfg MongoConfig) []UserInfo {
	//TODO insert user info when user log in at first time
	connect(cfg)
	defer curSession.Close()

	infos := []UserInfo{}
	value = strings.ToLower(value)
	selector := bson.M{key: value}

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

//SetUserInfo set user fee , phone, email and so on
func SetUserInfo(info UserInfo, cfg MongoConfig) error {
	connect(cfg)
	defer curSession.Close()

	selector := bson.M{}

	if info.UserAccount != "" {
		selector["account"] = info.UserAccount
	} else if info.Wallet != "" {
		selector["walletAddress"] = info.Wallet
	} else {
		seelog.Info("cannot find useraccount or wallet as set idex")
		return errors.New("cannot find useraccount or wallet as set index")
	}
	setdata := bson.M{}

	val := reflect.ValueOf(&info).Elem()

	for i := 0; i < val.NumField(); i++ {
		if val.Type().Field(i).Tag.Get("json") == "account" || val.Type().Field(i).Tag.Get("json") == "walletAddress" {
			continue
		}
		setdata[val.Type().Field(i).Tag.Get("json")] = val.Field(i).Interface()
	}
	setter := bson.M{"$set": setdata}

	db := curSession.DB("etc_pool")
	collection := db.C("user_info")
	changeinfo, err := collection.UpdateAll(selector, setter)
	seelog.Info("change info:", changeinfo)
	if err != nil {
		seelog.Error("set user info error:", err)
		return err
	}
	return nil
}

//SetFee use for settin common users fee or vip users fee
func SetFee(fee float64, vip int, cfg MongoConfig) error {
	connect(cfg)
	defer curSession.Close()

	selector := bson.M{}
	setdata := bson.M{}
	if vip == 1 { //set all vip user's fee
		selector["vip"] = 1
		setdata["fee"] = fee
	} else if vip == 0 { //set all normall user's fee
		selector["vip"] = 0
		setdata["fee"] = fee
	} else {
		return errors.New("error vip flag")
	}

	setter := bson.M{"$set": setdata}
	db := curSession.DB("etc_pool")
	collection := db.C("user_info")
	changeinfo, err := collection.UpdateAll(selector, setter)
	seelog.Info("change info:", changeinfo)
	if err != nil {
		seelog.Error("change fee error:", err)
		return err
	}
	return nil
}
