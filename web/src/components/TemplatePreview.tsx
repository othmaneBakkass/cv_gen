import { cn } from '#/lib/utils'
import type { TemplateId } from '#/lib/templates'

// A miniature, CSS-drawn stand-in for each template's actual PDF layout. It
// mirrors the real styling decisions each Go theme makes (internal/pdf/theme
// T1/T2/T3): accent hue, whether headings carry a rule, upper vs title case,
// the photo box, and the bullet marker — so the preview reads as that
// template rather than a generic placeholder.

interface Spec {
  accent: string
  rule: 'none' | 'thin' | 'full' // none = t3, thin = t1 grey, full = t2 accent
  uppercase: boolean
  photo: boolean
  marker: 'dot' | 'dash'
}

const specs: Record<TemplateId, Spec> = {
  t1: { accent: '#1F3864', rule: 'thin', uppercase: true, photo: false, marker: 'dot' },
  t2: { accent: '#1F5673', rule: 'full', uppercase: true, photo: true, marker: 'dot' },
  t3: { accent: '#3730A3', rule: 'none', uppercase: false, photo: false, marker: 'dash' },
}

export function TemplatePreview({ id, className }: { id: TemplateId; className?: string }) {
  const s = specs[id]
  return (
    <div
      className={cn(
        'relative mx-auto aspect-[210/297] w-full max-w-56 overflow-hidden rounded-[3px] bg-[#fdfdfc] p-4',
        'ring-1 ring-black/5',
        'shadow-[0_1px_1px_rgba(26,26,26,0.05),0_18px_36px_-14px_rgba(26,26,26,0.30)]',
        className,
      )}
    >
      {/* Header */}
      <div className="flex items-start justify-between">
        <div className="flex flex-col gap-[5px]">
          <Bar w="52px" h={5} color={s.accent} />
          <Bar w="70px" h={2.5} color="#1a1a1a" />
          <Bar w="58px" h={2} color="#8a969e" />
        </div>
        {s.photo && <div className="h-9 w-7 rounded-[2px] bg-[#e7eef1] ring-1 ring-black/5" />}
      </div>
      {s.rule !== 'none' && (
        <div
          className="mt-2 w-full"
          style={{ height: 1.5, background: s.accent }}
        />
      )}

      {/* Sections */}
      {[
        { label: 26, lines: [92, 74] },
        { label: 32, lines: [86, 62, 80] },
        { label: 22, lines: [70, 54] },
      ].map((section, i) => (
        <div key={i} className="mt-3">
          <Bar w={`${section.label}px`} h={2.5} color={s.accent} />
          {s.rule === 'full' && <div className="mt-1 h-px w-full" style={{ background: s.accent }} />}
          {s.rule === 'thin' && <div className="mt-1 h-px w-full bg-[#e2e5e7]" />}
          <div className="mt-2 flex flex-col gap-[5px]">
            {section.lines.map((pct, j) => (
              <div key={j} className="flex items-center gap-1">
                <span
                  className="inline-block shrink-0"
                  style={{
                    width: 3,
                    height: s.marker === 'dot' ? 3 : 1.5,
                    borderRadius: s.marker === 'dot' ? 9999 : 0,
                    background: '#8a969e',
                  }}
                />
                <Bar w={`${pct}%`} h={2} color="#5a6b75" />
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  )
}

function Bar({ w, h, color }: { w: string; h: number; color: string }) {
  return <div className="rounded-[1px]" style={{ width: w, height: h, background: color }} />
}
