import { Link } from '@tanstack/react-router'
import { FileText } from 'lucide-react'
import { Button } from './ui/button'
import ThemeToggle from './ThemeToggle'

export default function Header() {
  return (
    <header className="sticky top-0 z-50 border-b border-border bg-background/80 px-4 backdrop-blur-md">
      <nav className="mx-auto flex max-w-6xl items-center gap-6 py-3">
        <Link
          to="/"
          className="group flex items-center gap-2 no-underline"
          aria-label="cv_gen home"
        >
          <span className="flex size-7 items-center justify-center rounded-md bg-primary text-primary-foreground transition-transform duration-200 group-hover:scale-105">
            <FileText className="size-4" />
          </span>
          <span className="text-base font-extrabold tracking-tight text-foreground">cv_gen</span>
        </Link>

        <div className="flex flex-1 items-center gap-5 text-sm font-medium">
          <Link
            to="/templates"
            className="text-muted-foreground no-underline transition-colors hover:text-foreground"
            activeProps={{ className: 'text-foreground' }}
          >
            Templates
          </Link>
          <Link
            to="/build"
            className="text-muted-foreground no-underline transition-colors hover:text-foreground"
            activeProps={{ className: 'text-foreground' }}
          >
            Build
          </Link>
        </div>

        <ThemeToggle />
        <Button asChild size="sm" className="hidden sm:inline-flex">
          <Link to="/build">Build a CV</Link>
        </Button>
      </nav>
    </header>
  )
}
