import { z } from 'zod'

// Mirrors internal/schema/schema.go and schemas/v1.json. Kept in sync by
// hand since the frontend has no backend call (yet) to source it from.

export const headSchema = z.object({
  fullName: z.string().min(1, 'Required'),
  address: z.string().min(1, 'Required'),
  phone: z.string().min(1, 'Required'),
  email: z.string().min(1, 'Required').email('Must be a valid email'),
  jobTitle: z.string().optional(),
  specialties: z.array(z.string().min(1)).optional(),
  linkedin: z.string().optional(),
  github: z.string().optional(),
  photo: z.string().optional(),
})

export const educationSchema = z.object({
  school: z.string().min(1, 'Required'),
  location: z.string().min(1, 'Required'),
  startedAt: z.string().min(1, 'Required'),
  endedAt: z.string().min(1, 'Required'),
  degree: z.string().min(1, 'Required'),
  description: z.string().min(1, 'Required'),
})

export const jobSchema = z.object({
  company: z.string().min(1, 'Required'),
  location: z.string().min(1, 'Required'),
  position: z.string().min(1, 'Required'),
  startedAt: z.string().min(1, 'Required'),
  endedAt: z.string().min(1, 'Required'),
  tools: z.array(z.string().min(1)).min(1, 'Add at least one'),
  highlights: z.array(z.string().min(1)).min(1, 'Add at least one'),
})

export const languageSchema = z.object({
  language: z.string().min(1, 'Required'),
  level: z.string().min(1, 'Required'),
})

export const skillSchema = z.object({
  category: z.string().min(1, 'Required'),
  items: z.array(z.string().min(1)).min(1, 'Add at least one'),
})

export const projectSchema = z.object({
  name: z.string().min(1, 'Required'),
  tech: z.string().min(1, 'Required'),
  description: z.string().min(1, 'Required'),
})

// Template ids mirror internal/pdf/templates' registry (t1/t2/t3).
export const templateId = z.enum(['t1', 't2', 't3'])
export type TemplateId = z.infer<typeof templateId>

// A hex color, "#RRGGBB" or "RRGGBB" — the same rule as theme.ParseHex and
// schemas/v1.json. Empty string is allowed and means "no override".
const hexColor = z
  .string()
  .regex(/^$|^#?[0-9a-fA-F]{6}$/, 'Use a hex color like #1F3864')

// Semantic palette roles (see theme.Palette). Every role is optional; the
// four legacy aliases still work on the Go side but the UI only writes roles.
export const colorSettingsSchema = z.object({
  headline: hexColor.optional(),
  subheadline: hexColor.optional(),
  body: hexColor.optional(),
  metadata: hexColor.optional(),
  link: hexColor.optional(),
  border: hexColor.optional(),
  profile: hexColor.optional(),
  background: hexColor.optional(),
})
export type ColorSettings = z.infer<typeof colorSettingsSchema>

// The palette roles the color-picker UI exposes, with human labels. Kept in
// one place so the editor and any legend stay in sync.
export const colorRoles: Array<{ key: keyof ColorSettings; label: string; hint: string }> = [
  { key: 'headline', label: 'Headline', hint: 'Section headings, name, "Stack:" label' },
  { key: 'border', label: 'Border', hint: 'Header & section rules' },
  { key: 'body', label: 'Body', hint: 'Body copy, titles, bullets' },
  { key: 'metadata', label: 'Metadata', hint: 'Dates, company, specialties' },
  { key: 'subheadline', label: 'Subheadline', hint: 'Tagline under the name' },
  { key: 'link', label: 'Link', hint: 'LinkedIn / GitHub' },
  { key: 'profile', label: 'Profile', hint: 'Summary paragraph' },
  { key: 'background', label: 'Background', hint: 'Page background' },
]

export const sectionStyleSchema = z.object({
  divider: z.boolean().optional(),
  uppercase: z.boolean().optional(),
  spaceBefore: z.number().optional(),
  colors: colorSettingsSchema.optional(),
})

export const sectionSettingSchema = z.object({
  include: z.boolean().optional(),
  priority: z.number().optional(),
  style: sectionStyleSchema.optional(),
})

// The seven orderable/toggleable section keys, matching render.SectionKey.
export const sectionKeys = [
  'profile',
  'education',
  'experience',
  'skills',
  'projects',
  'certifications',
  'languages',
] as const
export type SectionKey = (typeof sectionKeys)[number]

export const settingsSchema = z.object({
  lang: z.enum(['en', 'fr']).optional(),
  density: z.enum(['dense', 'normal', 'airy']).optional(),
  sections: z.record(z.enum(sectionKeys), sectionSettingSchema).optional(),
  colors: colorSettingsSchema.optional(),
})
export type Settings = z.infer<typeof settingsSchema>

export const cvSchema = z.object({
  template: templateId,
  fileName: z.string().min(1, 'Required'),
  head: headSchema,
  profile: z.string().optional(),
  education: z.array(educationSchema).min(1, 'Add at least one'),
  jobs: z.array(jobSchema).min(1, 'Add at least one'),
  skills: z.array(skillSchema).optional(),
  projects: z.array(projectSchema).optional(),
  certifications: z.array(z.string().min(1)).optional(),
  languages: z.array(languageSchema).min(1, 'Add at least one'),
  settings: settingsSchema.optional(),
})

export type CV = z.infer<typeof cvSchema>

export const emptyCV: CV = {
  template: 't1',
  fileName: '',
  head: { fullName: '', address: '', phone: '', email: '', specialties: [] },
  profile: '',
  education: [
    { school: '', location: '', startedAt: '', endedAt: '', degree: '', description: '' },
  ],
  jobs: [
    {
      company: '',
      location: '',
      position: '',
      startedAt: '',
      endedAt: '',
      tools: [''],
      highlights: [''],
    },
  ],
  skills: [],
  projects: [],
  certifications: [],
  languages: [{ language: '', level: '' }],
}

// A filled-in example, downloadable from the "Upload JSON" tab so a user
// can edit it in a text editor and upload it back instead of using the form.
export const exampleCV: CV = {
  template: 't2',
  fileName: 'your-name',
  head: {
    fullName: 'Jane Doe',
    address: 'Casablanca, Morocco',
    phone: '+212 600 000 000',
    email: 'jane.doe@example.com',
    jobTitle: 'Senior Backend Engineer',
    specialties: ['Distributed Systems', 'Cloud Infrastructure'],
    linkedin: 'linkedin.com/in/jane-doe',
    github: 'github.com/janedoe',
  },
  profile: 'A short, 3-4 line summary of who you are and what you bring.',
  education: [
    {
      school: 'Your University',
      location: 'City, Country',
      startedAt: '2016',
      endedAt: '2019',
      degree: 'B.Sc. in Software Engineering',
      description: 'One line on your focus or dissertation.',
    },
  ],
  jobs: [
    {
      company: 'Company Name',
      location: 'Remote',
      position: 'Backend Engineer',
      startedAt: '2021',
      endedAt: 'Present',
      tools: ['Go', 'PostgreSQL', 'Kubernetes'],
      highlights: [
        'A concrete achievement with a measurable result.',
        'A second achievement, or a responsibility you owned.',
      ],
    },
  ],
  skills: [{ category: 'Languages', items: ['Go', 'TypeScript'] }],
  projects: [
    { name: 'Project name', tech: 'Tech used', description: 'One line on what it does.' },
  ],
  certifications: ['Certification title — Issuing organization'],
  languages: [{ language: 'English', level: 'Fluent' }],
}
