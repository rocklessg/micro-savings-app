package models

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
	Transfer   TransactionType = "transfer"
	Investment TransactionType = "investment"
)

// IsValid checks if a transaction type is valid
func (t TransactionType) IsValid() bool {
	switch t {
	case Deposit, Withdrawal, Transfer, Investment:
		return true
	default:
		return false
	}
}