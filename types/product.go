package types

type Product struct {
	Id          int     `json:"id,omitempty"`
	ProductName string  `json:"productName,omitempty"`
	Price       float32 `json:"price,omitempty"`
}

func NewProduct(name string, price float32) *Product {
	return &Product{
		ProductName: name,
		Price:       price,
	}
}
