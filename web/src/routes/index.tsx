import { useEffect, useState } from 'react'
import { Link, createFileRoute } from '@tanstack/react-router'
import {
  ArrowRight,
  Check,
  FileText,
  Languages,
  LayoutGrid,
  Palette,
  SlidersHorizontal,
  Sparkles,
  Wand2,
} from 'lucide-react'
import { cn } from '#/lib/utils'

import { Button } from '#/components/ui/button'
import { Reveal } from '#/components/Reveal'
import { TemplatePreview } from '#/components/TemplatePreview'
import { templates } from '#/lib/templates'

export const Route = createFileRoute('/')({ component: Landing })

function Landing() {
  return (
    <main className="flex flex-col">
      <Hero />
      <StatStrip />
      <HowItWorks />
      <TemplatesShowcase />
      <Features />
      <FinalCta />
    </main>
  )
}

function Hero() {
  return (
    <section className="relative overflow-hidden px-4 pb-20 pt-16 sm:pt-24">
      <div className="pointer-events-none absolute inset-0 -z-10 bg-grid" aria-hidden />
      <div className="pointer-events-none absolute inset-x-0 top-0 -z-10 h-[420px] glow" aria-hidden />

      <div className="mx-auto grid max-w-6xl items-center gap-12 lg:grid-cols-[1.05fr_0.95fr]">
        <div>
          <Reveal
            as="span"
            className="inline-flex items-center gap-2 rounded-full border border-border bg-background/70 px-3 py-1 text-xs font-semibold text-primary shadow-xs backdrop-blur-sm"
          >
            <Sparkles className="size-3.5" />
            Structured CVs, generated as real PDFs
          </Reveal>

          <Reveal
            as="h1"
            delay={60}
            className="mt-5 text-balance text-4xl font-extrabold leading-[1.05] tracking-tight text-foreground sm:text-5xl lg:text-6xl"
          >
            Your CV, typeset like a{' '}
            <span className="text-gradient">design system</span>, not a template.
          </Reveal>

          <Reveal as="p" delay={120} className="mt-5 max-w-xl text-lg leading-relaxed text-muted-foreground">
            Fill in a form or drop in JSON, pick a layout, and get a print-ready
            one-page PDF. Every section is an independent component — colors,
            spacing, dividers, and language all configurable, with the same engine
            behind the CLI and this editor.
          </Reveal>

          <Reveal delay={180} className="mt-8 flex flex-wrap items-center gap-3">
            <Button asChild size="lg" className="group">
              <Link to="/build">
                Build your CV
                <ArrowRight className="transition-transform duration-200 group-hover:translate-x-0.5" />
              </Link>
            </Button>
            <Button asChild size="lg" variant="outline">
              <Link to="/templates">Browse templates</Link>
            </Button>
          </Reveal>

          <Reveal delay={240} className="mt-6 flex flex-wrap gap-x-5 gap-y-2 text-sm text-muted-foreground">
            {['No account needed', 'ATS-safe output', 'Fits to one page'].map((t) => (
              <span key={t} className="inline-flex items-center gap-1.5">
                <Check className="size-4 text-primary" />
                {t}
              </span>
            ))}
          </Reveal>
        </div>

        <Reveal delay={120} className="relative">
          <HeroDeck />
        </Reveal>
      </div>
    </section>
  )
}

function usePrefersReducedMotion() {
  const [reduced, setReduced] = useState(false)
  useEffect(() => {
    const mq = window.matchMedia('(prefers-reduced-motion: reduce)')
    const update = () => setReduced(mq.matches)
    update()
    mq.addEventListener('change', update)
    return () => mq.removeEventListener('change', update)
  }, [])
  return reduced
}

