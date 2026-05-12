import {
    weiToEther,
    etherToWei,
    isFeeKey,
    validateConfig,
    formatGeneralConfig,
    postConfig,
    hasDuplicateConfirmationAmounts,
    isfeePercentageKey,
    isToggableFeeKey,
    formatCap
} from './configUtils.js';

const generalChanged = { value: false };
const peginChanged = { value: false };
const pegoutChanged = { value: false };

const setTextContent = (id, text) => document.getElementById(id).textContent = text;

const fetchData = async (url, elementId, csrfToken) => {
    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-CSRF-Token': csrfToken
            }
        });
        const responseData = await response.json();
        setTextContent(elementId, weiToEther(responseData.collateral) + " rBTC");
    } catch (error) {
        console.error(error);
    }
};

const populateProviderData = (providerData, rskAddress, btcAddress) => {
    setTextContent('providerRskAddress', rskAddress);
    setTextContent('providerBtcAddress', btcAddress);
    setTextContent('isOperational', providerData.status ? "Operational" : "Not Operational");
};

const createInput = (section, key, value) => {
    const div = document.createElement('div');
    div.classList.add('mb-3');

    const label = document.createElement('label');
    label.classList.add('form-label');
    label.textContent = key;

    const inputContainer = document.createElement('div');
    inputContainer.classList.add('input-container');

    if (typeof value === 'boolean') {
        createCheckboxInput(inputContainer, section, key, value);
    } else if (isToggableFeeKey(key)) {
        createToggableFeeInput(inputContainer, label, section, key, value);
    } else if (isFeeKey(key)) {
        createFeeInput(inputContainer, label, section, key, value);
    } else if (isfeePercentageKey(key)) {
        createFeePercentageInput(inputContainer, section, key, value);
    } else {
        createDefaultInput(inputContainer, section, key, value);
    }

    div.appendChild(label);
    div.appendChild(inputContainer);
    section.appendChild(div);
};

const createCheckboxInput = (inputContainer, section, key, value) => {
    const checkbox = document.createElement('input');
    checkbox.type = 'checkbox';
    checkbox.classList.add('form-check-input');
    checkbox.style.marginRight = '10px';
    checkbox.dataset.key = key;
    checkbox.setAttribute('data-testid', `config-${section.id.replace('Config','')}-${key}-checkbox`);
    checkbox.checked = value;
    checkbox.addEventListener('change', () => setChanged(section.id));
    inputContainer.appendChild(checkbox);
};

const createToggableFeeInput = (inputContainer, label, section, key, value) => {
    const checkbox = document.createElement('input');
    checkbox.type = 'checkbox';
    checkbox.classList.add('form-check-input');
    checkbox.style.marginRight = '10px';
    checkbox.dataset.key = `${key}_enabled`;
    checkbox.setAttribute('data-testid', `config-${section.id.replace('Config','')}-${key}-checkbox`);

    const input = document.createElement('input');
    input.type = 'text';
    input.style.width = '40%';
    input.classList.add('form-control');
    input.dataset.key = key;
    input.setAttribute('data-testid', `config-${section.id.replace('Config','')}-${key}-input`);
    input.dataset.originalValue = value;

    if (value === '0' || value === 0) {
        checkbox.checked = false;
        input.value = '0';
        input.disabled = true;
    } else {
        checkbox.checked = true;
        input.value = isFeeKey(key) ? weiToEther(value) : value;
        input.disabled = false;
    }

    checkbox.addEventListener('change', () => {
        if (checkbox.checked) {
            input.disabled = false;
            input.value = (input.dataset.originalValue === '0' || input.dataset.originalValue === 0) ? '' :
                            isFeeKey(key) ? weiToEther(input.dataset.originalValue) : input.dataset.originalValue;
        } else {
            input.disabled = true;
            input.value = '0';
        }
        setChanged(section.id);
        checkFeeWarnings();
    });

    input.addEventListener('input', () => setChanged(section.id));
    inputContainer.appendChild(checkbox);
    inputContainer.appendChild(input);
    const questionIcon = createQuestionIcon(getTooltipText(key));
    label.appendChild(questionIcon);
};

const createFeeInput = (inputContainer, label, section, key, value) => {
    const input = document.createElement('input');
    input.type = 'text';
    input.style.width = '40%';
    input.classList.add('form-control');
    input.dataset.key = key;
    input.setAttribute('data-testid', `config-${section.id.replace('Config','')}-${key}-input`);
    input.value = isFeeKey(key) ? weiToEther(value) : value;
    input.addEventListener('input', () => setChanged(section.id));
    inputContainer.appendChild(input);
    const questionIcon = createQuestionIcon(getTooltipText(key));
    label.appendChild(questionIcon);
};

