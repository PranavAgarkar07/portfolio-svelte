package analytics

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

func InsertSession(db *sql.DB, s SessionPayload) error {
	_, err := db.Exec(
		`INSERT INTO sessions (id, ip_hash, country, city, referrer, source, device, os, browser, theme, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
		 ON CONFLICT (id) DO NOTHING`,
		s.ID, s.IPHash, s.Country, s.City, s.Referrer, s.Source, s.Device, s.OS, s.Browser, s.Theme,
	)
	return err
}

const maxBatchSize = 50

func InsertEvents(db *sql.DB, events []EventPayload) error {
	if len(events) == 0 {
		return nil
	}
	if len(events) > maxBatchSize {
		events = events[:maxBatchSize]
	}

	query := `INSERT INTO events (session_id, type, target, value, ts) VALUES `
	params := make([]interface{}, 0, len(events)*5)
	for i, e := range events {
		if i > 0 {
			query += ", "
		}
		base := i * 5
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, to_timestamp($%d/1000.0))",
			base+1, base+2, base+3, base+4, base+5)
		params = append(params, e.SessionID, e.Type, e.Target, e.Value, e.Ts)
	}
	_, err := db.Exec(query, params...)
	return err
}

func GetDashboardStats(db *sql.DB, since time.Time) DashboardStats {
	var ds DashboardStats
	sinceTs := since.Unix()

	if err := db.QueryRow(`SELECT COUNT(*) FROM sessions`).Scan(&ds.TotalViews); err != nil {
		slog.Error("analytics: failed to count total views", "error", err)
	}
	if err := db.QueryRow(`SELECT COUNT(DISTINCT ip_hash) FROM sessions WHERE created_at >= to_timestamp($1)`, sinceTs).Scan(&ds.UniqueVisitors); err != nil {
		slog.Error("analytics: failed to count unique visitors", "error", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM events WHERE type = 'click' AND target = 'resume_pdf' AND ts >= to_timestamp($1)`, sinceTs).Scan(&ds.ResumeDownloads); err != nil {
		slog.Error("analytics: failed to count resume downloads", "error", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM events WHERE type = 'form' AND target = 'form_submit' AND ts >= to_timestamp($1)`, sinceTs).Scan(&ds.FormSubmissions); err != nil {
		slog.Error("analytics: failed to count form submissions", "error", err)
	}
	if err := db.QueryRow(`SELECT COALESCE(AVG(NULLIF(value, '')::numeric), 0) FROM events WHERE type = 'timing' AND target = 'session_duration' AND ts >= to_timestamp($1)`,
		sinceTs).Scan(&ds.AvgTimeOnSite); err != nil {
		slog.Error("analytics: failed to compute avg time on site", "error", err)
	}

	ds.TopReferrers = GetReferrerBreakdown(db, sinceTs)
	ds.TopTargets = GetTopTargets(db, sinceTs, 20)
	ds.CountryBreakdown = GetCountryBreakdown(db, sinceTs)
	ds.SourceBreakdown = GetSourceBreakdown(db, sinceTs)

	return ds
}

func GetTopTargets(db *sql.DB, sinceTs int64, limit int) []TargetCount {
	rows, err := db.Query(
		`SELECT target, COUNT(*) as cnt FROM events WHERE ts >= to_timestamp($1) AND target != '' GROUP BY target ORDER BY cnt DESC LIMIT $2`,
		sinceTs, limit,
	)
	if err != nil {
		slog.Error("analytics: failed to query top targets", "error", err)
		return []TargetCount{}
	}
	defer rows.Close()
	out := []TargetCount{}
	for rows.Next() {
		var tc TargetCount
		if err := rows.Scan(&tc.Target, &tc.Count); err != nil {
			slog.Error("analytics: failed to scan target row", "error", err)
			continue
		}
		out = append(out, tc)
	}
	return out
}

func GetCountryBreakdown(db *sql.DB, sinceTs int64) []CountryCount {
	rows, err := db.Query(
		`SELECT country, COUNT(*) as cnt FROM sessions WHERE created_at >= to_timestamp($1) AND country != '' GROUP BY country ORDER BY cnt DESC`,
		sinceTs,
	)
	if err != nil {
		slog.Error("analytics: failed to query country breakdown", "error", err)
		return []CountryCount{}
	}
	defer rows.Close()
	out := []CountryCount{}
	for rows.Next() {
		var cc CountryCount
		if err := rows.Scan(&cc.Country, &cc.Count); err != nil {
			slog.Error("analytics: failed to scan country row", "error", err)
			continue
		}
		out = append(out, cc)
	}
	return out
}

func GetReferrerBreakdown(db *sql.DB, sinceTs int64) []ReferrerCount {
	rows, err := db.Query(
		`SELECT referrer, COUNT(*) as cnt FROM sessions WHERE created_at >= to_timestamp($1) AND referrer != '' GROUP BY referrer ORDER BY cnt DESC`,
		sinceTs,
	)
	if err != nil {
		slog.Error("analytics: failed to query referrer breakdown", "error", err)
		return []ReferrerCount{}
	}
	defer rows.Close()
	out := []ReferrerCount{}
	for rows.Next() {
		var rc ReferrerCount
		if err := rows.Scan(&rc.Referrer, &rc.Count); err != nil {
			slog.Error("analytics: failed to scan referrer row", "error", err)
			continue
		}
		out = append(out, rc)
	}
	return out
}

func GetSourceBreakdown(db *sql.DB, sinceTs int64) []SourceCount {
	rows, err := db.Query(
		`SELECT COALESCE(NULLIF(source, ''), 'direct') as src, COUNT(*) as cnt FROM sessions WHERE created_at >= to_timestamp($1) GROUP BY src ORDER BY cnt DESC`,
		sinceTs,
	)
	if err != nil {
		slog.Error("analytics: failed to query source breakdown", "error", err)
		return []SourceCount{}
	}
	defer rows.Close()
	out := []SourceCount{}
	for rows.Next() {
		var sc SourceCount
		if err := rows.Scan(&sc.Source, &sc.Count); err != nil {
			slog.Error("analytics: failed to scan source row", "error", err)
			continue
		}
		out = append(out, sc)
	}
	return out
}

const dataRetentionDays = 90

func CleanupOldData(db *sql.DB) {
	cutoff := time.Now().AddDate(0, 0, -dataRetentionDays).Unix()
	result, err := db.Exec(`DELETE FROM events WHERE ts < to_timestamp($1)`, cutoff)
	if err != nil {
		slog.Error("analytics: failed to cleanup old events", "error", err)
	} else if n, _ := result.RowsAffected(); n > 0 {
		slog.Info("analytics: cleaned up old events", "count", n)
	}
	result, err = db.Exec(`DELETE FROM sessions WHERE created_at < to_timestamp($1)`, cutoff)
	if err != nil {
		slog.Error("analytics: failed to cleanup old sessions", "error", err)
	} else if n, _ := result.RowsAffected(); n > 0 {
		slog.Info("analytics: cleaned up old sessions", "count", n)
	}
}
