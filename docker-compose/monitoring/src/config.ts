import { type BitcoinAddress } from '@rsksmart/rsk-monitor'

export const MONITORED_ADDRESSES: BitcoinAddress[] = [
  { address: 'mwEc...', alias: 'testnet1' },
  { address: 'mvL2...', alias: 'testnet2' },
  { address: 'tb1q7...', alias: 'testnet3' }
]

export const MONITOR_CONFIG = {
  pollingIntervalSeconds: 10,
  monitorName: 'bitcoin-balance-monitor',
  network: 'testnet' as 'mainnet' | 'testnet'
} 