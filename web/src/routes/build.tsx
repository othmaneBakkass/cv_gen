import { useEffect, useRef, useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { useForm } from '@tanstack/react-form'
import { toast } from 'sonner'
import {
  ArrowDown,
  ArrowUp,
  Check,
  ChevronLeft,
  ChevronRight,
  Download,
  FileText,
  Loader2,
  RotateCcw,
} from 'lucide-react'
import { cn } from '#/lib/utils'

import { Button } from '#/components/ui/button'
import { Textarea } from '#/components/ui/textarea'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs'
import { Card, CardContent } from '#/components/ui/card'
import { Separator } from '#/components/ui/separator'
import { Switch } from '#/components/ui/switch'
import { Label } from '#/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select'
import { StringListField } from '#/components/cv-form/StringListField'
import {
  EducationSection,
  JobsSection,
  LanguagesSection,
  ProjectsSection,
  SkillsSection,
  TextInput,
} from '#/components/cv-form/sections'
import {
  colorRoles,
  cvSchema,
  emptyCV,
  exampleCV,
  sectionKeys,
  templateId,
} from '#/lib/cv-schema'
import type { CV, ColorSettings, SectionKey, Settings } from '#/lib/cv-schema'
import { templates } from '#/lib/templates'
import { generatePdf, GenerateError } from '#/lib/api'

export const Route = createFileRoute('/build')({
  validateSearch: (search) => ({
    template: templateId.optional().parse(search.template),
  }),
  component: BuildPage,
})

// --- Editor settings (separate from the form so reorder/toggles are fully
// controllable) ---------------------------------------------------------

const REQUIRED: ReadonlyArray<SectionKey> = ['education', 'experience', 'languages']

const sectionLabels: Record<SectionKey, string> = {
  profile: 'Profile',
  education: 'Education',
  experience: 'Experience',
  skills: 'Skills',
  projects: 'Projects',
  certifications: 'Certifications',
  languages: 'Languages',
}

interface SectionRow {
  key: SectionKey
  include: boolean
  divider: boolean // false hides this section's rule; true keeps the template default
}

interface EditorSettings {
  density: 'dense' | 'normal' | 'airy'
  lang: 'en' | 'fr'
  order: Array<SectionRow>
  colors: ColorSettings
}

function defaultEditorSettings(): EditorSettings {
  return {
    density: 'normal',
    lang: 'en',
    order: sectionKeys.map((key) => ({ key, include: true, divider: true })),
    colors: {},
  }
}

// Convert the editor's settings into the schema.Settings the API accepts.
function toSettings(s: EditorSettings): Settings {
  const sections: Settings['sections'] = {}
  s.order.forEach((row, i) => {
    const entry: NonNullable<Settings['sections']>[SectionKey] = { priority: (i + 1) * 10 }
    if (!row.include && !REQUIRED.includes(row.key)) entry.include = false
    if (!row.divider) entry.style = { divider: false }
    sections[row.key] = entry
  })

  const colors: ColorSettings = {}
  for (const [k, v] of Object.entries(s.colors)) {
    if (v && v.length > 0) colors[k as keyof ColorSettings] = v
  }

  return {
    lang: s.lang === 'fr' ? 'fr' : undefined,
    density: s.density !== 'normal' ? s.density : undefined,
    sections,
    colors: Object.keys(colors).length ? colors : undefined,
  }
}

// Best-effort hydrate of the editor settings from an uploaded CV's settings.
function fromSettings(s: Settings | undefined): EditorSettings {
  const base = defaultEditorSettings()
  if (!s) return base
  if (s.density) base.density = s.density
  if (s.lang) base.lang = s.lang
  if (s.colors) base.colors = { ...s.colors }
  if (s.sections) {
    const rows = [...base.order]
    rows.sort((a, b) => {
      const pa = s.sections?.[a.key]?.priority ?? Number.MAX_SAFE_INTEGER
      const pb = s.sections?.[b.key]?.priority ?? Number.MAX_SAFE_INTEGER
      return pa - pb
    })
    base.order = rows.map((r) => ({
      key: r.key,
      include: s.sections?.[r.key]?.include ?? true,
      divider: s.sections?.[r.key]?.style?.divider ?? true,
    }))
  }
  return base
}

// --- Generation state ---------------------------------------------------

type GenState =
  | { status: 'idle' }
  | { status: 'loading' }
  | { status: 'done'; url: string; fileName: string }
  | { status: 'error'; message: string; issues: Array<string> }

function download(filename: string, contents: string, type = 'application/json') {
  const blob = new Blob([contents], { type })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

function BuildPage() {
  const search = Route.useSearch()
  const [mode, setMode] = useState<'form' | 'upload'>('form')
  const [uploadError, setUploadError] = useState<string | null>(null)
  const [settings, setSettings] = useState<EditorSettings>(defaultEditorSettings)
  const [gen, setGen] = useState<GenState>({ status: 'idle' })
  const fileInputRef = useRef<HTMLInputElement>(null)
  const lastUrl = useRef<string | null>(null)

  const form = useForm({
    defaultValues: { ...emptyCV, template: search.template ?? 't1' } as CV,
    validators: { onChange: cvSchema },
    onSubmit: ({ value }) => generate(value),
    onSubmitInvalid: () => toast.error('Some fields need fixing before generating.'),
  })

  // Revoke the previous blob URL whenever we replace or unmount.
  useEffect(() => () => {
    if (lastUrl.current) URL.revokeObjectURL(lastUrl.current)
  }, [])

  function assemble(value: CV): CV {
    return { ...value, settings: toSettings(settings) }
  }

  async function generate(value: CV) {
    setGen({ status: 'loading' })
    try {
      const result = await generatePdf(assemble(value))
      if (lastUrl.current) URL.revokeObjectURL(lastUrl.current)
      lastUrl.current = result.url
      setGen({ status: 'done', url: result.url, fileName: result.fileName })
      toast.success('PDF generated.')
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Generation failed.'
      const issues = err instanceof GenerateError ? err.issues : []
      setGen({ status: 'error', message, issues })
      toast.error('Could not generate the PDF.')
    }
  }

  function handleFile(file: File) {
    setUploadError(null)
    file
      .text()
      .then((text) => {
        const raw = JSON.parse(text)
        const body = raw && typeof raw === 'object' && 'data' in raw ? raw.data?.[0] : raw
        const parsed = cvSchema.safeParse(body)
        if (!parsed.success) {
          setUploadError(
            parsed.error.issues.map((i) => `${i.path.join('.')}: ${i.message}`).join('; '),
          )
          return
        }
        form.reset(parsed.data)
        setSettings(fromSettings(parsed.data.settings))
        setMode('form')
        toast.success('Loaded. Review it, then generate.')
      })
      .catch(() => setUploadError('That file is not valid JSON.'))
  }

  return (
    <main className="mx-auto max-w-6xl px-4 py-10">
      <div className="mb-8">
        <h1 className="text-3xl font-extrabold tracking-tight text-foreground">Build your CV</h1>
        <p className="mt-2 max-w-2xl text-muted-foreground">
          Fill in the form or upload JSON, tune the look on the right, then
          generate a print-ready PDF — no account needed.
        </p>
      </div>

      <div className="grid gap-8 lg:grid-cols-[1fr_minmax(0,400px)]">
        {/* Left: content editor */}
        <div className="min-w-0">
          <Tabs value={mode} onValueChange={(v) => setMode(v as 'form' | 'upload')}>
            <TabsList>
              <TabsTrigger value="form">Fill in the form</TabsTrigger>
              <TabsTrigger value="upload">Upload JSON</TabsTrigger>
            </TabsList>

            <TabsContent value="upload" className="mt-6">
              <UploadTab
                fileInputRef={fileInputRef}
                onFile={handleFile}
                uploadError={uploadError}
                template={form.state.values.template}
              />
            </TabsContent>

            <TabsContent value="form" className="mt-6">
              <Wizard form={form} />
            </TabsContent>
          </Tabs>
        </div>

        {/* Right: settings + generate */}
        <aside className="flex flex-col gap-6 lg:sticky lg:top-20 lg:h-fit">
          <SettingsPanel form={form} settings={settings} setSettings={setSettings} />
          <GeneratePanel
            gen={gen}
            onGenerate={() => form.handleSubmit()}
            onExportJson={() =>
              download(
                `${form.state.values.fileName || 'cv'}.json`,
                JSON.stringify({ data: [assemble(form.state.values)] }, null, 2),
              )
            }
          />
        </aside>
      </div>
    </main>
  )
}

// --- Content form -------------------------------------------------------

const WIZARD_STEPS = [
  { title: 'Contact', hint: 'Who you are' },
  { title: 'Experience', hint: 'Profile & roles' },
  { title: 'Education', hint: 'Studies & skills' },
  { title: 'Extras', hint: 'Projects, certs, languages' },
] as const

// Wizard turns the CV form into a guided, multi-step flow. Values live in the
// TanStack form (not the step components), so navigating between steps never
// loses input, and Generate still validates the whole CV at once.
function Wizard({ form }: { form: any }) {
  const [step, setStep] = useState(0)
  const [dir, setDir] = useState<1 | -1>(1)
  const last = WIZARD_STEPS.length - 1

  function go(next: number) {
    if (next < 0 || next > last || next === step) return
    setDir(next > step ? 1 : -1)
    setStep(next)
    if (typeof window !== 'undefined') window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault()
        form.handleSubmit()
      }}
    >
      <Stepper step={step} onStep={go} />

      <div
        key={step}
        className={cn(
          'mt-8 space-y-10 duration-300 animate-in fade-in',
          dir === 1 ? 'slide-in-from-right-4' : 'slide-in-from-left-4',
        )}
      >
        {step === 0 && <ContactStep form={form} />}
        {step === 1 && <ExperienceStep form={form} />}
        {step === 2 && <EducationStep form={form} />}
        {step === 3 && <ExtrasStep form={form} />}
      </div>

      <div className="mt-10 flex items-center justify-between border-t border-border pt-6">
        <Button type="button" variant="ghost" onClick={() => go(step - 1)} disabled={step === 0}>
          <ChevronLeft /> Back
        </Button>
        {step < last ? (
          <Button type="button" onClick={() => go(step + 1)}>
            Next: {WIZARD_STEPS[step + 1].title} <ChevronRight />
          </Button>
        ) : (
          <Button type="submit">
            <FileText /> Generate PDF
          </Button>
        )}
      </div>
    </form>
  )
}

