import { type BitcoinAddress } from '@rsksmart/rsk-monitor'


/**
 * Uncomment this to monitor testnet addresses
 */
export const MONITORED_ADDRESSES: BitcoinAddress[] = [
  { address: 'mwEceC31MwWmF6hc5SSQ8FmbgdsSoBSnbm', alias: 'testnet1' },
  { address: 'mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe', alias: 'testnet2' },
  { address: 'tb1q7hec37mcmfk6hmqn67echzdf8zwg5n5pqnfzma', alias: 'testnet3' }
]

export const MONITOR_CONFIG = {
  pollingIntervalSeconds: 10,
  monitorName: 'bitcoin-balance-monitor',
  network: 'testnet' as 'mainnet' | 'testnet'
} 


/**
 * Uncomment this to monitor mainnet addresses
 */
// export const MONITORED_ADDRESSES: BitcoinAddress[] = [
//   { address: 'bc1qv7l4jgnzxyjgn598ee04l72lanudx50fqpdq0t', alias: 'mainnet1' },
//   { address: '3DGxAYYUA61WrrdbBac8Ra9eA9peAQwTJF', alias: 'mainnet2' }
// ]

// export const MONITOR_CONFIG = {
//   pollingIntervalSeconds: 2 * 60,
//   monitorName: 'bitcoin-balance-monitor',
//   network: 'mainnet' as 'mainnet' | 'testnet'
// } 