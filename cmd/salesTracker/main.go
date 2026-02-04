package main

func main() {
	storage := NewSalesService()
	Run(":8080", storage)
}
