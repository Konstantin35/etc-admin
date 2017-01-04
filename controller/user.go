package controller

import (
	"encoding/json"
	"etc-pool-admin/storage"
	"github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", http.MethodPut)
		res.WriteHeader(http.StatusOK)
		return
	}
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	modify := storage.UserInfo{}
	body, _ := ioutil.ReadAll(req.Body)
	seelog.Info("data body:", string(body))

	err = json.Unmarshal(body, &modify)
	if err != nil {
		seelog.Info("decode requst body error:", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	// decoder := json.NewDecoder(req.Body)
	// err = decoder.Decode(&modify)
	// if err != nil {
	// 	seelog.Info("decode requst body error:", err)
	// 	res.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	seelog.Info("after unmarshal:", modify)
	seelog.Info("json email:", modify.Email)
	seelog.Info("json wallet:", modify.Wallet)
	err = storage.SetUserInfo(modify, Conf.Mongo)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

//SetFee set common users or vip users fee
func SetFee(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", http.MethodPut)
		res.WriteHeader(http.StatusOK)
		return
	}
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}
	vip := mux.Vars(req)["vip"]
	isvip, err := strconv.Atoi(vip)
	if err != nil {
		seelog.Info("invalid vip flag:", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	req.ParseForm()

	v := req.PostFormValue("fee")

	fee, _ := strconv.ParseFloat(v, 8)
	err = storage.SetFee(fee, isvip, Conf.Mongo)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}

//GetAddressChartData get 24 housrs hashrate from redis
func GetAddressChartData(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", http.MethodPut)
		res.WriteHeader(http.StatusOK)
		return
	}
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}
	address := mux.Vars(req)["address"]
	address = strings.ToLower(address)
	chartdata, err := Backend.GetAccountChartValues(address)
	if err != nil {
		seelog.Info("get wallet address chart data error:", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(chartdata)
	if err != nil {
		seelog.Error("query chart data error, when serializing response data:", err)
	}
}

//GetAddressStaticData get revenue data include matured revenue, immature revenue, last paid, total paid
func GetAddressStaticData(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", http.MethodPut)
		res.WriteHeader(http.StatusOK)
		return
	}
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}
}

//GetMinersInfo get miner state about online offline and its hashrate data
func GetMinersInfo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", http.MethodPut)
		res.WriteHeader(http.StatusOK)
		return
	}
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}
}

func QueryPaymentHistory(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", http.MethodPut)
		res.WriteHeader(http.StatusOK)
		return
	}
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}
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
