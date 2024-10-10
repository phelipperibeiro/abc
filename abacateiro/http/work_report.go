package http

import (
	"application"
	"application/workreport"
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// func dump(data interface{}) {
// 	jsonData, err := json.MarshalIndent(data, "", "  ")
// 	if err != nil {
// 		fmt.Println("Erro ao serializar dados:", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println(string(jsonData))
// }

type findWorkReportResponse struct {
	WorkReports []*application.WorkReport `json:"work_reports"`
	Metadata    application.Metadata      `json:"metadata"`
}

type findWorkReportTopicResponse struct {
	WorkReportTopics []*application.WorkReportTopic `json:"work_report_topics"`
	Metadata         application.Metadata           `json:"metadata"`
}

type findWRAdvSearchResponse struct {
	Results  []*application.WRAdvSearchResult `json:"work_report_adv_search_results"`
	Metadata application.Metadata             `json:"metadata"`
}

func (s *Server) RegisterWorkReportRoutes(router chi.Router) {
	router.Get("/work-reports", s.handleWorkReportList)
	router.Post("/work-reports/{file_name}", s.handleCreateWorkReport)
	router.Get("/work-report-topics", s.handleWorkReportTopicList)
	router.Get("/work-report-topics/adv-search", s.handleWorkReportAdvSearch)

}

func (s *Server) handleCreateWorkReport(w http.ResponseWriter, r *http.Request) {

	fileName := chi.URLParam(r, "file_name")

	if filepath.Ext(fileName) != ".docx" {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "Formato inválido: apenas docx é aceito"))
		return
	}

	workReport, err := application.GetWorkReportFromFileName(fileName)

	if err != nil {
		s.Error(w, r, err)
		return
	}

	units, _, _ := s.unitService.FindUnits(r.Context(), application.UnitFilter{Name: &workReport.UnitName})
	if len(units) != 1 {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "Unidade não reconhecida: %s", workReport.UnitName))
		return
	}

	workReport.UnitID = units[0].ID
	err = r.ParseMultipartForm(200 << 20)
	if err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "maximum size allowed: 200MB"))
		return
	}

	mpFile, mpHeader, err := r.FormFile("file")
	if err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "arquivo inválido: %v", err))
		return
	}
	defer mpFile.Close()

	zr, err := zip.NewReader(mpFile, mpHeader.Size)
	if err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "arquivo inválido: %v", err))
		return
	}

	text, topics, err := workreport.ExtractText(zr)
	if err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "arquivo inválido: %v", err))
		return
	}
	workReport.Text = text

	if _, err := mpFile.Seek(0, io.SeekStart); err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "arquivo inválido: %v", err))
		return
	}

	workReport.Data, err = io.ReadAll(mpFile)
	if err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "arquivo inválido: %v", err))
		return
	}

	// verifica se o relatório já existe
	wrs, _, err := s.workReportService.FindWorkReports(r.Context(), application.WorkReportFilter{DocName: &workReport.DocName})
	if err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "Erro ao verificar relatório de trabalho: %v", err))
		return
	}

	if len(wrs) > 0 {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "Relatório de trabalho já existe: %s", workReport.DocName))
		return
	}

	err = s.workReportService.CreateWorkReport(r.Context(), workReport)
	if err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "Erro ao criar relatório de trabalho: %v", err))
		return
	}

	for _, topic := range topics {
		t := &application.WorkReportTopic{
			Title:        topic.Title,
			Text:         topic.Text,
			WorkReportID: workReport.ID,
		}
		if err := s.workReportService.CreateWorkReportTopic(r.Context(), t); err != nil {
			s.Error(w, r, err)
			return
		}
	}

	go s.workReportService.DeleteDuplicatesWorkReportTopics(r.Context())

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(workReport); err != nil {
		s.Error(w, r, err)
		return
	}
}

