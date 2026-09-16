import type { CV } from './cv-schema'

// Client for the Go engine's HTTP API (cmd/serve). In dev, Vite proxies
// /api to the Go server (see vite.config.ts); in a co-hosted build the API
// is same-origin, so a relative base works in both cases.

export interface ApiError {
  error: string
  detail?: string
  issues?: Array<string>
}

export class GenerateError extends Error {
  issues: Array<string>
  constructor(message: string, issues: Array<string> = []) {
    super(message)
    this.name = 'GenerateError'
    this.issues = issues
  }
}

export interface TemplateInfo {
  id: string
  name: string
  description: string
}

export async function fetchTemplates(signal?: AbortSignal): Promise<Array<TemplateInfo>> {
  const res = await fetch('/api/templates', { signal })
  if (!res.ok) throw new GenerateError('Could not load templates')
  const body = (await res.json()) as { templates: Array<TemplateInfo> }
  return body.templates
}

export interface GenerateResult {
  blob: Blob
  url: string
  fileName: string
}

// generatePdf posts a CV to the engine and returns the PDF as a blob URL
// (revoke it when you're done, e.g. on unmount / next generate). Throws a
// GenerateError carrying any field-level issues on a 4xx/5xx.
export async function generatePdf(cv: CV, signal?: AbortSignal): Promise<GenerateResult> {
  const res = await fetch('/api/generate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(cv),
    signal,
  })

  if (!res.ok) {
    let msg = `Generation failed (${res.status})`
    let issues: Array<string> = []
    try {
      const body = (await res.json()) as ApiError
      if (body.error) msg = body.detail ? `${body.error}: ${body.detail}` : body.error
      if (body.issues) issues = body.issues
    } catch {
      // non-JSON error body — keep the status message
    }
    throw new GenerateError(msg, issues)
  }

  const blob = await res.blob()
  const header = res.headers.get('X-Filename')
  const fileName = header && header.length > 0 ? header : `${cv.fileName || 'cv'}.pdf`
  return { blob, url: URL.createObjectURL(blob), fileName }
}
