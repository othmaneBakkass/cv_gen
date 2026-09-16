import { Link, createFileRoute } from '@tanstack/react-router'
import { ArrowRight, Check } from 'lucide-react'

import { Button } from '#/components/ui/button'
import { Reveal } from '#/components/Reveal'
import { TemplatePreview } from '#/components/TemplatePreview'
import { templates } from '#/lib/templates'

export const Route = createFileRoute('/templates')({ component: TemplatesPage })

function TemplatesPage() {
  return (
    <main className="mx-auto max-w-6xl px-4 py-16">
      <Reveal className="mx-auto max-w-2xl text-center">
        <p className="mb-3 text-xs font-bold uppercase tracking-[0.16em] text-primary">Templates</p>
        <h1 className="text-balance text-4xl font-extrabold tracking-tight text-foreground">
          Pick a starting point
        </h1>
        <p className="mt-3 text-muted-foreground">
          All three render from the same content and the same engine — choose the
          look, then fine-tune colors, density and sections in the editor.
        </p>
      </Reveal>

      <div className="mt-12 flex flex-col gap-8">
        {templates.map((t, i) => (
          <Reveal key={t.id} delay={i * 70}>
            <article className="grid items-center gap-8 rounded-2xl border border-border bg-card p-6 sm:p-8 md:grid-cols-[minmax(0,220px)_1fr]">
              <div className="rounded-xl bg-secondary/50 p-6">
                <TemplatePreview id={t.id} />
              </div>
              <div>
                <div className="flex flex-wrap items-center gap-3">
                  <span className="size-3 rounded-full" style={{ background: t.accent }} />
                  <h2 className="text-2xl font-extrabold tracking-tight text-foreground">{t.name}</h2>
                  <span className="rounded-full border border-border px-2.5 py-0.5 text-xs font-medium text-muted-foreground">
                    {t.tagline}
                  </span>
                </div>
                <p className="mt-3 max-w-2xl leading-relaxed text-muted-foreground">{t.description}</p>

                <div className="mt-4 flex flex-wrap gap-2">
                  {t.sections.map((s) => (
                    <span
                      key={s}
                      className="rounded-md bg-secondary px-2.5 py-1 text-xs font-medium text-secondary-foreground"
                    >
                      {s}
                    </span>
                  ))}
                </div>

                {t.recommended && (
                  <p className="mt-4 inline-flex items-start gap-2 text-sm text-muted-foreground">
                    <Check className="mt-0.5 size-4 shrink-0 text-primary" />
                    {t.recommended}
                  </p>
                )}

                <div className="mt-6">
                  <Button asChild className="group">
                    <Link to="/build" search={{ template: t.id }}>
                      Use {t.name}
                      <ArrowRight className="size-4 transition-transform duration-200 group-hover:translate-x-0.5" />
                    </Link>
                  </Button>
                </div>
              </div>
            </article>
          </Reveal>
        ))}
      </div>
    </main>
  )
}