const createFeePercentageInput = (inputContainer, section, key, value) => {
    const input = document.createElement('input');
    input.type = 'text';
    input.style.width = '40%';
    input.classList.add('form-control');
    input.dataset.key = key;
    input.setAttribute('data-testid', `config-${section.id.replace('Config','')}-${key}-input`);
    input.value = typeof value === 'number' ? value.toString() : value;
    input.addEventListener('input', () => setChanged(section.id));
    inputContainer.appendChild(input);
};

const createDefaultInput = (inputContainer, section, key, value) => {
    const input = document.createElement('input');
    input.type = 'text';
    input.style.width = '40%';
    input.classList.add('form-control');
    input.dataset.key = key;
    input.setAttribute('data-testid', `config-${section.id.replace('Config','')}-${key}-input`);
    input.value = value;
    input.addEventListener('input', () => setChanged(section.id));
    inputContainer.appendChild(input);
};

const createQuestionIcon = (tooltipText) => {
    const questionIcon = document.createElement('span');
    questionIcon.classList.add('question-mark');
    const img = document.createElement('img');
    img.src = '../static/questionIcon.svg';
    img.width = 13;
    img.height = 13;
    img.alt = 'Question Mark';
    img.classList.add('bi', 'bi-question-circle');
    questionIcon.appendChild(img);
    const tooltip = document.createElement('div');
    tooltip.classList.add('custom-tooltip');
    tooltip.textContent = tooltipText;
    questionIcon.appendChild(tooltip);
    return questionIcon;
};

const getTooltipText = (key) => {
    const tooltips = {
        timeForDeposit: 'The time (in seconds) for which a deposit is considered valid.',
        expireTime: 'The time (in seconds) after which a quote is considered expired.',
        penaltyFee: 'The penalty fee (in BTC) charged for invalid transactions.',
        callFee: 'The fee (in BTC) charged by the LP for processing a transaction.',
        maxValue: 'The maximum value (in BTC) allowed for a transaction.',
        minValue: 'The minimum value (in BTC) allowed for a transaction.',
        expireBlocks: 'The number of blocks after which a quote is considered expired.',
        bridgeTransactionMin: 'The amount of rBTC that needs to be gathered in peg out refunds before executing a native peg out.',
        fixedFee: 'A fixed fee charged for transactions.',
        feePercentage: 'A percentage fee charged based on the transaction amount.'
    };
    return tooltips[key] || 'No description available';
};

const createConfirmationConfig = (section, configKey, confirmations) => {
    const container = document.createElement('div');
    container.classList.add('confirmation-config');
    container.dataset.configKey = configKey;

    const header = document.createElement('h5');
    header.textContent = configKey;
    container.appendChild(header);

    const entriesContainer = document.createElement('div');
    entriesContainer.classList.add('entries-container');

    const sortedConfirmations = Object.entries(confirmations).sort(([amountWeiA], [amountWeiB]) => {
        const valA = parseFloat(weiToEther(amountWeiA));
        const valB = parseFloat(weiToEther(amountWeiB));
        return valA - valB;
    });

    sortedConfirmations.forEach(([amountWei, confirmation], index) => {
        createConfirmationEntry(entriesContainer, configKey, index, amountWei, confirmation);
    });

    container.appendChild(entriesContainer);

    const addButton = document.createElement('button');
    addButton.type = 'button';
    addButton.classList.add('btn', 'btn-secondary', 'mt-2');
    addButton.textContent = 'Add Entry';
    addButton.addEventListener('click', () => {
        const index = entriesContainer.querySelectorAll('.input-group').length;
        createConfirmationEntry(entriesContainer, configKey, index);
        setChanged(section.id);
    });

    container.appendChild(addButton);
    section.appendChild(container);
};

const setChanged = (sectionId) => {
    if (sectionId === 'generalConfig') generalChanged.value = true;
    else if (sectionId === 'peginConfig') peginChanged.value = true;
    else if (sectionId === 'pegoutConfig') pegoutChanged.value = true;
    else if (sectionId === 'rskConfirmations' || sectionId === 'btcConfirmations') generalChanged.value = true;
};