function Stepper({ step, onStep }: { step: number; onStep: (n: number) => void }) {
  return (
    <div>
      <p className="mb-3 text-sm font-medium text-muted-foreground sm:hidden">
        Step {step + 1} of {WIZARD_STEPS.length} —{' '}
        <span className="text-foreground">{WIZARD_STEPS[step].title}</span>
      </p>
      <ol className="hidden gap-2 sm:grid sm:grid-cols-4">
        {WIZARD_STEPS.map((s, i) => {
          const state = i < step ? 'done' : i === step ? 'current' : 'todo'
          return (
            <li key={s.title}>
              <button
                type="button"
                onClick={() => onStep(i)}
                className="group w-full text-left focus-visible:outline-none"
              >
                <div className="flex items-center gap-2">
                  <span
                    className={cn(
                      'flex size-6 shrink-0 items-center justify-center rounded-full text-xs font-bold transition-colors',
                      state === 'todo'
                        ? 'border border-border text-muted-foreground group-hover:border-primary/50'
                        : 'bg-primary text-primary-foreground',
                    )}
                  >
                    {state === 'done' ? <Check className="size-3.5" /> : i + 1}
                  </span>
                  <span
                    className={cn(
                      'text-sm font-semibold transition-colors',
                      state === 'todo' ? 'text-muted-foreground' : 'text-foreground',
                    )}
                  >
                    {s.title}
                  </span>
                </div>
                <div
                  className={cn(
                    'mt-2 h-1 rounded-full transition-colors',
                    state === 'todo' ? 'bg-border' : 'bg-primary',
                  )}
                />
              </button>
            </li>
          )
        })}
      </ol>
    </div>
  )
}

