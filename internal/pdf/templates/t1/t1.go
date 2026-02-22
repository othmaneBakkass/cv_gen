package tone

import (
	"fmt"
	"strings"

	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	apperror "github.com/othmaneBakkass/cv_gen/internal/common/appError"
	"github.com/othmaneBakkass/cv_gen/internal/common/stringc"
	"github.com/othmaneBakkass/cv_gen/internal/pdf"
)

type InputData struct {
	Data []JSONSchema `json:"data"`
}

type fontSize struct {
	lg   float64
	base float64
	sm   float64
}

var fontSizes = fontSize{
	lg:   16,
	base: 13,
	sm:   11,
}

func t1(data JSONSchema) core.Maroto {
	engine := pdf.GetEngine()
	head(engine, data.Head)
	header(engine, "Education")
	educationSection(engine, data.Education)
	header(engine, "Experience")
	jobSection(engine, data.Jobs)
	header(engine, "Languages")
	langSection(engine, data.Languages)
	return engine
}

func head(engine core.Maroto, params HeadSchema) {
	engine.AddAutoRow(text.NewCol(12, params.FullName, props.Text{Align: align.Center, Style: fontstyle.Bold, Size: fontSizes.lg}))
	engine.AddRow(1)
	engine.AddAutoRow(
		text.NewCol(12, fmt.Sprintf("%s | %s | %s", params.Address, params.Email, params.Phone), props.Text{Align: align.Center, Size: 13}),
	)
	engine.AddRow(6)
}

func header(engine core.Maroto, label string) {
	engine.AddAutoRow(text.NewCol(12, strings.ToUpper(label), props.Text{Align: align.Left, Style: fontstyle.Bold, Size: 13}))
	engine.AddRow(1)
	engine.AddAutoRow(line.NewCol(12, props.Line{OffsetPercent: 0, SizePercent: 100}))
	engine.AddRow(2)
}

func educationSection(engine core.Maroto, education []EducationSchema) {
	for _, v := range education {
		engine.AddAutoRow(
			text.NewCol(6, stringc.ToCapital(v.School), props.Text{Align: align.Left, Style: fontstyle.Bold, Size: fontSizes.base}),
			text.NewCol(6, fmt.Sprintf("%s | %s - %s", v.Location, v.StartedAt, v.EndedAt), props.Text{Align: align.Right, Size: fontSizes.base}),
		)
		engine.AddRow(1)
		engine.AddAutoRow(text.NewCol(12, v.Degree, props.Text{Align: align.Left, Size: fontSizes.base}))
		engine.AddRow(1)
		engine.AddAutoRow(text.NewCol(12, v.Description, props.Text{Align: align.Left, Size: fontSizes.base}))
		engine.AddRow(4)
	}
}

func jobSection(engine core.Maroto, jobs []JobSchema) {
	for _, v := range jobs {
		engine.AddAutoRow(
			text.NewCol(6, stringc.ToCapital(v.Company), props.Text{Align: align.Left, Style: fontstyle.Bold, Size: fontSizes.base}),
			text.NewCol(6, v.Location, props.Text{Align: align.Right, Size: fontSizes.base}),
		)
		engine.AddRow(1)
		engine.AddAutoRow(
			text.NewCol(6, stringc.ToCapital(v.Position), props.Text{Align: align.Left, Style: fontstyle.Bold, Size: fontSizes.base}),
			text.NewCol(6, fmt.Sprintf("%s - %s", v.StartedAt, v.EndedAt), props.Text{Align: align.Right, Style: fontstyle.Italic, Size: fontSizes.base}),
		)
		engine.AddRow(1)
		engine.AddAutoRow(text.NewCol(12, strings.Join(v.Tools, ", "), props.Text{Align: align.Left, Size: fontSizes.sm}))
		var highlights []core.Row
		for _, highlight := range v.Highlights {
			highlights = append(highlights, row.New().Add(
				text.NewCol(0, "  •  ", props.Text{Size: 13}),
				text.NewCol(12, highlight, props.Text{Size: 13, Left: 6}),
			), row.New(2))
		}
		engine.AddRow(2)
		engine.AddRows(highlights...)
		engine.AddRow(4)
	}
}

func langSection(engine core.Maroto, langs []LanguageSchema) {
	var parts = []string{}
	for _, v := range langs {
		parts = append(parts, fmt.Sprintf("%s - %s", stringc.ToCapital(v.Language), stringc.ToCapital(v.Level)))
	}
	engine.AddAutoRow(
		text.NewCol(12, strings.Join(parts, " | "), props.Text{Align: align.Left, Size: 13}),
	)
}

func GenerateT1PDF(saveAt string, data JSONSchema) error {
	var engine = t1(data)
	doc, err := engine.Generate()

	if err != nil {
		return apperror.New("PDF generation", "PDF could not be generated", apperror.ErrorCodeUnknown, apperror.ErrorSensitivityPublic)
	}

	err = doc.Save(saveAt)
	if err != nil {
		return apperror.New("PDF save", "PDF could not be saved", apperror.ErrorCodeUnknown, apperror.ErrorSensitivityPublic)
	}
	return nil
}
