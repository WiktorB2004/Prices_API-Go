package app

// Product model
type Product struct {
	ID         string `json:"id" bson:"_id"`
	Name       string `json:"productName" bson:"productName"`
	SupplierID string `json:"supplierId" bson:"supplierId"`
	Price      int    `json:"price" bson:"price"`
}
