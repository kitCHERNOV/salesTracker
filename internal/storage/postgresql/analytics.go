package postgresql

import (
	"fmt"
	_ "github.com/lib/pq"
	"math"
	"time"
)

// variable
const packageOp = "storage.postgresql.analytics."

// ====================================================================
// ANALYTICS - Аналитические функции
// ====================================================================

// PeriodSummary — суммарные показатели за период
type PeriodSummary struct {
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	TotalRevenue float64   `json:"total_revenue"`
	OrderCount   int       `json:"order_count"`
}

// TotalRevenueByPeriod — получить сумму заказов за определенный период
func (s *Storage) TotalRevenueByPeriod(start, end time.Time) (*PeriodSummary, error) {
	const op = packageOp + "TotalRevenueByPeriod"
	query := `SELECT SUM(total_amount), COUNT(*)
			FROM orders
			WHERE order_date BETWEEN $1 AND $2`

	var (
		totalRevenue float64
		ordersAmount int
	)

	err := s.DB.QueryRow(query, start, end).Scan(&totalRevenue, &ordersAmount)
	if err != nil {
		return nil, fmt.Errorf("getting total revenue by period error; %s, %v", op, err)
	}

	if ordersAmount < 1 || totalRevenue < 0 {
		return nil, fmt.Errorf("invalid meaning; %s", op)
	}
	return &PeriodSummary{
		StartDate:    start,
		EndDate:      end,
		TotalRevenue: totalRevenue,
		OrderCount:   ordersAmount,
	}, nil
}

// DailyOrders — количество заказов по дням
type DailyOrders struct {
	Date        string  `json:"date"`
	OrderCount  int     `json:"order_count"`
	TotalAmount float64 `json:"total_amount"`
}

// OrdersPerDay — количество заказов в день за период
func (s *Storage) OrdersPerDay(start, end time.Time) ([]DailyOrders, error) {
	const op = packageOp + "OrdersPerDay"
	var dailyOrders []DailyOrders

	query := `SELECT COUNT(*), SUM(total_amount)
			FROM orders
			WHERE order_date $1`

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s,%v", op, err)
	}

	for curDate := start; !curDate.After(end); curDate = curDate.AddDate(0, 0, 1) {
		var (
			orderCount  int
			totalAmount float64
		)
		row := stmt.QueryRow(curDate)
		err = row.Scan(&orderCount, &totalAmount)
		if err != nil {
			return nil, fmt.Errorf("%s,%v", op, err)
		}
		dailyOrders = append(dailyOrders, DailyOrders{
			Date:        curDate.Format("2006-01-02"),
			OrderCount:  orderCount,
			TotalAmount: totalAmount},
		)
	}

	return dailyOrders, nil
	// Mock implementation
	//return []DailyOrders{
	//	{Date: "2024-01-15", OrderCount: 12, TotalAmount: 125000.50},
	//	{Date: "2024-01-16", OrderCount: 15, TotalAmount: 187500.00},
	//	{Date: "2024-01-17", OrderCount: 8, TotalAmount: 95000.00},
	//	{Date: "2024-01-18", OrderCount: 20, TotalAmount: 250000.00},
	//	{Date: "2024-01-19", OrderCount: 18, TotalAmount: 210000.00},
	//}, nil
}

// AverageCheck — средний чек за период
type AverageCheckStats struct {
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	AverageCheck float64   `json:"average_check"`
	MinCheck     float64   `json:"min_check"`
	MaxCheck     float64   `json:"max_check"`
}

