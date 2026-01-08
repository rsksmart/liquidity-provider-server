function weiToEther(wei) {
    if (wei === null || wei === undefined) return '0';
    try {
        const decimalValue = new Decimal(wei);
        return decimalValue.dividedBy(new Decimal(1e18)).toString();
    } catch (error) {
        throw new Error(`Failed to convert wei to ether. Input: "${wei}". Error: ${error.message}`);
    }
}

function etherToWei(ether) {
    if (typeof ether !== 'string' && typeof ether !== 'number') {
        throw new TypeError(`Invalid input type for ether: ${typeof ether}. Expected a number or string.`);
    }
    try {
        const num = new Decimal(ether);
        if (num.isNegative()) {
            throw new RangeError(`The input "${ether}" is not a valid number.`);
        }
        return num.times(new Decimal(1e18)).toFixed();
    } catch (error) {
        throw new Error(`Failed to convert ether to wei. Input: "${ether}". Error: ${error.message}`);
    }
}

function isFeeKey(key) {
    return ['penaltyFee', 'callFee', 'maxValue', 'minValue', 'bridgeTransactionMin','fixedFee', 'maxLiquidity'].includes(key);
}

function isMaxLiquidityKey(key) {
    return key === 'maxLiquidity';
}

/**
 * Validates a maxLiquidity value.
 * @param {string|number} value - The value to validate (in RBTC)
 * @returns {{isValid: boolean, error: string|null}} Validation result
 */
function validateMaxLiquidity(value) {
    if (value === null || value === undefined || value === '') {
        return { isValid: false, error: 'Max liquidity is required' };
    }

    const strValue = String(value).trim();
    
    // Check if it's a valid number
    const num = parseFloat(strValue);
    if (isNaN(num)) {
        return { isValid: false, error: 'Max liquidity must be a valid number' };
    }

    // Check if positive
    if (num <= 0) {
        return { isValid: false, error: 'Max liquidity must be a positive number' };
    }

    // Check decimal places (maximum 18)
    const decimalPart = strValue.split('.')[1];
    if (decimalPart && decimalPart.length > 18) {
        return { isValid: false, error: 'Max liquidity cannot have more than 18 decimal places' };
    }

    return { isValid: true, error: null };
}

function isfeePercentageKey(key) {
    return key === 'feePercentage';
}

function isExcessToleranceKey(key) {
    return key === 'excessToleranceFixed' || key === 'excessTolerancePercentage';
}

function isExcessToleranceFixedKey(key) {
    return key === 'excessToleranceFixed';
}

function isExcessTolerancePercentageKey(key) {
    return key === 'excessTolerancePercentage';
}

/**
 * Validates an excessToleranceFixed value.
 * @param {string|number} value - The value to validate (in wei as bigint string)
 * @returns {{isValid: boolean, error: string|null}} Validation result
 */
function validateExcessToleranceFixed(value) {
    if (value === null || value === undefined || value === '') {
        return { isValid: true, error: null }; // Optional field, empty is OK
    }

    const strValue = String(value).trim();
    
    // Check if it's a valid number
    const num = parseFloat(strValue);
    if (isNaN(num)) {
        return { isValid: false, error: 'Excess tolerance fixed must be a valid number' };
    }

    // Check if non-negative
    if (num < 0) {
        return { isValid: false, error: 'Excess tolerance fixed must be a non-negative number' };
    }

    // Check decimal places (maximum 18 for wei conversion)
    const decimalPart = strValue.split('.')[1];
    if (decimalPart && decimalPart.length > 18) {
        return { isValid: false, error: 'Excess tolerance fixed cannot have more than 18 decimal places' };
    }

    return { isValid: true, error: null };
}

/**
 * Validates an excessTolerancePercentage value.
 * @param {string|number} value - The value to validate (0-100 percentage)
 * @returns {{isValid: boolean, error: string|null}} Validation result
 */
function validateExcessTolerancePercentage(value) {
    if (value === null || value === undefined || value === '') {
        return { isValid: true, error: null }; // Optional field, empty is OK (will default to 0)
    }

    const strValue = String(value).trim();
    const num = parseFloat(strValue);

    if (isNaN(num)) {
        return { isValid: false, error: 'Excess tolerance percentage must be a valid number' };
    }

    if (num < 0) {
        return { isValid: false, error: 'Excess tolerance percentage must be non-negative' };
    }

    if (num > 100) {
        return { isValid: false, error: 'Excess tolerance percentage cannot exceed 100%' };
    }

    return { isValid: true, error: null };
}

