package main

import "time"

type PingResponse struct {
	Ping string `json:"ping"`
	User string `json:"user"`
}

type GetTableResponse struct {
	Description         string `json:"description,omitempty"`
	FullName            string `json:"full_name,omitempty"`
	ID                  int    `json:"id,omitempty"`
	Monitored           bool   `json:"monitored,omitempty"`
	NotificationChannel struct {
		ChannelType string `json:"channel_type,omitempty"`
		Description string `json:"description,omitempty"`
		ID          int    `json:"id,omitempty"`
	} `json:"notification_channel,omitempty"`
	RecentStatus struct {
		IntervalID           int       `json:"interval_id,omitempty"`
		LatestRunChecksJobID string    `json:"latest_run_checks_job_id,omitempty"`
		Status               string    `json:"status,omitempty"`
		StatusDisplay        string    `json:"status_display,omitempty"`
		TimePeriodEnd        time.Time `json:"time_period_end,omitempty"`
		TimePeriodStart      time.Time `json:"time_period_start,omitempty"`
	} `json:"recent_status,omitempty"`
	Warehouse struct {
		ID   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"warehouse,omitempty"`
	Config struct {
		TableID                   int           `json:"table_id,omitempty"`
		CheckCadenceType          string        `json:"check_cadence_type,omitempty"`
		Definition                string        `json:"definition,omitempty"`
		TimeColumnType            string        `json:"time_column_type,omitempty"`
		NotifyAfter               string        `json:"notify_after,omitempty"`
		NotificationChannelID     int           `json:"notification_channel_id,omitempty"`
		TimeColumns               []interface{} `json:"time_columns,omitempty"`
		FreshAfter                string        `json:"fresh_after,omitempty"`
		CheckCadenceRunAtDuration string        `json:"check_cadence_run_at_duration,omitempty"`
		IntervalSkipExpr          string        `json:"interval_skip_expr,omitempty"`
		Created                   time.Time     `json:"created,omitempty"`
		CreatedBy                 struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"created_by,omitempty"`
		LastEditedAt string `json:"last_edited_at,omitempty"`
		LastEditedBy struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"last_edited_by,omitempty"`
		SlackUsers []interface{} `json:"slack_users,omitempty"`
	} `json:"config,omitempty"`
}

type ConfigureTableResponse struct {
	TableID                   int           `json:"table_id,omitempty"`
	CheckCadenceType          string        `json:"check_cadence_type,omitempty"`
	Definition                string        `json:"definition,omitempty"`
	TimeColumnType            string        `json:"time_column_type,omitempty"`
	NotifyAfter               string        `json:"notify_after,omitempty"`
	NotificationChannelID     int           `json:"notification_channel_id,omitempty"`
	TimeColumns               []interface{} `json:"time_columns,omitempty"`
	FreshAfter                string        `json:"fresh_after,omitempty"`
	CheckCadenceRunAtDuration string        `json:"check_cadence_run_at_duration,omitempty"`
	IntervalSkipExpr          string        `json:"interval_skip_expr,omitempty"`
	Created                   time.Time     `json:"created,omitempty"`
	CreatedBy                 struct {
		ID   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"created_by,omitempty"`
	LastEditedAt string `json:"last_edited_at,omitempty"`
	LastEditedBy struct {
		ID   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"last_edited_by,omitempty"`
	SlackUsers []interface{} `json:"slack_users,omitempty"`
}
