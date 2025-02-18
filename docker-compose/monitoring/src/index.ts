import { BitcoinMempoolMonitorBuilder, ConsoleExporter, MonitorConfig, PrometheusExporter } from '@rsksmart/rsk-monitor'
import { MONITORED_ADDRESSES, MONITOR_CONFIG } from './config'

async function main (): Promise<void> {
  await BitcoinMempoolMonitorBuilder.create({
    pollingIntervalSeconds: MONITOR_CONFIG.pollingIntervalSeconds,
    monitorName: MONITOR_CONFIG.monitorName,
    network: MONITOR_CONFIG.network
  } as MonitorConfig, MONITORED_ADDRESSES)
    .withBalanceMetric()
    .withExporters(new ConsoleExporter(), new PrometheusExporter(MONITOR_CONFIG.port, 'bitcoinbalancemonitor'))
    .build()
    .run()

  console.log('Starting Bitcoin balance monitor...')
  console.log('Monitoring addresses:', MONITORED_ADDRESSES)

  process.on('SIGINT', () => {
    console.log('Stopping monitor...')
    process.exit(0)
  })
}

main().catch(console.error)
