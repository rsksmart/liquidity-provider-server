// Manual Approval UI JavaScript
// Uses vanilla JavaScript following the existing management.js pattern

// Import utility functions from configUtils
import { weiToEther } from './configUtils.js';

// State management
let pendingTransactions = [];
let historyRecords = [];
let selectedTxIds = new Set();
let currentPendingPage = 1;
let currentHistoryPage = 1;
const perPage = 20;
let pendingSortDesc = true;
let historySortDesc = true;

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    initializeEventListeners();
    loadPendingTransactions();
});

// Initialize all event listeners
function initializeEventListeners() {
    // Tab change listeners
    document.getElementById('pending-tab').addEventListener('click', () => {
        loadPendingTransactions();
    });
    document.getElementById('history-tab').addEventListener('click', () => {
        loadHistoryRecords();
    });

    // Pending table listeners
    document.getElementById('selectAllPending').addEventListener('change', handleSelectAll);
    document.getElementById('bulkApproveButton').addEventListener('click', () => handleBulkAction('approve'));
    document.getElementById('bulkDenyButton').addEventListener('click', () => handleBulkAction('deny'));
    document.getElementById('clearSelectionButton').addEventListener('click', clearSelection);
    document.getElementById('pendingSearchButton').addEventListener('click', () => {
        currentPendingPage = 1;
        loadPendingTransactions();
    });
    document.getElementById('pendingSearch').addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            currentPendingPage = 1;
            loadPendingTransactions();
        }
    });
    document.getElementById('pendingDateHeader').addEventListener('click', () => {
        pendingSortDesc = !pendingSortDesc;
        loadPendingTransactions();
    });

    // History table listeners
    document.getElementById('historyFilterButton').addEventListener('click', () => {
        currentHistoryPage = 1;
        loadHistoryRecords();
    });
    document.getElementById('historyDateHeader').addEventListener('click', () => {
        historySortDesc = !historySortDesc;
        loadHistoryRecords();
    });

    // Modal confirm listeners
    document.getElementById('confirmApproveButton').addEventListener('click', confirmApprove);
    document.getElementById('confirmDenyButton').addEventListener('click', confirmDeny);
}

// Fetch with CSRF token
async function fetchWithCsrf(url, method = 'GET', body = null) {
    const options = {
        method,
        headers: {
            'Content-Type': 'application/json',
            'X-CSRF-Token': csrfToken
        }
    };
    
    if (body) {
        options.body = JSON.stringify(body);
    }
    
    const response = await fetch(url, options);
    
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return response.json();
}

// Load pending transactions
async function loadPendingTransactions() {
    try {
        const searchTerm = document.getElementById('pendingSearch').value.trim();
        const params = new URLSearchParams({
            page: currentPendingPage,
            perPage: perPage,
            sortDesc: pendingSortDesc
        });
        
        if (searchTerm) {
            params.append('search', searchTerm);
        }
        
        const data = await fetchWithCsrf(`/management/manual-approval/pending?${params}`);
        pendingTransactions = data.transactions || [];
        
        renderPendingTable(pendingTransactions);
        renderPendingPagination(data.totalCount);
        updatePendingCount(data.totalCount);
        
        // Update sort indicator
        document.getElementById('pendingDateSort').textContent = pendingSortDesc ? '↓' : '↑';
    } catch (error) {
        console.error('Error loading pending transactions:', error);
        showToast('Failed to load pending transactions', false);
        renderPendingTable([]);
    }
}