// HeroDeck is the hero's interactive product shot: the three template
// previews as a 3D card deck. It auto-cycles when idle, pauses on hover,
// brings any card to the front on click, and gives the active card a subtle
// pointer-parallax tilt. All motion is disabled under prefers-reduced-motion,
// where it degrades to a static, still-clickable fan.
function HeroDeck() {
  const [active, setActive] = useState(1) // t2 (Navy Rule) centered by default
  const [tilt, setTilt] = useState({ x: 0, y: 0 })
  const [hovering, setHovering] = useState(false)
  const reduced = usePrefersReducedMotion()
  const count = templates.length

  useEffect(() => {
    if (reduced || hovering) return
    const id = setInterval(() => setActive((a) => (a + 1) % count), 3800)
    return () => clearInterval(id)
  }, [reduced, hovering, count])

  function handleMove(e: React.PointerEvent<HTMLDivElement>) {
    if (reduced) return
    const r = e.currentTarget.getBoundingClientRect()
    setTilt({
      x: (e.clientX - r.left) / r.width - 0.5,
      y: (e.clientY - r.top) / r.height - 0.5,
    })
  }

  return (
    <div className="mx-auto max-w-md">
      <div
        className="relative flex h-[380px] items-center justify-center [perspective:1200px] sm:h-[460px]"
        onPointerMove={handleMove}
        onPointerEnter={() => setHovering(true)}
        onPointerLeave={() => {
          setHovering(false)
          setTilt({ x: 0, y: 0 })
        }}
      >
        {/* soft glow that tracks the active card */}
        <div
          className="pointer-events-none absolute h-56 w-56 rounded-full bg-primary/20 blur-3xl transition-opacity duration-500"
          style={{ opacity: hovering ? 0.9 : 0.5 }}
          aria-hidden
        />
        {templates.map((t, i) => {
          const rel = (i - active + count) % count // 0 center, 1 right, 2 left
          const layout =
            rel === 0
              ? { x: 0, rot: 0, scale: 1, z: 30, op: 1, blur: 0 }
              : rel === 1
                ? { x: 40, rot: 8, scale: 0.82, z: 10, op: 0.5, blur: 1.2 }
                : { x: -40, rot: -8, scale: 0.82, z: 10, op: 0.5, blur: 1.2 }
          const tiltStr =
            rel === 0 && !reduced
              ? ` rotateX(${(-tilt.y * 7).toFixed(2)}deg) rotateY(${(tilt.x * 9).toFixed(2)}deg)`
              : ''
          return (
            <button
              key={t.id}
              type="button"
              aria-label={`Show ${t.name}`}
              aria-pressed={rel === 0}
              onClick={() => setActive(i)}
              className="absolute w-44 cursor-pointer origin-center rounded-[3px] transition-[transform,opacity,filter] duration-500 ease-[cubic-bezier(0.16,1,0.3,1)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 focus-visible:ring-offset-background sm:w-52"
              style={{
                transform: `translateX(${layout.x}%) rotate(${layout.rot}deg) scale(${layout.scale})${tiltStr}`,
                zIndex: layout.z,
                opacity: layout.op,
                filter: layout.blur ? `blur(${layout.blur}px)` : undefined,
                transformStyle: 'preserve-3d',
              }}
            >
              <TemplatePreview id={t.id} />
            </button>
          )
        })}
      </div>

      {/* Controls: dots + a contextual link to build with the active template. */}
      <div className="mt-5 flex items-center justify-center gap-4">
        <div className="flex items-center gap-1.5">
          {templates.map((t, i) => (
            <button
              key={t.id}
              type="button"
              aria-label={`Show ${t.name}`}
              onClick={() => setActive(i)}
              className={cn(
                'h-1.5 rounded-full transition-all duration-300',
                i === active ? 'w-6 bg-primary' : 'w-1.5 bg-border hover:bg-muted-foreground',
              )}
            />
          ))}
        </div>
        <Link
          to="/build"
          search={{ template: templates[active].id }}
          className="group inline-flex items-center gap-1 text-sm font-semibold text-primary no-underline"
        >
          Use {templates[active].name}
          <ArrowRight className="size-4 transition-transform duration-200 group-hover:translate-x-0.5" />
        </Link>
      </div>
    </div>
  )
}

