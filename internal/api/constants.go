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

type lastPriceType string

const (
	LastPriceUnspecified lastPriceType = "LAST_PRICE_UNSPECIFIED"
	LastPriceExchange    lastPriceType = "LAST_PRICE_EXCHANGE"
	LastPriceDealer      lastPriceType = "LAST_PRICE_DEALER"
)

type instrumentStatus string

const (
	InstrumentStatusUnspecified instrumentStatus = "INSTRUMENT_STATUS_UNSPECIFIED"
	InstrumentStatusBase        instrumentStatus = "INSTRUMENT_STATUS_BASE"
	InstrumentStatusAll         instrumentStatus = "INSTRUMENT_STATUS_ALL"
)
