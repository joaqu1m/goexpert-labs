package product

type ProductUsecase struct {
	repo ProductRepositoryInterface
}

func NewProductUsecase(repo ProductRepositoryInterface) *ProductUsecase {
	return &ProductUsecase{repo: repo}
}

func (u *ProductUsecase) GetProductByID(id int) (*Product, error) {
	product, err := u.repo.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}