const createConfirmationEntry = (container, configKey, index, amount = '', confirmation = '') => {
    const div = document.createElement('div');
    div.classList.add('d-flex', 'align-items-center', 'mb-2');
    const fieldWidth = '180px';

    const amountGroup = document.createElement('div');
    amountGroup.classList.add('input-group', 'me-2');
    const amountInput = document.createElement('input');
    amountInput.type = 'text';
    amountInput.value = amount ? weiToEther(amount) : '';
    amountInput.classList.add('form-control', 'form-control-sm');
    amountInput.placeholder = 'Amount';
    amountInput.dataset.configKey = configKey;
    amountInput.dataset.field = 'amount';
    amountInput.dataset.index = index;
    amountInput.style.maxWidth = fieldWidth;

    const amountInputAppend = document.createElement('span');
    amountInputAppend.classList.add('input-group-text', 'input-group-text-sm');
    amountInputAppend.textContent = configKey === 'btcConfirmations' ? 'BTC' : 'rBTC';
    amountGroup.appendChild(amountInput);
    amountGroup.appendChild(amountInputAppend);

    const confirmationGroup = document.createElement('div');
    confirmationGroup.classList.add('input-group', 'me-2');

    const confirmationInput = document.createElement('input');
    confirmationInput.type = 'number';
    confirmationInput.value = confirmation;
    confirmationInput.classList.add('form-control', 'form-control-sm');
    confirmationInput.placeholder = 'Confirmations';
    confirmationInput.dataset.configKey = configKey;
    confirmationInput.dataset.field = 'confirmation';
    confirmationInput.dataset.index = index;
    confirmationInput.style.maxWidth = fieldWidth;
    confirmationInput.setAttribute('data-testid', `config-${configKey}-${index}`);

    const confirmationInputAppend = document.createElement('span');
    confirmationInputAppend.classList.add('input-group-text', 'input-group-text-sm');
    confirmationInputAppend.textContent = 'confirmations';

    confirmationGroup.appendChild(confirmationInput);
    confirmationGroup.appendChild(confirmationInputAppend);

    const removeButton = document.createElement('button');
    removeButton.type = 'button';
    removeButton.classList.add('btn', 'btn-danger', 'btn-sm');
    removeButton.textContent = 'Remove';

    removeButton.addEventListener('click', () => {
        div.remove();
        setChanged(configKey);
    });

    amountInput.addEventListener('input', () => setChanged(configKey));
    confirmationInput.addEventListener('input', () => setChanged(configKey));

    div.appendChild(amountGroup);
    div.appendChild(confirmationGroup);
    div.appendChild(removeButton);
    container.appendChild(div);
};

const populateConfigSection = (sectionId, config) => {
    const section = document.getElementById(sectionId);
    section.innerHTML = '';
    Object.entries(config).forEach(([key, value]) => {
        if (key === 'rskConfirmations' || key === 'btcConfirmations') {
            createConfirmationConfig(section, key, value);
        } else {
            createInput(section, key, value);
        }
    });
};

const showSuccessToast = () => {
    const toastElement = document.getElementById('successToast');
    const toast = new bootstrap.Toast(toastElement);
    toast.show();
};

const showErrorToast = (errorMessage) => {
    const toastElement = document.createElement('div');
    toastElement.classList.add('toast');
    toastElement.setAttribute('role', 'alert');
    toastElement.setAttribute('aria-live', 'assertive');
    toastElement.setAttribute('aria-atomic', 'true');
    toastElement.innerHTML = `
        <div class="toast-header">
            <strong class="me-auto">Error</strong>
            <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
        </div>
        <div class="toast-body"></div>
    `;
    toastElement.querySelector('.toast-body').textContent = errorMessage;
    document.querySelector('.toast-container').appendChild(toastElement);
    const toast = new bootstrap.Toast(toastElement);
    toast.show();
};

const showWarningToast = (warningMessage) => {
    const existingToast = document.getElementById('warningToast');
    if (existingToast) existingToast.parentNode.removeChild(existingToast);

    const toastElement = document.createElement('div');
    toastElement.id = 'warningToast';
    toastElement.classList.add('toast', 'text-bg-warning');
    toastElement.setAttribute('role', 'alert');
    toastElement.setAttribute('aria-live', 'assertive');
    toastElement.setAttribute('aria-atomic', 'true');
    toastElement.innerHTML = `
        <div class="toast-header">
            <strong class="me-auto">Warning</strong>
            <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
        </div>
        <div class="toast-body">
            ${warningMessage}
        </div>
    `;
    document.querySelector('.toast-container').appendChild(toastElement);
    const toast = new bootstrap.Toast(toastElement);
    toast.show();
};

