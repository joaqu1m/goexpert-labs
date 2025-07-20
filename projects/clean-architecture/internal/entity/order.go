package entity

import "errors"

type Order struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

func NewOrder(id string, price, tax float64) (*Order, error) {
	order := Order{
		ID:    id,
		Price: price,
		Tax:   tax,
	}
	if err := order.IsValid(); err != nil {
		return nil, err
	}
	return &order, nil
}

func (o *Order) IsValid() error {
	if o.ID == "" {
		return errors.New("order ID cannot be empty")
	}
	if o.Price <= 0 {
		return errors.New("order price must be greater than zero")
	}
	if o.Tax < 0 {
		return errors.New("order tax cannot be negative")
	}
	return nil
}

func (o *Order) CalculateFinalPrice() {
	o.FinalPrice = o.Price + o.Tax
}
