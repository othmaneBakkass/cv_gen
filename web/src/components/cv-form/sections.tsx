import { Button } from '#/components/ui/button'
import { Input } from '#/components/ui/input'
import { Label } from '#/components/ui/label'
import { Textarea } from '#/components/ui/textarea'
import { Separator } from '#/components/ui/separator'
import { StringListField } from './StringListField'

// Same note as StringListField: `form` stays loosely typed on purpose —
// these are internal section editors, not part of the public form API.

function TextInput({
  form,
  name,
  label,
  placeholder,
}: {
  form: any
  name: string
  label: string
  placeholder?: string
}) {
  const id = name.replace(/[[\].]/g, '-')
  return (
    <form.Field name={name}>
      {(field: any) => (
        <div className="space-y-1.5">
          <Label htmlFor={id}>{label}</Label>
          <Input
            id={id}
            value={field.state.value}
            placeholder={placeholder}
            onBlur={field.handleBlur}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              field.handleChange(e.target.value)
            }
          />
          {(field.state.meta.isTouched || form.state.submissionAttempts > 0) &&
            field.state.meta.errors.length > 0 && (
            <p className="text-xs text-destructive">
              {field.state.meta.errors.map((e: any) => e.message ?? e).join(', ')}
            </p>
          )}
        </div>
      )}
    </form.Field>
  )
}

export function EducationSection({ form }: { form: any }) {
  return (
    <form.Field name="education" mode="array">
      {(field: any) => (
        <div className="space-y-6">
          {field.state.value.map((_: unknown, i: number) => (
            <div key={i} className="space-y-4">
              {i > 0 && <Separator />}
              <div className="grid gap-4 sm:grid-cols-2">
                <TextInput form={form} name={`education[${i}].degree`} label="Degree" />
                <TextInput form={form} name={`education[${i}].school`} label="School" />
                <TextInput form={form} name={`education[${i}].location`} label="Location" />
                <div className="grid grid-cols-2 gap-2">
                  <TextInput form={form} name={`education[${i}].startedAt`} label="Started" />
                  <TextInput form={form} name={`education[${i}].endedAt`} label="Ended" />
                </div>
              </div>
              <form.Field name={`education[${i}].description`}>
                {(field2: any) => (
                  <div className="space-y-1.5">
                    <Label htmlFor={`education-${i}-description`}>Description</Label>
                    <Textarea
                      id={`education-${i}-description`}
                      rows={2}
                      value={field2.state.value}
                      onBlur={field2.handleBlur}
                      onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
                        field2.handleChange(e.target.value)
                      }
                    />
                  </div>
                )}
              </form.Field>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() => field.removeValue(i)}
                disabled={field.state.value.length <= 1}
              >
                Remove this school
              </Button>
            </div>
          ))}
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() =>
              field.pushValue({
                school: '',
                location: '',
                startedAt: '',
                endedAt: '',
                degree: '',
                description: '',
              })
            }
          >
            Add school
          </Button>
        </div>
      )}
    </form.Field>
  )
}

export function JobsSection({ form }: { form: any }) {
  return (
    <form.Field name="jobs" mode="array">
      {(field: any) => (
        <div className="space-y-6">
          {field.state.value.map((_: unknown, i: number) => (
            <div key={i} className="space-y-4">
              {i > 0 && <Separator />}
              <div className="grid gap-4 sm:grid-cols-2">
                <TextInput form={form} name={`jobs[${i}].position`} label="Position" />
                <TextInput form={form} name={`jobs[${i}].company`} label="Company" />
                <TextInput form={form} name={`jobs[${i}].location`} label="Location" />
                <div className="grid grid-cols-2 gap-2">
                  <TextInput form={form} name={`jobs[${i}].startedAt`} label="Started" />
                  <TextInput form={form} name={`jobs[${i}].endedAt`} label="Ended" />
                </div>
              </div>
              <StringListField
                form={form}
                name={`jobs[${i}].tools`}
                label="Tools"
                placeholder="e.g. Go"
                addLabel="Add tool"
              />
              <StringListField
                form={form}
                name={`jobs[${i}].highlights`}
                label="Highlights"
                placeholder="What you did, and the result"
                addLabel="Add highlight"
              />
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() => field.removeValue(i)}
                disabled={field.state.value.length <= 1}
              >
                Remove this job
              </Button>
            </div>
          ))}
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() =>
              field.pushValue({
                company: '',
                location: '',
                position: '',
                startedAt: '',
                endedAt: '',
                tools: [''],
                highlights: [''],
              })
            }
          >
            Add job
          </Button>
        </div>
      )}
    </form.Field>
  )
}

export function LanguagesSection({ form }: { form: any }) {
  return (
    <form.Field name="languages" mode="array">
      {(field: any) => (
        <div className="space-y-3">
          {field.state.value.map((_: unknown, i: number) => (
            <div key={i} className="flex gap-2">
              <div className="flex-1">
                <TextInput form={form} name={`languages[${i}].language`} label="Language" />
              </div>
              <div className="flex-1">
                <TextInput form={form} name={`languages[${i}].level`} label="Level" />
              </div>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                className="mt-6"
                onClick={() => field.removeValue(i)}
                disabled={field.state.value.length <= 1}
              >
                Remove
              </Button>
            </div>
          ))}
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => field.pushValue({ language: '', level: '' })}
          >
            Add language
          </Button>
        </div>
      )}
    </form.Field>
  )
}

export function SkillsSection({ form }: { form: any }) {
  return (
    <form.Field name="skills" mode="array">
      {(field: any) => (
        <div className="space-y-6">
          {field.state.value.map((_: unknown, i: number) => (
            <div key={i} className="space-y-3">
              {i > 0 && <Separator />}
              <TextInput
                form={form}
                name={`skills[${i}].category`}
                label="Category"
                placeholder="e.g. Languages"
              />
              <StringListField
                form={form}
                name={`skills[${i}].items`}
                label="Items"
                placeholder="e.g. Go"
                addLabel="Add item"
              />
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() => field.removeValue(i)}
              >
                Remove this category
              </Button>
            </div>
          ))}
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => field.pushValue({ category: '', items: [''] })}
          >
            Add skill category
          </Button>
        </div>
      )}
    </form.Field>
  )
}

export function ProjectsSection({ form }: { form: any }) {
  return (
    <form.Field name="projects" mode="array">
      {(field: any) => (
        <div className="space-y-6">
          {field.state.value.map((_: unknown, i: number) => (
            <div key={i} className="space-y-3">
              {i > 0 && <Separator />}
              <div className="grid gap-4 sm:grid-cols-2">
                <TextInput form={form} name={`projects[${i}].name`} label="Name" />
                <TextInput form={form} name={`projects[${i}].tech`} label="Tech" />
              </div>
              <form.Field name={`projects[${i}].description`}>
                {(field2: any) => (
                  <div className="space-y-1.5">
                    <Label htmlFor={`projects-${i}-description`}>Description</Label>
                    <Textarea
                      id={`projects-${i}-description`}
                      rows={2}
                      value={field2.state.value}
                      onBlur={field2.handleBlur}
                      onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
                        field2.handleChange(e.target.value)
                      }
                    />
                  </div>
                )}
              </form.Field>
              <Button type="button" variant="ghost" size="sm" onClick={() => field.removeValue(i)}>
                Remove this project
              </Button>
            </div>
          ))}
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => field.pushValue({ name: '', tech: '', description: '' })}
          >
            Add project
          </Button>
        </div>
      )}
    </form.Field>
  )
}

export { TextInput }
