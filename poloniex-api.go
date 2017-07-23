package poloniexapi

import (
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

const (
	URL_PUBLIC  = "https://poloniex.com/public"
	URL_PRIVATE = "https://poloniex.com/tradingApi"

	CMD_PUBLIC_TICKER                      = "returnTicker"
	CMD_PUBLIC_24HVOLUME                   = "return24hVolume"
	CMD_PUBLIC_ORDER_BOOK                  = "returnOrderBook"
	CMD_PUBLIC_TRADE_HISTORY               = "returnTradeHistory"
	CMD_PUBLIC_CHART_DATA                  = "returnChartData"
	CMD_PUBLIC_CURRENCIES                  = "returnCurrencies"
	CMD_PUBLIC_LOAN_ORDERS                 = "returnLoanOrders"
	CMD_PRIVATE_BALANCES                   = "returnBalances"
	CMD_PRIVATE_COMPLETE_BALANCES          = "returnCompleteBalances"
	CMD_PRIVATE_DEPOSIT_ADDRESSES          = "returnDepositAddresses"
	CMD_PRIVATE_NEW_ADDRESS                = "generateNewAddress"
	CMD_PRIVATE_DEPOSIT_WITHDRAWALS        = "returnDepositsWithdrawals"
	CMD_PRIVATE_OPEN_ORDERS                = "returnOpenOrders"
	CMD_PRIVATE_TRADE_HISTORY              = "returnTradeHistory"
	CMD_PRIVATE_ORDER_TRADES               = "returnOrderTrades"
	CMD_PRIVATE_BUY                        = "buy"
	CMD_PRIVATE_SELL                       = "sell"
	CMD_PRIVATE_CANCEL_ORDER               = "cancelOrder"
	CMD_PRIVATE_MOVE_ORDER                 = "moveOrder"
	CMD_PRIVATE_WITHDRAW                   = "withdraw"
	CMD_PRIVATE_FEE_INFO                   = "returnFeeInfo"
	CMD_PRIVATE_AVAILABLE_ACCOUNT_BALANCES = "returnAvailableAccountBalances"
	// Margin
	CMD_PRIVATE_TRADABLE_BALANCES      = "returnTradableBalances"
	CMD_PRIVATE_TRANSFER_BALANCES      = "transferBalance"
	CMD_PRIVATE_MARGIN_ACCOUNT_SUMMARY = "returnMarginAccountSummary" // Todo
	CMD_PRIVATE_MARGIN_BUY             = "marginBuy"                  // Todo
	CMD_PRIVATE_MARGIN_SELL            = "marginSell"                 // Todo
	CMD_PRIVATE_MARGIN_POSITION        = "getMarginPosition"          // Todo
	CMD_PRIVATE_CLOSE_MARGIN_POSITION  = "closeMarginPosition"        // Todo
	// Loan
	CMD_PRIVATE_CREATE_LOAD_OFFER = "createLoanOffer"      // Todo
	CMD_PRIVATE_CANCEL_LOAD_OFFER = "cancelLoanOffer"      // Todo
	CMD_PRIVATE_OPEN_LOAD_OFFER   = "returnOpenLoanOffers" // Todo
	CMD_PRIVATE_ACTIVE_LOANS      = "returnActiveLoans"    // Todo
	CMD_PRIVATE_LENDING_HISTORY   = "returnLendingHistory" // Todo
	CMD_PRIVATE_TOGGLE_AUTO_RENEW = "toggleAutoRenew"      // Todo
)

type PoloniexApi struct {
	Key       string
	secret    string
	UserAgent string
	Client    *http.Client
}

func New(key string, secret string) *PoloniexApi {
	client := &http.Client{}
	user_agent := "poloniex-api"

	return &PoloniexApi{key, secret, user_agent, client}
}

/*
returnTicker
Returns the ticker for all markets. Sample output:

{"BTC_LTC":{"last":"0.0251","lowestAsk":"0.02589999","highestBid":"0.0251","percentChange":"0.02390438",
"baseVolume":"6.16485315","quoteVolume":"245.82513926"},"BTC_NXT":{"last":"0.00005730","lowestAsk":"0.00005710",
"highestBid":"0.00004903","percentChange":"0.16701570","baseVolume":"0.45347489","quoteVolume":"9094"}, ... }

Call: https://poloniex.com/public?command=returnTicker
*/
func (api *PoloniexApi) ApiPublicTicker() (map[string]Ticker, error) {
	params := url.Values{}
	params.Set("command", CMD_PUBLIC_TICKER)

	out := make(map[string]Ticker)

	_, err := api.queryparse(URL_PUBLIC, params, false, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
return24Volume
Returns the 24-hour volume for all markets, plus totals for primary currencies. Sample output:

{"BTC_LTC":{"BTC":"2.23248854","LTC":"87.10381314"},"BTC_NXT":{"BTC":"0.981616","NXT":"14145"},
 "totalBTC":"81.89657704","totalLTC":"78.52083806"}

Call: https://poloniex.com/public?command=return24hVolume
*/
func (api *PoloniexApi) ApiPublic24hVolume() (map[string]float64, map[string]map[string]float64, error) {
	params := url.Values{}
	params.Set("command", CMD_PUBLIC_24HVOLUME)

	content, err := api.queryparse(URL_PUBLIC, params, false, nil)
	if err != nil {
		return nil, nil, err
	}

	outTotal := make(map[string]float64)
	out := make(map[string]map[string]float64)

	for key, value := range content.(map[string]interface{}) {
		outsub := make(map[string]float64)

		// find out if string (total) or pair
		if reflect.TypeOf(value).Name() == "string" {
			volume, err := strconv.ParseFloat(value.(string), 64)
			if err != nil {
				return nil, nil, err
			}

			outTotal[key] = volume
		} else {
			for subkey, subvalue := range value.(map[string]interface{}) {
				volume, err := strconv.ParseFloat(subvalue.(string), 64)
				if err != nil {
					return nil, nil, err
				}

				outsub[subkey] = volume
			}
			out[key] = outsub
		}
	}

	return outTotal, out, nil
}

func getOrderBookFromInterface(in interface{}) (OrderBookEntry, error) {
	var err error
	var isFrozen int
	var seq float64
	var asks, bids [][2]float64

	for key, value := range in.(map[string]interface{}) {
		switch key {
		case "asks":
			asks, err = interfaceTo2FloatArray(value)
			if err != nil {
				return OrderBookEntry{}, err
			}

		case "bids":
			bids, err = interfaceTo2FloatArray(value)
			if err != nil {
				return OrderBookEntry{}, err
			}

		case "isFrozen":
			isFrozen, err = strconv.Atoi(value.(string))
			if err != nil {
				return OrderBookEntry{}, err
			}

		case "seq":
			seq = value.(float64)
		}
	}

	order := OrderBookEntry{
		IsFrozen: isFrozen,
		Seq:      seq,
		Asks:     asks,
		Bids:     bids,
	}

	return order, nil

}

/*
returnOrderBook

Returns the order book for a given market, as well as a sequence number for use with the Push API and an
indicator specifying whether the market is frozen. You may set currencyPair to "all" to get the order books
of all markets. Sample output:

{"asks":[[0.00007600,1164],[0.00007620,1300], ... ], "bids":[[0.00006901,200],[0.00006900,408], ... ],
 "isFrozen": 0, "seq": 18849}

Or, for all markets:

{"BTC_NXT":{"asks":[[0.00007600,1164],[0.00007620,1300], ... ], "bids":[[0.00006901,200],[0.00006900,408], ... ],
  "isFrozen": 0, "seq": 149},"BTC_XMR":...}

Call: https://poloniex.com/public?command=returnOrderBook&currencyPair=BTC_NXT&depth=10
*/
func (api *PoloniexApi) ApiPublicOrderBook(pair string, depth int) (map[string]OrderBookEntry, error) {
	params := url.Values{}
	params.Set("command", CMD_PUBLIC_ORDER_BOOK)
	params.Set("currencyPair", pair)

	if depth != 0 {
		params.Set("depth", strconv.Itoa(depth))
	}

	resp, err := api.query(URL_PUBLIC, params, false)
	if err != nil {
		return nil, err
	}

	content, err := api.parse(resp, nil)
	if err != nil {
		return nil, err
	}

	out := make(map[string]OrderBookEntry)

	if pair != "all" {
		order, err := getOrderBookFromInterface(content)
		if err != nil {
			return nil, err
		}
		out[pair] = order
	} else {
		for k, v := range content.(map[string]interface{}) {
			order, err := getOrderBookFromInterface(v)
			if err != nil {
				return nil, err
			}
			out[k] = order
		}
	}

	return out, nil
}

/*
returnTradeHistory
Returns the past 200 trades for a given market, or up to 50,000 trades between a range specified in UNIX
timestamps by the "start" and "end" GET parameters. Sample output:

[{"date":"2014-02-10 04:23:23","type":"buy","rate":"0.00007600","amount":"140","total":"0.01064"},
 {"date":"2014-02-10 01:19:37","type":"buy","rate":"0.00007600","amount":"655","total":"0.04978"}, ... ]

Call: https://poloniex.com/public?command=returnTradeHistory&currencyPair=BTC_NXT&start=1410158341&end=1410499372
*/
func (api *PoloniexApi) ApiPublicTradeHistory(pair string, start, end int) ([]Trade, error) {
	params := url.Values{}
	params.Set("command", CMD_PUBLIC_TRADE_HISTORY)
	params.Set("currencyPair", pair)

	if start != 0 {
		params.Set("start", strconv.Itoa(start))
	}

	if end != 0 {
		params.Set("end", strconv.Itoa(end))
	}

	out := make([]Trade, 0)

	_, err := api.queryparse(URL_PUBLIC, params, false, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
returnChartData
Returns candlestick chart data. Required GET parameters are "currencyPair", "period" (candlestick period in
seconds; valid values are 300, 900, 1800, 7200, 14400, and 86400), "start", and "end". "Start" and "end" are
given in UNIX timestamp format and used to specify the date range for the data returned. Sample output:

[{"date":1405699200,"high":0.0045388,"low":0.00403001,"open":0.00404545,"close":0.00427592,
 "volume":44.11655644,"quoteVolume":10259.29079097,"weightedAverage":0.00430015}, ...]

Call: https://poloniex.com/public?command=returnChartData&currencyPair=BTC_XMR&start=1405699200&end=9999999999&period=14400
*/
func (api *PoloniexApi) ApiChartData(pair string, start, end, period int64) ([]ChartEntry, error) {
	params := url.Values{}
	params.Set("command", CMD_PUBLIC_CHART_DATA)
	params.Set("currencyPair", pair)

	if start != 0 {
		params.Set("start", strconv.FormatInt(start, 10))
	}

	if end != 0 {
		params.Set("end", strconv.FormatInt(end, 10))
	}

	if period != 0 {
		params.Set("period", strconv.FormatInt(period, 10))
	}

	out := make([]ChartEntry, 0)

	_, err := api.queryparse(URL_PUBLIC, params, false, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
Returns information about currencies. Sample output:

{"1CR":{"maxDailyWithdrawal":10000,"txFee":0.01,"minConf":3,"disabled":0},
 "ABY":{"maxDailyWithdrawal":10000000,"txFee":0.01,"minConf":8,"disabled":0}, ... }

Call: https://poloniex.com/public?command=returnCurrencies
*/
func (api *PoloniexApi) ApiCurrencies() (map[string]Currency, error) {
	params := url.Values{}
	params.Set("command", CMD_PUBLIC_CURRENCIES)

	out := make(map[string]Currency)

	_, err := api.queryparse(URL_PUBLIC, params, false, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

/*
Returns the list of loan offers and demands for a given currency, specified by the "currency" GET parameter.
 Sample output:

{"offers":[{"rate":"0.00200000","amount":"64.66305732","rangeMin":2,"rangeMax":8}, ... ],
 "demands":[{"rate":"0.00170000","amount":"26.54848841","rangeMin":2,"rangeMax":2}, ... ]}

Call: https://poloniex.com/public?command=returnLoanOrders&currency=BTC
*/
func (api *PoloniexApi) ApiLoanOrders(currency string) (*LoanOrders, error) {
	params := url.Values{}
	params.Set("command", CMD_PUBLIC_LOAN_ORDERS)
	params.Set("currency", currency)

	out := new(LoanOrders)

	_, err := api.queryparse(URL_PUBLIC, params, false, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
