package api

type AccountStatus string

const (
	AccountStatusUnspecified AccountStatus = "ACCOUNT_STATUS_UNSPECIFIED"
	AccountStatusNew         AccountStatus = "ACCOUNT_STATUS_NEW"
	AccountStatusOpen        AccountStatus = "ACCOUNT_STATUS_OPEN"
	AccountStatusClosed      AccountStatus = "ACCOUNT_STATUS_CLOSED"
	AccountStatusAll         AccountStatus = "ACCOUNT_STATUS_ALL"
)

type InstrumentType string

const (
	InstrumentTypeUnspecified         InstrumentType = "INSTRUMENT_TYPE_UNSPECIFIED"
	InstrumentTypeBond                InstrumentType = "INSTRUMENT_TYPE_BOND"
	InstrumentTypeShare               InstrumentType = "INSTRUMENT_TYPE_SHARE"
	InstrumentTypeCurrency            InstrumentType = "INSTRUMENT_TYPE_CURRENCY"
	InstrumentTypeETF                 InstrumentType = "INSTRUMENT_TYPE_ETF"
	InstrumentTypeFutures             InstrumentType = "INSTRUMENT_TYPE_FUTURES"
	InstrumentTypeSP                  InstrumentType = "INSTRUMENT_TYPE_SP"
	InstrumentTypeOption              InstrumentType = "INSTRUMENT_TYPE_OPTION"
	InstrumentTypeClearingCertificate InstrumentType = "INSTRUMENT_TYPE_CLEARING_CERTIFICATE"
	InstrumentTypeIndex               InstrumentType = "INSTRUMENT_TYPE_INDEX"
	InstrumentTypeCommodity           InstrumentType = "INSTRUMENT_TYPE_COMMODITY"
)

type LastPriceType string

const (
	LastPriceUnspecified LastPriceType = "LAST_PRICE_UNSPECIFIED"
	LastPriceExchange    LastPriceType = "LAST_PRICE_EXCHANGE"
	LastPriceDealer      LastPriceType = "LAST_PRICE_DEALER"
)

type InstrumentStatus string

const (
	InstrumentStatusUnspecified InstrumentStatus = "INSTRUMENT_STATUS_UNSPECIFIED"
	InstrumentStatusBase        InstrumentStatus = "INSTRUMENT_STATUS_BASE"
	InstrumentStatusAll         InstrumentStatus = "INSTRUMENT_STATUS_ALL"
)

type OrderDirection string

const (
	OrderDirectionUnspecified OrderDirection = "ORDER_DIRECTION_UNSPECIFIED"
	OrderDirectionBuy         OrderDirection = "ORDER_DIRECTION_BUY"
	OrderDirectionSell        OrderDirection = "ORDER_DIRECTION_SELL"
)

type OrderType string

const (
	OrderTypeUnspecified OrderType = "ORDER_TYPE_UNSPECIFIED"
	OrderTypeLimit       OrderType = "ORDER_TYPE_LIMIT"
	OrderTypeMarket      OrderType = "ORDER_TYPE_MARKET"
	OrderTypeBestPrice   OrderType = "ORDER_TYPE_BESTPRICE"
)

type TimeInForce string

const (
	TimeInForceUnspecified TimeInForce = "TIME_IN_FORCE_UNSPECIFIED"
	TimeInForceDay         TimeInForce = "TIME_IN_FORCE_DAY"
	TimeInForceFillAndKill TimeInForce = "TIME_IN_FORCE_FILL_AND_KILL"
	TimeInForceFillOrKill  TimeInForce = "TIME_IN_FORCE_FILL_OR_KILL"
)

type PriceType string

const (
	PriceTypeUnspecified PriceType = "PRICE_TYPE_UNSPECIFIED"
	PriceTypePoint       PriceType = "PRICE_TYPE_POINT"
	PriceTypeCurrency    PriceType = "PRICE_TYPE_CURRENCY"
)
