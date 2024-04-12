package repository

import (
	"cashier/model"
	"cashier/utils"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type ProductRepository interface {
	Create(newProduct model.ProductModel) error
	ProductCheckName(name string) (bool, error)
	GetAll(ctx context.Context, page, pageSize int, search string) ([]*model.ProductModel, int, error)
	FindById(id string) (*model.ProductModel, error)
	CheckNameForUpdate(nama, id string) (bool, error)
	Update(id string, updated model.ProductModel) error
	Delete(ID string) error
	UploadProduct(productImage model.ProductImage) error
}

type productRepository struct {
	db *sql.DB
}

func (pr *productRepository) Create(newProduct model.ProductModel) (error) {
	newProduct.Id = utils.UuidGenerate()

	insertQuery := "INSERT INTO products(id, product_name, thumbnail, price, status, categoty_id, created_at, updated_at) VALUES($1,$2,$3,$4,$5, $6, $7, $8) RETURNING id"
	var id string
	err := pr.db.QueryRow(insertQuery, newProduct.Id, newProduct.ProductName, newProduct.Thumbnail, newProduct.Price, newProduct.Status, newProduct.CategoryId, newProduct.CreatedAt, newProduct.UpdatedAt).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("err on productRepository.Create: %v", err)
	}
	return nil
}

func (pr *productRepository) GetAll(ctx context.Context, page, pageSize int, search string) ([]*model.ProductModel, int, error) {
	offset := (page - 1) * pageSize
	getQuery :=
		`SELECT 
		p.id, 
		p.product_name, 
		p.price, 
		p.status, 
		p.category_id, 
		p.thumbnail , 
		p.category_id,
		c.category_name
	FROM 
		products p
	INNER JOIN 
		categories c ON p.category_id = c.id
	WHERE 
		p.product_name LIKE '%' || $1 || '%'
	ORDER BY 
		p.id 
	LIMIT 
		$2 
	OFFSET 
		$3;
	`

	rows, err := pr.db.QueryContext(ctx, getQuery, search, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error on productRepository.GetAll: %w", err)
	}
	defer rows.Close()

	var arrProduct []*model.ProductModel

	for rows.Next() {
		product := &model.ProductModel{}
		category := &model.CategoryModel{}
		if err := rows.Scan(
			&product.Id, 
			&product.ProductName, 
			&product.Price, 
			&product.Status, 
			&product.CategoryId,
			&product.Thumbnail,
			&category.Id, 
			&category.CategoryName, 
		); err != nil {
			return nil, 0, fmt.Errorf("error scanning campaign row: %w", err)
		}
		product.Categories = append(product.Categories, category)
		arrProduct = append(arrProduct, product)
	}

	countQuery := "SELECT COUNT(id) FROM products" // Menyesuaikan query COUNT dengan offset
	var totalDocs int
	if err := pr.db.QueryRowContext(ctx, countQuery).Scan(&totalDocs); err != nil {
		return nil, 0, fmt.Errorf("error on pr.db.QueryRowContext.countQuery: %w", err)
	}
	return arrProduct, totalDocs, nil
}

func (pr *productRepository) ProductCheckName(name string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM products WHERE product_name = $1)"

	var exists bool
	err := pr.db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error on productRepository.ProductCheckName: %w", err)
	}

	return exists, nil
}

func (pr *productRepository) FindById(id string) (*model.ProductModel, error) {
	getQuery := "SELECT id, product_name, thumbnail, price, status, category_id, created_at, updated_at FROM products WHERE id=$1"

	row := pr.db.QueryRow(getQuery, id)

	product := &model.ProductModel{}
	err := row.Scan(&product.Id, &product.ProductName, &product.Thumbnail, &product.Price, &product.CategoryId, &product.Status, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		fmt.Println(err)
		return product, fmt.Errorf("error on productRepository. FindById: %w", err)
	}

	return product, nil
}

func (pr *productRepository) CheckNameForUpdate(nama, id string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM products WHERE product_name = $1 AND id != $2)"

	var exists bool
	err := pr.db.QueryRow(query, nama, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking existing product name for update: %w", err)
	}

	return exists, nil
}

func (pr *productRepository) Update(id string, updated model.ProductModel) error {
	updateQuery :=
		`UPDATE products SET 
		product_name = $1, 
		price = $2,
		status = $3,
		category_id = $4,
		updated_at = $5,
		thumbnail = $6
		WHERE id = $7`

	_, err := pr.db.Exec(
		updateQuery,
		updated.ProductName,
		updated.Price,
		updated.Status,
		updated.CategoryId,
		time.Now(),
		updated.Thumbnail,
		id)
	if err != nil {
		return fmt.Errorf("error updating products: %w", err)
	}

	return nil
}

func (pr *productRepository) Delete(ID string) error {
	deleteQuery := "DELETE FROM products WHERE id = $1"

	_, err := pr.db.Query(deleteQuery, ID)
	if err != nil {
		return fmt.Errorf("error fetching products by Id: %w", err)
	}
	return nil
}

func (pr *productRepository) UploadProduct(productImage model.ProductImage) error{
	productImage.Id = utils.UuidGenerate()
	insertQuery := "INSERT INTO product_images (id , product_id, product_name) VALUES ($1, $2, $3)"
	_, err := pr.db.Exec(insertQuery, productImage.Id, productImage.ProductId, productImage.ProductFile)
	if err != nil {
		return fmt.Errorf("failed to insert product image in database: %v", err)
	}
	return nil
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}
