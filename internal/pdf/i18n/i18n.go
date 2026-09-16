// Package i18n holds the template's own static labels (section headings,
// "Stack:", date separators) in each supported language. It does not
// translate user-supplied content (job titles, bullet text, etc.) — that
// stays in whatever language the input JSON was written in.
package i18n

// Lang is a supported output language for template-owned strings.
type Lang string

const (
	EN Lang = "en"
	FR Lang = "fr"
)

// Labels is the set of strings a template needs, in one language.
type Labels struct {
	Profile        string
	Education      string
	Experience     string
	Skills         string
	Projects       string
	Certifications string
	Languages      string
	Stack          string // e.g. "Stack:" — the trailing space/colon is included
	DateSep        string // joins a start and end date, e.g. "2021 – Present"
}

var english = Labels{
	Profile:        "Profile",
	Education:      "Education",
	Experience:     "Experience",
	Skills:         "Skills",
	Projects:       "Projects",
	Certifications: "Certifications",
	Languages:      "Languages",
	Stack:          "Stack: ",
	DateSep:        " – ",
}

var french = Labels{
	Profile:        "Profil",
	Education:      "Formation",
	Experience:     "Expérience professionnelle",
	Skills:         "Compétences techniques",
	Projects:       "Projets",
	Certifications: "Certifications",
	Languages:      "Langues",
	Stack:          "Stack : ",
	DateSep:        " – ",
}

// For returns the label set for lang, defaulting to English for anything
// other than French.
func For(lang Lang) Labels {
	if lang == FR {
		return french
	}
	return english
}
