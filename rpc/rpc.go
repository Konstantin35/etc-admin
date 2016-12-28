package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/cihub/seelog"
	"manageserver/utils/common"
	"math/big"
	"net/http"
	"time"
)

type Upstream struct {
	Url     string `json:"url"`
	Timeout string `json:"timeout"`
}

type RPCClient struct {
	Url    string
	client *http.Client
}

type JSONRpcResp struct {
	Id     *json.RawMessage       `json:"id"`
	Result *json.RawMessage       `json:"result"`
	Error  map[string]interface{} `json:"error"`
}

//NewRPCClient make a new rpc lient for get data from node
func NewRPCClient(url string, timeout string) *RPCClient {
	rpcClient := &RPCClient{Url: url}
	timeoutIntv, err := time.ParseDuration(timeout)
	if err != nil {
		seelog.Info("parse timeout when new rpc clien error:", err)
		panic(err)
	}
	rpcClient.client = &http.Client{
		Timeout: timeoutIntv,
	}
	return rpcClient
}

//GetAccountBalance get given address balance by rpc
func (r *RPCClient) GetAccountBalance(account string) (*big.Int, error) {
	rpcResp, err := r.doPost(r.Url, "eth_getBalance", []string{account, "latest"})
	if err != nil {
		return nil, err
	}
	var reply string
	err = json.Unmarshal(*rpcResp.Result, &reply)
	if err != nil {
		return nil, err
	}
	return common.String2Big(reply), err
}

func (r *RPCClient) doPost(url string, method string, params interface{}) (*JSONRpcResp, error) {
	jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 0}
	data, _ := json.Marshal(jsonReq)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Length", (string)(len(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp *JSONRpcResp
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		return nil, err
	}
	if rpcResp.Error != nil {
		return nil, errors.New(rpcResp.Error["message"].(string))
	}
	return rpcResp, err
}
