package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	log "github.com/sirupsen/logrus"
)

// Mock data structures for manual approval API

type ManualApprovalTransaction struct {
	TxId        string `json:"txId"`
	QuoteHash   string `json:"quoteHash"`
	Date        string `json:"date"`
	Amount      string `json:"amount"` // wei as string
	Type        string `json:"type"`   // "pegin" or "pegout"
	State       string `json:"state"`
	UserAddress string `json:"userAddress"`
	CallFee     string `json:"callFee"` // wei as string
	GasFee      string `json:"gasFee"`  // wei as string
}

type PendingTransactionsResponse struct {
	Transactions []ManualApprovalTransaction `json:"transactions"`
	TotalCount   int                         `json:"totalCount"`
	Page         int                         `json:"page"`
	PerPage      int                         `json:"perPage"`
}

type HistoryRecord struct {
	ManualApprovalTransaction
	ApprovedOrDeniedBy string `json:"approvedOrDeniedBy"`
	Decision           string `json:"decision"`     // "approved" or "denied"
	DecisionDate       string `json:"decisionDate"` // ISO 8601 date string
}

type HistoryResponse struct {
	History    []HistoryRecord `json:"history"`
	TotalCount int             `json:"totalCount"`
	Page       int             `json:"page"`
	PerPage    int             `json:"perPage"`
}

type ApprovalActionRequest struct {
	TxIds []string `json:"txIds"`
}

type ApprovalActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Mock data generators. Will be replaced with actual data from the database.

func generateMockPendingTransactions() []ManualApprovalTransaction {
	// Hardcoded threshold: 1 BTC = 1,000,000,000,000,000,000 wei
	transactions := []ManualApprovalTransaction{
		{
			TxId:        "0x1a2b3c4d5e6f7890abcdef1234567890abcdef1234567890abcdef1234567890",
			QuoteHash:   "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			Date:        time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			Amount:      "1500000000000000000", // 1.5 BTC
			Type:        "pegin",
			State:       "WaitingForDeposit",
			UserAddress: "0x9876543210987654321098765432109876543210",
			CallFee:     "10000000000000000",
			GasFee:      "5000000000000000",
		},
		{
			TxId:        "0x2b3c4d5e6f7890abcdef1234567890abcdef1234567890abcdef1234567891",
			QuoteHash:   "0xbcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567891",
			Date:        time.Now().Add(-5 * time.Hour).Format(time.RFC3339),
			Amount:      "2000000000000000000", // 2 BTC
			Type:        "pegout",
			State:       "WaitingForDepositConfirmations",
			UserAddress: "0x8765432109876543210987654321098765432101",
			CallFee:     "15000000000000000",
			GasFee:      "7000000000000000",
		},
		{
			TxId:        "0x3c4d5e6f7890abcdef1234567890abcdef1234567890abcdef1234567892",
			QuoteHash:   "0xcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567892",
			Date:        time.Now().Add(-10 * time.Hour).Format(time.RFC3339),
			Amount:      "3500000000000000000", // 3.5 BTC
			Type:        "pegin",
			State:       "WaitingForDeposit",
			UserAddress: "0x7654321098765432109876543210987654321012",
			CallFee:     "20000000000000000",
			GasFee:      "8000000000000000",
		},
		{
			TxId:        "0x4d5e6f7890abcdef1234567890abcdef1234567890abcdef1234567893",
			QuoteHash:   "0xdef1234567890abcdef1234567890abcdef1234567890abcdef1234567893",
			Date:        time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			Amount:      "5000000000000000000", // 5 BTC
			Type:        "pegout",
			State:       "WaitingForDeposit",
			UserAddress: "0x6543210987654321098765432109876543210123",
			CallFee:     "25000000000000000",
			GasFee:      "10000000000000000",
		},
		{
			TxId:        "0x5e6f7890abcdef1234567890abcdef1234567890abcdef1234567894",
			QuoteHash:   "0xef1234567890abcdef1234567890abcdef1234567890abcdef1234567894",
			Date:        time.Now().Add(-48 * time.Hour).Format(time.RFC3339),
			Amount:      "1200000000000000000", // 1.2 BTC
			Type:        "pegin",
			State:       "WaitingForDepositConfirmations",
			UserAddress: "0x5432109876543210987654321098765432101234",
			CallFee:     "12000000000000000",
			GasFee:      "6000000000000000",
		},
	}
	return transactions
}

