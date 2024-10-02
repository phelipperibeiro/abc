package http

import (
	"application"
	"application/workreport"
	"archive/zip"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

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

// type findWRAdvSearchResponse struct {
// 	Results  []*application.WRAdvSearchResult `json:"work_report_adv_search_results"`
// 	Metadata application.Metadata             `json:"metadata"`
// }

func (s *Server) RegisterWorkReportRoutes(router chi.Router) {
	router.Get("/work-reports", s.handleWorkReportList)
	router.Post("/work-reports/{file_name}", s.handleCreateWorkReport)
	router.Get("/work-report-topics", s.handleWorkReportTopicList)

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