func (s *Server) handleWorkReportList(w http.ResponseWriter, r *http.Request) {

	var filter application.WorkReportFilter

	switch r.Header.Get("Content-type") {
	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			s.Error(w, r, application.Errorf(application.ERRINVALID, "Invalid JSON body"))
			return
		}
	default:
		filter.Pagination.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
		filter.Pagination.PageSize = 20
	}

	filter.LimitPagination()

	workReports, meta, err := s.workReportService.FindWorkReports(r.Context(), filter)
	if err != nil {
		s.Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(findWorkReportResponse{
		WorkReports: workReports,
		Metadata:    meta,
	}); err != nil {
		s.Error(w, r, err)
		return
	}
}

func (s *Server) handleWorkReportTopicList(w http.ResponseWriter, r *http.Request) {

	var filter application.WorkReportTopicFilter

	switch r.Header.Get("Content-type") {
	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			s.Error(w, r, application.Errorf(application.ERRINVALID, "Invalid JSON body"))
			return
		}
	default:
		filter.Pagination.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
		filter.Pagination.PageSize, _ = strconv.Atoi(r.URL.Query().Get("page_size"))
	}

	filter.LimitPagination()

	search := r.URL.Query().Get("search")

	filter.GlobalSearch = &search

	topics, meta, err := s.workReportService.FindWorkReportTopics(r.Context(), filter)
	if err != nil {
		s.Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(findWorkReportTopicResponse{
		WorkReportTopics: topics,
		Metadata:         meta,
	}); err != nil {
		s.Error(w, r, err)
		return
	}
}

func dump(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Erro ao serializar dados:", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonData))
}

// type WRAdvSearchFilter struct {
// 	UnitID *int       `json:"unit_id"`
// 	From   *time.Time `json:"work_report_from"`
// 	To     *time.Time `json:"work_report_to"`
// 	Years  []int      `json:"years"`

// 	GlobalSearch *string `json:"global_search"`

// 	Pagination
// }

func (s *Server) handleWorkReportAdvSearch(w http.ResponseWriter, r *http.Request) {

	var filter application.WRAdvSearchFilter

	// Define content type and decode JSON body if applicable using switch
	switch r.Header.Get("Content-type") {
	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			s.Error(w, r, application.Errorf(application.ERRINVALID, "Invalid JSON body"))
			return
		}
	default:
		filter.Pagination.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
		filter.Pagination.PageSize, _ = strconv.Atoi(r.URL.Query().Get("page_size"))
	}

	filter.LimitPagination()

	getYears := func(yearsParams interface{}) ([]int, error) {
		var years []int

		// Verifica se é um array de strings e maior que 0
		params, ok := yearsParams.([]string)
		if !ok || len(params) == 0 {
			return nil, fmt.Errorf("year must be a non-empty array of strings")
		}

		for _, yearStr := range params {
			year, err := strconv.Atoi(yearStr)
			if err != nil {
				return nil, fmt.Errorf("invalid year format: %s", yearStr)
			}
			years = append(years, year)
		}
		return years, nil
	}

	// Extract query parameters
	years, _ := getYears(r.URL.Query()["year[]"])
	unitID, _ := strconv.Atoi(r.URL.Query().Get("unit_id"))
	search := r.URL.Query().Get("search")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	// Parse dates
	fromDate, _ := time.Parse("2006-01-02", from)
	toDate, _ := time.Parse("2006-01-02", to)

	// Set filter fields
	filter.From = &fromDate
	filter.To = &toDate
	filter.UnitID = &unitID
	filter.GlobalSearch = &search
	filter.Years = years

	// Call the service to get results
	results, meta, err := s.workReportService.FindWorkReportTopicsAdvSearch(r.Context(), filter)
	if err != nil {
		s.Error(w, r, err)
		return
	}

	// Return the response as JSON
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(findWRAdvSearchResponse{
		Results:  results,
		Metadata: meta,
	}); err != nil {
		s.Error(w, r, err)
		return
	}
}
