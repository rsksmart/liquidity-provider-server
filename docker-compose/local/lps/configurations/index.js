const { CONSTANTS } = require('./constants');
const { CONFIG_REQUESTS } = require('./requests');


const cookieJar = {};

async function fetchWithCookies(url, cookies, options = {}) {
    // Attach stored cookies to the request
    options.headers = {
        ...options.headers,
        ...(Object.keys(cookies).length && {
            Cookie: Object.entries(cookies)
                .map(([k, v]) => `${k}=${v}`)
                .join("; "),
        }),
    };
    const res = await fetch(url, options);

    // Save cookies from the response
    for (const [name, value] of (res.headers.getSetCookie?.() ?? []).map((c) =>
        c.split(";")[0].split("=")
    )) {
        cookies[name.trim()] = value.trim();
    }
    return res;
}

function extractCsrfFromHtml(html) {
    const m = html.match(/name=["']csrf["'][^>]*\bvalue=["']([^"']+)["']/i);
    return m?.[1] ?? null;
}

function getCsrfToken(cookies) {
    return fetchWithCookies(CONSTANTS.LPS_URL+"/management", cookies)
        .then(res => {
            if (!res.ok) {
                throw new Error(`Failed to fetch management page: ${res.status} ${res.statusText}`);
            }
            return res.text();
        })
        .then(html => {
            const csrf = extractCsrfFromHtml(html);
            if (!csrf) {
                throw new Error("CSRF token not found in management page HTML");
            }
            return csrf.replaceAll('&#43;', '+')
        })
}

(async function(){
   const csrf = await getCsrfToken(cookieJar)
   for (const [key, value] of Object.entries(CONFIG_REQUESTS)) {
       console.log(`Sending ${key} request...`)
       const res = await fetchWithCookies(CONSTANTS.LPS_URL+value.path, cookieJar, {
           method: value.method,
           headers: {
               'Content-Type': 'application/json',
               'X-CSRF-Token': csrf
           },
           body: JSON.stringify(value.body)
       })
       if (!res.ok) {
           console.error(`Failed to send ${key} request: ${res.status} ${res.statusText}`);
           const errorText = await res.text();
           console.error("Response body:", errorText);
       } else {
           console.log(`${key} request successful`);
       }
   }
})()