// AverageCheckByPeriod — средний чек за определенный период
func (s *Storage) AverageCheckByPeriod(start, end time.Time) (*AverageCheckStats, error) {
	const op = packageOp + "AverageCheckByPeriod"
	query := `SELECT COUNT(*), SUM(total_amount), MIN(total_amount), MAX(total_amount)
			FROM orders
			WHERE order_date $1`

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s,%v", op, err)
	}

	var MinCheckTotally = math.MaxFloat64
	var MaxCheckTotally = -1.
	var SumOfBills = 0.
	var ordersAmount = 0

	for curDate := start; !curDate.After(end); curDate = curDate.AddDate(0, 0, 1) {
		var (
			orderCount  int
			totalAmount float64
			minCheck    float64
			maxCheck    float64
		)
		row := stmt.QueryRow(curDate)
		err = row.Scan(&orderCount, &totalAmount, &minCheck, &maxCheck)
		if err != nil {
			return nil, fmt.Errorf("%s,%v", op, err)
		}

		if minCheck < MinCheckTotally {
			MinCheckTotally = minCheck
		}
		if maxCheck > MaxCheckTotally {
			MaxCheckTotally = maxCheck
		}

		SumOfBills = SumOfBills + totalAmount
		ordersAmount += orderCount
	}

	return &AverageCheckStats{
		StartDate:    start,
		EndDate:      end,
		AverageCheck: SumOfBills / float64(ordersAmount),
		MinCheck:     MinCheckTotally,
		MaxCheck:     MaxCheckTotally,
	}, nil

	// Mock implementation
	//return &AverageCheckStats{
	//	StartDate:    start,
	//	EndDate:      end,
	//	AverageCheck: 8500.00,
	//	MinCheck:     450.00,
	//	MaxCheck:     125000.00,
	//}, nil
}

// MedianStats — результаты расчета медианы
type MedianStats struct {
	Metric     string  `json:"metric"`
	Median     float64 `json:"median"`
	SampleSize int     `json:"sample_size"`
}

// OrdersMedian — медиана суммы заказов за период по дням
func (s *Storage) OrdersMedian(start, end time.Time) (*MedianStats, error) {
	const op = packageOp + "OrdersMedian"
	query := `SELECT total_amount
		FROM orders
		WHERE order_date = $1`

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s, %v", op, err)
	}

	total_amounts := make([]float64, 0)
	days_shift := int(end.Sub(start).Hours() / 24)

	for currentDate := start; !currentDate.After(end); currentDate = currentDate.AddDate(0, 0, 1) {
		rows, err := stmt.Query(currentDate)
		if err != nil {
			return nil, fmt.Errorf("%s, %v", op, err)
		}

		defer rows.Close()

		for rows.Next() {
			var orderTotalAmount float64
			err = rows.Scan(&orderTotalAmount)
			if err != nil {
				return nil, fmt.Errorf("%s, %v", op, err)
			}
			total_amounts = append(total_amounts, orderTotalAmount)
		}

	}

	return &MedianStats{
		Metric:     "order_total",
		Median:     7200.00,
		SampleSize: days_shift,
	}, nil

	// Mock implementation
	//return &MedianStats{
	//	Metric:     "order_total",
	//	Median:     7200.00,
	//	SampleSize: 150,
	//}, nil
}
/*
TODO: implement CustomerSpendingMedian:
- add sql request to get data
- add math logic to compute customer spending median
*/
// CustomerSpendingMedian — медиана трат покупателей за период
func (s *Storage) CustomerSpendingMedian(start, end time.Time) (*MedianStats, error) {
	// Mock implementation
	return &MedianStats{
		Metric:     "customer_spending",
		Median:     12500.00,
		SampleSize: 85,
	}, nil
}

// PercentileStats — результаты расчета перцентиля
type PercentileStats struct {
	Metric     string  `json:"metric"`
	Percentile int     `json:"percentile"`
	Value      float64 `json:"value"`
	SampleSize int     `json:"sample_size"`
}