func generateMockHistory() []HistoryRecord {
	baseTime := time.Now().Add(-72 * time.Hour)
	history := []HistoryRecord{
		{
			ManualApprovalTransaction: ManualApprovalTransaction{
				TxId:        "0xf1234567890abcdef1234567890abcdef1234567890abcdef1234567895",
				QuoteHash:   "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567895",
				Date:        baseTime.Format(time.RFC3339),
				Amount:      "1800000000000000000",
				Type:        "pegin",
				State:       "CallForUserSucceeded",
				UserAddress: "0x4321098765432109876543210987654321012345",
				CallFee:     "18000000000000000",
				GasFee:      "9000000000000000",
			},
			ApprovedOrDeniedBy: "admin@example.com",
			Decision:           "approved",
			DecisionDate:       baseTime.Add(1 * time.Hour).Format(time.RFC3339),
		},
		{
			ManualApprovalTransaction: ManualApprovalTransaction{
				TxId:        "0xe234567890abcdef1234567890abcdef1234567890abcdef1234567896",
				QuoteHash:   "0x234567890abcdef1234567890abcdef1234567890abcdef1234567896",
				Date:        baseTime.Add(-24 * time.Hour).Format(time.RFC3339),
				Amount:      "2500000000000000000",
				Type:        "pegout",
				State:       "SendPegoutSucceeded",
				UserAddress: "0x3210987654321098765432109876543210123456",
				CallFee:     "20000000000000000",
				GasFee:      "11000000000000000",
			},
			ApprovedOrDeniedBy: "admin@example.com",
			Decision:           "approved",
			DecisionDate:       baseTime.Add(-23 * time.Hour).Format(time.RFC3339),
		},
		{
			ManualApprovalTransaction: ManualApprovalTransaction{
				TxId:        "0xd34567890abcdef1234567890abcdef1234567890abcdef1234567897",
				QuoteHash:   "0x34567890abcdef1234567890abcdef1234567890abcdef1234567897",
				Date:        baseTime.Add(-48 * time.Hour).Format(time.RFC3339),
				Amount:      "1000000000000000000",
				Type:        "pegin",
				State:       "TimeForDepositElapsed",
				UserAddress: "0x2109876543210987654321098765432101234567",
				CallFee:     "10000000000000000",
				GasFee:      "5000000000000000",
			},
			ApprovedOrDeniedBy: "admin@example.com",
			Decision:           "denied",
			DecisionDate:       baseTime.Add(-47 * time.Hour).Format(time.RFC3339),
		},
	}
	return history
}

// NewGetPendingTransactionsHandler returns pending transactions requiring manual approval
// @Title Get Pending Transactions
// @Description Returns pending transactions above the approval threshold
// @Success 200 object PendingTransactionsResponse
// @Route /management/manual-approval/pending [get]
func NewGetPendingTransactionsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Parse query parameters
		page := parseIntParam(req, "page", 1)
		perPage := parseIntParam(req, "perPage", 20)
		search := req.URL.Query().Get("search")
		sortDescStr := req.URL.Query().Get("sortDesc")
		sortDesc := sortDescStr == "true"

		// Get mock data
		allTransactions := generateMockPendingTransactions()

		// Filter by search term
		var filteredTransactions []ManualApprovalTransaction
		if search != "" {
			searchLower := strings.ToLower(search)
			for _, tx := range allTransactions {
				if strings.Contains(strings.ToLower(tx.TxId), searchLower) ||
					strings.Contains(strings.ToLower(tx.QuoteHash), searchLower) {
					filteredTransactions = append(filteredTransactions, tx)
				}
			}
		} else {
			filteredTransactions = allTransactions
		}

		// Sort by date
		if !sortDesc {
			// Reverse for ascending
			for i, j := 0, len(filteredTransactions)-1; i < j; i, j = i+1, j-1 {
				filteredTransactions[i], filteredTransactions[j] = filteredTransactions[j], filteredTransactions[i]
			}
		}

		// Paginate
		totalCount := len(filteredTransactions)
		start := (page - 1) * perPage
		end := start + perPage

		if start >= totalCount {
			filteredTransactions = []ManualApprovalTransaction{}
		} else {
			if end > totalCount {
				end = totalCount
			}
			filteredTransactions = filteredTransactions[start:end]
		}

		response := PendingTransactionsResponse{
			Transactions: filteredTransactions,
			TotalCount:   totalCount,
			Page:         page,
			PerPage:      perPage,
		}

		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}

// NewGetHistoryHandler returns approval/denial history
// @Title Get Approval History
// @Description Returns history of approved and denied transactions
// @Success 200 object HistoryResponse
// @Route /management/manual-approval/history [get]
func NewGetHistoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Parse query parameters
		page := parseIntParam(req, "page", 1)
		perPage := parseIntParam(req, "perPage", 20)
		search := req.URL.Query().Get("search")
		status := req.URL.Query().Get("status")
		startDate := req.URL.Query().Get("startDate")
		endDate := req.URL.Query().Get("endDate")
		sortDescStr := req.URL.Query().Get("sortDesc")
		sortDesc := sortDescStr == "true"

		// Get and filter mock data
		allHistory := generateMockHistory()
		filteredHistory := filterHistoryRecords(allHistory, search, status, startDate, endDate)

		// Sort by date
		if !sortDesc {
			reverseHistoryRecords(filteredHistory)
		}

		// Paginate
		totalCount := len(filteredHistory)
		paginatedHistory := paginateHistoryRecords(filteredHistory, page, perPage)

		response := HistoryResponse{
			History:    paginatedHistory,
			TotalCount: totalCount,
			Page:       page,
			PerPage:    perPage,
		}

		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}

