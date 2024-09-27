package main

import (
	"application"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"application/postgres"

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

// func dd(data interface{}) {
// 	dump(data)
// 	os.Exit(0)
// }

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
	workReportService := postgres.NewWorkReportService(dbPool)

	ligado := false

	if ligado {
		testCreateWorkReport(logger, workReportService)
		testFindWorkReportByID(logger, workReportService)
		testFindWorkReports(logger, workReportService)
		testUpdateWorkReport(logger, workReportService)
		testDeleteWorkReport(logger, workReportService)
		testCreateWorkReportTopic(logger, workReportService)
		testFindWorkReportTopicByID(logger, workReportService)
	}

	testFindWorkReportTopics(logger, workReportService)
}

func testFindWorkReportByID(logger *log.Logger, workReportService *postgres.WorkReportService) application.WorkReport {
	workReport, err := workReportService.FindWorkReportByID(context.Background(), 2)
	if err != nil {
		logger.Fatalf("Erro ao buscar relatório de trabalho: %v", err)
	}

	dump(workReport)
	logger.Printf("Relatório de trabalho encontrado: %v", workReport)
	return *workReport
}

func testCreateWorkReport(logger *log.Logger, workReportService *postgres.WorkReportService) {
	err := workReportService.CreateWorkReport(context.Background(), application.WorkReport{
		DocName: "docname",
		UnitID:  5,
		From:    time.Now(),
		To:      time.Now(),
		Text:    "text",
		Data:    []byte("data"),
	})
	if err != nil {
		logger.Fatalf("Erro ao criar relatório de trabalho: %v", err)
	}
	logger.Println("Relatório de trabalho criado com sucesso")
}

func testUpdateWorkReport(logger *log.Logger, workReportService *postgres.WorkReportService) {

	workReport := testFindWorkReportByID(logger, workReportService)
	workReport.DocName = "docname updated"
	workReport.UnitID = 6
	workReport.From = workReport.From.AddDate(0, 0, 1)
	workReport.To = workReport.To.AddDate(0, 0, 1)

	_, err := workReportService.UpdateWorkReport(context.Background(), workReport.ID, application.WorkReportUpdate{
		DocName: &workReport.DocName,
		UnitID:  &workReport.UnitID,
		From:    &workReport.From,
		To:      &workReport.To,
	})

	if err != nil {
		logger.Fatalf("Erro ao atualizar relatório de trabalho: %v", err)
	}

	logger.Println("Relatório de trabalho atualizado com sucesso")
}

func testFindWorkReports(logger *log.Logger, workReportService *postgres.WorkReportService) {
	workReports, _, err := workReportService.FindWorkReports(context.Background(), application.WorkReportFilter{})
	if err != nil {
		logger.Fatalf("Erro ao buscar relatórios de trabalho: %v", err)
	}
	dump(workReports)
	logger.Printf("Relatórios de trabalho encontrados: %v", workReports)
}

func testDeleteWorkReport(logger *log.Logger, workReportService *postgres.WorkReportService) {
	err := workReportService.DeleteWorkReport(context.Background(), 1)
	if err != nil {
		logger.Fatalf("Erro ao deletar relatório de trabalho: %v", err)
	}

	logger.Println("Relatório de trabalho deletado com sucesso")
}

func testCreateWorkReportTopic(logger *log.Logger, workReportService *postgres.WorkReportService) {

	workReport, err := workReportService.FindWorkReportByID(context.Background(), 2)

	if err != nil {
		logger.Fatalf("Erro ao buscar relatório de trabalho: %v", err)
	}

	err = workReportService.CreateWorkReportTopic(context.Background(), application.WorkReportTopic{
		WorkReportID: workReport.ID,
		Title:        "title",
		Text:         "text",
	})
	if err != nil {
		logger.Fatalf("Erro ao criar tópico de relatório de trabalho: %v", err)
	}

	logger.Println("Tópico de relatório de trabalho criado com sucesso")
}

func testFindWorkReportTopicByID(logger *log.Logger, workReportService *postgres.WorkReportService) {

	workReportTopic, err := workReportService.FindWorkReportTopicByID(context.Background(), 1)
	if err != nil {
		logger.Fatalf("Erro ao buscar tópico de relatório de trabalho: %v", err)
	}

	dump(workReportTopic)
	logger.Printf("Tópico de relatório de trabalho encontrado: %v", workReportTopic)
}

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func testFindWorkReportTopics(logger *log.Logger, workReportService *postgres.WorkReportService) {

	workReportTopics, _, err := workReportService.FindWorkReportTopics(context.Background(), application.WorkReportTopicFilter{
		ID:           intPtr(1),
		Title:        strPtr("title"),
		Text:         strPtr("text"),
		UnitID:       intPtr(1),
		From:         timePtr(time.Now().AddDate(0, 0, -1)),
		To:           timePtr(time.Now().AddDate(0, 0, 1)),
		GlobalSearch: strPtr("title"),
		Pagination: application.Pagination{
			Page:     1,
			PageSize: 10,
		},
	})
	if err != nil {
		logger.Fatalf("Erro ao buscar tópicos de relatório de trabalho: %v", err)
	}
	dump(workReportTopics)
	logger.Printf("Tópicos de relatório de trabalho encontrados: %v", workReportTopics)
}

// func testFindWorkReportTopicsAdvSearch(logger *log.Logger, workReportService *postgres.WorkReportService) {
// 	// implementar logica de busca avançada de tópicos de relatório
// }
