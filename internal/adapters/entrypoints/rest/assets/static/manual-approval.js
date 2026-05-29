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

function createEmptyRow(colspan, message) {
    const tr = document.createElement('tr');
    const td = document.createElement('td');
    td.colSpan = colspan;
    td.className = 'text-center text-muted py-4';
    td.textContent = message;
    tr.appendChild(td);
    return tr;
}

function createBadgeCell(text, badgeClass) {
    const td = document.createElement('td');
    const span = document.createElement('span');
    span.className = `badge ${badgeClass}`;
    span.textContent = text;
    td.appendChild(span);
    return td;
}

function createMonospaceCell(text) {
    const td = document.createElement('td');
    const small = document.createElement('small');
    small.className = 'font-monospace';
    small.textContent = text;
    td.appendChild(small);
    return td;
}

function createTextCell(text) {
    const td = document.createElement('td');
    td.textContent = text;
    return td;
}

function renderPendingTable(transactions) {
    const tbody = document.getElementById('pendingTableBody');
    tbody.replaceChildren();
    
    if (!transactions || transactions.length === 0) {
        tbody.appendChild(createEmptyRow(8, 'No pending transactions found'));
        return;
    }
    
    for (const tx of transactions) {
        const tr = document.createElement('tr');

        const checkboxTd = document.createElement('td');
        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.className = 'form-check-input tx-checkbox';
        checkbox.value = tx.txId;
        checkbox.checked = selectedTxIds.has(tx.txId);
        checkbox.addEventListener('change', handleCheckboxChange);
        checkboxTd.appendChild(checkbox);
        tr.appendChild(checkboxTd);

        tr.appendChild(createBadgeCell(tx.type, tx.type === 'pegin' ? 'bg-primary' : 'bg-info'));
        tr.appendChild(createMonospaceCell(truncateHash(tx.quoteHash)));
        tr.appendChild(createTextCell(formatDate(tx.date)));
        tr.appendChild(createTextCell(`${weiToEther(tx.amount)} ${tx.type === 'pegin' ? 'BTC' : 'rBTC'}`));
        tr.appendChild(createBadgeCell(tx.state, 'bg-secondary status-badge'));
        tr.appendChild(createMonospaceCell(truncateHash(tx.userAddress)));

        const actionsTd = document.createElement('td');
        const approveBtn = document.createElement('button');
        approveBtn.className = 'btn btn-success btn-sm';
        approveBtn.textContent = 'Approve';
        approveBtn.addEventListener('click', () => showActionModal('approve', [tx.txId]));
        const denyBtn = document.createElement('button');
        denyBtn.className = 'btn btn-danger btn-sm';
        denyBtn.textContent = 'Deny';
        denyBtn.addEventListener('click', () => showActionModal('deny', [tx.txId]));
        actionsTd.appendChild(approveBtn);
        actionsTd.appendChild(denyBtn);
        tr.appendChild(actionsTd);

        tbody.appendChild(tr);
    }
}

