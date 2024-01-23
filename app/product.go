package app

// Product model
type Product struct {
	ID         string `json:"id" bson:"_id"`
	Name       string `json:"productName" bson:"productName"`
	SupplierID string `json:"supplierId" bson:"supplierId"`
	Price      int    `json:"price" bson:"price"`
}

type ExtProduct struct {
	ID           string `json:"id" bson:"_id"`
	Name         string `json:"productName" bson:"productName"`
	SupplierName string `json:"supplierName" bson:"supplierName"`
	Price        int    `json:"price" bson:"price"`
}

type ProductResponse struct {
	Message         string  `json:"message"`
	ExistingProduct Product `json:"existingProduct"`
}