function checkFeeWarnings() {
    const activeTabPane = document.querySelector('#configTabContent .tab-pane.active');
    let activeSectionId;
    if (activeTabPane) {
        switch (activeTabPane.id) {
            case 'general':
                activeSectionId = 'generalConfig';
                break;
            case 'peginConfig':
                activeSectionId = 'peginConfig';
                break;
            case 'pegoutConfig':
                activeSectionId = 'pegoutConfig';
                break;
            default:
                activeSectionId = undefined;
        }
    }
    if (!activeSectionId) return;
    const sectionElement = document.getElementById(activeSectionId);
    if (!sectionElement) return;
    const fixedFeeCheckbox = sectionElement.querySelector('input[data-key="fixedFee_enabled"]');
    const feePercentageCheckbox = sectionElement.querySelector('input[data-key="feePercentage_enabled"]');
    const shouldWarn = (
        fixedFeeCheckbox &&
        feePercentageCheckbox &&
        !fixedFeeCheckbox.checked &&
        !feePercentageCheckbox.checked
    );
    const existingToast = document.getElementById('warningToast');
    if (shouldWarn) {
        if (!existingToast) {
            showWarningToast('You have configured a zero-fee setting. This means you won\'t earn fees from bridging transactions.');
        } else {
            bootstrap.Toast.getOrCreateInstance(existingToast).show();
        }
    } else if (existingToast) {
        existingToast.parentNode.removeChild(existingToast);
    }
}

function getConfirmationConfig(sectionId) {
    const entries = document.querySelectorAll(`#${sectionId} .confirmation-config`);
    const config = {};
    entries.forEach(entry => {
        const configKey = entry.dataset.configKey;
        const inputGroups = entry.querySelectorAll('input');
        let tempArray = [];

        inputGroups.forEach(input => {
            const idx = input.dataset.index;
            if (!tempArray[idx]) tempArray[idx] = {};

            if (input.value.trim() === '') {
                const amountLabel = configKey === 'btcConfirmations' ? 'BTC amount' : 'rBTC amount';
                showErrorToast(`Please enter a non-empty value for "${input.dataset.field === 'amount' ? amountLabel : 'confirmations'}."`);
                throw new Error(`Empty ${input.dataset.field} input`);
            }

            if (input.dataset.field === 'amount') {
                try {
                    tempArray[idx].amount = etherToWei(input.value).toString();
                } catch (error) {
                    const amountLabel = configKey === 'btcConfirmations' ? 'BTC amount' : 'rBTC amount';
                    showErrorToast(`Invalid input "${input.value}" for ${amountLabel}. Please enter a valid non-negative number.`);
                    throw error;
                }
            } else if (input.dataset.field === 'confirmation') {
                const val = Number(input.value);
                if (isNaN(val) || !Number.isInteger(val) || val < 0) {
                    showErrorToast(`Invalid input "${input.value}" for confirmations. Please enter a valid non-negative integer.`);
                    throw new Error('Invalid confirmation number');
                }
                tempArray[idx].confirmation = val;
            }
        });

        tempArray = tempArray.filter( entryObj =>
            entryObj !== undefined &&
            entryObj.amount !== undefined &&
            entryObj.confirmation !== undefined
        );
        config[configKey] = tempArray;
    });
    return config;
}

function getRegularConfig(sectionId) {
    const inputs = document.querySelectorAll(`#${sectionId} input:not(.form-check-input):not([data-field="amount"]):not([data-field="confirmation"])`);
    const checkboxes = document.querySelectorAll(`#${sectionId} input.form-check-input`);
    const config = {};

    inputs.forEach(input => {
        const key = input.dataset.key;
        let value;

        if (input.disabled) {
            if (isfeePercentageKey(key)) {
                value = 0;
            } else {
                value = '0';
            }
        } else {
            if (isFeeKey(key)) {
                try {
                    value = etherToWei(input.value).toString();
                } catch (error) {
                    showErrorToast(`"${sectionId}": Invalid input "${input.value}" for field "${key}". Please enter a valid number.`);
                    throw error;
                }
            } else if (isfeePercentageKey(key)) {
                const rawInput = input.value.trim();
                const percentagePattern = /^\d+(\.\d+)?%?$/;
                if (!percentagePattern.test(rawInput)) {
                    showErrorToast(`"${sectionId}": Invalid percentage entered "${rawInput}". Please provide a numeric value between 0% and 100%.`);
                    throw new Error('Invalid feePercentage');
                }
                const numericPart = rawInput.endsWith('%') ? rawInput.slice(0, -1) : rawInput;
                value = parseFloat(numericPart);
                if (isNaN(value)) {
                    showErrorToast(`"${sectionId}": Invalid percentage entered "${rawInput}". Please provide a valid value between 0% and 100%.`);
                    throw new Error('Invalid feePercentage');
                }
                if (value < 0) {
                    showErrorToast(`"${sectionId}": Fee percentage cannot be negative. Please enter a value between 0% and 100%.`);
                    throw new Error('Invalid feePercentage');
                }
                if (value > 100) {
                    showErrorToast(`"${sectionId}": Fee percentage cannot exceed 100%. Please enter a value between 0% and 100%.`);
                    throw new Error('Invalid feePercentage');
                }
            } else {
                value = input.value;
                if (!isNaN(value) && value !== '') value = Number(value);
            }
        }
        config[key] = value;
    });

    checkboxes.forEach(input => {
        const key = input.dataset.key;
        if (!key.endsWith('_enabled')) config[key] = input.checked;
    });

    return config;
}

