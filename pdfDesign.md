# CV Style Guide — "Navy Rule" Template

A minimal, editorial one-page résumé style built for ATS-safe Word/Google Docs documents. Structured, high-contrast, no decorative elements beyond thin rules and a single accent color.

---

## 1. Design Intent

- Print-safe, single-page, works identically in Microsoft Word and Google Docs.
- Information hierarchy is carried by **weight, color and thin rules** — never by boxes, icons, or background fills.
- Reads as "corporate technical" — appropriate for engineering/IT roles, not a creative portfolio.

---

## 2. Color Palette

| Token              | Hex       | RGB           | Usage                                                                             |
| ------------------ | --------- | ------------- | --------------------------------------------------------------------------------- |
| **Accent (Navy)**  | `#1F3864` | 31, 56, 100   | Name, section heading labels, section rule lines, skill/language table row labels |
| **Subtle (Slate)** | `#44546A` | 68, 84, 106   | Subtitle line under name, company names, dates, "Stack:" labels and values        |
| **Text (Ink)**     | `#1A1A1A` | 26, 26, 26    | Body copy, bullet text, table values                                              |
| **Background**     | `#FFFFFF` | 255, 255, 255 | Page background — always pure white, no off-white/cream                           |

Rules:

- Only ever three text colors + white background. No fourth color, no tints/gradients.
- Accent is reserved for **structural markers** (name, headings, rule lines, table row labels) — never for body prose.
- Subtle is for **metadata** (dates, company names, stack lists) — anything secondary to the main content.
- Accent-on-white contrast is AA-compliant for normal text; safe for both screen and print.

---

## 3. Typography

**Primary font:** Calibri
**Fallback stack (for cross-platform / Google Docs rendering):** Carlito → Arial → sans-serif
_(Carlito is metric-compatible with Calibri and is what Google Docs/LibreOffice substitute automatically — always test the layout after fallback substitution since line-wrapping can shift by a few pixels.)_

| Role                        | Size (pt) | Weight  | Style  | Color  | Case          | Letter-spacing         |
| --------------------------- | --------- | ------- | ------ | ------ | ------------- | ---------------------- |
| Name (header)               | 20pt      | Bold    | —      | Accent | As typed      | Wide (+0.5pt tracking) |
| Subtitle (role tagline)     | 11.5pt    | Regular | —      | Subtle | As typed      | Normal                 |
| Contact line                | 9.5pt     | Regular | —      | Text   | As typed      | Normal                 |
| Section heading             | 10.5pt    | Bold    | —      | Accent | UPPERCASE     | Wide (+0.5pt tracking) |
| Role title                  | 10.5pt    | Bold    | —      | Text   | As typed      | Normal                 |
| Company name                | 10.5pt    | Regular | —      | Subtle | As typed      | Normal                 |
| Date range                  | 9.5pt     | Regular | Italic | Subtle | As typed      | Normal                 |
| "Stack:" label              | 9.5pt     | Bold    | —      | Subtle | As typed      | Normal                 |
| Stack value                 | 9.5pt     | Regular | Italic | Subtle | As typed      | Normal                 |
| Body / bullet text          | 10pt      | Regular | —      | Text   | Sentence case | Normal                 |
| Table row label (left col)  | 9.5pt     | Bold    | —      | Accent | As typed      | Normal                 |
| Table row value (right col) | 9.5pt     | Regular | —      | Text   | As typed      | Normal                 |

No more than **two weights** (Regular, Bold) and **one italic style** are used anywhere. No underlines except the structural rule lines described below.

---

## 4. Page Setup

| Property            | Value                                                        |
| ------------------- | ------------------------------------------------------------ |
| Page size           | A4 (210 × 297 mm)                                            |
| Top / bottom margin | 0.33 in (8.5 mm / 24 pt)                                     |
| Left / right margin | 0.59 in (15 mm / 42.5 pt)                                    |
| Content width       | ~7.1 in (181 mm)                                             |
| Target length       | Exactly 1 page — tighten spacing before adding a second page |

---

## 5. Spacing System

A compact vertical rhythm — everything is tuned tighter than typical Word defaults to force one-page density. Use these as a scale, not fixed pixel values:

| Token      | Value      | Used for                                            |
| ---------- | ---------- | --------------------------------------------------- |
| `space-xs` | 1–2 pt     | Between a bullet and the next bullet                |
| `space-sm` | 2.5–3.5 pt | After a "Stack:" line; after a section heading rule |
| `space-md` | 5.5–9 pt   | Before a section heading; before a new role entry   |
| `space-lg` | ~12 pt     | Between the header block and the first section      |

