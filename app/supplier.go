package app

// Supplier model
type Supplier struct {
	ID            string   `json:"id" bson:"_id"`
	Name          string   `json:"supplierName" bson:"supplierName"`
	Phone         string   `json:"phoneNumber" bson:"phoneNumber"`
	Email         string   `json:"email" bson:"email"`
	Products      []string `json:"products" bson:"products"`
	ProductsCount int      `json:"productCount" bson:"productCount"`
}

type SupplierResponse struct {
	Message         string   `json:"message"`
	ExistingProduct Supplier `json:"existingSupplier"`
}