function StatStrip() {
  const stats = [
    { n: '3', label: 'templates' },
    { n: '7', label: 'orderable sections' },
    { n: 'EN / FR', label: 'labels' },
    { n: '1', label: 'page, auto-fit' },
  ]
  return (
    <section className="border-y border-border bg-secondary/40 px-4 py-8">
      <div className="mx-auto grid max-w-5xl grid-cols-2 gap-6 sm:grid-cols-4">
        {stats.map((s, i) => (
          <Reveal key={s.label} delay={i * 60} className="text-center">
            <div className="text-2xl font-extrabold tracking-tight text-gradient sm:text-3xl">
              {s.n}
            </div>
            <div className="mt-1 text-xs font-medium uppercase tracking-wide text-muted-foreground">
              {s.label}
            </div>
          </Reveal>
        ))}
      </div>
    </section>
  )
}

function HowItWorks() {
  const steps = [
    {
      icon: FileText,
      title: 'Add your content',
      body: 'Type into a guided form or upload a JSON file shaped like the schema. Validation happens as you go.',
    },
    {
      icon: SlidersHorizontal,
      title: 'Configure the look',
      body: 'Choose a template, density and language, reorder or hide sections, and override any color role per section.',
    },
    {
      icon: Wand2,
      title: 'Generate the PDF',
      body: 'The same Go engine the CLI uses renders a print-ready, one-page A4 PDF — previewed inline, ready to download.',
    },
  ]
  return (
    <Section id="how" eyebrow="How it works" title="From details to a polished PDF in three steps">
      <div className="grid gap-6 md:grid-cols-3">
        {steps.map((s, i) => (
          <Reveal key={s.title} delay={i * 80}>
            <div className="group h-full rounded-xl border border-border bg-card p-6 transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-primary/5">
              <div className="flex size-11 items-center justify-center rounded-lg bg-accent text-accent-foreground transition-colors group-hover:bg-primary group-hover:text-primary-foreground">
                <s.icon className="size-5" />
              </div>
              <div className="mt-4 flex items-center gap-2">
                <span className="text-xs font-bold text-muted-foreground">
                  {String(i + 1).padStart(2, '0')}
                </span>
                <h3 className="text-base font-bold text-foreground">{s.title}</h3>
              </div>
              <p className="mt-2 text-sm leading-relaxed text-muted-foreground">{s.body}</p>
            </div>
          </Reveal>
        ))}
      </div>
    </Section>
  )
}

function TemplatesShowcase() {
  return (
    <Section
      id="templates"
      eyebrow="Templates"
      title="Three layouts, one engine"
      subtitle="Each is a theme over the same components — so your content stays identical and only the look changes."
    >
      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {templates.map((t, i) => (
          <Reveal key={t.id} delay={i * 80}>
            <Link
              to="/build"
              search={{ template: t.id }}
              className="group block h-full rounded-xl border border-border bg-card p-5 no-underline transition-all duration-300 hover:-translate-y-1 hover:border-primary/40 hover:shadow-xl hover:shadow-primary/5"
            >
              <div className="rounded-lg bg-secondary/50 p-5">
                <div className="transition-transform duration-300 group-hover:scale-[1.03]">
                  <TemplatePreview id={t.id} />
                </div>
              </div>
              <div className="mt-4 flex items-center gap-2">
                <span className="size-2.5 rounded-full" style={{ background: t.accent }} />
                <h3 className="text-base font-bold text-foreground">{t.name}</h3>
                <span className="ml-auto text-xs font-medium text-muted-foreground">{t.tagline}</span>
              </div>
              <p className="mt-2 text-sm leading-relaxed text-muted-foreground">{t.description}</p>
              <span className="mt-4 inline-flex items-center gap-1 text-sm font-semibold text-primary">
                Use this layout
                <ArrowRight className="size-4 transition-transform duration-200 group-hover:translate-x-0.5" />
              </span>
            </Link>
          </Reveal>
        ))}
      </div>
    </Section>
  )
}

