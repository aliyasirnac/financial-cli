package main

import (
	"context"
	"fmt"
	"github.com/aliyasirnac/financialManagement/internal/db"
	"github.com/aliyasirnac/financialManagement/internal/model"
	"log"
)

func main() {
	dbConfig := db.NewDatabase()
	db := dbConfig.OpenDatabase()
	defer db.Close()

	// Create the products table.
	ctx := context.Background()
	_, err := db.NewCreateTable().Model((*model.Product)(nil)).Exec(ctx)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	fmt.Println("Table created successfully")
}
