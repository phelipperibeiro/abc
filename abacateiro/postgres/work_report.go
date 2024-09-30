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
	query := `SELECT work_report_id, unit_id, work_report_docname, work_report_from, work_report_to, work_report_text, work_report_data FROM work_reports`

	var wrs []*application.WorkReport
	rows, err := s.db.Query(context.Background(), query)

	if err != nil {
		return nil, application.Metadata{}, fmt.Errorf("failed to find work reports: %w", err)
	}

	for rows.Next() {
		var wr application.WorkReport

		err = rows.Scan(&wr.ID, &wr.UnitID, &wr.DocName, &wr.From, &wr.To, &wr.Text, &wr.Data)

		if err != nil {
			return nil, application.Metadata{}, fmt.Errorf("failed to find work reports: %w", err)
		}

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

func (s *WorkReportService) FindWorkReportTopics(ctx context.Context, filter application.WorkReportTopicFilter) ([]*application.WorkReportTopic, application.Metadata, error) {

	query := `SELECT
				work_report_topic_id,
				work_report_topic_title,
				work_report_topic_text,
				work_report_id
			FROM work_report_topics
			JOIN work_reports USING(work_report_id)`

	var conditions []string
	var args []interface{}
	argID := 1 // $1

	if filter.ID != nil {
		conditions = append(conditions, fmt.Sprintf("work_report_topic_id = $%d", argID))
		args = append(args, *filter.ID)
		argID++ // $2
	}

	if filter.Title != nil {
		conditions = append(conditions, fmt.Sprintf("work_report_topic_title ILIKE $%d", argID))
		args = append(args, "%"+*filter.Title+"%")
		argID++ // $3
	}

	if filter.Text != nil {
		conditions = append(conditions, fmt.Sprintf("work_report_topic_text ILIKE $%d", argID))
		args = append(args, "%"+*filter.Text+"%")
		argID++ // $4
	}

	if filter.UnitID != nil {
		conditions = append(conditions, fmt.Sprintf("work_report_id = $%d", argID))
		args = append(args, *filter.UnitID)
		argID++ // $5
	}

	if filter.From != nil {
		conditions = append(conditions, fmt.Sprintf("work_report_from = $%d", argID))
		args = append(args, *filter.From)
		argID++ // $6
	}

	if filter.To != nil {
		conditions = append(conditions, fmt.Sprintf("work_report_to = $%d", argID))
		args = append(args, *filter.To)
		argID++ // $7
	}

	if filter.GlobalSearch != nil {
		search := "%" + *filter.GlobalSearch + "%"
		conditions = append(conditions, fmt.Sprintf("(work_report_topic_title ILIKE $%d OR work_report_topic_text ILIKE $%d)", argID, argID))
		args = append(args, search)
		argID++ // 8
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var wrts []*application.WorkReportTopic
	rows, err := s.db.Query(context.Background(), query, args...)

	if err != nil {
		return nil, application.Metadata{}, fmt.Errorf("failed to find work report topics: %w", err)
	}

	for rows.Next() {

		var wrt application.WorkReportTopic

		err = rows.Scan(&wrt.ID, &wrt.Title, &wrt.Text, &wrt.WorkReportID)

		if err != nil {
			return nil, application.Metadata{}, fmt.Errorf("failed to find work report topics: %w", err)
		}

		wrts = append(wrts, &wrt)
	}
	if rows.Err() != nil {
		return nil, application.Metadata{}, fmt.Errorf("row iteration error: %w", rows.Err())
	}

	return wrts, application.Metadata{}, nil
}

func (s *WorkReportService) FindWorkReportTopicsAdvSearch(ctx context.Context, filter application.WRAdvSearchFilter) ([]*application.WRAdvSearchResult, application.Metadata, error) {
	// implementar lógica de busca avançada de tópicos de relatório (full text search)
	return nil, application.Metadata{}, nil
}
