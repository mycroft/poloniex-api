package poloniexapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type BalancesJson map[string]json.Number

/*
returnBalances
Returns all of your available balances. Sample output:

{"BTC":"0.59098578","LTC":"3.31117268", ... }
*/
func (api *PoloniexApi) ApiPrivateBalances() (map[string]float64, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_BALANCES)

	out_json := new(BalancesJson)

	_, err := api.queryparse(URL_PRIVATE, params, true, &out_json)
	if err != nil {
		return nil, err
	}

	out := make(map[string]float64)

	for k, v := range *out_json {
		out[k], err = v.Float64()
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

/*
returnCompleteBalances
Returns all of your balances, including available balance, balance on orders,
and the estimated BTC value of your balance. By default, this call is limited
to your exchange account; set the "account" POST parameter to "all" to include
your margin and lending accounts. Sample output:

{"LTC":{"available":"5.015","onOrders":"1.0025","btcValue":"0.078"},"NXT:{...} ... }
*/
func (api *PoloniexApi) ApiPrivateCompleteBalances(complete bool) (map[string]Balance, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_COMPLETE_BALANCES)

	if complete {
		params.Set("account", "all")
	}

	out := make(map[string]Balance)

	_, err := api.queryparse(URL_PRIVATE, params, true, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
returnDepositAddresses
Returns all of your deposit addresses. Sample output:

{"BTC":"19YqztHmspv2egyD6jQM3yn81x5t5krVdJ","LTC":"LPgf9kjv9H1Vuh4XSaKhzBe8JHdou1WgUB", ...
 "ITC":"Press Generate.." ... }
*/
func (api *PoloniexApi) ApiPrivateDepositAddresses() (map[string]string, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_DEPOSIT_ADDRESSES)

	out := make(map[string]string)

	_, err := api.queryparse(URL_PRIVATE, params, true, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
generateNewAddress
Generates a new deposit address for the currency specified by the "currency" POST parameter.
Sample output:

{"success":1,"response":"CKXbbs8FAVbtEa397gJHSutmrdrBrhUMxe"}

Only one address per currency per day may be generated, and a new address may not be generated
before the previously-generated one has been used.
*/
func (api *PoloniexApi) ApiPrivateGenerateNewAddress(currency string) (*GenerateAddressResponse, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_NEW_ADDRESS)
	params.Set("currency", currency)

	out := new(GenerateAddressResponse)

	_, err := api.queryparse(URL_PRIVATE, params, true, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
returnDepositsWithdrawals
Returns your deposit and withdrawal history within a range, specified by the "start" and "end" POST parameters,
both of which should be given as UNIX timestamps. Sample output:

{"deposits":
[{"currency":"BTC","address":"...","amount":"0.01006132","confirmations":10,
"txid":"17f819a91369a9ff6c4a34216d434597cfc1b4a3d0489b46bd6f924137a47701","timestamp":1399305798,"status":"COMPLETE"},
{"currency":"BTC","address":"...","amount":"0.00404104","confirmations":10,
"txid":"7acb90965b252e55a894b535ef0b0b65f45821f2899e4a379d3e43799604695c","timestamp":1399245916,"status":"COMPLETE"}],
"withdrawals":
[{"withdrawalNumber":134933,"currency":"BTC","address":"1N2i5n8DwTGzUq2Vmn9TUL8J1vdr1XBDFg","amount":"5.00010000",
"timestamp":1399267904,"status":"COMPLETE: 36e483efa6aff9fd53a235177579d98451c4eb237c210e66cd2b9a2d4a988f8e","ipAddress":"..."}]}
*/
func (api *PoloniexApi) ApiPrivateDepositWithdrawals(start, end int64) (*DepositWithdrawal, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_DEPOSIT_WITHDRAWALS)

	params.Set("start", strconv.FormatInt(start, 10))
	params.Set("end", strconv.FormatInt(end, 10))

	out := new(DepositWithdrawal)
	_, err := api.queryparse(URL_PRIVATE, params, true, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
returnOpenOrders
Returns your open orders for a given market, specified by the "currencyPair" POST parameter, e.g. "BTC_XCP".
Set "currencyPair" to "all" to return open orders for all markets. Sample output for single market:

[{"orderNumber":"120466","type":"sell","rate":"0.025","amount":"100","total":"2.5"},
 {"orderNumber":"120467","type":"sell","rate":"0.04","amount":"100","total":"4"}, ... ]
*/
func (api *PoloniexApi) ApiPrivateOpenOrders(currencyPair string) (map[string][]OpenOrder, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_OPEN_ORDERS)
	params.Set("currencyPair", currencyPair)

	out := make(map[string][]OpenOrder)
	if currencyPair == "all" {
		_, err := api.queryparse(URL_PRIVATE, params, true, &out)
		if err != nil {
			return nil, err
		}
	} else {
		out_tmp := make([]OpenOrder, 0)
		_, err := api.queryparse(URL_PRIVATE, params, true, &out_tmp)
		if err != nil {
			return nil, err
		}
		out[currencyPair] = out_tmp
	}

	return out, nil
}

/*
returnTradeHistory
Returns your trade history for a given market, specified by the "currencyPair"
POST parameter. You may specify "all" as the currencyPair to receive your trade
history for all markets. You may optionally specify a range via "start" and/or
"end" POST parameters, given in UNIX timestamp format; if you do not specify a
range, it will be limited to one day. Sample output:

[{ "globalTradeID": 25129732, "tradeID": "6325758", "date": "2016-04-05 08:08:40",
   "rate": "0.02565498", "amount": "0.10000000", "total": "0.00256549", "fee": "0.00200000",
   "orderNumber": "34225313575", "type": "sell", "category": "exchange" },
 { "globalTradeID": 25129628, "tradeID": "6325741", "date": "2016-04-05 08:07:55",
   "rate": "0.02565499", "amount": "0.10000000", "total": "0.00256549", "fee": "0.00200000",
   "orderNumber": "34225195693", "type": "buy", "category": "exchange" }, ... ]
*/
func (api *PoloniexApi) ApiPrivateTradeHistory(currencyPair string, start, end int) (map[string][]Trade, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_TRADE_HISTORY)
	params.Set("currencyPair", currencyPair)

	if start != 0 {
		params.Set("start", strconv.Itoa(start))
	}

	if end != 0 {
		params.Set("end", strconv.Itoa(end))
	}

	out := make(map[string][]Trade)

	if currencyPair != "all" {
		out_tmp := make([]Trade, 0)
		_, err := api.queryparse(URL_PRIVATE, params, true, &out_tmp)
		if err != nil {
			return nil, err
		}
		out[currencyPair] = out_tmp
	} else {
		_, err := api.queryparse(URL_PRIVATE, params, true, &out)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

/*
returnOrderTrades
Returns all trades involving a given order, specified by the "orderNumber" POST parameter. If no trades for the order have occurred or you specify an order that does not belong to you, you will receive an error. Sample output:

[{"globalTradeID": 20825863, "tradeID": 147142, "currencyPair": "BTC_XVC", "type": "buy", "rate": "0.00018500", "amount": "455.34206390", "total": "0.08423828", "fee": "0.00200000", "date": "2016-03-14 01:04:36"}, ...]
*/
func (api *PoloniexApi) ApiPrivateOrderTrades(orderNumber string) ([]Trade, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_ORDER_TRADES)
	params.Set("orderNumber", orderNumber)

	out := make([]Trade, 0)

	_, err := api.queryparse(URL_PRIVATE, params, true, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
buy
Places a limit buy order in a given market. Required POST parameters are
"currencyPair", "rate", and "amount". If successful, the method will return the
order number. Sample output:

{"orderNumber":31226040,"resultingTrades":[
  {"amount":"338.8732","date":"2014-10-18 23:03:21","rate":"0.00000173",
   "total":"0.00058625","tradeID":"16164","type":"buy"}]}

You may optionally set "fillOrKill", "immediateOrCancel", "postOnly" to 1. A
fill-or-kill order will either fill in its entirety or be completely aborted.
An immediate-or-cancel order can be partially or completely filled, but any
portion of the order that cannot be filled immediately will be canceled rather
than left on the order book. A post-only order will only be placed if no
portion of it fills immediately; this guarantees you will never pay the taker
fee on any part of the order that fills.
*/
func (api *PoloniexApi) ApiPrivateBuy(currencyPair string, rate, amount float64, opts map[string]bool) (*Order, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_BUY)
	params.Set("currencyPair", currencyPair)
	params.Set("rate", strconv.FormatFloat(rate, 'f', -1, 64))
	params.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))

	if _, ok := opts["fillOrKill"]; ok {
		params.Set("fillOrKill", "1")
	}
	if _, ok := opts["immediateOrCancel"]; ok {
		params.Set("immediateOrCancel", "1")
	}
	if _, ok := opts["postOnly"]; ok {
		params.Set("postOnly", "1")
	}

	out := new(Order)

	_, err := api.queryparse(URL_PRIVATE, params, true, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
sell
Places a sell order in a given market. Parameters and output are the same as for the buy method.
*/
func (api *PoloniexApi) ApiPrivateSell(currencyPair string, rate, amount float64, opts map[string]bool) (*Order, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_SELL)
	params.Set("currencyPair", currencyPair)
	params.Set("rate", strconv.FormatFloat(rate, 'f', -1, 64))
	params.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))

	if _, ok := opts["fillOrKill"]; ok {
		params.Set("fillOrKill", "1")
	}
	if _, ok := opts["immediateOrCancel"]; ok {
		params.Set("immediateOrCancel", "1")
	}
	if _, ok := opts["postOnly"]; ok {
		params.Set("postOnly", "1")
	}

	out := new(Order)

	_, err := api.queryparse(URL_PRIVATE, params, true, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
cancelOrder
Cancels an order you have placed in a given market. Required POST parameter is "orderNumber". If successful, the method will return:
{"success":1}
*/
func (api *PoloniexApi) ApiPrivateCancel(orderNumber int64) (bool, *CancelOrder, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_CANCEL_ORDER)
	params.Set("orderNumber", strconv.FormatInt(orderNumber, 10))

	out := new(CancelOrder)

	_, err := api.queryparse(URL_PRIVATE, params, true, &out)
	if err != nil {
		return false, nil, err
	}

	if out.Error != "" {
		return false, nil, fmt.Errorf(out.Error)
	}

	return 1 == out.Success, out, nil
}

/*
moveOrder
Cancels an order and places a new one of the same type in a single atomic
transaction, meaning either both operations will succeed or both will fail.
Required POST parameters are "orderNumber" and "rate"; you may optionally
specify "amount" if you wish to change the amount of the new order. "postOnly"
or "immediateOrCancel" may be specified for exchange orders, but will have no
effect on margin orders. Sample output:

{"success":1,"orderNumber":"239574176","resultingTrades":{"BTC_BTS":[]}}
*/
func (api *PoloniexApi) ApiPrivateMoveOrder(orderNumber int64, rate, amount float64, opts map[string]bool) (*Order, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_MOVE_ORDER)
	params.Set("orderNumber", strconv.FormatInt(orderNumber, 10))
	params.Set("rate", strconv.FormatFloat(rate, 'f', -1, 64))

	if amount != 0.0 {
		params.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	}

	if _, ok := opts["immediateOrCancel"]; ok {
		params.Set("immediateOrCancel", "1")
	}
	if _, ok := opts["postOnly"]; ok {
		params.Set("postOnly", "1")
	}

	out := new(Order)

	_, err := api.queryparse(URL_PRIVATE, params, true, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
withdraw
Immediately places a withdrawal for a given currency, with no email
confirmation. In order to use this method, the withdrawal privilege must be
enabled for your API key. Required POST parameters are "currency", "amount",
and "address". For XMR withdrawals, you may optionally specify "paymentId".
Sample output:

{"response":"Withdrew 2398 NXT."}

XXX: Untested method
*/
func (api *PoloniexApi) ApiPrivateWithdraw(currency, address string, amount float64) (string, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_WITHDRAW)
	params.Set("currency", currency)
	params.Set("address", address)
	params.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))

	var out string

	_, err := api.queryparse(URL_PRIVATE, params, true, &out)
	if err != nil {
		return "", err
	}

	return out, nil
}

/*
returnFeeInfo
If you are enrolled in the maker-taker fee schedule, returns your current
trading fees and trailing 30-day volume in BTC. This information is updated
once every 24 hours.

{"makerFee": "0.00140000", "takerFee": "0.00240000", "thirtyDayVolume": "612.00248891", "nextTier": "1200.00000000"}
*/
func (api *PoloniexApi) ApiPrivateFeeInfo() (*FeeInfo, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_FEE_INFO)

	out := new(FeeInfo)

	_, err := api.queryparse(URL_PRIVATE, params, true, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func convertMapMapJsonMapMapFloat64(in map[string]map[string]json.Number) (map[string]map[string]float64, error) {
	var err error
	out := make(map[string]map[string]float64)

	for k, v := range in {
		sub_temp := make(map[string]float64)
		for subk, subv := range v {
			sub_temp[subk], err = subv.Float64()
			if err != nil {
				return nil, err
			}
		}
		out[k] = sub_temp
	}

	return out, nil
}

/*
returnAvailableAccountBalances
Returns your balances sorted by account. You may optionally specify the
"account" POST parameter if you wish to fetch only the balances of one account.
Please note that balances in your margin account may not be accessible if you
have any open margin positions or orders. Sample output:

{"exchange":{"BTC":"1.19042859","BTM":"386.52379392","CHA":"0.50000000","DASH":"120.00000000",
 "STR":"3205.32958001", "VNL":"9673.22570147"},"margin":{"BTC":"3.90015637","DASH":"250.00238240",
 "XMR":"497.12028113"},"lending":{"DASH":"0.01174765","LTC":"11.99936230"}}
*/
func (api *PoloniexApi) ApiPrivateAvailableAccountBalances(account string) (map[string]map[string]float64, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_AVAILABLE_ACCOUNT_BALANCES)
	if account != "" {
		params.Set("account", account)
	}

	out_tmp := make(map[string]map[string]json.Number)

	_, err := api.queryparse(URL_PRIVATE, params, true, &out_tmp)
	if err != nil {
		return nil, err
	}

	out, err := convertMapMapJsonMapMapFloat64(out_tmp)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
returnTradableBalances
Returns your current tradable balances for each currency in each market for
which margin trading is enabled. Please note that these balances may vary
continually with market conditions. Sample output:

{"BTC_DASH":{"BTC":"8.50274777","DASH":"654.05752077"},
 "BTC_LTC":{"BTC":"8.50274777","LTC":"1214.67825290"},
 "BTC_XMR":{"BTC":"8.50274777","XMR":"3696.84685650"}}
*/
func (api *PoloniexApi) ApiPrivateTradableBalances() (map[string]map[string]float64, error) {
	params := url.Values{}
	params.Set("command", CMD_PRIVATE_TRADABLE_BALANCES)

	out_tmp := make(map[string]map[string]json.Number)

	_, err := api.queryparse(URL_PRIVATE, params, true, &out_tmp)
	if err != nil {
		return nil, err
	}

	out, err := convertMapMapJsonMapMapFloat64(out_tmp)
	if err != nil {
		return nil, err
	}

	return out, nil
}
