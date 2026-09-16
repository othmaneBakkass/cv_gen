import { Button } from '#/components/ui/button'
import { Input } from '#/components/ui/input'
import { Label } from '#/components/ui/label'

// Editor for a plain string array (job tools, highlights, certifications,
// specialties). `form` is intentionally loosely typed: this is a small
// internal helper mounted at many different array paths, and threading
// TanStack Form's full generic signature through every call site here
// would cost more than it buys.
export function StringListField({
  form,
  name,
  label,
  placeholder,
  addLabel,
}: {
  form: any
  name: string
  label: string
  placeholder?: string
  addLabel: string
}) {
  const labelId = `${name}-label`.replace(/[[\].]/g, '-')
  return (
    <form.Field name={name} mode="array">
      {(field: any) => (
        <div className="space-y-2">
          <Label id={labelId}>{label}</Label>
          {field.state.value.map((_: string, i: number) => (
            <div key={i} className="flex gap-2">
              <form.Field name={`${name}[${i}]`}>
                {(sub: any) => (
                  <Input
                    aria-labelledby={labelId}
                    value={sub.state.value}
                    placeholder={placeholder}
                    onBlur={sub.handleBlur}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                      sub.handleChange(e.target.value)
                    }
                  />
                )}
              </form.Field>
              <Button
                type="button"
                variant="ghost"
                size="sm"
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
            onClick={() => field.pushValue('')}
          >
            {addLabel}
          </Button>
        </div>
      )}
    </form.Field>
  )
}
