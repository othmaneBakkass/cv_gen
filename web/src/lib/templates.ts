import type { TemplateId } from './cv-schema'

export type { TemplateId }

export interface TemplateMeta {
  id: TemplateId
  name: string
  tagline: string
  description: string
  sections: Array<string>
  /** Accent color used by the Go theme, for the preview swatch/mock. */
  accent: string
  recommended?: string
}

// Mirrors internal/pdf/templates' registry and internal/pdf/theme (T1/T2/T3).
// Static (not fetched) so the marketing pages render instantly; the running
// engine exposes the same list at GET /api/templates for parity checks.
export const templates: Array<TemplateMeta> = [
  {
    id: 't1',
    name: 'ATS Classic',
    tagline: 'Sober, single-column, ATS-safe',
    description:
      'Conservative navy on a clean single column with upright dates. Built to parse cleanly through applicant-tracking software and read as a sober professional CV.',
    sections: ['Profile', 'Experience', 'Education', 'Skills', 'Projects', 'Certifications', 'Languages'],
    accent: '#1F3864',
    recommended: 'Best when the CV goes through an ATS or a conservative recruiter first.',
  },
  {
    id: 't2',
    name: 'Navy Rule',
    tagline: 'Editorial, teal-accented',
    description:
      'Editorial and teal-accented, with tracked uppercase headings, thin rules under every section, and italic right-aligned dates.',
    sections: ['Profile', 'Education', 'Experience', 'Skills', 'Projects', 'Certifications', 'Languages'],
    accent: '#1F5673',
    recommended: 'A polished default when you want structure without looking like a template.',
  },
  {
    id: 't3',
    name: 'Modern Minimal',
    tagline: 'Compact, indigo, rule-free',
    description:
      'Compact and indigo-accented, with title-case headings, no rule lines, and an en-dash bullet — a contemporary, quieter look that leans on space and weight.',
    sections: ['Profile', 'Education', 'Experience', 'Skills', 'Projects', 'Certifications', 'Languages'],
    accent: '#3730A3',
    recommended: 'For product, design, and modern engineering roles where a lighter look fits.',
  },
]

export function getTemplate(id: string | undefined): TemplateMeta | undefined {
  return templates.find((t) => t.id === id)
}
