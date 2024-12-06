package docs

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type ProductInstance struct {
	Id           string `json:"id"`
	EnterpriseId string `json:"enterpriseId"`
	ProductId    int    `json:"productId"`
}

func main() {
	// Open the JSON file
	jsonFile, err := os.Open("product_instances_raw.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer jsonFile.Close()

	// Decode the JSON data
	var products []ProductInstance
	err = json.NewDecoder(jsonFile).Decode(&products)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Create a CSV file
	csvFile, err := os.Create("product_instances_csv_string.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer csvFile.Close()

	// Create a CSV writer
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{"id", "productId", "enterpriseId"})
	if err != nil {
		fmt.Println("Error writing CSV header:", err)
		return
	}

	// Write data to CSV
	for _, product := range products {
		productId := fmt.Sprintf("%d", product.ProductId)
		err := writer.Write([]string{product.Id, productId, product.EnterpriseId})
		//err := writer.Write([]string{product.Id, product.Name, product.RoleKey})
		if err != nil {
			fmt.Println("Error writing CSV data:", err)
			return
		}
	}

	fmt.Println("CSV file created successfully")
}
