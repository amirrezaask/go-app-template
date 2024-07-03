package sqldb

type Mysql struct {
	Conn             *sql.DB
	requestDurations prometheus.Histogram
	readFailures     prometheus.Gauge
	updateFailures   prometheus.Gauge
	n                string
}

// ExecContext implements DB.
func (s Mysql) ExecContext(ctx context.Context, stmt string, args ...any) (sql.Result, error) {
	now := time.Now()
	res, err := s.Conn.ExecContext(ctx, stmt, args...)
	if err == nil {
		if s.requestDurations != nil {
			s.requestDurations.(prometheus.ExemplarObserver).ObserveWithExemplar(
				float64(time.Since(now).Milliseconds())/1000.0, make(map[string]string, 0),
			)
		}
		if s.updateFailures != nil {
			s.updateFailures.Set(0)
		}
	} else {
		if s.updateFailures != nil {
			s.updateFailures.Inc()
		}
	}
	return res, err
}

// QueryContext implements DB.
func (s Mysql) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	now := time.Now()
	rows, err := s.Conn.QueryContext(ctx, query, args...)
	if err == nil {
		if s.requestDurations != nil {
			s.requestDurations.(prometheus.ExemplarObserver).ObserveWithExemplar(
				float64(time.Since(now).Milliseconds())/1000.0, make(map[string]string, 0),
			)
		}
		if s.readFailures != nil {
			s.readFailures.Set(0)
		}
	} else {
		if s.readFailures != nil {
			s.readFailures.Inc()
		}
	}
	return rows, err
}

// Tx implements DB.
func (s Mysql) Tx(ctx context.Context) (Tx, error) {
	return s.Conn.BeginTx(ctx, nil)
}

func NewMySQL(connectionString string, dbName string, reg prometheus.Registerer) (DB, error) {
	slog.Debug("connecting to mysql", "db", connectionString)
	conn, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	readFailures := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("mysql_%s_read_failures", dbName),
			Help: fmt.Sprintf("Mysql %s read failures", dbName),
		},
	)
	updateFailures := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("mysql_%s_update_failures", dbName),
			Help: fmt.Sprintf("Mysql %s update failures", dbName),
		},
	)
	requestDurations := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    fmt.Sprintf("mysql_%s_request_duration_seconds", dbName),
		Help:    fmt.Sprintf("A histogram of the %s MySQL request durations in milliseconds.", dbName),
		Buckets: []float64{0.001, 0.01, 0.1, 1, 5},
	})
	reg.MustRegister(
		readFailures,
		updateFailures,
		requestDurations,
	)
	return Mysql{Conn: conn, requestDurations: requestDurations, readFailures: readFailures, updateFailures: updateFailures, n: dbName}, nil
}
