package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"orders-api/events" // ajuste para o nome do seu módulo
	"orders-api/store"

	ulid "github.com/oklog/ulid/v2"
)

// ──────────────────────────────────────────────────────────────────────────────
// Tipos de request/response específicos da API
// ──────────────────────────────────────────────────────────────────────────────

type createReq struct {
	Customer string   `json:"customer"`
	Items    []string `json:"items"`
}

type updateStatusReq struct {
	Status string `json:"status"`
}

// ──────────────────────────────────────────────────────────────────────────────
// ULID helper (domínio de orders da API)
// ──────────────────────────────────────────────────────────────────────────────

var ulidEntropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)

func newID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy).String()
}

// ──────────────────────────────────────────────────────────────────────────────
// Server
// ──────────────────────────────────────────────────────────────────────────────

type Server struct {
	db        *sql.DB
	publisher *events.Publisher
	mux       *http.ServeMux
}

// NewServer recebe as dependências (DB e Kafka publisher) e monta as rotas.
func NewServer(db *sql.DB, publisher *events.Publisher) *Server {
	s := &Server{
		db:        db,
		publisher: publisher,
		mux:       http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

// ServeHTTP implementa http.Handler e delega para o mux interno.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// ──────────────────────────────────────────────────────────────────────────────
// Rotas
// ──────────────────────────────────────────────────────────────────────────────

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/orders", s.handleOrders)
	s.mux.HandleFunc("/orders/", s.handleOrderByID)
}

// ──────────────────────────────────────────────────────────────────────────────
// Handlers
// ──────────────────────────────────────────────────────────────────────────────

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"ok":true}`))
}

// /orders → POST (criar) / GET (listar)
func (s *Server) handleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.handleCreateOrder(w, r)
	case http.MethodGet:
		s.handleListOrders(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// /orders/{id}           → GET
// /orders/{id}/status    → PUT
func (s *Server) handleOrderByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/orders/")
	if path == "" {
		http.NotFound(w, r)
		return
	}

	// PUT /orders/{id}/status
	if strings.HasSuffix(path, "/status") {
		if r.Method != http.MethodPut {
			http.Error(w, "use PUT", http.StatusMethodNotAllowed)
			return
		}
		id := strings.TrimSuffix(path, "/status")
		id = strings.TrimSuffix(id, "/")
		s.handleUpdateStatus(w, r, id)
		return
	}

	// GET /orders/{id}
	if r.Method != http.MethodGet {
		http.Error(w, "use GET", http.StatusMethodNotAllowed)
		return
	}
	if strings.Contains(path, "/") {
		http.Error(w, "invalid path; expected /orders/{id}", http.StatusBadRequest)
		return
	}
	s.handleGetOrder(w, r, path)
}

// ──────────────────────────────────────────────────────────────────────────────
// Handlers específicos
// ──────────────────────────────────────────────────────────────────────────────

func (s *Server) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()
	itemsJSON, _ := json.Marshal(req.Items)

	id := newID()

	if _, err := s.db.Exec(`INSERT INTO orders (id, customer, status, items_json, created_at, updated_at)
		VALUES (?,?,?,?,?,?)`,
		id, req.Customer, "OPEN", string(itemsJSON), now, now); err != nil {
		log.Printf("ERROR insert order: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	evt := map[string]any{
		"type":     "OrderCreated",
		"id":       id,
		"customer": req.Customer,
		"status":   "OPEN",
		"items":    req.Items,
		"ts":       now.Format(time.RFC3339Nano),
	}

	if _, err := s.publisher.PublishWithDigest(r.Context(), id, evt,
		map[string]string{"x-event": "OrderCreated"}); err != nil {
		log.Printf("WARN publish OrderCreated failed: %v", err)
		http.Error(w, "kafka unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":       id,
		"customer": req.Customer,
		"items":    req.Items,
		"status":   "OPEN",
	})
}

func (s *Server) handleListOrders(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	status := q.Get("status")
	customer := q.Get("customer")
	since := q.Get("since")
	until := q.Get("until")

	limit, offset := 50, 0
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 500 {
			limit = n
		}
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	var (
		conds []string
		args  []any
	)
	if status != "" {
		conds = append(conds, "status = ?")
		args = append(args, status)
	}
	if customer != "" {
		conds = append(conds, "customer LIKE ?")
		args = append(args, "%"+customer+"%")
	}
	if since != "" {
		if t, err := time.Parse(time.RFC3339, since); err == nil {
			conds = append(conds, "created_at >= ?")
			args = append(args, t)
		}
	}
	if until != "" {
		if t, err := time.Parse(time.RFC3339, until); err == nil {
			conds = append(conds, "created_at <= ?")
			args = append(args, t)
		}
	}

	var sb strings.Builder
	sb.WriteString("SELECT id, customer, status, items_json, created_at, updated_at FROM orders")
	if len(conds) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(conds, " AND "))
	}
	sb.WriteString(" ORDER BY created_at DESC, id DESC")
	sb.WriteString(" LIMIT ? OFFSET ?")
	args = append(args, limit, offset)

	rows, err := s.db.Query(sb.String(), args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	list, err := store.ScanOrders(rows)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"items":  list,
		"limit":  limit,
		"offset": offset,
		"count":  len(list),
	})
}

func (s *Server) handleUpdateStatus(w http.ResponseWriter, r *http.Request, id string) {
	var req updateStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()

	res, err := s.db.Exec(`UPDATE orders SET status=?, updated_at=? WHERE id=?`, req.Status, now, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	evt := map[string]any{
		"type":   "OrderStatusUpdated",
		"id":     id,
		"status": req.Status,
		"ts":     now.Format(time.RFC3339Nano),
	}
	if _, err := s.publisher.PublishWithDigest(r.Context(), id, evt,
		map[string]string{"x-event": "OrderStatusUpdated"}); err != nil {
		log.Printf("WARN publish OrderStatusUpdated failed: %v", err)
		http.Error(w, "kafka unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":     id,
		"status": req.Status,
	})
}

func (s *Server) handleGetOrder(w http.ResponseWriter, r *http.Request, id string) {
	row := s.db.QueryRow(`SELECT id, customer, status, items_json, created_at, updated_at FROM orders WHERE id=?`, id)
	o, err := store.ScanOrder(row)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(o)
}