function getConfig(sectionId) {
    const confirmationInputs = document.querySelectorAll(`#${sectionId} .confirmation-config`);
    let config = {};

    if (confirmationInputs.length > 0) {
        const confirmationConfig = getConfirmationConfig(sectionId);
        config = { ...config, ...confirmationConfig };
    }

    const regularConfig = getRegularConfig(sectionId);
    config = { ...config, ...regularConfig };
    return config;
}

const saveConfig = async (csrfToken, configurations) => {
    let saveSuccess = true;

    let generalConfig, peginConfigData, pegoutConfigData;
    try {
        generalConfig = getConfig('generalConfig');
        peginConfigData = getConfig('peginConfig');
        pegoutConfigData = getConfig('pegoutConfig');
    } catch (error) {
        return;
    }

    ['rskConfirmations', 'btcConfirmations'].forEach(key => {
        if (!generalConfig[key]?.length) {
            showErrorToast(`Please provide at least one fully filled out entry for ${key}.`);
            throw new Error('Missing confirmations');
        }
    });

    for (const key of ['rskConfirmations', 'btcConfirmations']) {
        if (generalConfig[key] && hasDuplicateConfirmationAmounts(generalConfig[key])) {
            showErrorToast(`Duplicate rBTC amounts found in ${key}. Please remove duplicates before saving.`);
            return;
        }
    }

    const { isValid: isGeneralValid, errors: generalErrors } = validateConfig(formatGeneralConfig(generalConfig), configurations.general);
    if (!isGeneralValid) {
        showErrorToast(generalErrors.join('<br>'));
        saveSuccess = false;
    } else if (generalChanged.value) {
        try {
            await postConfig('generalConfig', '/configuration', formatGeneralConfig(generalConfig), csrfToken);
        } catch (error) {
            showErrorToast(error.message);
            saveSuccess = false;
        }
    }

    const { isValid: isPeginValid, errors: peginErrors } = validateConfig(peginConfigData, configurations.pegin);
    if (!isPeginValid) {
        showErrorToast(peginErrors.join('<br>'));
        saveSuccess = false;
    } else if (peginChanged.value) {
        try {
            await postConfig('peginConfig', '/pegin/configuration', peginConfigData, csrfToken);
        } catch (error) {
            showErrorToast(error.message);
            saveSuccess = false;
        }
    }

    const { isValid: isPegoutValid, errors: pegoutErrors } = validateConfig(pegoutConfigData, configurations.pegout);
    if (!isPegoutValid) {
        showErrorToast(pegoutErrors.join('<br>'));
        saveSuccess = false;
    } else if (pegoutChanged.value) {
        try {
            await postConfig('pegoutConfig', '/pegout/configuration', pegoutConfigData, csrfToken);
        } catch (error) {
            showErrorToast(error.message);
            saveSuccess = false;
        }
    }

    if (saveSuccess) showSuccessToast();
};

const addCollateral = async (amountId, endpoint, elementId, loadingBarId, buttonId, csrfToken) => {
    const amountInEther = document.getElementById(amountId).value;
    const loadingBar = document.getElementById(loadingBarId);
    const button = document.getElementById(buttonId);
    loadingBar.style.display = 'block';
    button.disabled = true;
    try {
        const amountInWei = Number(etherToWei(amountInEther));
        const response = await fetch(endpoint, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', 'X-CSRF-Token': csrfToken },
            body: JSON.stringify({ amount: amountInWei })
        });
        if (response.ok) {
            fetchData(endpoint.replace('/addCollateral', '/collateral'), elementId, csrfToken);
        } else {
            const errorData = await response.json();
            showErrorToast(`Error adding collateral: ${errorData.message || 'Unknown error'}`);
        }
    } catch (error) {
        showErrorToast(`Invalid input "${amountInEther}" for collateral amount. Please enter a valid number.`);
    } finally {
        loadingBar.style.display = 'none';
        button.disabled = false;
    }
};