function ContactStep({ form }: { form: any }) {
  return (
    <FormBlock title="Contact" hint="The header block at the top of your CV.">
      <div className="grid gap-4 sm:grid-cols-2">
        <TextInput form={form} name="fileName" label="Output file name" placeholder="jane-doe" />
        <TextInput form={form} name="head.fullName" label="Full name" />
        <TextInput form={form} name="head.jobTitle" label="Job title" />
        <TextInput form={form} name="head.email" label="Email" />
        <TextInput form={form} name="head.phone" label="Phone" />
        <TextInput form={form} name="head.address" label="Address" />
        <TextInput form={form} name="head.linkedin" label="LinkedIn" />
        <TextInput form={form} name="head.github" label="GitHub" />
      </div>
      <StringListField
        form={form}
        name="head.specialties"
        label="Specialties"
        placeholder="e.g. Distributed Systems"
        addLabel="Add specialty"
      />
    </FormBlock>
  )
}

function ExperienceStep({ form }: { form: any }) {
  return (
    <>
      <section className="space-y-3">
        <div>
          <h2 id="profile-heading" className="text-lg font-bold text-foreground">
            Profile
          </h2>
          <p className="text-sm text-muted-foreground">A short summary — optional, but a strong opener.</p>
        </div>
        <form.Field name="profile">
          {(field: any) => (
            <Textarea
              aria-labelledby="profile-heading"
              rows={3}
              placeholder="A short summary of who you are and what you bring."
              value={field.state.value ?? ''}
              onBlur={field.handleBlur}
              onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => field.handleChange(e.target.value)}
            />
          )}
        </form.Field>
      </section>
      <FormBlock title="Experience" hint="Your roles, most recent first.">
        <JobsSection form={form} />
      </FormBlock>
    </>
  )
}

