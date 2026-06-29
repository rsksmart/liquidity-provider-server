import type { Plugin, ViteDevServer } from 'vite'

const LPS_DEV_TARGET = process.env.LPS_DEV_TARGET ?? 'http://localhost:8080'

function extractCsrfFromHtml(html: string): string {
  const match = html.match(/<meta name="csrf-token" content="([^"]*)"/)
  if (!match?.[1]?.trim()) {
    throw new Error('csrf-token meta missing or empty in LPS shell')
  }
  return match[1].replaceAll('&#43;', '+')
}

function extractInitialDataFromHtml(html: string): unknown {
  const match = html.match(/<script\b[^>]*\bid="initial-data"[^>]*>([\s\S]*?)<\/script>/)
  if (!match?.[1]?.trim()) {
    throw new Error('initial-data script missing or empty in LPS shell')
  }
  return JSON.parse(match[1]) as unknown
}

export function lpsDevBootstrapPlugin(): Plugin {
  return {
    name: 'lps-dev-bootstrap',
    apply: 'serve',
    configureServer(server: ViteDevServer) {
      server.middlewares.use('/__dev/lps-shell', async (req, res, next) => {
        if (req.method !== 'GET') {
          next()
          return
        }

        try {
          const cookie = req.headers.cookie ?? ''
          const lpsRes = await fetch(`${LPS_DEV_TARGET}/management/next/login`, {
            headers: cookie ? { cookie } : undefined,
          })
          if (!lpsRes.ok) {
            res.statusCode = 502
            res.setHeader('Content-Type', 'application/json')
            res.end(JSON.stringify({ error: `LPS shell GET failed: ${lpsRes.status}` }))
            return
          }

          const html = await lpsRes.text()
          const csrf = extractCsrfFromHtml(html)
          const initialData = extractInitialDataFromHtml(html)
          const setCookies = lpsRes.headers.getSetCookie?.() ?? []

          if (setCookies.length > 0) {
            for (const cookie of setCookies) {
              res.appendHeader('Set-Cookie', cookie)
            }
          } else {
            const rawSetCookie = lpsRes.headers.get('set-cookie')
            if (rawSetCookie) {
              res.setHeader('Set-Cookie', rawSetCookie)
            }
          }

          res.statusCode = 200
          res.setHeader('Content-Type', 'application/json')
          res.end(JSON.stringify({ csrf, initialData }))
        } catch (err) {
          res.statusCode = 502
          res.setHeader('Content-Type', 'application/json')
          res.end(
            JSON.stringify({
              error: err instanceof Error ? err.message : 'LPS dev bootstrap failed',
            }),
          )
        }
      })
    },
  }
}

export const lpsDevProxyTarget = LPS_DEV_TARGET