const displaySummaryData = (container, data) => {
    container.innerHTML = '';
    const table = document.createElement('table');
    table.classList.add('table', 'table-striped');
    const rows = [
        { label: 'Total Quotes', value: data.totalQuotesCount },
        { label: 'Accepted Quotes', value: data.acceptedQuotesCount },
        { label: 'Paid Quotes', value: data.paidQuotesCount },
        { label: 'Paid Quotes Amount', value: data.paidQuotesAmount },
        { label: 'Total Accepted Amount', value: data.totalAcceptedQuotedAmount },
        { label: 'Total Fees Collected', value: data.totalFeesCollected },
        { label: 'Refunded Quotes', value: data.refundedQuotesCount },
        { label: 'Total Penalty Amount', value: data.totalPenaltyAmount },
        { label: 'LP Earnings', value: data.lpEarnings }
    ];
    rows.forEach(row => {
        const tr = document.createElement('tr');
        const th = document.createElement('th');
        th.textContent = row.label;
        const td = document.createElement('td');
        td.textContent = row.value;
        tr.appendChild(th);
        tr.appendChild(td);
        table.appendChild(tr);
    });
    container.appendChild(table);
};

const fetchSummariesReport = async (csrfToken) => {
    const startDate = document.getElementById('summaryStartDate').value;
    const endDate = document.getElementById('summaryEndDate').value;
    if (!startDate || !endDate) {
        showErrorToast('Please select both start and end dates');
        return;
    }
    try {
        const response = await fetch(`/reports/summaries?startDate=${startDate}&endDate=${endDate}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-CSRF-Token': csrfToken
            }
        });
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || 'Failed to fetch summaries');
        }
        const data = await response.json();
        document.getElementById('summariesResult').style.display = 'block';
        displaySummaryData(document.getElementById('peginSummary'), data.peginSummary);
        displaySummaryData(document.getElementById('pegoutSummary'), data.pegoutSummary);
    } catch (error) {
        showErrorToast(`Error fetching summaries: ${error.message}`);
    } finally {
        loadingBar.style.display = 'none';
    }
};

const fetchTrustedAccounts = async (csrfToken) => {
    const loadingBar = document.getElementById('trustedAccountsLoadingBar');
    loadingBar.style.display = 'block';
    try {
        const response = await fetch('/management/trusted-accounts', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-CSRF-Token': csrfToken
            }
        });
        const data = await response.json();
        if (!response.ok) {
            const errorMessage = data.details?.error || data.message || 'Unknown error';
            populateTrustedAccountsTable([], csrfToken, errorMessage);
            showErrorToast(`Failed to load trusted accounts: ${errorMessage}`);
        } else {
            const accountsData = data.accounts || [];
            populateTrustedAccountsTable(accountsData, csrfToken);
        }
    } catch (error) {
        console.error('Error fetching trusted accounts:', error);
        populateTrustedAccountsTable([], csrfToken, error.message);
        showErrorToast(`Failed to load trusted accounts: ${error.message}`);
    } finally {
        loadingBar.style.display = 'none';
    }
};

const populateTrustedAccountsTable = (accounts, csrfToken, errorMessage = null) => {
    const tableBody = document.getElementById('trustedAccountsTable');
    tableBody.innerHTML = '';

    if (errorMessage) {
        const row = document.createElement('tr');
        const cell = document.createElement('td');
        cell.colSpan = 5;
        cell.classList.add('text-center', 'text-danger');
        cell.textContent = `Error: ${errorMessage}`;
        row.appendChild(cell);
        tableBody.appendChild(row);
        return;
    }

    if (!accounts || accounts.length === 0) {
        const row = document.createElement('tr');
        const cell = document.createElement('td');
        cell.colSpan = 5;
        cell.classList.add('text-center');
        cell.textContent = 'No trusted accounts found.';
        row.appendChild(cell);
        tableBody.appendChild(row);
        return;
    }

    accounts.forEach(account => {
        const row = document.createElement('tr');
        const nameCell = document.createElement('td');
        nameCell.textContent = account.name || 'Unknown';
        const addressCell = document.createElement('td');
        addressCell.classList.add('address-cell');
        addressCell.textContent = account.address;

        const btcCapCell = document.createElement('td');
        btcCapCell.classList.add('cap-cell');
        const btcValue = weiToEther(account.btcLockingCap);
        btcCapCell.textContent = formatCap(btcValue, 'BTC');

        const rbtcCapCell = document.createElement('td');
        rbtcCapCell.classList.add('cap-cell');
        const rbtcValue = weiToEther(account.rbtcLockingCap);
        rbtcCapCell.textContent = formatCap(rbtcValue, 'rBTC');

        const actionsCell = document.createElement('td');
        const deleteButton = document.createElement('button');
        deleteButton.type = 'button';
        deleteButton.classList.add('btn', 'btn-danger', 'btn-sm');
        deleteButton.textContent = 'Remove';
        deleteButton.addEventListener('click', () => removeTrustedAccount(account.address, csrfToken));
        actionsCell.appendChild(deleteButton);

        row.appendChild(nameCell);
        row.appendChild(addressCell);
        row.appendChild(btcCapCell);
        row.appendChild(rbtcCapCell);
        row.appendChild(actionsCell);
        tableBody.appendChild(row);
    });
};

// Helper function to clear form validation states
const clearFormValidation = () => {
    const formFields = ['accountName', 'accountAddress', 'btc_locking_cap', 'rbtc_locking_cap'];
    formFields.forEach(fieldId => {
        const field = document.getElementById(fieldId);
        field.classList.remove('is-invalid', 'is-valid');
        // Remove any existing error message
        const feedback = field.parentElement.querySelector('.invalid-feedback');
        if (feedback) {
            feedback.remove();
        }
    });
};

// Helper function to show field-specific validation errors
const showFieldError = (fieldId, errorMessage) => {
    const field = document.getElementById(fieldId);
    field.classList.add('is-invalid');
    field.classList.remove('is-valid');

    // Remove existing error message if any
    const existingFeedback = field.parentElement.querySelector('.invalid-feedback');
    if (existingFeedback) {
        existingFeedback.remove();
    }

    // Add new error message
    const errorDiv = document.createElement('div');
    errorDiv.className = 'invalid-feedback';
    errorDiv.textContent = errorMessage;
    field.parentElement.appendChild(errorDiv);
};

// Helper function to validate numeric input
const validatePositiveNumber = (value, fieldName) => {
    if (!value || value.trim() === '') {
        return `${fieldName} is required`;
    }

    const numValue = parseFloat(value);
    if (isNaN(numValue) || numValue <= 0) {
        return `${fieldName} must be a positive number`;
    }

    return null;
};

const addTrustedAccount = async (csrfToken) => {
    // Clear any previous validation states
    clearFormValidation();

    const name = document.getElementById('accountName').value.trim();
    const address = document.getElementById('accountAddress').value.trim();
    let btcLockingCap = document.getElementById('btc_locking_cap').value.trim();
    let rbtcLockingCap = document.getElementById('rbtc_locking_cap').value.trim();

    let hasValidationErrors = false;

    if (!name) {
        showFieldError('accountName', 'Account name is required');
        hasValidationErrors = true;
    }

    if (!address) {
        showFieldError('accountAddress', 'Account address is required');
        hasValidationErrors = true;
    }

    const btcCapError = validatePositiveNumber(btcLockingCap, 'BTC Locking Cap');
    if (btcCapError) {
        showFieldError('btc_locking_cap', btcCapError);
        hasValidationErrors = true;
    }

    const rbtcCapError = validatePositiveNumber(rbtcLockingCap, 'rBTC Locking Cap');
    if (rbtcCapError) {
        showFieldError('rbtc_locking_cap', rbtcCapError);
        hasValidationErrors = true;
    }

    if (hasValidationErrors) {
        return;
    }

    try {
        // Convert to wei only for valid positive numbers
        btcLockingCap = etherToWei(btcLockingCap);
        rbtcLockingCap = etherToWei(rbtcLockingCap);

        const response = await fetch('/management/trusted-accounts', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-CSRF-Token': csrfToken
            },
            body: `{
                "name": ${JSON.stringify(name)},
                "address": ${JSON.stringify(address)},
                "btcLockingCap": ${btcLockingCap},
                "rbtcLockingCap": ${rbtcLockingCap}
            }`
        });

        if (response.ok) {
            const modal = bootstrap.Modal.getInstance(document.getElementById('addTrustedAccountModal'));
            modal.hide();
            document.getElementById('accountName').value = '';
            document.getElementById('accountAddress').value = '';
            document.getElementById('btc_locking_cap').value = '';
            document.getElementById('rbtc_locking_cap').value = '';
            clearFormValidation();
            showSuccessToast();
            fetchTrustedAccounts(csrfToken);
        } else {
            const errorData = await response.json();

            // Handle field-specific validation errors from backend
            if (errorData.message === 'validation error' && errorData.details) {
                let hasFieldErrors = false;

                // Map backend field names to frontend field IDs
                const fieldMapping = {
                    'Name': 'accountName',
                    'Address': 'accountAddress',
                    'BtcLockingCap': 'btc_locking_cap',
                    'RbtcLockingCap': 'rbtc_locking_cap'
                };

                for (const [backendField, errorMessage] of Object.entries(errorData.details)) {
                    const frontendFieldId = fieldMapping[backendField];
                    if (frontendFieldId) {
                        showFieldError(frontendFieldId, errorMessage);
                        hasFieldErrors = true;
                    }
                }

                if (hasFieldErrors) {
                    return; // Don't show toast if we have field-specific errors
                }
            }

            // Show generic error if no field-specific errors
            const errorMessage = errorData.details?.error || errorData.message || 'Unknown error';
            showErrorToast(`Error adding trusted account: ${errorMessage}`);
        }
    } catch (error) {
        showErrorToast(`Error adding trusted account: ${error.message}`);
    }
};

const removeTrustedAccount = async (address, csrfToken) => {
    if (!confirm(`Are you sure you want to remove the trusted account with address ${address}?`)) return;
    try {
        const response = await fetch(`/management/trusted-accounts?address=${encodeURIComponent(address)}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'X-CSRF-Token': csrfToken
            }
        });
        if (response.ok) {
            showSuccessToast();
            fetchTrustedAccounts(csrfToken);
        } else {
            const errorData = await response.json();
            const errorMessage = errorData.details?.error || errorData.message || 'Unknown error';
            showErrorToast(`Error removing trusted account: ${errorMessage}`);
        }
    } catch (error) {
        showErrorToast(`Error removing trusted account: ${error.message}`);
    }
};

