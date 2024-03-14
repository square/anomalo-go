package anomalo

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
		RecentIntervals []struct {
			IntervalID           int       `json:"interval_id,omitempty"`
			LatestRunChecksJobID string    `json:"latest_run_checks_job_id,omitempty"`
			Status               string    `json:"status,omitempty"`
			StatusDisplay        string    `json:"status_display,omitempty"`
			TimePeriodEnd        time.Time `json:"time_period_end,omitempty"`
			TimePeriodStart      time.Time `json:"time_period_start,omitempty"`
		} `json:"recent_intervals,omitempty"`
	} `json:"recent_status,omitempty"`
	Warehouse struct {
		ID   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"warehouse,omitempty"`
	Config struct {
		TableID                   int       `json:"table_id,omitempty"`
		CheckCadenceType          string    `json:"check_cadence_type,omitempty"`
		Definition                string    `json:"definition,omitempty"`
		TimeColumnType            string    `json:"time_column_type,omitempty"`
		NotifyAfter               string    `json:"notify_after,omitempty"`
		NotificationChannelID     int       `json:"notification_channel_id,omitempty"`
		TimeColumns               []string  `json:"time_columns,omitempty"`
		FreshAfter                string    `json:"fresh_after,omitempty"`
		CheckCadenceRunAtDuration string    `json:"check_cadence_run_at_duration,omitempty"`
		IntervalSkipExpr          string    `json:"interval_skip_expr,omitempty"`
		AlwaysAlertOnErrors       bool      `json:"always_alert_on_errors,omitempty"`
		Created                   time.Time `json:"created,omitempty"`
		CreatedBy                 struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"created_by,omitempty"`
		LastEditedAt string `json:"last_edited_at,omitempty"`
		LastEditedBy struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"last_edited_by,omitempty"`
	} `json:"config,omitempty"`
}

type ConfigureTableRequest struct {
	TableID                   int      `json:"table_id,omitempty"`
	CheckCadenceType          *string  `json:"check_cadence_type"`
	Definition                string   `json:"definition,omitempty"`
	TimeColumnType            string   `json:"time_column_type,omitempty"`
	NotifyAfter               string   `json:"notify_after,omitempty"`
	NotificationChannelID     int      `json:"notification_channel_id,omitempty"`
	TimeColumns               []string `json:"time_columns,omitempty"`
	FreshAfter                string   `json:"fresh_after,omitempty"`
	CheckCadenceRunAtDuration string   `json:"check_cadence_run_at_duration,omitempty"`
	IntervalSkipExpr          string   `json:"interval_skip_expr,omitempty"`
	AlwaysAlertOnErrors       bool     `json:"always_alert_on_errors,omitempty"`
}

type ConfigureTableResponse struct {
	ID            int       `json:"id,omitempty"`
	Created       time.Time `json:"created,omitempty"`
	Modified      time.Time `json:"modified,omitempty"`
	SchemaID      int       `json:"schema_id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Definition    string    `json:"definition,omitempty"`
	LastRefreshed time.Time `json:"last_refreshed,omitempty"`
	Config        struct {
	} `json:"config,omitempty"`
	UpdateModified bool `json:"update_modified,omitempty"`
}

type Check struct {
	CheckID       int    `json:"check_id,omitempty"`
	CheckStaticID int    `json:"check_static_id,omitempty"`
	Ref           string `json:"ref,omitempty"`
	CheckType     string `json:"check_type,omitempty"`
	Config        struct {
		Metadata struct {
			CheckMessage     string `json:"check_message,omitempty"`
			CheckMessageHTML string `json:"check_message_html,omitempty"`
			CheckType        string `json:"check_type,omitempty"`
			Description      string `json:"description,omitempty"`
			IsSystemCheck    bool   `json:"is_system_check,omitempty"`
			PriorityLevel    string `json:"priority_level,omitempty"`
		} `json:"_metadata,omitempty"`
		Check  string                 `json:"check,omitempty"`
		Params map[string]interface{} `json:"params,omitempty"`
	} `json:"config,omitempty"`
	Created   time.Time `json:"created,omitempty"`
	CreatedBy struct {
		ID   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"created_by,omitempty"`
	LastEditedAt string `json:"last_edited_at,omitempty"`
	LastEditedBy struct {
		ID   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"last_edited_by,omitempty"`
	TriageStatus                    string `json:"triage_status,omitempty"`
	AdditionalNotificationChannelID int    `json:"additional_notification_channel_id,omitempty"`
}

type GetChecksResponse struct {
	Checks []Check `json:"checks,omitempty"`
}

type CreateCheckRequest struct {
	CheckType string            `json:"check_type,omitempty"`
	Params    map[string]string `json:"params,omitempty"`
	TableID   int               `json:"table_id,omitempty"`
}

type CreateCheckResponse struct {
	CheckID       int    `json:"check_id,omitempty"`
	CheckRef      string `json:"ref,omitempty"`
	CheckStaticId int    `json:"check_static_id,omitempty"`
}

type DeleteCheckRequest struct {
	TableID int `json:"table_id,omitempty"`
	CheckID int `json:"check_id,omitempty"`
}

type DeleteCheckResponse struct {
	DeletedCount int `json:"deleted_count,omitempty"`
}

type NotificationChannel struct {
	ChannelType string `json:"channel_type,omitempty"`
	Description string `json:"description,omitempty"`
	ID          int    `json:"id,omitempty"`
}

type GetNotificationChannelsResponse struct {
	NotificationChannels []NotificationChannel `json:"notification_channels,omitempty"`
}

// TODO not actually an API object
type Organization struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// TODO not actually an API object
type GetOrganizationsResponse struct {
	Organizations []Organization `json:"organizations,omitempty"`
}
