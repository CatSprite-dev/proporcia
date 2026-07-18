package api

type AccountStatus string

const (
	AccountStatusUnspecified AccountStatus = "ACCOUNT_STATUS_UNSPECIFIED"
	AccountStatusNew         AccountStatus = "ACCOUNT_STATUS_NEW"
	AccountStatusOpen        AccountStatus = "ACCOUNT_STATUS_OPEN"
	AccountStatusClosed      AccountStatus = "ACCOUNT_STATUS_CLOSED"
	AccountStatusAll         AccountStatus = "ACCOUNT_STATUS_ALL"
)