// Render pending transactions table
function renderPendingTable(transactions) {
    const tbody = document.getElementById('pendingTableBody');
    
    if (!transactions || transactions.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="8" class="text-center text-muted py-4">
                    No pending transactions found
                </td>
            </tr>
        `;
        return;
    }
    
    tbody.innerHTML = transactions.map(tx => `
        <tr>
            <td>
                <input type="checkbox" class="form-check-input tx-checkbox" value="${tx.txId}" 
                    ${selectedTxIds.has(tx.txId) ? 'checked' : ''}>
            </td>
            <td>
                <span class="badge ${tx.type === 'pegin' ? 'bg-primary' : 'bg-info'}">${tx.type}</span>
            </td>
            <td>
                <small class="font-monospace">${truncateHash(tx.quoteHash)}</small>
            </td>
            <td>${formatDate(tx.date)}</td>
            <td>${weiToEther(tx.amount)} ${tx.type === 'pegin' ? 'BTC' : 'rBTC'}</td>
            <td>
                <span class="badge bg-secondary status-badge">${tx.state}</span>
            </td>
            <td>
                <small class="font-monospace">${truncateHash(tx.userAddress)}</small>
            </td>
            <td>
                <button class="btn btn-success btn-sm" onclick="window.approveOne('${tx.txId}')">
                    Approve
                </button>
                <button class="btn btn-danger btn-sm" onclick="window.denyOne('${tx.txId}')">
                    Deny
                </button>
            </td>
        </tr>
    `).join('');
    
    // Re-attach checkbox event listeners
    document.querySelectorAll('.tx-checkbox').forEach(checkbox => {
        checkbox.addEventListener('change', handleCheckboxChange);
    });
}

// Render pending pagination
function renderPendingPagination(totalCount) {
    const totalPages = Math.ceil(totalCount / perPage);
    const paginationContainer = document.getElementById('pendingPagination');
    
    const start = (currentPendingPage - 1) * perPage + 1;
    const end = Math.min(currentPendingPage * perPage, totalCount);
    document.getElementById('pendingPaginationInfo').textContent = 
        `Showing ${start}-${end} of ${totalCount}`;
    
    paginationContainer.innerHTML = renderPaginationButtons(currentPendingPage, totalPages, (page) => {
        currentPendingPage = page;
        loadPendingTransactions();
    });
}

// Load history records
async function loadHistoryRecords() {
    try {
        const searchTerm = document.getElementById('historySearch').value.trim();
        const statusFilter = document.getElementById('historyStatusFilter').value;
        const startDate = document.getElementById('historyStartDate').value;
        const endDate = document.getElementById('historyEndDate').value;
        
        const params = new URLSearchParams({
            page: currentHistoryPage,
            perPage: perPage,
            sortDesc: historySortDesc
        });
        
        if (searchTerm) params.append('search', searchTerm);
        if (statusFilter) params.append('status', statusFilter);
        if (startDate) params.append('startDate', startDate);
        if (endDate) params.append('endDate', endDate);
        
        const data = await fetchWithCsrf(`/management/manual-approval/history?${params}`);
        historyRecords = data.history || [];
        
        renderHistoryTable(historyRecords);
        renderHistoryPagination(data.totalCount);
        
        // Update sort indicator
        document.getElementById('historyDateSort').textContent = historySortDesc ? '↓' : '↑';
    } catch (error) {
        console.error('Error loading history:', error);
        showToast('Failed to load history', false);
        renderHistoryTable([]);
    }
}

// Render history table
function renderHistoryTable(records) {
    const tbody = document.getElementById('historyTableBody');
    
    if (!records || records.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="7" class="text-center text-muted py-4">
                    No history records found
                </td>
            </tr>
        `;
        return;
    }
    
    tbody.innerHTML = records.map(record => `
        <tr>
            <td>
                <span class="badge ${record.type === 'pegin' ? 'bg-primary' : 'bg-info'}">${record.type}</span>
            </td>
            <td>
                <small class="font-monospace">${truncateHash(record.quoteHash)}</small>
            </td>
            <td>${formatDate(record.date)}</td>
            <td>${weiToEther(record.amount)} ${record.type === 'pegin' ? 'BTC' : 'rBTC'}</td>
            <td>
                <span class="badge ${record.decision === 'approved' ? 'bg-success' : 'bg-danger'} status-badge">
                    ${record.decision}
                </span>
            </td>
            <td>${record.approvedOrDeniedBy}</td>
            <td>${formatDate(record.decisionDate)}</td>
        </tr>
    `).join('');
}

