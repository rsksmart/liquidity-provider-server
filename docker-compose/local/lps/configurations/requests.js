const {readPassword} = require("./credentials");
const {CONSTANTS} = require("./constants");

const CONFIG_REQUESTS = {
    login: {
        path: '/management/login',
        method: 'POST',
        body: {
            username: "admin",
            password: readPassword(),
        }
    },
    generalConfiguration: {
        path: '/configuration',
        method: 'POST',
        body: {
            configuration: {
                rskConfirmations: {
                    "100000000000000000": 4,
                    "2000000000000000000": 20,
                    "400000000000000000": 12,
                    "4000000000000000000": 40,
                    "8000000000000000000": 80
                },
                btcConfirmations: {
                    "100000000000000000": 2,
                    "2000000000000000000": 10,
                    "400000000000000000": 6,
                    "4000000000000000000": 20,
                    "8000000000000000000": 40
                },
                publicLiquidityCheck: true,
                maxLiquidity: "2000000000000000000000",
                reimbursementWindowBlocks: 100,
                excessTolerance: {
                    isFixed: false,
                    percentageValue: 15,
                    fixedValue: "0"
                }
            }
        }
    },
    peginConfiguration: {
        path: '/pegin/configuration',
        method: 'POST',
        body: {
            configuration: {
                timeForDeposit: 3600,
                callTime: 7200,
                penaltyFee: "1000000000000000",
                maxValue: "10000000000000000000",
                minValue: "600000000000000000",
                feePercentage: 0.33,
                fixedFee: "200000000000000"
            }
        }
    },
    pegoutConfiguration: {
        path: '/pegout/configuration',
        method: 'POST',
        body: {
            configuration: {
                timeForDeposit: 3600,
                // expireTime/expireBlocks must exceed the estimated confirmation time for maxValue.
                // At maxValue=10 ETH the general config resolves to 80 RSK + 40 BTC confirmations:
                // 80*30 + 40*600 = 26400s. So expireTime > 26400s and expireBlocks*30 > 26400s (>880 blocks).
                expireTime: 28800,
                penaltyFee: "1000000000000000",
                maxValue: "10000000000000000000",
                minValue: "600000000000000000",
                expireBlocks: 1000,
                bridgeTransactionMin: "1500000000000000000",
                feePercentage: 0.33,
                fixedFee: "200000000000000"
            }
        }
    },
    trustedAccountsConfiguration: {
        path: '/management/trusted-accounts',
        method: 'POST',
        body: {
            address: process.env.TRUSTED_ACCOUNT_ADDRESS,
            name: "Boletaz",
            btcLockingCap: 3000000000000000000,
            rbtcLockingCap: 3000000000000000000
        }
    },
    updatePassword: {
        path: "/management/credentials",
        method: "POST",
        body: {
            oldUsername: "admin",
            newUsername: "admin",
            oldPassword: readPassword(),
            newPassword: CONSTANTS.MANAGEMENT_DEFAULT_PASSWORD
        }
    }
}

module.exports = { CONFIG_REQUESTS }