document.addEventListener('DOMContentLoaded', () => {
    const csrfToken = data.CsrfToken;
    const configurations = data.Configuration;
    const providerData = data.ProviderData;
    const rskAddress = data.RskAddress;
    const btcAddress = data.BtcAddress;

    document.getElementById('addPeginCollateralButton').addEventListener('click', () => addCollateral('addPeginCollateralAmount', '/pegin/addCollateral', 'peginCollateral', 'peginLoadingBar', 'addPeginCollateralButton', csrfToken));
    document.getElementById('addPegoutCollateralButton').addEventListener('click', () => addCollateral('addPegoutCollateralAmount', '/pegout/addCollateral', 'pegoutCollateral', 'pegoutLoadingBar', 'addPegoutCollateralButton', csrfToken));
    document.getElementById('saveConfig').addEventListener('click', () => saveConfig(csrfToken, configurations));
    document.getElementById('fetchSummariesButton').addEventListener('click', () => fetchSummariesReport(csrfToken));
    document.getElementById('saveAccountButton').addEventListener('click', () => addTrustedAccount(csrfToken));

    // Clear validation states when modal is opened
    document.getElementById('addTrustedAccountModal').addEventListener('show.bs.modal', () => {
        clearFormValidation();
    });

    // Clear field-specific validation errors when user starts typing
    ['accountName', 'accountAddress', 'btc_locking_cap', 'rbtc_locking_cap'].forEach(fieldId => {
        document.getElementById(fieldId).addEventListener('input', function() {
            this.classList.remove('is-invalid');
            const feedback = this.parentElement.querySelector('.invalid-feedback');
            if (feedback) {
                feedback.remove();
            }
        });
    });

    document.querySelectorAll('#configTabs a[data-bs-toggle="tab"]').forEach(tabEl => {
        tabEl.addEventListener('shown.bs.tab', () => checkFeeWarnings());
    });

    populateConfigSection('generalConfig', configurations.general);
    populateConfigSection('peginConfig', configurations.pegin);
    populateConfigSection('pegoutConfig', configurations.pegout);
    populateProviderData(providerData, rskAddress, btcAddress);

    fetchData('/pegin/collateral', 'peginCollateral', csrfToken);
    fetchData('/pegout/collateral', 'pegoutCollateral', csrfToken);
    checkFeeWarnings();

    const today = new Date();
    const lastMonth = new Date(today);
    lastMonth.setMonth(today.getMonth() - 1);
    fetchTrustedAccounts(csrfToken);
    document.getElementById('summaryStartDate').value = lastMonth.toISOString().split('T')[0];
    document.getElementById('summaryEndDate').value = today.toISOString().split('T')[0];
});
