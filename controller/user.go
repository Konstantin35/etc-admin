package controller

import (
	"encoding/json"
	"etc-pool-admin/storage"
	"github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
)

const (
	queryByAccount = "account"
	queryByAddress = "walletAddress"
	queryByPhone   = "phone"
	queryByEmail   = "email"
)

type userAllInfo struct {
	BasicInfo   storage.UserInfo
	LastRevenue int64
	AllRevenue  int64
	OfflineTime string
}

//QueryUsers get givin pattern correspodding users and info
func QueryUsers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", "GET")
		res.WriteHeader(http.StatusOK)
		return
	}
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	param := mux.Vars(req)["value"]
	querykey := regexpParam(param)
	basicInfo := storage.GetUserInfo(querykey, param, Conf.Mongo)

	allUserInfo := make([]userAllInfo, len(basicInfo))
	for idx, basic := range basicInfo {
		allUserInfo[idx].BasicInfo = basic
		allUserInfo[idx].LastRevenue, allUserInfo[idx].AllRevenue = Backend.GetRevenue(basic.Wallet)
		allUserInfo[idx].OfflineTime = storage.GetOnlineStat(basic.Wallet, Conf.Mongo)
	}
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(allUserInfo)
	if err != nil {
		seelog.Error("query users error, when serializing response data:", err)
	}
}

func SetUserInfo(res http.ResponseWriter, req *http.Request) {

}

func GetAddressChartData(res http.ResponseWriter, req *http.Request) {

}

func GetAddressBenefitData(res http.ResponseWriter, req *http.Request) {

}

func GetMinersInfo(res http.ResponseWriter, req *http.Request) {

}

func QueryPayment(res http.ResponseWriter, req *http.Request) {

}

func regexpParam(param string) string {
	//regexp wallet address
	regAddr := regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	regPhone := regexp.MustCompile(`^1[3|4|5|8][0-9]{9}$`)
	regEmail := regexp.MustCompile(`^[a-z_0-9.-]{1,64}@([a-z0-9-]{1,200}.){1,5}[a-z]{1,6}$`)

	if regAddr.MatchString(param) {
		return queryByAddress
	} else if regPhone.MatchString(param) {
		return queryByPhone
	} else if regEmail.MatchString(param) {
		return queryByEmail
	}

	return queryByAccount
}
