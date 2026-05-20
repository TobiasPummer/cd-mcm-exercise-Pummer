package store

import (
	"errors"
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()

	product := model.Product{
		Name:  "Laptop",
		Price: 999.99,
	}

	created := s.Create(product)

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if got.ID != created.ID {
		t.Errorf("ID = %v, want %v", got.ID, created.ID)
	}
	if got.Name != product.Name {
		t.Errorf("Name = %v, want %v", got.Name, product.Name)
	}
	if got.Price != product.Price {
		t.Errorf("Price = %v, want %v", got.Price, product.Price)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()

	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestUpdateProduct(t *testing.T) {
	s := NewMemoryStore()

	product := model.Product{
		Name:  "Old name",
		Price: 10.0,
	}

	created := s.Create(product)

	updatedInput := model.Product{
		ID:    created.ID,
		Name:  "New name",
		Price: 20.0,
	}

	updated, err := s.Update(created.ID, updatedInput)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if updated.Name != "New name" {
		t.Errorf("Name = %v, want %v", updated.Name, "New name")
	}
	if updated.Price != 20.0 {
		t.Errorf("Price = %v, want %v", updated.Price, 20.0)
	}
}

func TestDeleteProduct(t *testing.T) {
	s := NewMemoryStore()

	product := model.Product{
		Name:  "To delete",
		Price: 15.0,
	}

	created := s.Create(product)

	err := s.Delete(created.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = s.GetByID(created.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()

	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := NewMemoryStore()

	tests := []struct {
		name string
		id   int
	}{
		{"unknown id", 999},
		{"zero id", 0},
		{"negative id", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.GetByID(tt.id)
			if !errors.Is(err, ErrNotFound) {
				t.Errorf("GetByID(%d) error = %v, want %v", tt.id, err, ErrNotFound)
			}
		})
	}
}
