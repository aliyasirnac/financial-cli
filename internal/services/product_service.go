package services

import (
	"context"
	"github.com/aliyasirnac/financialManagement/internal/db"
	"github.com/aliyasirnac/financialManagement/internal/model"
	"github.com/uptrace/bun"
)

type ProductService struct {
	db *bun.DB
}

func NewProductService() (*ProductService, error) {
	databaseConfig := db.NewDatabase()
	database := databaseConfig.OpenDatabase()

	service := &ProductService{
		db: database,
	}

	// Create products table if not exists
	ctx := context.Background()
	_, err := service.db.NewCreateTable().Model((*model.Product)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ProductService) CreateProduct(product *model.Product) error {
	_, err := s.db.NewInsert().Model(product).Exec(context.Background())
	return err
}

func (s *ProductService) GetProducts() ([]model.Product, error) {
	var products []model.Product
	err := s.db.NewSelect().Model(&products).Scan(context.Background())
	return products, err
}

func (s *ProductService) UpdateProduct(product *model.Product) error {
	_, err := s.db.NewUpdate().Model(product).WherePK().Exec(context.Background())
	return err
}

func (s *ProductService) DeleteProduct(id uint64) error {
	product := &model.Product{ID: id}
	_, err := s.db.NewDelete().Model(product).WherePK().Exec(context.Background())
	return err
}

func (s *ProductService) IncrementProductCount(id uint64, count uint) error {
	product := &model.Product{ID: id}
	err := s.db.NewSelect().Model(product).WherePK().Scan(context.Background())
	if err != nil {
		return err
	}
	product.Count += count
	_, err = s.db.NewUpdate().Model(product).WherePK().Exec(context.Background())
	return err
}

func (s *ProductService) CalculateTotalSpending() (uint16, error) {
	var products []model.Product
	err := s.db.NewSelect().Model(&products).Scan(context.Background())
	if err != nil {
		return 0, err
	}
	var total uint16
	for _, product := range products {
		total += product.Price * uint16(product.Count)
	}
	return total, nil
}