function EducationStep({ form }: { form: any }) {
  return (
    <>
      <FormBlock title="Education" hint="Degrees and programs.">
        <EducationSection form={form} />
      </FormBlock>
      <FormBlock title="Skills" hint="Group related skills into categories.">
        <SkillsSection form={form} />
      </FormBlock>
    </>
  )
}

function ExtrasStep({ form }: { form: any }) {
  return (
    <>
      <FormBlock title="Projects" hint="Optional — notable things you've built.">
        <ProjectsSection form={form} />
      </FormBlock>
      <section className="space-y-3">
        <h2 className="text-lg font-bold text-foreground">Certifications</h2>
        <StringListField
          form={form}
          name="certifications"
          label="Certifications"
          placeholder="Title — Issuing organization"
          addLabel="Add certification"
        />
      </section>
      <FormBlock title="Languages" hint="Languages and proficiency.">
        <LanguagesSection form={form} />
      </FormBlock>
    </>
  )
}

function FormBlock({
  title,
  hint,
  children,
}: {
  title: string
  hint?: string
  children: React.ReactNode
}) {
  return (
    <section className="space-y-4">
      <div>
        <h2 className="text-lg font-bold text-foreground">{title}</h2>
        {hint && <p className="text-sm text-muted-foreground">{hint}</p>}
      </div>
      {children}
    </section>
  )
}

function UploadTab({
  fileInputRef,
  onFile,
  uploadError,
  template,
}: {
  fileInputRef: React.RefObject<HTMLInputElement | null>
  onFile: (f: File) => void
  uploadError: string | null
  template: CV['template']
}) {
  return (
    <Card>
      <CardContent className="space-y-4 pt-6">
        <p className="text-sm text-muted-foreground">
          Download a JSON file shaped like the schema, fill it in, then upload it
          back here — settings and all.
        </p>
        <div className="flex flex-wrap gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() =>
              download('cv-blank.json', JSON.stringify({ ...emptyCV, template }, null, 2))
            }
          >
            Download blank JSON
          </Button>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => download('cv-example.json', JSON.stringify(exampleCV, null, 2))}
          >
            Download filled example
          </Button>
        </div>
        <Separator />
        <div>
          <input
            ref={fileInputRef}
            type="file"
            accept="application/json"
            className="hidden"
            onChange={(e) => {
              const file = e.target.files?.[0]
              if (file) onFile(file)
              e.target.value = ''
            }}
          />
          <Button type="button" onClick={() => fileInputRef.current?.click()}>
            Choose a JSON file
          </Button>
          {uploadError && <p className="mt-3 text-sm text-destructive">{uploadError}</p>}
        </div>
      </CardContent>
    </Card>
  )
}

// --- Settings panel -----------------------------------------------------

