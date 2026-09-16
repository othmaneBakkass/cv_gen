import { Link } from '@tanstack/react-router'
import { FileText } from 'lucide-react'

export default function Footer() {
  return (
    <footer className="mt-20 border-t border-border px-4 py-10">
      <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 sm:flex-row">
        <Link to="/" className="flex items-center gap-2 no-underline">
          <span className="flex size-6 items-center justify-center rounded-md bg-primary text-primary-foreground">
            <FileText className="size-3.5" />
          </span>
          <span className="text-sm font-bold text-foreground">cv_gen</span>
        </Link>
        <nav className="flex items-center gap-5 text-sm text-muted-foreground">
          <Link to="/templates" className="no-underline transition-colors hover:text-foreground">
            Templates
          </Link>
          <Link to="/build" className="no-underline transition-colors hover:text-foreground">
            Build
          </Link>
        </nav>
        <p className="m-0 text-sm text-muted-foreground">
          Generate CVs from JSON, in three configurable templates.
        </p>
      </div>
    </footer>
  )
}
