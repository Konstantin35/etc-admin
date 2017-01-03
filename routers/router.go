package router

import (
	"etc-pool-admin/controller"
	"github.com/gorilla/mux"
	"net/http"
)

var Routes *mux.Router

func init() {
	Routes = mux.NewRouter()
	Routes.NotFoundHandler = http.HandlerFunc(controller.NotFound)
	Routes.HandleFunc("/login", controller.Login).Methods("POST")
	Routes.HandleFunc("/main/poolchart", controller.PoolChartData).Methods("GET", "OPTIONS")
	Routes.HandleFunc("/main/statistic", controller.StatisticData).Methods("GET", "OPTIONS")
	/****************routers for user manage********************/
	Routes.HandleFunc("/user/query/{value}/{vip}", controller.QueryUsers).Methods("GET", "OPTIONS") //query by wallet address or login account or email or phone number
	Routes.HandleFunc("/user/info/settings", controller.SetUserInfo).Methods("PUT")
	Routes.HandleFunc("/user/walletaddress/data/chart", controller.GetAddressChartData).Methods("GET", "OPTIONS")
	Routes.HandleFunc("/user/walletaddress/data/benefit", controller.GetAddressBenefitData).Methods("GET", "OPTIONS")
	Routes.HandleFunc("/user/walletaddress/info/miners", controller.GetMinersInfo).Methods("GET", "OPTIONS")
	Routes.HandleFunc("/user/walletaddress/history/payment", controller.QueryPayment).Methods("GET", "OPTIONS")
}
