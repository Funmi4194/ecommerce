package enum

type OrderStatus string

func (o OrderStatus) String() string {
	return string(o)
}

// Currency describes the currency of a transaction
const (
	// pending
	Pending OrderStatus = "PENDING"

	// Cancelled denotes an order cancelled                                       by an admin and user
	Cancelled OrderStatus = "CANCELLED"

	// Approved denotes an order completed by an admin
	Completed OrderStatus = "COMPLETED"

	// Rejected demotes an order rejected by an admin
	Rejected OrderStatus = "REJECTED"

	// Approved denotes an order approved by an admin
	Approved OrderStatus = "APPROVED"
)