function SettingsPanel({
  form,
  settings,
  setSettings,
}: {
  form: any
  settings: EditorSettings
  setSettings: React.Dispatch<React.SetStateAction<EditorSettings>>
}) {
  function move(index: number, dir: -1 | 1) {
    setSettings((s) => {
      const order = [...s.order]
      const j = index + dir
      if (j < 0 || j >= order.length) return s
      ;[order[index], order[j]] = [order[j], order[index]]
      return { ...s, order }
    })
  }

  return (
    <Card>
      <CardContent className="space-y-6 pt-6">
        <div>
          <h2 className="text-sm font-bold uppercase tracking-wide text-muted-foreground">Design</h2>
        </div>

        {/* Template */}
        <div className="space-y-2">
          <Label>Template</Label>
          <form.Field name="template">
            {(field: any) => (
              <div className="grid grid-cols-3 gap-2">
                {templates.map((t) => (
                  <button
                    key={t.id}
                    type="button"
                    onClick={() => field.handleChange(t.id)}
                    className={
                      'flex flex-col items-start gap-1 rounded-lg border p-2.5 text-left transition-all ' +
                      (field.state.value === t.id
                        ? 'border-primary bg-accent/60 ring-1 ring-primary'
                        : 'border-border hover:border-primary/40')
                    }
                  >
                    <span className="size-2.5 rounded-full" style={{ background: t.accent }} />
                    <span className="text-xs font-semibold text-foreground">{t.name}</span>
                  </button>
                ))}
              </div>
            )}
          </form.Field>
        </div>

        {/* Density + language */}
        <div className="grid grid-cols-2 gap-3">
          <div className="space-y-2">
            <Label>Density</Label>
            <Select
              value={settings.density}
              onValueChange={(v) => setSettings((s) => ({ ...s, density: v as EditorSettings['density'] }))}
            >
              <SelectTrigger className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="dense">Dense</SelectItem>
                <SelectItem value="normal">Normal</SelectItem>
                <SelectItem value="airy">Airy</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Labels</Label>
            <Select
              value={settings.lang}
              onValueChange={(v) => setSettings((s) => ({ ...s, lang: v as EditorSettings['lang'] }))}
            >
              <SelectTrigger className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="en">English</SelectItem>
                <SelectItem value="fr">Français</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <Separator />

        {/* Section order / visibility / divider */}
        <div className="space-y-2">
          <Label>Sections</Label>
          <p className="text-xs text-muted-foreground">
            Reorder, hide, or drop the divider per section.
          </p>
          <ul className="space-y-1.5">
            {settings.order.map((row, i) => (
              <li
                key={row.key}
                className="flex items-center gap-2 rounded-md border border-border bg-background px-2 py-1.5"
              >
                <div className="flex flex-col">
                  <button
                    type="button"
                    aria-label={`Move ${sectionLabels[row.key]} up`}
                    disabled={i === 0}
                    onClick={() => move(i, -1)}
                    className="text-muted-foreground transition-colors hover:text-foreground disabled:opacity-30"
                  >
                    <ArrowUp className="size-3" />
                  </button>
                  <button
                    type="button"
                    aria-label={`Move ${sectionLabels[row.key]} down`}
                    disabled={i === settings.order.length - 1}
                    onClick={() => move(i, 1)}
                    className="text-muted-foreground transition-colors hover:text-foreground disabled:opacity-30"
                  >
                    <ArrowDown className="size-3" />
                  </button>
                </div>
                <span className="flex-1 text-sm font-medium text-foreground">
                  {sectionLabels[row.key]}
                </span>
                <label className="flex items-center gap-1 text-[11px] text-muted-foreground">
                  rule
                  <Switch
                    size="sm"
                    checked={row.divider}
                    onCheckedChange={(v) =>
                      setSettings((s) => ({
                        ...s,
                        order: s.order.map((r) => (r.key === row.key ? { ...r, divider: v } : r)),
                      }))
                    }
                  />
                </label>
                <Switch
                  size="sm"
                  aria-label={`Show ${sectionLabels[row.key]}`}
                  checked={REQUIRED.includes(row.key) ? true : row.include}
                  disabled={REQUIRED.includes(row.key)}
                  onCheckedChange={(v) =>
                    setSettings((s) => ({
                      ...s,
                      order: s.order.map((r) => (r.key === row.key ? { ...r, include: v } : r)),
                    }))
                  }
                />
              </li>
            ))}
          </ul>
        </div>

        <Separator />

        {/* Colors */}
        <div className="space-y-2">
          <Label>Colors</Label>
          <p className="text-xs text-muted-foreground">
            Override any role; leave blank to keep the template default.
          </p>
          <div className="grid grid-cols-2 gap-2">
            {colorRoles.map((role) => (
              <ColorField
                key={role.key}
                label={role.label}
                hint={role.hint}
                value={settings.colors[role.key] ?? ''}
                onChange={(v) =>
                  setSettings((s) => ({ ...s, colors: { ...s.colors, [role.key]: v } }))
                }
              />
            ))}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

function ColorField({
  label,
  hint,
  value,
  onChange,
}: {
  label: string
  hint: string
  value: string
  onChange: (v: string) => void
}) {
  const swatch = value && /^#?[0-9a-fA-F]{6}$/.test(value) ? (value.startsWith('#') ? value : `#${value}`) : '#ffffff'
  return (
    <div className="flex items-center gap-2 rounded-md border border-border bg-background px-2 py-1.5" title={hint}>
      <input
        type="color"
        aria-label={`${label} color`}
        value={swatch}
        onChange={(e) => onChange(e.target.value)}
        className="size-6 shrink-0 cursor-pointer rounded border-0 bg-transparent p-0"
      />
      <div className="min-w-0 flex-1">
        <div className="truncate text-[11px] font-medium text-foreground">{label}</div>
      </div>
      {value ? (
        <button
          type="button"
          aria-label={`Clear ${label}`}
          onClick={() => onChange('')}
          className="text-[11px] text-muted-foreground hover:text-destructive"
        >
          clear
        </button>
      ) : null}
    </div>
  )
}

// --- Generate panel -----------------------------------------------------

function GeneratePanel({
  gen,
  onGenerate,
  onExportJson,
}: {
  gen: GenState
  onGenerate: () => void
  onExportJson: () => void
}) {
  return (
    <Card className="overflow-hidden">
      <CardContent className="space-y-4 pt-6">
        <div className="flex items-center gap-2">
          <FileText className="size-4 text-primary" />
          <h2 className="text-sm font-bold uppercase tracking-wide text-muted-foreground">Preview</h2>
        </div>

        <div className="relative aspect-[210/297] w-full overflow-hidden rounded-lg border border-border bg-secondary/40">
          {gen.status === 'done' ? (
            <object data={gen.url} type="application/pdf" className="h-full w-full" aria-label="CV preview">
              <div className="flex h-full items-center justify-center p-4 text-center text-sm text-muted-foreground">
                Preview unavailable in this browser — use Download.
              </div>
            </object>
          ) : gen.status === 'loading' ? (
            <div className="flex h-full flex-col items-center justify-center gap-3 text-muted-foreground">
              <Loader2 className="size-6 animate-spin text-primary" />
              <p className="text-sm">Rendering your PDF…</p>
            </div>
          ) : gen.status === 'error' ? (
            <div className="flex h-full flex-col items-center justify-center gap-2 px-5 text-center">
              <p className="text-sm font-semibold text-destructive">{gen.message}</p>
              {gen.issues.length > 0 && (
                <ul className="max-h-40 space-y-1 overflow-auto text-xs text-muted-foreground">
                  {gen.issues.map((issue, i) => (
                    <li key={i}>{issue}</li>
                  ))}
                </ul>
              )}
            </div>
          ) : (
            <div className="flex h-full flex-col items-center justify-center gap-2 px-6 text-center text-muted-foreground">
              <FileText className="size-8 opacity-40" />
              <p className="text-sm">Your PDF preview will appear here.</p>
            </div>
          )}
        </div>

        <Button type="button" className="w-full" size="lg" onClick={onGenerate} disabled={gen.status === 'loading'}>
          {gen.status === 'loading' ? (
            <>
              <Loader2 className="animate-spin" /> Generating…
            </>
          ) : gen.status === 'done' ? (
            <>
              <RotateCcw /> Regenerate
            </>
          ) : (
            'Generate PDF'
          )}
        </Button>

        {gen.status === 'done' && (
          <Button asChild variant="outline" className="w-full">
            <a href={gen.url} download={gen.fileName}>
              <Download /> Download {gen.fileName}
            </a>
          </Button>
        )}

        <Button type="button" variant="ghost" size="sm" className="w-full" onClick={onExportJson}>
          Export JSON for the CLI
        </Button>
      </CardContent>
    </Card>
  )
}
