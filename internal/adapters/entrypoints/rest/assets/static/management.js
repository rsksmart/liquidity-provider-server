import {
    weiToEther,
    etherToWei,
    isFeeKey,
    validateConfig,
    formatGeneralConfig,
    postConfig,
    hasDuplicateConfirmationAmounts,
    isfeePercentageKey,
    isToggableFeeKey
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

    const input = document.createElement('input');
    input.type = 'text';
    input.style.width = '40%';
    input.classList.add('form-control');
    input.dataset.key = key;
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

document.addEventListener('DOMContentLoaded', () => {
    const csrfToken = data.CsrfToken;
    const configurations = data.Configuration;
    const providerData = data.ProviderData;
    const rskAddress = data.RskAddress;
    const btcAddress = data.BtcAddress;

    document.getElementById('addPeginCollateralButton').addEventListener('click', () => addCollateral('addPeginCollateralAmount', '/pegin/addCollateral', 'peginCollateral', 'peginLoadingBar', 'addPeginCollateralButton', csrfToken));
    document.getElementById('addPegoutCollateralButton').addEventListener('click', () => addCollateral('addPegoutCollateralAmount', '/pegout/addCollateral', 'pegoutCollateral', 'pegoutLoadingBar', 'addPegoutCollateralButton', csrfToken));
    document.getElementById('saveConfig').addEventListener('click', () => saveConfig(csrfToken, configurations));

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
});