function Features() {
  const features = [
    {
      icon: LayoutGrid,
      title: 'Sections as components',
      body: 'Reorder, show or hide any of seven sections. Each inherits global defaults and can override them on its own.',
    },
    {
      icon: Palette,
      title: 'Semantic color roles',
      body: 'Recolor headline, body, metadata, link, border and more — globally or per section — never a fixed palette.',
    },
    {
      icon: SlidersHorizontal,
      title: 'Density & margins',
      body: 'Dense, normal or airy scales both the spacing and the page margins together, matching the reference designs.',
    },
    {
      icon: Languages,
      title: 'English or French labels',
      body: 'Template labels switch language without touching your content — headings, "Stack:", date ranges and more.',
    },
    {
      icon: FileText,
      title: 'ATS-safe, one page',
      body: 'Single-column, real text in reading order, embedded fonts — and spacing that tightens to hold one page.',
    },
    {
      icon: Wand2,
      title: 'CLI or browser',
      body: 'The identical engine powers the command line and this editor. Export JSON for CI, or generate a PDF here.',
    },
  ]
  return (
    <Section
      id="features"
      eyebrow="Why it's different"
      title="Configurable down to the section, not just the template"
    >
      <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
        {features.map((f, i) => (
          <Reveal key={f.title} delay={(i % 3) * 70}>
            <div className="h-full rounded-xl border border-border bg-card p-6 transition-colors duration-300 hover:border-primary/30">
              <f.icon className="size-5 text-primary" />
              <h3 className="mt-3 text-base font-bold text-foreground">{f.title}</h3>
              <p className="mt-2 text-sm leading-relaxed text-muted-foreground">{f.body}</p>
            </div>
          </Reveal>
        ))}
      </div>
    </Section>
  )
}

function FinalCta() {
  return (
    <section className="px-4 py-20">
      <Reveal className="relative mx-auto max-w-4xl overflow-hidden rounded-2xl border border-border bg-card px-6 py-14 text-center">
        <div className="pointer-events-none absolute inset-0 -z-10 glow" aria-hidden />
        <h2 className="mx-auto max-w-2xl text-balance text-3xl font-extrabold tracking-tight text-foreground sm:text-4xl">
          Build a CV you'd actually want to send.
        </h2>
        <p className="mx-auto mt-3 max-w-xl text-muted-foreground">
          No sign-up, no watermark. Fill it in, tune the look, and download the PDF.
        </p>
        <div className="mt-8 flex flex-wrap items-center justify-center gap-3">
          <Button asChild size="lg" className="group">
            <Link to="/build">
              Start building
              <ArrowRight className="transition-transform duration-200 group-hover:translate-x-0.5" />
            </Link>
          </Button>
          <Button asChild size="lg" variant="outline">
            <Link to="/templates">See the templates</Link>
          </Button>
        </div>
      </Reveal>
    </section>
  )
}

// Section is the shared vertical-rhythm wrapper for the landing page's bands.
function Section({
  id,
  eyebrow,
  title,
  subtitle,
  children,
}: {
  id?: string
  eyebrow: string
  title: string
  subtitle?: string
  children: React.ReactNode
}) {
  return (
    <section id={id} className="px-4 py-20">
      <div className="mx-auto max-w-6xl">
        <Reveal className="mx-auto mb-12 max-w-2xl text-center">
          <p className="mb-3 text-xs font-bold uppercase tracking-[0.16em] text-primary">{eyebrow}</p>
          <h2 className="text-balance text-3xl font-extrabold tracking-tight text-foreground sm:text-4xl">
            {title}
          </h2>
          {subtitle && <p className="mt-3 text-muted-foreground">{subtitle}</p>}
        </Reveal>
        {children}
      </div>
    </section>
  )
}
