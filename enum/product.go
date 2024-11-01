package enum

import "strings"

type ProductStatus string

// Upper converts the type to uppercase
func (p ProductStatus) Upper() string {
	return strings.ToUpper(string(p))
}

// Lower converts the type to lowercase
func (p ProductStatus) Lower() string {
	return strings.ToLower(string(p))
}

// String converts the type to string
func (p ProductStatus) String() string {
	return string(p)
}

// Product Statuses
const (
	// Archived denotes a product archived by an admin
	Archived ProductStatus = "ARCHIVED"

	// Published denotes a product accepted by an admin
	Published ProductStatus = "PUBLISHED"

	// Delisted denotes a product delisted by an admin
	Delisted ProductStatus = "DELISTED"
)
