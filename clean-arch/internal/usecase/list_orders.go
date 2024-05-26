package usecase

import (
	"github.com/devfullcycle/20-CleanArch/internal/entity"
)

//type OrderInputDTO struct {
//ID    string  `json:"id"`
//Price float64 `json:"price"`
//Tax   float64 `json:"tax"`
//}

//type OrderOutputDTO struct {
//ID         string  `json:"id"`
//Price      float64 `json:"price"`
//Tax        float64 `json:"tax"`
//FinalPrice float64 `json:"final_price"`
//}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	//OrdersListed    events.EventInterface
	//EventDispatcher events.EventDispatcherInterface
}

func NewListOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	//OrdersListed events.EventInterface,
	//EventDispatcher events.EventDispatcherInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
		//OrdersListed:    OrdersListed,
		//EventDispatcher: EventDispatcher,
	}
}

func (c *ListOrdersUseCase) Execute() (*[]OrderOutputDTO, error) {
	l, err := c.OrderRepository.GetOrders()
	if err != nil {
		return nil, err
	}
	//dto := make(OrderOutputDTO,1)
	var dto []OrderOutputDTO
	for _, v := range *l {
		dto = append(dto, OrderOutputDTO{
			ID:         v.ID,
			Price:      v.Price,
			Tax:        v.Tax,
			FinalPrice: v.Price + v.Tax,
		},
		)
	}

	//c.OrdersListed.SetPayload(dto)
	//c.EventDispatcher.Dispatch(c.OrdersListed)

	return &dto, nil
}