// NewApproveTransactionsHandler approves transactions
// @Title Approve Transactions
// @Description Approves one or more pending transactions
// @Success 200 object ApprovalActionResponse
// @Route /management/manual-approval/approve [post]
func NewApproveTransactionsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var request ApprovalActionRequest
		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("Invalid request body", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		if len(request.TxIds) == 0 {
			jsonErr := rest.NewErrorResponse("No transaction IDs provided", false)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		log.Infof("Approved %d transaction(s)", len(request.TxIds))

		// Mock response - in real implementation, this would:
		// 1. Validate transaction IDs exist
		// 2. Execute approval logic (update state, interact with blockchain)
		// 3. Create audit trail entry
		// 4. Return success/failure

		response := ApprovalActionResponse{
			Success: true,
			Message: fmt.Sprintf("Successfully approved %d transaction(s)", len(request.TxIds)),
		}

		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}

// NewDenyTransactionsHandler denies transactions
// @Title Deny Transactions
// @Description Denies one or more pending transactions
// @Success 200 object ApprovalActionResponse
// @Route /management/manual-approval/deny [post]
func NewDenyTransactionsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var request ApprovalActionRequest
		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("Invalid request body", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		if len(request.TxIds) == 0 {
			jsonErr := rest.NewErrorResponse("No transaction IDs provided", false)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		log.Infof("Denied %d transaction(s)", len(request.TxIds))

		// Mock response - in real implementation, this would:
		// 1. Validate transaction IDs exist
		// 2. Execute denial logic (update state, possibly trigger refund)
		// 3. Create audit trail entry
		// 4. Return success/failure

		response := ApprovalActionResponse{
			Success: true,
			Message: fmt.Sprintf("Successfully denied %d transaction(s)", len(request.TxIds)),
		}

		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}

// Helper function to parse integer query parameters
func parseIntParam(req *http.Request, param string, defaultValue int) int {
	valueStr := req.URL.Query().Get(param)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// filterHistoryRecords applies search, status, and date filters to history records
func filterHistoryRecords(records []HistoryRecord, search, status, startDate, endDate string) []HistoryRecord {
	filtered := make([]HistoryRecord, 0, len(records))
	for _, record := range records {
		if !matchesSearchFilter(record, search) {
			continue
		}
		if !matchesStatusFilter(record, status) {
			continue
		}
		if !matchesDateRangeFilter(record, startDate, endDate) {
			continue
		}
		filtered = append(filtered, record)
	}
	return filtered
}

// matchesSearchFilter checks if record matches search term
func matchesSearchFilter(record HistoryRecord, search string) bool {
	if search == "" {
		return true
	}
	searchLower := strings.ToLower(search)
	return strings.Contains(strings.ToLower(record.TxId), searchLower) ||
		strings.Contains(strings.ToLower(record.QuoteHash), searchLower)
}

// matchesStatusFilter checks if record matches status filter
func matchesStatusFilter(record HistoryRecord, status string) bool {
	if status == "" {
		return true
	}
	return record.Decision == status
}

// matchesDateRangeFilter checks if record is within date range
func matchesDateRangeFilter(record HistoryRecord, startDate, endDate string) bool {
	if startDate == "" && endDate == "" {
		return true
	}

	recordDate, err := time.Parse(time.RFC3339, record.Date)
	if err != nil {
		return false
	}

	if startDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err == nil && recordDate.Before(start) {
			return false
		}
	}

	if endDate != "" {
		end, err := time.Parse("2006-01-02", endDate)
		if err == nil && recordDate.After(end.Add(24*time.Hour)) {
			return false
		}
	}

	return true
}

// reverseHistoryRecords reverses the order of history records for ascending sort
func reverseHistoryRecords(records []HistoryRecord) {
	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}
}

// paginateHistoryRecords returns a paginated slice of history records
func paginateHistoryRecords(records []HistoryRecord, page, perPage int) []HistoryRecord {
	totalCount := len(records)
	start := (page - 1) * perPage
	end := start + perPage

	if start >= totalCount {
		return []HistoryRecord{}
	}

	if end > totalCount {
		end = totalCount
	}

	return records[start:end]
}
