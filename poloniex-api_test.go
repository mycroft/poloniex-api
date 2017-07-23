package poloniexapi

import (
	"encoding/json"
	"io/ioutil"

	"log"
	"testing"
	"time"
)

type Config struct {
	Key    string
	Secret string
}

var api = CreatePrivateApiClient()

func LoadConfiguration(path string, config interface{}) (interface{}, error) {
	content, err := ioutil.ReadFile(path)
	CheckErr(err)

	err = json.Unmarshal(content, &config)
	CheckErr(err)

	return config, nil
}

func CreatePrivateApiClient() *PoloniexApi {
	var config Config

	_, err := LoadConfiguration("config.json", &config)
	CheckErr(err)

	return New(config.Key, config.Secret)
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func TestApiPublicTicker(t *testing.T) {
	resp, err := api.ApiPublicTicker()
	CheckErr(err)

	for key, value := range resp {
		log.Printf("%10s (%3d) %8.3f b:%8.3f a:%8.3f p:%.3f v:%8.3f q:%.3f h:%.3f l:%.3f\n",
			key,
			value.Id,
			value.Last,
			value.LowestAsk,
			value.HighestBid,
			value.PercentChange,
			value.BaseVolume,
			value.QuoteVolume,
			value.High24hr,
			value.Low24hr,
		)
	}
}

func TestApiPublic24hVolume(t *testing.T) {
	resp_total, resp, err := api.ApiPublic24hVolume()
	CheckErr(err)

	for key, value := range resp_total {
		log.Printf("%s: %f\n", key, value)
	}

	for key, value := range resp {
		for subkey, value := range value {
			log.Printf("%s: %s: %f\n", key, subkey, value)
		}
	}
}

func TestApiPublicOrderBook(t *testing.T) {
	order_books, err := api.ApiPublicOrderBook("BTC_NXT", 10)
	CheckErr(err)

	for pair, order_book := range order_books {
		log.Printf("Pair: %s - IsFrozen: %v SeqNum: %f\n", pair, order_book.IsFrozen, order_book.Seq)
		log.Println("Asks:", order_book.Asks)
		log.Println("Bids", order_book.Bids)
	}

	order_books, err = api.ApiPublicOrderBook("all", 10)
	CheckErr(err)

	for pair, order_book := range order_books {
		log.Printf("Pair: %s - IsFrozen: %v SeqNum: %f\n", pair, order_book.IsFrozen, order_book.Seq)
		log.Println("Asks:", order_book.Asks)
		log.Println("Bids:", order_book.Bids)
	}
}

func TestApiPublicTradeHistory(t *testing.T) {
	history, err := api.ApiPublicTradeHistory("BTC_NXT", 0, 0)
	CheckErr(err)

	log.Println(history)
}

func TestApiChartData(t *testing.T) {
	now := time.Now().Unix()
	onedayago := now - 24*60*60

	data, err := api.ApiChartData("BTC_NXT", onedayago, now, 300)
	CheckErr(err)

	for k, v := range data {
		log.Printf("%d: date:%d %f -> [o:%f / c:%f] -> %f (vol:%f quote:%f wa:%f)\n", k, v.Date, v.High, v.Open, v.Close, v.Low, v.Volume, v.QuoteVolume, v.WeightedAverage)
	}
}

func TestApiCurrencies(t *testing.T) {
	data, err := api.ApiCurrencies()
	CheckErr(err)

	for k, v := range data {
		log.Printf("%d %s %s txfee: %f min conf: %d disabled: %d\n",
			v.Id,
			k,
			v.Name,
			v.TxFee,
			v.MinConf,
			v.Disabled,
		)
	}
}

func TestApiLoanOrders(t *testing.T) {
	data, err := api.ApiLoanOrders("BTC")
	CheckErr(err)

	log.Println("Demands for BTC...")
	for _, v := range data.Demands {
		log.Printf("%f %d %d %f\n", v.Amount, v.RangeMax, v.RangeMin, v.Rate)
	}

	log.Println("Offers for BTC...")
	for _, v := range data.Offers {
		log.Printf("%f %d %d %f\n", v.Amount, v.RangeMax, v.RangeMin, v.Rate)
	}
}

func TestApiPrivateBalances(t *testing.T) {
	balances, err := api.ApiPrivateBalances()
	CheckErr(err)

	for key, value := range balances {
		if value != 0.0 {
			log.Printf("%s %f\n", key, value)
		}
	}
}

func TestApiPrivateCompleteBalances(t *testing.T) {
	balances, err := api.ApiPrivateCompleteBalances(true)
	CheckErr(err)

	for k, v := range balances {
		if v.Available == 0.0 && v.OnOrders == 0.0 && v.BtcValue == 0.0 {
			continue
		}
		log.Printf("%s: Available:%f OnOrder:%f BtcValue:%f\n",
			k,
			v.Available,
			v.OnOrders,
			v.BtcValue,
		)
	}
}

func TestApiPrivateDepositAddresses(t *testing.T) {
	addresses, err := api.ApiPrivateDepositAddresses()
	CheckErr(err)

	for k, v := range addresses {
		log.Printf("%s: %s\n", k, v)
	}
}

func TestApiPrivateGenerateNewAddress(t *testing.T) {
	resp, err := api.ApiPrivateGenerateNewAddress("XMR")
	CheckErr(err)

	log.Println(resp)
}

func TestApiPrivateDepositWithDrawals(t *testing.T) {
	depositsWithdrawals, err := api.ApiPrivateDepositWithdrawals(0, time.Now().Unix())
	CheckErr(err)

	log.Println(depositsWithdrawals)
}

func TestApiPrivateOpenOrdersSingle(t *testing.T) {
	out, err := api.ApiPrivateOpenOrders("XMR_LTC")
	CheckErr(err)

	log.Println(out)
}

func TestApiPrivateOpenOrdersAll(t *testing.T) {
	out, err := api.ApiPrivateOpenOrders("all")
	CheckErr(err)

	for k, v := range out {
		if len(v) == 0 {
			continue
		}

		log.Println(k)
		for _, order := range v {
			log.Println(order)
		}
	}
}

func TestApiPrivateTradeHistorySingle(t *testing.T) {
	out, err := api.ApiPrivateTradeHistory("XMR_LTC", 1497450844, 1497623644)
	CheckErr(err)

	log.Println(out)
}

func TestApiPrivateTradeHistoryAll(t *testing.T) {
	out, err := api.ApiPrivateTradeHistory("all", 1497450844, 1497623644)
	CheckErr(err)

	log.Println(out)
}

func TestApiPrivateOrderTrades(t *testing.T) {
	out, err := api.ApiPrivateOrderTrades("28029530549")
	CheckErr(err)

	log.Println(out)
}

func TestApiPrivateBuy(t *testing.T) {
	out, err := api.ApiPrivateBuy("XMR_LTC", 0.001, 42, map[string]bool{})
	CheckErr(err)

	log.Println(out)
}

func TestApiPrivateSell(t *testing.T) {
	out, err := api.ApiPrivateSell("BTC_XMR", 4999, 0.001, map[string]bool{})
	CheckErr(err)

	log.Println(out)
}

func TestApiPrivateCancel(t *testing.T) {
	success, out, err := api.ApiPrivateCancel(179587130217)
	CheckErr(err)

	log.Println(success)
	log.Println(out)
}

func TestApiPrivateMoveOrder(t *testing.T) {
	out, err := api.ApiPrivateBuy("XMR_LTC", 0.001, 42, map[string]bool{})
	CheckErr(err)

	log.Println(out)

	out, err = api.ApiPrivateMoveOrder(out.OrderNumber, 0.0005, 42.5, map[string]bool{})
	CheckErr(err)

	log.Println(out)
}

func TestApiPrivateFeeInfo(t *testing.T) {
	out, err := api.ApiPrivateFeeInfo()
	CheckErr(err)

	log.Println(out)
}

func TestApiPrivateAvailableAccountBalances(t *testing.T) {
	out, err := api.ApiPrivateAvailableAccountBalances("")
	CheckErr(err)

	log.Println(out)
}

func TestApiPrivateTradableBalances(t *testing.T) {
	out, err := api.ApiPrivateTradableBalances()
	CheckErr(err)

	log.Println(out)
}