Rule of thumb: **headings get more space above than below** (they should feel attached to the content that follows, separated from what precedes).

---

## 6. Structural Components

### 6.1 Header block

```
NAME                                            ← 20pt bold, accent, uppercase-ish tracking
Role Tagline                                    ← 11.5pt regular, subtle
City, Country  |  email@domain.com  |  phone    ← 9.5pt regular, text color
──────────────────────────────────────────────  ← 1.25pt solid rule, accent color, full content width
```

The contact line sits directly above a single horizontal rule (1.25pt / accent) that spans the full content width. This rule is the only "heavy" line on the page — everything else is thinner.

### 6.2 Section heading

```
SECTION LABEL                                   ← 10.5pt bold, accent, ALL CAPS, +tracking
──────────────────────────────────────────────  ← 0.75pt solid rule, accent color, full content width
```

Every section heading is immediately followed by a thin (0.75pt) accent rule spanning the full content width. This is the recurring signature element of the template — repeat it identically for every section, no exceptions and no variation in weight.

### 6.3 Experience / education entry

```
Role Title  —  Company, City                                    Month Year – Month Year
Stack: Technology, Technology, Technology
• Bullet describing an achievement, starting with a past-tense action verb.
• Second bullet.
```

- Title, company, and date sit on **one line**, using a right-aligned tab stop so the date range is flush with the right margin of the content area — never let it wrap to its own line.
- Title: bold, text color. Company: regular, subtle color, separated from the title by an em dash with spaces (`  —  `).
- Date range: italic, subtle color, right-aligned.
- "Stack:" line is optional per entry — include only where a specific tech stack applies (e.g., work experience), omit for education entries.
- Bullets: 3 per role is the target (2 minimum, 4 maximum) — each one sentence, starting with an action verb, no trailing period is fine either way but stay consistent within the document.

### 6.4 Two-column key-value table (skills / languages)

- Two-column table, **no borders, no shading, no gridlines**.
- Left column ≈ 27% of content width, right column ≈ 73%.
- Left column: bold, accent color (the "key").
- Right column: regular, text color (the "value").
- Row vertical padding: minimal (~2pt top/bottom) — this is a dense list, not a spaced-out table.
- Used for exactly two sections in this template: **Compétences techniques / Skills** and **Langues / Languages**. Do not use a table anywhere else.

### 6.5 Bullet list

- Marker: simple round bullet (`•`), not a dash or arrow.
- Hanging indent: text aligns ~13pt from the left margin; bullet sits ~3pt further left.
- No nested/sub-bullets anywhere in this template — flatten everything to one level.

---

## 7. Section Order (top to bottom)

1. Header block (name, tagline, contact, rule)
2. Expérience professionnelle / Experience — **reverse chronological**
3. Compétences techniques / Skills (table)
4. Formation / Education — reverse chronological, same entry component as experience (title/institution/dates line), no bullets needed unless a note is useful
5. Langues / Languages (table)

There is intentionally **no summary/profile paragraph** — the template leads with proof (experience), not a self-description.

---

## 8. Do / Don't

**Do**

- Keep the accent color to exactly one hue across the whole document.
- Keep every section heading's rule line the same weight (0.75pt) — only the header's rule is heavier (1.25pt).
- Right-align dates using a tab stop, not manual spaces.
- Compress line spacing aggressively before shrinking font sizes — font sizes below 9.5pt should be a last resort.

**Don't**

- Don't add icons, colored backgrounds, sidebars, or a photo placeholder box — this template is text/rule-based only.
- Don't mix more than one accent color, even for "categories" — use bold/case/position for hierarchy instead of more colors.
- Don't let a table use percentage-based column widths — use fixed widths (inches/points/twips) so the layout survives round-tripping between Word and Google Docs without reflowing.
- Don't justify body text — left-align everything.

---

## 9. File Format Notes (for reproducing in Word / Google Docs)

- Build tables with **fixed-width columns** (not "auto" or percentage-based) — percentage table widths render inconsistently between Word and Google Docs.
- Avoid text boxes, floating shapes, or absolute-positioned elements — everything should be in the normal document flow so it survives conversion between formats.
- Use native paragraph borders for the rule lines (bottom border on a paragraph), not inserted horizontal-line objects or drawn shapes.
- Use a right tab stop for the date alignment rather than a table — it's more robust across editors than a borderless table row.