function renderPendingPagination(totalCount) {
    const totalPages = Math.ceil(totalCount / perPage);
    const paginationContainer = document.getElementById('pendingPagination');
    
    const start = (currentPendingPage - 1) * perPage + 1;
    const end = Math.min(currentPendingPage * perPage, totalCount);
    document.getElementById('pendingPaginationInfo').textContent = 
        `Showing ${start}-${end} of ${totalCount}`;
    
    paginationContainer.replaceChildren(renderPaginationButtons(currentPendingPage, totalPages, (page) => {
        currentPendingPage = page;
        loadPendingTransactions();
    }));
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

function renderHistoryTable(records) {
    const tbody = document.getElementById('historyTableBody');
    tbody.replaceChildren();
    
    if (!records || records.length === 0) {
        tbody.appendChild(createEmptyRow(7, 'No history records found'));
        return;
    }
    
    for (const record of records) {
        const tr = document.createElement('tr');
        tr.appendChild(createBadgeCell(record.type, record.type === 'pegin' ? 'bg-primary' : 'bg-info'));
        tr.appendChild(createMonospaceCell(truncateHash(record.quoteHash)));
        tr.appendChild(createTextCell(formatDate(record.date)));
        tr.appendChild(createTextCell(`${weiToEther(record.amount)} ${record.type === 'pegin' ? 'BTC' : 'rBTC'}`));
        tr.appendChild(createBadgeCell(record.decision, `${record.decision === 'approved' ? 'bg-success' : 'bg-danger'} status-badge`));
        tr.appendChild(createTextCell(record.approvedOrDeniedBy));
        tr.appendChild(createTextCell(formatDate(record.decisionDate)));
        tbody.appendChild(tr);
    }
}

function renderHistoryPagination(totalCount) {
    const totalPages = Math.ceil(totalCount / perPage);
    const paginationContainer = document.getElementById('historyPagination');
    
    const start = (currentHistoryPage - 1) * perPage + 1;
    const end = Math.min(currentHistoryPage * perPage, totalCount);
    document.getElementById('historyPaginationInfo').textContent = 
        `Showing ${start}-${end} of ${totalCount}`;
    
    paginationContainer.replaceChildren(renderPaginationButtons(currentHistoryPage, totalPages, (page) => {
        currentHistoryPage = page;
        loadHistoryRecords();
    }));
}

function renderPaginationButtons(currentPage, totalPages, onPageChange) {
    const fragment = document.createDocumentFragment();
    if (totalPages <= 1) return fragment;

    function createPageItem(label, page, options = {}) {
        const li = document.createElement('li');
        li.className = 'page-item';
        if (options.disabled) li.classList.add('disabled');
        if (options.active) li.classList.add('active');

        const link = document.createElement(options.disabled && !page ? 'span' : 'a');
        link.className = 'page-link';
        link.textContent = label;
        if (link.tagName === 'A') {
            link.href = '#';
            link.addEventListener('click', (e) => {
                e.preventDefault();
                if (!options.disabled && page != null) onPageChange(page);
            });
        }
        li.appendChild(link);
        return li;
    }

    fragment.appendChild(createPageItem('Previous', currentPage - 1, { disabled: currentPage === 1 }));

    const maxButtons = 5;
    let startPage = Math.max(1, currentPage - Math.floor(maxButtons / 2));
    let endPage = Math.min(totalPages, startPage + maxButtons - 1);
    if (endPage - startPage < maxButtons - 1) {
        startPage = Math.max(1, endPage - maxButtons + 1);
    }

    if (startPage > 1) {
        fragment.appendChild(createPageItem('1', 1));
        if (startPage > 2) {
            fragment.appendChild(createPageItem('...', null, { disabled: true }));
        }
    }

    for (let i = startPage; i <= endPage; i++) {
        fragment.appendChild(createPageItem(String(i), i, { active: i === currentPage }));
    }

    if (endPage < totalPages) {
        if (endPage < totalPages - 1) {
            fragment.appendChild(createPageItem('...', null, { disabled: true }));
        }
        fragment.appendChild(createPageItem(String(totalPages), totalPages));
    }

    fragment.appendChild(createPageItem('Next', currentPage + 1, { disabled: currentPage === totalPages }));

    return fragment;
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

function showActionModal(action, txIds) {
    const modalId = action === 'approve' ? 'approveModal' : 'denyModal';
    const countEl = document.getElementById(`${action}Count`);
    const listEl = document.getElementById(`${action}TxList`);
    
    countEl.textContent = txIds.length;
    listEl.replaceChildren();
    for (const txId of txIds) {
        const li = document.createElement('li');
        const small = document.createElement('small');
        small.className = 'font-monospace';
        small.textContent = txId;
        li.appendChild(small);
        listEl.appendChild(li);
    }
    
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

function setConfirmButtonLoading(button) {
    button.disabled = true;
    button.replaceChildren();
    const spinner = document.createElement('span');
    spinner.className = 'spinner-border spinner-border-sm';
    spinner.setAttribute('role', 'status');
    button.appendChild(spinner);
    button.appendChild(document.createTextNode(' Processing...'));
}

function resetConfirmButton(button, action) {
    button.disabled = false;
    button.replaceChildren();
    const icon = document.createElement('i');
    icon.className = `bi bi-${action === 'approve' ? 'check' : 'x'}-circle`;
    button.appendChild(icon);
    button.appendChild(document.createTextNode(` Confirm ${action === 'approve' ? 'Approval' : 'Denial'}`));
}

async function executeAction(action, txIds) {
    const confirmButton = document.getElementById(`confirm${action === 'approve' ? 'Approve' : 'Deny'}Button`);
    try {
        setConfirmButtonLoading(confirmButton);
        
        const endpoint = `/management/manual-approval/${action}`;
        await fetchWithCsrf(endpoint, 'POST', { txIds });
        
        const modalId = action === 'approve' ? 'approveModal' : 'denyModal';
        const modal = bootstrap.Modal.getInstance(document.getElementById(modalId));
        modal.hide();
        
        showToast(`Successfully ${action}ed ${txIds.length} transaction(s)`, true);
        
        clearSelection();
        loadPendingTransactions();
        
        resetConfirmButton(confirmButton, action);
    } catch (error) {
        console.error(`Error ${action}ing transactions:`, error);
        showToast(`Failed to ${action} transactions`, false);
        resetConfirmButton(confirmButton, action);
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
