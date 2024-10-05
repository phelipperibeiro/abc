package postgres

import (
	"application"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkReportService struct {
	db *pgxpool.Pool
}

// implements WorkReportServiceInterface
func NewWorkReportService(db *pgxpool.Pool) *WorkReportService {
	return &WorkReportService{
		db: db,
	}
}

func (s *WorkReportService) CreateWorkReport(ctx context.Context, wr *application.WorkReport) error {

	if err := wr.Validate(); err != nil {
		return fmt.Errorf("invalid work report: %w", err)
	}

	query := `INSERT INTO work_reports (work_report_id, unit_id, work_report_docname, work_report_from, work_report_to, work_report_text, work_report_data) VALUES (DEFAULT, $1, $2, $3, $4, $5, $6) RETURNING work_report_id`

	err := s.db.QueryRow(context.Background(), query, wr.UnitID, wr.DocName, wr.From, wr.To, wr.Text, wr.Data).Scan(&wr.ID)

	if err != nil {
		return fmt.Errorf("failed to create work report: %w", err)
	}

	return nil
}

func (s *WorkReportService) FindWorkReportByID(ctx context.Context, id int) (*application.WorkReport, error) {
	query := `SELECT work_report_id, unit_id, work_report_docname, work_report_from, work_report_to, work_report_text, work_report_data FROM work_reports WHERE work_report_id = $1`

	var wr application.WorkReport
	err := s.db.QueryRow(context.Background(), query, id).Scan(&wr.ID, &wr.UnitID, &wr.DocName, &wr.From, &wr.To, &wr.Text, &wr.Data)

	if err != nil {
		return nil, fmt.Errorf("failed to find work report: %w", err)
	}

	return &wr, nil
}

func (s *WorkReportService) FindWorkReports(ctx context.Context, filter application.WorkReportFilter) ([]*application.WorkReport, application.Metadata, error) {

	query := `SELECT
				work_report_id,
				unit_id as work_report_unit_id,
				work_report_docname,
				work_report_from,
				work_report_to,
				work_report_text, 
				work_report_data,

				unit_id,
				unit_name,
				storage_path as unit_storage_path,
				prometheus_server_address as unit_prometheus_server_address
			FROM work_reports
			JOIN units_v2 USING(unit_id)`

	var wrs []*application.WorkReport
	rows, err := s.db.Query(context.Background(), query)

	if err != nil {
		return nil, application.Metadata{}, fmt.Errorf("failed to find work reports: %w", err)
	}

	for rows.Next() {

		var (
			wr   application.WorkReport
			unit application.Unit
		)

		err = rows.Scan(
			&wr.ID, &wr.UnitID, &wr.DocName, &wr.From, &wr.To, &wr.Text, &wr.Data, // work_report
			&unit.ID, &unit.Name, &unit.StoragePath, &unit.PrometheusServerAddr, // unit
		)

		if err != nil {
			return nil, application.Metadata{}, fmt.Errorf("failed to find work reports: %w", err)
		}

		wr.Unit = &unit

		wrs = append(wrs, &wr)
	}

	meta := application.CalculateMetadata(len(wrs), filter.Page, filter.PageSize)

	return wrs, meta, nil
}

func (s *WorkReportService) UpdateWorkReport(ctx context.Context, id int, upd application.WorkReportUpdate) (*application.WorkReport, error) {
	query := `UPDATE work_reports SET work_report_docname = $1, unit_id = $2, work_report_from = $3, work_report_to = $4 WHERE work_report_id = $5 RETURNING work_report_id`

	var wr application.WorkReport
	err := s.db.QueryRow(context.Background(), query, upd.DocName, upd.UnitID, upd.From, upd.To, id).Scan(&wr.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to update work report: %w", err)
	}

	return &wr, nil
}

func (s *WorkReportService) DeleteWorkReport(ctx context.Context, id int) error {

	query := `DELETE FROM work_reports WHERE work_report_id = $1`

	_, err := s.db.Exec(context.Background(), query, id)

	if err != nil {
		return fmt.Errorf("failed to delete work report: %w", err)
	}

	return nil
}

func (s *WorkReportService) CreateWorkReportTopic(ctx context.Context, wr *application.WorkReportTopic) error {

	query := `INSERT INTO work_report_topics (work_report_topic_title, work_report_topic_text, work_report_id, work_report_topic_text_hash) VALUES ($1, $2, $3, hashtext($2)) RETURNING work_report_topic_id`

	err := s.db.QueryRow(context.Background(), query, wr.Title, wr.Text, wr.WorkReportID).Scan(&wr.ID)

	if err != nil {
		return fmt.Errorf("failed to create work report topic: %w", err)
	}

	return nil
}

func (s *WorkReportService) DeleteDuplicatesWorkReportTopics(ctx context.Context) error {

	query := `WITH Duplicates AS (

					SELECT
						work_report_topic_id,
						ROW_NUMBER() OVER (
							PARTITION BY work_report_topic_title, unit_id, work_report_topic_text_hash
							ORDER BY work_report_from ASC)
							AS rn
					FROM
						work_report_topics JOIN work_reports USING(work_report_id)

				)
				DELETE FROM work_report_topics
				WHERE

					work_report_topic_id IN (
						SELECT
							work_report_topic_id
						FROM
							Duplicates
						WHERE
							rn > 1

				)`

	_, err := s.db.Exec(context.Background(), query)

	if err != nil {
		return fmt.Errorf("failed to delete work report: %w", err)
	}

	return nil
}

func (s *WorkReportService) FindWorkReportTopicByID(ctx context.Context, id int) (*application.WorkReportTopic, error) {
	query := `SELECT work_report_topic_id, work_report_id, work_report_topic_title, work_report_topic_text FROM work_report_topics WHERE work_report_topic_id = $1`

	var wrt application.WorkReportTopic
	err := s.db.QueryRow(context.Background(), query, id).Scan(&wrt.ID, &wrt.WorkReportID, &wrt.Title, &wrt.Text)

	if err != nil {
		return nil, fmt.Errorf("failed to find work report topic: %w", err)
	}

	return &wrt, nil
}

// func dump(data interface{}) {
// 	jsonData, err := json.MarshalIndent(data, "", "  ")
// 	if err != nil {
// 		fmt.Println("Erro ao serializar dados:", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println(string(jsonData))
// }

func (s *WorkReportService) FindWorkReportTopics(ctx context.Context, filter application.WorkReportTopicFilter) ([]*application.WorkReportTopic, application.Metadata, error) {
	var (
		conditions []string
		args       []interface{}
		argID      = 1
	)

	// Função auxiliar genérica para adicionar condições
	addCondition := func(condition string, values ...interface{}) {
		placeholders := make([]interface{}, len(values))
		for i := range values {
			placeholders[i] = argID + i
		}
		conditions = append(conditions, fmt.Sprintf(condition, placeholders...))
		args = append(args, values...)
		argID += len(values)
	}

	// Aplicar os filtros
	if filter.ID != nil {
		addCondition("work_report_topic_id = $%d", *filter.ID)
	}
	if filter.Title != nil {
		addCondition("work_report_topic_title ILIKE $%d", "%"+*filter.Title+"%")
	}
	if filter.Text != nil {
		addCondition("work_report_topic_text ILIKE $%d", "%"+*filter.Text+"%")
	}
	if filter.UnitID != nil {
		addCondition("work_report_id = $%d", *filter.UnitID)
	}
	if filter.From != nil {
		addCondition("work_report_from = $%d", *filter.From)
	}
	if filter.To != nil {
		addCondition("work_report_to = $%d", *filter.To)
	}
	if filter.GlobalSearch != nil {
		search := "%" + *filter.GlobalSearch + "%"
		addCondition("(work_report_topic_title ILIKE $%d OR work_report_topic_text ILIKE $%d)", search, search)
	}

	// Construir a query de contagem, para calcular a metadata
	countQuery := `
		SELECT COUNT(1)
			FROM work_report_topics
		JOIN work_reports USING(work_report_id)
		JOIN units_v2 USING(unit_id)
	`

	// Construir a query principal
	query := `
		SELECT
			work_report_topic_id,
			work_report_topic_title,
			work_report_topic_text,
			work_report_id AS work_report_topic_work_report_id,
			work_report_id,
			work_report_from,
			work_report_to,
			work_report_docname,
			unit_id AS work_report_unit_id,
			unit_id,
			unit_name,
			storage_path as unit_storage_path,
			prometheus_server_address as unit_prometheus_server_address
		FROM work_report_topics
		JOIN work_reports USING(work_report_id)
		JOIN units_v2 USING(unit_id)
	`

	// Adicionar condições se existirem
	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		countQuery += whereClause
		query += whereClause
	}

	var count int
	if err := s.db.QueryRow(ctx, countQuery, args...).Scan(&count); err != nil {
		return nil, application.Metadata{}, fmt.Errorf("failed to count work report topics: %w", err)
	}

	// Adicionar ordenação, limite e offset
	query += formatOrderBy(
		filter.SortBy,
		filter.SortDescending,
		"work_report_topic_title",
		[]string{
			"work_report_topic_title",
			"work_report_docname",
			"work_report_from",
			"work_report_to",
		},
	)
	// Adicionar limit e offset
	query += formatLimitOffset(filter.Limit(), filter.Offset())

	// Executar a query
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, application.Metadata{}, fmt.Errorf("failed to find work report topics: %w", err)
	}
	defer rows.Close()

	// Processar os resultados
	var wrts []*application.WorkReportTopic
	for rows.Next() {
		var (
			wrt  application.WorkReportTopic
			wr   application.WorkReport
			unit application.Unit
		)
		// Fazer o scan das colunas
		if err := rows.Scan(
			&wrt.ID, &wrt.Title, &wrt.Text, &wrt.WorkReportID, // work_report_topic
			&wr.ID, &wr.From, &wr.To, &wr.DocName, &wr.UnitID, // work_report
			&unit.ID, &unit.Name, &unit.StoragePath, &unit.PrometheusServerAddr, // unit
		); err != nil {
			return nil, application.Metadata{}, fmt.Errorf("failed to scan work report topics: %w", err)
		}
		// Associar o relatório de trabalho e a unidade ao tópico
		wr.Unit = &unit
		wrt.WorkReport = &wr
		wrts = append(wrts, &wrt)
	}

	// Verificar se houve algum erro durante a iteração
	if err = rows.Err(); err != nil {
		return nil, application.Metadata{}, fmt.Errorf("row iteration error: %w", err)
	}

	// Calcular metadata de paginação
	meta := application.CalculateMetadata(count, filter.Page, filter.PageSize)

	return wrts, meta, nil
}

func (s *WorkReportService) FindWorkReportTopicsAdvSearch(ctx context.Context, filter application.WRAdvSearchFilter) ([]*application.WRAdvSearchResult, application.Metadata, error) {
	// implementar lógica de busca avançada de tópicos de relatório (full text search)
	return nil, application.Metadata{}, nil
}
