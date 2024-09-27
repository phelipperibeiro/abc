package main

import (
	"application"
	"application/postgres"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func dump(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Erro ao serializar dados:", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonData))
}

func initializeLogger() *log.Logger {
	return log.New(log.Writer(), "app: ", log.LstdFlags)
}

func initializeDatabase(logger *log.Logger) *pgxpool.Pool {
	connString := "host=localhost port=5432 user=abacateiro password=abacateiro dbname=abacateiro sslmode=disable"
	dbPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		logger.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	return dbPool
}

func main() {
	logger := initializeLogger()
	dbPool := initializeDatabase(logger)
	unitService := postgres.NewUnitService(dbPool)
	ligado := false

	if ligado {
		testCreateUnit(logger, unitService)
		testFindUnitByID(logger, unitService)
		testFindUnits(logger, unitService)
		testFindAllUnits(logger, unitService)
		testUpdateUnit(logger, unitService)
		testDeleteUnit(logger, unitService)
	}
}

func testCreateUnit(logger *log.Logger, unitService *postgres.UnitService) {
	unit := &application.Unit{
		Name:                 "P72",
		StoragePath:          "/srv/sbsautfs01/05_Replica_Repositorio_P-72",
		PrometheusServerAddr: "https://sp70autlx02:9090",
		ID:                   10,
		TotalDevices:         0,
	}

	err := unitService.CreateUnit(context.Background(), unit)

	if err != nil {
		logger.Fatalf("Erro ao criar unidade: %v", err)
	}
	logger.Printf("Unidade criada com sucesso: %v", unit)
}

func testFindUnitByID(logger *log.Logger, unitService *postgres.UnitService) {

	unit, err := unitService.FindUnitByID(context.Background(), 10)

	if err != nil {
		logger.Fatalf("Erro ao buscar unidade: %v", err)
	}

	dump(unit)

	logger.Printf("Unidade encontrada: %v", unit)
}

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}

func testFindUnits(logger *log.Logger, unitService *postgres.UnitService) {
	units, _, err := unitService.FindUnits(context.Background(), application.UnitFilter{
		ID:           intPtr(10),
		IDs:          []int{9, 8},
		Name:         strPtr("P72"),
		GlobalSearch: strPtr("P72"),
		Pagination: application.Pagination{
			Sort:              "unit_id",
			SortBy:            "unit_id",
			SortDescending:    false,
			DisablePagination: false,
			SortSafeList: []string{
				"unit_id",
				"unit_name",
				"storage_path",
				"prometheus_server_address",
			},
			Page:     1,
			PageSize: 10,
		},
	})

	if err != nil {
		logger.Fatalf("Erro ao buscar unidades: %v", err)
	}

	dump(units)

	logger.Printf("Unidades encontradas: %v", units)
}

func testFindAllUnits(logger *log.Logger, unitService *postgres.UnitService) {
	units, err := unitService.FindAllUnits(context.Background())

	if err != nil {
		logger.Fatalf("Erro ao buscar todas as unidades: %v", err)
	}

	dump(units)

	logger.Printf("Unidades encontradas: %v", units)
}

func testUpdateUnit(logger *log.Logger, unitService *postgres.UnitService) {
	unit, _ := unitService.FindUnitByID(context.Background(), 10)
	unit.Name = "P72 updated"
	unit.PrometheusServerAddr = "https://sp70autlx02:9090 updated"

	_, err := unitService.UpdateUnit(context.Background(), unit.ID, application.UnitUpdate{
		Name:                 &unit.Name,
		PrometheusServerAddr: &unit.PrometheusServerAddr,
	})

	if err != nil {
		logger.Fatalf("Erro ao atualizar unidade: %v", err)
	}

	logger.Println("Unidade atualizada com sucesso")
}

func testDeleteUnit(logger *log.Logger, unitService *postgres.UnitService) {
	err := unitService.DeleteUnit(context.Background(), 10)

	if err != nil {
		logger.Fatalf("Erro ao deletar unidade: %v", err)
	}

	logger.Println("Unidade deletada com sucesso")
}