// TODO: add orders percentili implemention
// OrdersPercentile — перцентиль суммы заказов за период
func (s *Storage) OrdersPercentile(start, end time.Time, percentile int) (*PercentileStats, error) {
	const op = packageOp + "OrdersPercentile"
	
	// Validate percentile range
	if percentile < 0 || percentile > 100 {
		return nil, fmt.Errorf("%s, percentile must be between 0 and 100", op)
	}

	query := `SELECT total_amount
		FROM orders
		WHERE order_date BETWEEN $1 AND $2
		ORDER BY total_amount ASC`

	rows, err := s.DB.Query(query, start, end)
	if err != nil {
		return nil, fmt.Errorf("%s, %v", op, err)
	}
	defer rows.Close()

	// Collect all total_amount values
	totalAmounts := make([]float64, 0)
	for rows.Next() {
		var totalAmount float64
		err := rows.Scan(&totalAmount)
		if err != nil {
			return nil, fmt.Errorf("%s, %v", op, err)
		}
		totalAmounts = append(totalAmounts, totalAmount)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s, %v", op, err)
	}

	sampleSize := len(totalAmounts)
	if sampleSize == 0 {
		return nil, fmt.Errorf("%s, no orders found in the specified period", op)
	}

	// Calculate percentile using linear interpolation method
	// Formula: percentile_value = values[floor((n-1)*p/100)] + ((n-1)*p/100 - floor((n-1)*p/100)) * (values[floor((n-1)*p/100)+1] - values[floor((n-1)*p/100)])
	
	// Values are already sorted by ORDER BY in the query
	n := float64(sampleSize - 1)
	pos := n * float64(percentile) / 100.0
	lowerIndex := int(math.Floor(pos))
	upperIndex := int(math.Ceil(pos))
	
	// Handle edge case where upperIndex equals sampleSize
	if upperIndex >= sampleSize {
		upperIndex = sampleSize - 1
	}

	// Interpolate between values
	fraction := pos - float64(lowerIndex)
	percentileValue := totalAmounts[lowerIndex] + fraction*(totalAmounts[upperIndex]-totalAmounts[lowerIndex])

	return &PercentileStats{
		Metric:     "order_total",
		Percentile: percentile,
		Value:      percentileValue,
		SampleSize: sampleSize,
	}, nil
}

// CustomerSpendingPercentile — перцентиль трат покупателей за период
func (s *Storage) CustomerSpendingPercentile(start, end time.Time, percentile int) (*PercentileStats, error) {
	// Mock implementation
	mockValues := map[int]float64{
		50: 12500.00,
		75: 22000.00,
		90: 35000.00,
		95: 50000.00,
		99: 95000.00,
	}

	value := mockValues[percentile]
	if value == 0 {
		value = 15000.00 // default
	}

	return &PercentileStats{
		Metric:     "customer_spending",
		Percentile: percentile,
		Value:      value,
		SampleSize: 85,
	}, nil
}

// ====================================================================
// COMBINED ANALYTICS — Комбинированные аналитические отчеты
// ====================================================================

// SalesReport — полный отчет по продажам за период
type SalesReport struct {
	Period       PeriodSummary     `json:"period"`
	DailyStats   []DailyOrders     `json:"daily_stats"`
	AverageCheck AverageCheckStats `json:"average_check"`
	Median       MedianStats       `json:"median"`
	Percentile75 PercentileStats   `json:"percentile_75"`
	Percentile95 PercentileStats   `json:"percentile_95"`
}

// GenerateSalesReport — сгенерировать полный отчет по продажам
func (s *Storage) GenerateSalesReport(start, end time.Time) (*SalesReport, error) {
	// Mock implementation
	period, _ := s.TotalRevenueByPeriod(start, end)
	dailyStats, _ := s.OrdersPerDay(start, end)
	avgCheck, _ := s.AverageCheckByPeriod(start, end)
	median, _ := s.OrdersMedian(start, end)
	p75, _ := s.OrdersPercentile(start, end, 75)
	p95, _ := s.OrdersPercentile(start, end, 95)

	return &SalesReport{
		Period:       *period,
		DailyStats:   dailyStats,
		AverageCheck: *avgCheck,
		Median:       *median,
		Percentile75: *p75,
		Percentile95: *p95,
	}, nil
}
