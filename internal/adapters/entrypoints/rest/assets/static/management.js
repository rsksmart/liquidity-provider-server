import { 
    weiToEther, 
    etherToWei,
    isFeeKey,
    validateConfig,
    formatGeneralConfig,
    postConfig
} from './configUtils.js';

document.addEventListener('DOMContentLoaded', () => {
    const generalChanged = { value: false };
    const peginChanged = { value: false };
    const pegoutChanged = { value: false };

    const setTextContent = (id, text) => document.getElementById(id).textContent = text;

    const fetchData = async (url, elementId) => {
        try {
            const response = await fetch(url, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'X-CSRF-Token': data.CsrfToken
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

        const input = document.createElement('input');
        input.style.marginLeft = "10px";
        input.dataset.key = key;
        input.dataset.originalValue = value;

        if (typeof value === 'boolean') {
            input.type = 'checkbox';
            input.classList.add('form-check-input');
            input.checked = value;
            input.addEventListener('change', () => setChanged(section.id));
        } else {
            input.type = 'text';
            input.style.width = "40%";
            input.classList.add('form-control');
            input.value = isFeeKey(key) ? weiToEther(value) : value;
            input.addEventListener('input', () => setChanged(section.id));
            if (isFeeKey(key)) {
                const questionIcon = createQuestionIcon(getTooltipText(key));
                label.appendChild(questionIcon);
            }
        }

        inputContainer.appendChild(input);
        div.appendChild(label);
        div.appendChild(inputContainer);
        section.appendChild(div);
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
            bridgeTransactionMin: 'The amount of rBTC that needs to be gathered in peg out refunds before executing a native peg out.'
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

        const sortedEntries = Object.entries(confirmations).sort(([amountWeiA], [amountWeiB]) => {
            const amountA = new Decimal(weiToEther(amountWeiA));
            const amountB = new Decimal(weiToEther(amountWeiB));
            return amountA.minus(amountB).toNumber();
        });
        sortedEntries.forEach(([amountWei, confirmation], index) => {
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
            setChanged(configKey);
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
        amountInputAppend.textContent = 'rBTC';
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
            <div class="toast-body">
                ${errorMessage}
            </div>
        `;
        document.querySelector('.toast-container').appendChild(toastElement);
        const toast = new bootstrap.Toast(toastElement);
        toast.show();
    };

    function getConfirmationConfig(sectionId) {
        const entries = document.querySelectorAll(`#${sectionId} .confirmation-config`);
        const config = {};
        entries.forEach(entry => {
            const configKey = entry.dataset.configKey;
            const inputGroups = entry.querySelectorAll('input');
            const tempArray = [];
            inputGroups.forEach(input => {
                const idx = input.dataset.index;
                if (!tempArray[idx]) tempArray[idx] = {};
                if (input.dataset.field === 'amount') {
                    tempArray[idx].amount = etherToWei(input.value).toString();
                } else if (input.dataset.field === 'confirmation') {
                    const val = Number(input.value);
                    if (isNaN(val) || !Number.isInteger(val) || val < 0) {
                        showErrorToast(`Invalid input "${input.value}" for field "Confirmation". Please enter a valid non-negative integer.`);
                        throw new Error('Invalid confirmation number');
                    }
                    tempArray[idx].confirmation = val;
                }
            });
            config[configKey] = tempArray;
        });
        return config;
    }

    function getRegularConfig(sectionId) {
        const inputs = document.querySelectorAll(`#${sectionId} input:not(.form-check-input):not([data-field="amount"]):not([data-field="confirmation"])`);
        const checkboxes = document.querySelectorAll(`#${sectionId} input.form-check-input`);
        const config = {};

        inputs.forEach(input => {
            let value;
            if (isFeeKey(input.dataset.key)) {
                try {
                    value = etherToWei(input.value).toString();
                } catch (error) {
                    showErrorToast(`Invalid input "${input.value}" for field "${input.dataset.key}". Please enter a valid number.`);
                    throw error;
                }
            } else {
                value = input.value;
                if (!isNaN(value) && !isNaN(parseFloat(value))) {
                    value = Number(value);
                }
            }
            config[input.dataset.key] = value;
        });

        checkboxes.forEach(input => {
            config[input.dataset.key] = input.checked;
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

    const saveConfig = async () => {
        const csrfToken = data.CsrfToken;
        let saveSuccess = true;

        let generalConfig, peginConfigData, pegoutConfigData;
        try {
            generalConfig = getConfig('generalConfig');
            peginConfigData = getConfig('peginConfig');
            pegoutConfigData = getConfig('pegoutConfig');
        } catch (error) {
            return;
        }

        const { isValid: isGeneralValid, errors: generalErrors } = validateConfig(formatGeneralConfig(generalConfig), data.Configuration.general);
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

        const { isValid: isPeginValid, errors: peginErrors } = validateConfig(peginConfigData, data.Configuration.pegin);
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

        const { isValid: isPegoutValid, errors: pegoutErrors } = validateConfig(pegoutConfigData, data.Configuration.pegout);
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

        if (saveSuccess) {
            showSuccessToast();
            sortConfirmationEntries();
        }
    };
    
    const sortConfirmationEntries = () => {
        const confirmationConfigs = document.querySelectorAll('.confirmation-config');
        confirmationConfigs.forEach(container => {
            const entriesContainer = container.querySelector('.entries-container');
            const entries = Array.from(entriesContainer.children);

            entries.sort((a, b) => {
                const amountA = getAmountFromEntry(a);
                const amountB = getAmountFromEntry(b);
                return amountA.minus(amountB).toNumber();
            });

            entriesContainer.innerHTML = '';
            entries.forEach(entry => entriesContainer.appendChild(entry));
        });
    };

    const getAmountFromEntry = (entry) => {
        const amountInput = entry.querySelector('input[data-field="amount"]');
        const amountValue = amountInput.value;
        try {
            return new Decimal(amountValue);
        } catch (error) {
            return new Decimal(0);
        }
    };

    const addCollateral = async (amountId, endpoint, elementId, loadingBarId, buttonId) => {
        const amountInEther = document.getElementById(amountId).value;
        const loadingBar = document.getElementById(loadingBarId);
        const button = document.getElementById(buttonId);
        loadingBar.style.display = 'block';
        button.disabled = true;
        try {
            const amountInWei = Number(etherToWei(amountInEther));
            const response = await fetch(endpoint, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'X-CSRF-Token': data.CsrfToken },
                body: JSON.stringify({ amount: amountInWei })
            });
            if (response.ok) {
                fetchData(endpoint.replace('/addCollateral', '/collateral'), elementId);
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

    document.getElementById('addPeginCollateralButton').addEventListener('click', () => addCollateral('addPeginCollateralAmount', '/pegin/addCollateral', 'peginCollateral', 'peginLoadingBar', 'addPeginCollateralButton'));
    document.getElementById('addPegoutCollateralButton').addEventListener('click', () => addCollateral('addPegoutCollateralAmount', '/pegout/addCollateral', 'pegoutCollateral', 'pegoutLoadingBar', 'addPegoutCollateralButton'));
    document.getElementById('saveConfig').addEventListener('click', saveConfig);

    populateConfigSection('generalConfig', data.Configuration.general);
    populateConfigSection('peginConfig', data.Configuration.pegin);
    populateConfigSection('pegoutConfig', data.Configuration.pegout);
    populateProviderData(data.ProviderData, data.RskAddress, data.BtcAddress);

    fetchData('/pegin/collateral', 'peginCollateral');
    fetchData('/pegout/collateral', 'pegoutCollateral');
});