// Render history pagination
function renderHistoryPagination(totalCount) {
    const totalPages = Math.ceil(totalCount / perPage);
    const paginationContainer = document.getElementById('historyPagination');
    
    const start = (currentHistoryPage - 1) * perPage + 1;
    const end = Math.min(currentHistoryPage * perPage, totalCount);
    document.getElementById('historyPaginationInfo').textContent = 
        `Showing ${start}-${end} of ${totalCount}`;
    
    paginationContainer.innerHTML = renderPaginationButtons(currentHistoryPage, totalPages, (page) => {
        currentHistoryPage = page;
        loadHistoryRecords();
    });
}

// Generic pagination button renderer
function renderPaginationButtons(currentPage, totalPages, onPageChange) {
    if (totalPages <= 1) return '';
    
    let html = '';
    
    // Previous button
    html += `
        <li class="page-item ${currentPage === 1 ? 'disabled' : ''}">
            <a class="page-link" href="#" data-page="${currentPage - 1}">Previous</a>
        </li>
    `;
    
    // Page numbers
    const maxButtons = 5;
    let startPage = Math.max(1, currentPage - Math.floor(maxButtons / 2));
    let endPage = Math.min(totalPages, startPage + maxButtons - 1);
    
    if (endPage - startPage < maxButtons - 1) {
        startPage = Math.max(1, endPage - maxButtons + 1);
    }
    
    if (startPage > 1) {
        html += `<li class="page-item"><a class="page-link" href="#" data-page="1">1</a></li>`;
        if (startPage > 2) {
            html += `<li class="page-item disabled"><span class="page-link">...</span></li>`;
        }
    }
    
    for (let i = startPage; i <= endPage; i++) {
        html += `
            <li class="page-item ${i === currentPage ? 'active' : ''}">
                <a class="page-link" href="#" data-page="${i}">${i}</a>
            </li>
        `;
    }
    
    if (endPage < totalPages) {
        if (endPage < totalPages - 1) {
            html += `<li class="page-item disabled"><span class="page-link">...</span></li>`;
        }
        html += `<li class="page-item"><a class="page-link" href="#" data-page="${totalPages}">${totalPages}</a></li>`;
    }
    
    // Next button
    html += `
        <li class="page-item ${currentPage === totalPages ? 'disabled' : ''}">
            <a class="page-link" href="#" data-page="${currentPage + 1}">Next</a>
        </li>
    `;
    
    // Attach click handlers
    setTimeout(() => {
        document.querySelectorAll('#pendingPagination .page-link, #historyPagination .page-link').forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                const page = parseInt(e.target.getAttribute('data-page'));
                if (!isNaN(page)) {
                    onPageChange(page);
                }
            });
        });
    }, 0);
    
    return html;
}

// Handle select all checkbox
function handleSelectAll(e) {
    const isChecked = e.target.checked;
    selectedTxIds.clear();
    
    if (isChecked) {
        pendingTransactions.forEach(tx => selectedTxIds.add(tx.txId));
    }
    
    document.querySelectorAll('.tx-checkbox').forEach(checkbox => {
        checkbox.checked = isChecked;
    });
    
    updateBulkActionsToolbar();
}

// Handle individual checkbox change
function handleCheckboxChange(e) {
    const txId = e.target.value;
    
    if (e.target.checked) {
        selectedTxIds.add(txId);
    } else {
        selectedTxIds.delete(txId);
    }
    
    // Update select all checkbox
    const allSelected = pendingTransactions.length > 0 && 
        pendingTransactions.every(tx => selectedTxIds.has(tx.txId));
    document.getElementById('selectAllPending').checked = allSelected;
    
    updateBulkActionsToolbar();
}

// Update bulk actions toolbar visibility
function updateBulkActionsToolbar() {
    const toolbar = document.getElementById('bulkActionsToolbar');
    const count = selectedTxIds.size;
    
    if (count > 0) {
        toolbar.classList.add('show');
        document.getElementById('selectedCount').textContent = count;
    } else {
        toolbar.classList.remove('show');
    }
}

