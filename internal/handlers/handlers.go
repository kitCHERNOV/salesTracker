package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"salesTracker/internal/storage/postgresql"
)

// ====================================================================
// REQUEST/RESPONSE DTOs
// ====================================================================

// CategoryRequest - DTO для создания/обновления категории
type CategoryRequest struct {
	CategoryName string `json:"category_name"`
	Description  string `json:"description"`
}

// ProductRequest - DTO для создания/обновления товара
type ProductRequest struct {
	ProductName   string  `json:"product_name"`
	CategoryID    int     `json:"category_id"`
	Price         float64 `json:"price"`
	Cost          float64 `json:"cost"`
	StockQuantity int     `json:"stock_quantity"`
}

// CustomerRequest - DTO для создания/обновления покупателя
type CustomerRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	City      string `json:"city"`
}

// OrderRequest - DTO для создания/обновления заказа
type OrderRequest struct {
	CustomerID    int     `json:"customer_id"`
	OrderDate     string  `json:"order_date"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
	TotalAmount   float64 `json:"total_amount"`
}

// OrderItemRequest - DTO для создания/обновления позиции заказа
type OrderItemRequest struct {
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Discount  float64 `json:"discount"`
}

// ====================================================================
// HELPERS
// ====================================================================

func parseURLParamID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	return strconv.Atoi(idStr)
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func respondError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, map[string]string{"error": message})
}

// ====================================================================
// CATEGORIES HANDLERS
// ====================================================================

// CreateCategory - создать категорию
func CreateCategory(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CategoryRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid request body")
			return
		}

		id, err := storage.AddCategory(req.CategoryName, req.Description)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		category, err := storage.GetCategory(id)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, category)
	}
}

// GetCategory - получить категорию по ID
func GetCategory(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid category id")
			return
		}

		category, err := storage.GetCategory(id)
		if err != nil {
			respondError(w, r, http.StatusNotFound, "category not found")
			return
		}

		render.JSON(w, r, category)
	}
}

// ListCategories - получить список всех категорий
func ListCategories(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := storage.ListCategories()
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, categories)
	}
}

// UpdateCategory - обновить категорию
func UpdateCategory(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid category id")
			return
		}

		var req CategoryRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := storage.UpdateCategory(id, req.CategoryName, req.Description); err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		category, err := storage.GetCategory(id)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, category)
	}
}

// DeleteCategory - удалить категорию
func DeleteCategory(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid category id")
			return
		}

		if err := storage.DeleteCategory(id); err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.Status(r, http.StatusNoContent)
	}
}

// ====================================================================
// PRODUCTS HANDLERS
// ====================================================================

// CreateProduct - создать товар
func CreateProduct(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ProductRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid request body")
			return
		}

		id, err := storage.AddProduct(req.ProductName, req.CategoryID, req.Price, req.Cost, req.StockQuantity)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		product, err := storage.GetProduct(id)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, product)
	}
}

// GetProduct - получить товар по ID
func GetProduct(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid product id")
			return
		}

		product, err := storage.GetProduct(id)
		if err != nil {
			respondError(w, r, http.StatusNotFound, "product not found")
			return
		}

		render.JSON(w, r, product)
	}
}

// ListProducts - получить список всех товаров
func ListProducts(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := storage.ListProducts()
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, products)
	}
}

// ListProductsByCategory - получить товары по категории
func ListProductsByCategory(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryID, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid category id")
			return
		}

		products, err := storage.ListProductsByCategory(categoryID)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, products)
	}
}

// UpdateProduct - обновить товар
func UpdateProduct(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid product id")
			return
		}

		var req ProductRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := storage.UpdateProduct(id, req.ProductName, req.CategoryID, req.Price, req.Cost, req.StockQuantity); err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		product, err := storage.GetProduct(id)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, product)
	}
}

// DeleteProduct - удалить товар
func DeleteProduct(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid product id")
			return
		}

		if err := storage.DeleteProduct(id); err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.Status(r, http.StatusNoContent)
	}
}

// ====================================================================
// CUSTOMERS HANDLERS
// ====================================================================

// CreateCustomer - создать покупателя
func CreateCustomer(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CustomerRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid request body")
			return
		}

		id, err := storage.AddCustomer(req.FirstName, req.LastName, req.Email, req.Phone, req.City, time.Now())
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		customer, err := storage.GetCustomer(id)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, customer)
	}
}

// GetCustomer - получить покупателя по ID
func GetCustomer(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid customer id")
			return
		}

		customer, err := storage.GetCustomer(id)
		if err != nil {
			respondError(w, r, http.StatusNotFound, "customer not found")
			return
		}

		render.JSON(w, r, customer)
	}
}

// ListCustomers - получить список всех покупателей
func ListCustomers(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customers, err := storage.ListCustomers()
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, customers)
	}
}

// UpdateCustomer - обновить покупателя
func UpdateCustomer(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid customer id")
			return
		}

		var req CustomerRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := storage.UpdateCustomer(id, req.FirstName, req.LastName, req.Email, req.Phone, req.City); err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		customer, err := storage.GetCustomer(id)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, customer)
	}
}

// DeleteCustomer - удалить покупателя
func DeleteCustomer(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid customer id")
			return
		}

		if err := storage.DeleteCustomer(id); err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.Status(r, http.StatusNoContent)
	}
}

// ====================================================================
// ORDERS HANDLERS
// ====================================================================

// CreateOrder - создать заказ
func CreateOrder(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req OrderRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid request body")
			return
		}

		orderDate, err := parseDate(req.OrderDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid order date format, use YYYY-MM-DD")
			return
		}

		id, err := storage.AddOrder(req.CustomerID, orderDate, req.Status, req.PaymentMethod, req.TotalAmount)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		order, err := storage.GetOrder(id)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, order)
	}
}

// GetOrder - получить заказ по ID
func GetOrder(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid order id")
			return
		}

		order, err := storage.GetOrder(id)
		if err != nil {
			respondError(w, r, http.StatusNotFound, "order not found")
			return
		}

		render.JSON(w, r, order)
	}
}

// ListOrders - получить список всех заказов
func ListOrders(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := storage.ListOrders()
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, orders)
	}
}

// ListOrdersByCustomer - получить заказы покупателя
func ListOrdersByCustomer(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerID, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid customer id")
			return
		}

		orders, err := storage.ListOrdersByCustomer(customerID)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, orders)
	}
}

// UpdateOrder - обновить заказ
func UpdateOrder(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid order id")
			return
		}

		var req struct {
			Status      string  `json:"status"`
			TotalAmount float64 `json:"total_amount"`
		}
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := storage.UpdateOrder(id, req.Status, req.TotalAmount); err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		order, err := storage.GetOrder(id)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, order)
	}
}

// DeleteOrder - удалить заказ
func DeleteOrder(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid order id")
			return
		}

		if err := storage.DeleteOrder(id); err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.Status(r, http.StatusNoContent)
	}
}

// ====================================================================
// ORDER ITEMS HANDLERS
// ====================================================================

// CreateOrderItem - создать позицию заказа
func CreateOrderItem(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req OrderItemRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid request body")
			return
		}

		id, err := storage.AddOrderItem(req.OrderID, req.ProductID, req.Quantity, req.Price, req.Discount)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		item, err := storage.GetOrderItem(id)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, item)
	}
}

// GetOrderItem - получить позицию заказа по ID
func GetOrderItem(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid order item id")
			return
		}

		item, err := storage.GetOrderItem(id)
		if err != nil {
			respondError(w, r, http.StatusNotFound, "order item not found")
			return
		}

		render.JSON(w, r, item)
	}
}

// ListOrderItems - получить позиции заказа
func ListOrderItems(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid order id")
			return
		}

		items, err := storage.ListOrderItems(orderID)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, items)
	}
}

// UpdateOrderItem - обновить позицию заказа
func UpdateOrderItem(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid order item id")
			return
		}

		var req struct {
			Quantity int     `json:"quantity"`
			Price    float64 `json:"price"`
			Discount float64 `json:"discount"`
		}
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := storage.UpdateOrderItem(id, req.Quantity, req.Price, req.Discount); err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		item, err := storage.GetOrderItem(id)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, item)
	}
}

// DeleteOrderItem - удалить позицию заказа
func DeleteOrderItem(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseURLParamID(r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid order item id")
			return
		}

		if err := storage.DeleteOrderItem(id); err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.Status(r, http.StatusNoContent)
	}
}