function inferType(value) {
    if (value === null || value === undefined) return 'undefined';
    if (Array.isArray(value)) return 'array';
    return typeof value;
}

function validateConfig(config, originalConfig) {
    const errors = [];
    const confirmationKeys = ['rskConfirmations', 'btcConfirmations'];

    for (const [key, value] of Object.entries(config)) {
        const expectedValue = originalConfig[key];
        const expectedType = inferType(expectedValue);
        let actualType = inferType(value);
        if (
            (isFeeKey(key) || isfeePercentageKey(key) || isExcessToleranceKey(key)) &&
            ((expectedType === 'number' && actualType === 'string') ||
             (expectedType === 'string' && actualType === 'number'))
        ) {
            actualType = expectedType;
        }
        if (expectedType === 'undefined') continue;
        if (confirmationKeys.includes(key)) {
            if (actualType !== 'object') {
                errors.push(`Invalid type for ${key}: expected object, got ${actualType}`);
                continue;
            }
            for (const [subKey, subValue] of Object.entries(value)) {
                if (inferType(subValue) !== 'number') {
                    errors.push(`Invalid type for ${key} confirmation value of amount ${subKey}: expected number, got ${inferType(subValue)}`);
                }
            }
        } else if (actualType !== expectedType) {
            errors.push(`Invalid type for ${key}: expected ${expectedType}, got ${actualType}`);
        }
    }
    return { isValid: errors.length === 0, errors };
}

function formatGeneralConfig(config) {
    const formattedConfig = {};
    Object.keys(config).forEach(key => {
        if (key === 'rskConfirmations' || key === 'btcConfirmations') {
            formattedConfig[key] = {};
            config[key].forEach(entry => {
                if (entry.amount && entry.confirmation !== undefined) {
                    formattedConfig[key][entry.amount] = entry.confirmation;
                }
            });
        } else {
            formattedConfig[key] = config[key];
        }
    });
    return formattedConfig;
}

async function postConfig(sectionId, endpoint, config, csrfToken) {
    try {
        const response = await fetch(endpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-CSRF-Token': csrfToken
            },
            body: JSON.stringify({ configuration: config })
        });
        if (!response.ok) {
            const errorData = await response.json();            
            if (errorData.details && typeof errorData.details === 'object') {
                const detailMessages = Object.entries(errorData.details)
                    .map(([field, message]) => `${field}: ${message}`)
                    .join(', ');
                throw new Error(detailMessages);
            }
            throw new Error(`Error saving ${sectionId} configuration: ${errorData.message || 'Unknown error'}`);
        }
        return true;
    } catch (error) {
        throw new Error(`${sectionId}: ${error.message}`);
    }
}

/**
 * Checks if there are any duplicated rBTC amounts in the given confirmation array.
 */
function hasDuplicateConfirmationAmounts(confirmationArray) {
    const amounts = confirmationArray.map(entry => entry.amount);
    const uniqueAmounts = new Set(amounts);
    return uniqueAmounts.size < amounts.length;
}

function isToggableFeeKey(key) {
    return key === 'fixedFee' || key === 'feePercentage';
}
const formatCap = (value, unit) => {
    try {
        const num = parseFloat(value);
        return parseFloat(num.toFixed(4)).toString() + ' ' + unit;
    } catch (e) {
        console.error('Error formatting cap:', e);
        return `Error: ${e.message || 'Failed to format value'}`;
    }
};

export {
    weiToEther,
    etherToWei,
    isFeeKey,
    isMaxLiquidityKey,
    validateMaxLiquidity,
    inferType,
    validateConfig,
    formatGeneralConfig,
    postConfig,
    hasDuplicateConfirmationAmounts,
    isfeePercentageKey,
    isToggableFeeKey,
    formatCap,
    isExcessToleranceKey,
    isExcessToleranceFixedKey,
    isExcessTolerancePercentageKey,
    validateExcessToleranceFixed,
    validateExcessTolerancePercentage
};