// Clear selection
function clearSelection() {
    selectedTxIds.clear();
    document.getElementById('selectAllPending').checked = false;
    document.querySelectorAll('.tx-checkbox').forEach(checkbox => {
        checkbox.checked = false;
    });
    updateBulkActionsToolbar();
}

// Handle bulk action (approve/deny)
function handleBulkAction(action) {
    if (selectedTxIds.size === 0) return;
    
    const txIds = Array.from(selectedTxIds);
    showActionModal(action, txIds);
}

// Approve single transaction
window.approveOne = (txId) => {
    showActionModal('approve', [txId]);
};

// Deny single transaction
window.denyOne = (txId) => {
    showActionModal('deny', [txId]);
};

// Show approve/deny modal
function showActionModal(action, txIds) {
    const modalId = action === 'approve' ? 'approveModal' : 'denyModal';
    const countEl = document.getElementById(`${action}Count`);
    const listEl = document.getElementById(`${action}TxList`);
    
    countEl.textContent = txIds.length;
    listEl.innerHTML = txIds.map(txId => `<li><small class="font-monospace">${txId}</small></li>`).join('');
    
    // Store txIds for confirmation
    window.currentActionTxIds = txIds;
    window.currentAction = action;
    
    const modal = new bootstrap.Modal(document.getElementById(modalId));
    modal.show();
}

// Confirm approve
async function confirmApprove() {
    await executeAction('approve', window.currentActionTxIds);
}

// Confirm deny
async function confirmDeny() {
    await executeAction('deny', window.currentActionTxIds);
}

// Execute approve/deny action
async function executeAction(action, txIds) {
    try {
        // Disable buttons
        const confirmButton = document.getElementById(`confirm${action === 'approve' ? 'Approve' : 'Deny'}Button`);
        confirmButton.disabled = true;
        confirmButton.innerHTML = '<span class="spinner-border spinner-border-sm" role="status"></span> Processing...';
        
        const endpoint = `/management/manual-approval/${action}`;
        await fetchWithCsrf(endpoint, 'POST', { txIds });
        
        // Close modal
        const modalId = action === 'approve' ? 'approveModal' : 'denyModal';
        const modal = bootstrap.Modal.getInstance(document.getElementById(modalId));
        modal.hide();
        
        // Show success toast
        showToast(`Successfully ${action}ed ${txIds.length} transaction(s)`, true);
        
        // Clear selection and reload
        clearSelection();
        loadPendingTransactions();
        
        // Reset button
        confirmButton.disabled = false;
        confirmButton.innerHTML = `<i class="bi bi-${action === 'approve' ? 'check' : 'x'}-circle"></i> Confirm ${action === 'approve' ? 'Approval' : 'Denial'}`;
    } catch (error) {
        console.error(`Error ${action}ing transactions:`, error);
        showToast(`Failed to ${action} transactions`, false);
        
        // Reset button
        const confirmButton = document.getElementById(`confirm${action === 'approve' ? 'Approve' : 'Deny'}Button`);
        confirmButton.disabled = false;
        confirmButton.innerHTML = `<i class="bi bi-${action === 'approve' ? 'check' : 'x'}-circle"></i> Confirm ${action === 'approve' ? 'Approval' : 'Denial'}`;
    }
}

// Update pending count badge
function updatePendingCount(count) {
    document.getElementById('pendingCount').textContent = count;
}

// Show toast notification
function showToast(message, isSuccess) {
    const toastId = isSuccess ? 'successToast' : 'errorToast';
    const toastEl = document.getElementById(toastId);
    const toastBody = toastEl.querySelector('.toast-body');
    
    toastBody.textContent = message;
    
    const toast = new bootstrap.Toast(toastEl);
    toast.show();
}

// Utility: Truncate hash
function truncateHash(hash) {
    if (!hash) return '';
    if (hash.length <= 12) return hash;
    return `${hash.substring(0, 6)}...${hash.substring(hash.length - 4)}`;
}

// Utility: Format date
function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}
