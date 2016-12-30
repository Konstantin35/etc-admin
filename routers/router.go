package router

import (
	"etc-pool-admin/controller"
	"github.com/gorilla/mux"
)

var Routes *mux.Router

func init() {
	Routes = mux.NewRouter()
	Routes.HandleFunc("/login", controller.Login).Methods("POST")
	Routes.HandleFunc("/main/poolchart", controller.PoolChartData).Methods("GET", "OPTIONS")
	Routes.HandleFunc("/main/statistic", controller.StatisticData).Methods("GET", "OPTIONS")
	/****************routers for user manage********************/
	Routes.HandleFunc("/user/query", controller.QueryUsers).Methods("GET") //query by wallet address or login account or email or phone number
	Routes.HandleFunc("/user/info/settings", controller.SetUserInfo).Methods("PUT")
	Routes.HandleFunc("/user/walletaddress/data/chart", controller.GetAddressChartData).Methods("GET")
	Routes.HandleFunc("/user/walletaddress/data/benefit", controller.GetAddressBenefitData).Methods("GET")
	Routes.HandleFunc("/user/walletaddress/info/miners", controller.GetMinersInfo).Methods("GET")
	Routes.HandleFunc("/user/walletaddress/history/payment", controller.QueryPayment).Methods("GET")
}
