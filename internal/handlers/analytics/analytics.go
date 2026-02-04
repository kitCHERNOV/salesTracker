package analytics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"

	"salesTracker/internal/storage/postgresql"
)

// ====================================================================
// HELPERS
// ====================================================================

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func respondError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, map[string]string{"error": message})
}

// ====================================================================
// TOTAL REVENUE
// ====================================================================

// TotalRevenueByPeriod - выручка за период
// GET /analytics/revenue?start=2024-01-01&end=2024-01-31
func TotalRevenueByPeriod(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start")
		endDate := r.URL.Query().Get("end")

		start, err := parseDate(startDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid start date format, use YYYY-MM-DD")
			return
		}

		end, err := parseDate(endDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid end date format, use YYYY-MM-DD")
			return
		}

		summary, err := storage.TotalRevenueByPeriod(start, end)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, summary)
	}
}

// ====================================================================
// DAILY ORDERS
// ====================================================================

// OrdersPerDay - заказы по дням
// GET /analytics/daily-orders?start=2024-01-01&end=2024-01-31
func OrdersPerDay(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start")
		endDate := r.URL.Query().Get("end")

		start, err := parseDate(startDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid start date format, use YYYY-MM-DD")
			return
		}

		end, err := parseDate(endDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid end date format, use YYYY-MM-DD")
			return
		}

		dailyOrders, err := storage.OrdersPerDay(start, end)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, dailyOrders)
	}
}

// ====================================================================
// AVERAGE CHECK
// ====================================================================

// AverageCheckByPeriod - средний чек за период
// GET /analytics/average-check?start=2024-01-01&end=2024-01-31
func AverageCheckByPeriod(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start")
		endDate := r.URL.Query().Get("end")

		start, err := parseDate(startDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid start date format, use YYYY-MM-DD")
			return
		}

		end, err := parseDate(endDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid end date format, use YYYY-MM-DD")
			return
		}

		avgCheck, err := storage.AverageCheckByPeriod(start, end)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, avgCheck)
	}
}

// ====================================================================
// MEDIAN
// ====================================================================

// OrdersMedian - медиана заказов
// GET /analytics/orders-median?start=2024-01-01&end=2024-01-31
func OrdersMedian(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start")
		endDate := r.URL.Query().Get("end")

		start, err := parseDate(startDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid start date format, use YYYY-MM-DD")
			return
		}

		end, err := parseDate(endDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid end date format, use YYYY-MM-DD")
			return
		}

		median, err := storage.OrdersMedian(start, end)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, median)
	}
}

// CustomerSpendingMedian - медиана трат покупателей
// GET /analytics/customer-median?start=2024-01-01&end=2024-01-31
func CustomerSpendingMedian(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start")
		endDate := r.URL.Query().Get("end")

		start, err := parseDate(startDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid start date format, use YYYY-MM-DD")
			return
		}

		end, err := parseDate(endDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid end date format, use YYYY-MM-DD")
			return
		}

		median, err := storage.CustomerSpendingMedian(start, end)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, median)
	}
}

// ====================================================================
// PERCENTILES
// ====================================================================

// OrdersPercentile - перцентиль заказов
// GET /analytics/orders-percentile?start=2024-01-01&end=2024-01-31&percentile=75
func OrdersPercentile(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start")
		endDate := r.URL.Query().Get("end")
		percentileStr := r.URL.Query().Get("percentile")

		start, err := parseDate(startDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid start date format, use YYYY-MM-DD")
			return
		}

		end, err := parseDate(endDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid end date format, use YYYY-MM-DD")
			return
		}

		percentile, err := strconv.Atoi(percentileStr)
		if err != nil || percentile < 0 || percentile > 100 {
			respondError(w, r, http.StatusBadRequest, "invalid percentile, must be between 0 and 100")
			return
		}

		result, err := storage.OrdersPercentile(start, end, percentile)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, result)
	}
}

// CustomerSpendingPercentile - перцентиль трат покупателей
// GET /analytics/customer-percentile?start=2024-01-01&end=2024-01-31&percentile=75
func CustomerSpendingPercentile(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start")
		endDate := r.URL.Query().Get("end")
		percentileStr := r.URL.Query().Get("percentile")

		start, err := parseDate(startDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid start date format, use YYYY-MM-DD")
			return
		}

		end, err := parseDate(endDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid end date format, use YYYY-MM-DD")
			return
		}

		percentile, err := strconv.Atoi(percentileStr)
		if err != nil || percentile < 0 || percentile > 100 {
			respondError(w, r, http.StatusBadRequest, "invalid percentile, must be between 0 and 100")
			return
		}

		result, err := storage.CustomerSpendingPercentile(start, end, percentile)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, result)
	}
}

// ====================================================================
// COMBINED REPORTS
// ====================================================================

// GenerateSalesReport - полный отчет по продажам
// GET /analytics/sales-report?start=2024-01-01&end=2024-01-31
func GenerateSalesReport(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("start")
		endDate := r.URL.Query().Get("end")

		start, err := parseDate(startDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid start date format, use YYYY-MM-DD")
			return
		}

		end, err := parseDate(endDate)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, "invalid end date format, use YYYY-MM-DD")
			return
		}

		report, err := storage.GenerateSalesReport(start, end)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, r, report)
	}
}

// ====================================================================
// ROUTE SETUP HELPER
// ====================================================================

// Routes возвращает пути для маршрутизации аналитики
type Route struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

// GetAnalyticsRoutes - получить все роуты аналитики
func GetAnalyticsRoutes(storage *postgresql.Storage) []Route {
	return []Route{
		// Revenue
		{"GET", "/analytics/revenue", TotalRevenueByPeriod(storage)},

		// Daily stats
		{"GET", "/analytics/daily-orders", OrdersPerDay(storage)},

		// Average check
		{"GET", "/analytics/average-check", AverageCheckByPeriod(storage)},

		// Median
		{"GET", "/analytics/orders-median", OrdersMedian(storage)},
		{"GET", "/analytics/customer-median", CustomerSpendingMedian(storage)},

		// Percentiles
		{"GET", "/analytics/orders-percentile", OrdersPercentile(storage)},
		{"GET", "/analytics/customer-percentile", CustomerSpendingPercentile(storage)},

		// Combined reports
		{"GET", "/analytics/sales-report", GenerateSalesReport(storage)},
	}
}
