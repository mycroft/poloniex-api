package poloniexapi

import (
	"encoding/json"
)

type Ticker struct {
	Id            int64
	Last          float64 `json:",string"`
	LowestAsk     float64 `json:",string"`
	HighestBid    float64 `json:",string"`
	PercentChange float64 `json:",string"`
	BaseVolume    float64 `json:",string"`
	QuoteVolume   float64 `json:",string"`
	IsFrozen      int     `json:",string"`
	High24hr      float64 `json:",string"`
	Low24hr       float64 `json:",string"`
}

type OrderBookEntry struct {
	IsFrozen int
	Seq      float64
	Asks     [][2]float64
	Bids     [][2]float64
}

type Trade struct {
	GlobalTradeID int64   `json:"globalTradeID"`
	TradeID       int64   `json:"tradeID"`
	Date          string  `json:"date"`
	Type          string  `json:"type"`
	Rate          float64 `json:"rate,string"`
	Amount        float64 `json:"amount,string"`
	Total         float64 `json:"total,string"`
	Fee           float64 `json:"fee,string"`
	OrderNumber   int64   `json:"orderNumber,string"`
	Category      string  `json:"category"`
}

type ChartEntry struct {
	Date            int64
	High            float64
	Low             float64
	Open            float64
	Close           float64
	Volume          float64
	QuoteVolume     float64
	WeightedAverage float64
}

type Currency struct {
	Delisted int64
	Frozen   int64
	Disabled int64
	Id       int64
	Name     string
	TxFee    float64 `json:",string"`
	MinConf  int64
}

type LoanOrder struct {
	Amount   float64 `json:",string"`
	RangeMax int64
	RangeMin int64
	Rate     float64 `json:",string"`
}

type LoanOrders struct {
	Demands []LoanOrder
	Offers  []LoanOrder
}

type Balance struct {
	Available float64 `json:"available,string"`
	OnOrders  float64 `json:"onOrders,string"`
	BtcValue  float64 `json:"btcValue,string"`
}

type GenerateAddressResponse struct {
	Success  int
	Response string
}

type Deposit struct {
	Currency      string
	Address       string
	Amount        float64 `json:"amount,string"`
	Confirmations int64
	Txid          string
	Timestamp     int64
	Status        string
}

type Withdrawal struct {
	WithdrawalNumber int64
	Currency         string
	Address          string
	Amount           float64 `json:"amount,string"`
	Timestamp        int64
	Status           string
	IpAddress        string
}

type DepositWithdrawal struct {
	Deposits    []Deposit
	Withdrawals []Withdrawal
}

type OpenOrder struct {
	OrderNumber    string  `json:"orderNumber"`
	Type           string  `json:"type"`
	Rate           float64 `json:"rate,string"`
	StartingAmount float64 `json:"startingAmount,string"`
	Amount         float64 `json:"amount,string"`
	Total          float64 `json:"total,string"`
	Date           string  `json:"date"`
	Margin         int64   `json:"margin"`
}

type Order struct {
	Success         int64              `json:"success"` // Use for moveOrder
	OrderNumber     int64              `json:"orderNumber"`
	ResultingTrades map[string][]Trade `json:"resultingTrades"`
}

type CancelOrder struct {
	Success int64   `json:"success"`
	Error   string  `json:"error"`
	Amount  float64 `json:"amount,string"`
	Message string  `json:"message"`
}

type FeeInfo struct {
	MakerFee        float64 `json:"makerFee,string"`
	TakerFee        float64 `json:"takerFee,string"`
	ThirtyDayVolume float64 `json:"thirtyDayVolume,string"`
	NextTier        float64 `json:"nextTier,string"`
}

func (t *Trade) UnmarshalJSON(data []byte) error {
	var err error

	type Alias Trade
	aux := &struct {
		TradeID json.Number `json:"tradeID"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.TradeID, err = aux.TradeID.Int64()
	if err != nil {
		return err
	}

	return nil
}

func (t *Order) UnmarshalJSON(data []byte) error {
	var err error

	type Alias Order
	aux := &struct {
		OrderNumber json.Number `json:"orderNumber"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.OrderNumber, err = aux.OrderNumber.Int64()
	if err != nil {
		return err
	}

	return nil
}
