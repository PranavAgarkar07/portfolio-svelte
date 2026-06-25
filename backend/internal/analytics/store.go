package analytics

import (
	"database/sql"
	"fmt"
	"time"
)

func InsertSession(db *sql.DB, s SessionPayload) error {
	_, err := db.Exec(
		`INSERT INTO sessions (id, ip_hash, country, city, referrer, device, os, browser, theme, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		 ON CONFLICT (id) DO NOTHING`,
		s.ID, s.IPHash, s.Country, s.City, s.Referrer, s.Device, s.OS, s.Browser, s.Theme,
	)
	return err
}

func InsertEvents(db *sql.DB, events []EventPayload) error {
	if len(events) == 0 {
		return nil
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

	db.QueryRow(`SELECT COUNT(*) FROM sessions`).Scan(&ds.TotalViews)
	db.QueryRow(`SELECT COUNT(DISTINCT ip_hash) FROM sessions WHERE created_at >= $1`, since).Scan(&ds.UniqueVisitors)
	db.QueryRow(`SELECT COUNT(*) FROM events WHERE type = 'click' AND target = 'resume_pdf' AND ts >= $1`, since).Scan(&ds.ResumeDownloads)
	db.QueryRow(`SELECT COUNT(*) FROM events WHERE type = 'form' AND target = 'form_submit' AND ts >= $1`, since).Scan(&ds.FormSubmissions)

	_ = db.QueryRow(`SELECT COALESCE(AVG(NULLIF(value, '')::numeric), 0) FROM events WHERE type = 'timing' AND target = 'session_duration' AND ts >= $1`,
		since).Scan(&ds.AvgTimeOnSite)

	ds.TopReferrers = GetReferrerBreakdown(db, since)
	ds.TopTargets = GetTopTargets(db, since, 20)
	ds.CountryBreakdown = GetCountryBreakdown(db, since)

	return ds
}

func GetTopTargets(db *sql.DB, since time.Time, limit int) []TargetCount {
	rows, err := db.Query(
		`SELECT target, COUNT(*) as cnt FROM events WHERE ts >= $1 AND target != '' GROUP BY target ORDER BY cnt DESC LIMIT $2`,
		since, limit,
	)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []TargetCount
	for rows.Next() {
		var tc TargetCount
		if err := rows.Scan(&tc.Target, &tc.Count); err != nil {
			continue
		}
		out = append(out, tc)
	}
	return out
}

func GetCountryBreakdown(db *sql.DB, since time.Time) []CountryCount {
	rows, err := db.Query(
		`SELECT country, COUNT(*) as cnt FROM sessions WHERE created_at >= $1 AND country != '' GROUP BY country ORDER BY cnt DESC`,
		since,
	)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []CountryCount
	for rows.Next() {
		var cc CountryCount
		if err := rows.Scan(&cc.Country, &cc.Count); err != nil {
			continue
		}
		out = append(out, cc)
	}
	return out
}

func GetReferrerBreakdown(db *sql.DB, since time.Time) []ReferrerCount {
	rows, err := db.Query(
		`SELECT referrer, COUNT(*) as cnt FROM sessions WHERE created_at >= $1 AND referrer != '' GROUP BY referrer ORDER BY cnt DESC`,
		since,
	)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []ReferrerCount
	for rows.Next() {
		var rc ReferrerCount
		if err := rows.Scan(&rc.Referrer, &rc.Count); err != nil {
			continue
		}
		out = append(out, rc)
	}
	return out
}
