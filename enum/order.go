package enum

type OrderStatus string

func (o OrderStatus) String() string {
	return string(o)
}

// Currency describes the currency of a transaction
const (
	// pending
	Pending OrderStatus = "PENDING"

	// canceled
	Canceled OrderStatus = "CANCELED"

	// completed
	Completed OrderStatus = "COMPLETED"
)
