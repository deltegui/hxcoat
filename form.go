package coat

import (
	"strconv"
	"strings"

	"github.com/deltegui/hx"
)

// Field is a row of a form: a Label, a control, a Description and an error.
func Field(children ...hx.INode) hx.INode { return hx.Div().Class("field").Body(children...) }

// Label is the label of a control. forID is the id of the control.
func Label(text, forID string) hx.INode {
	l := hx.Label().Class("label").Text(text)
	if forID != "" {
		l.Attribute("for", forID)
	}
	return l
}

// Description is the small text under a control.
func Description(text string) hx.INode { return hx.P().Text(text) }

// FormField is the common case: a label, a control and the errors of that
// field in red (nothing if errs is nil or empty). The error text follows the
// signal. The control gets the id that the label points to.
func FormField(label, id string, control hx.INode, errs hx.Gettable[[]string]) hx.INode {
	control.Id(id)
	parts := []hx.INode{Label(label, id), control}
	if errs != nil {
		alert := hx.P().Attribute("role", "alert").Style("color", "var(--destructive)")
		alert.Effect(func() { alert.Text(strings.Join(errs.Get(), " · ")) })
		parts = append(parts, alert)
	}
	return Field(parts...)
}

// Fieldset groups controls under a legend.
func Fieldset(legend string, children ...hx.INode) hx.INode {
	return hx.El("fieldset").Class("fieldset").Body(append([]hx.INode{hx.El("legend").Text(legend)}, children...)...)
}

// FieldSeparator is a line with a text in the middle ("or").
func FieldSeparator(text string) hx.INode {
	return hx.Div().Class("field-separator").Body(hx.Hr().Attribute("role", "separator"), hx.Span().Text(text))
}

// Input is an input of the given type ("text", "email", "password", "number",
// "date", "search"...), bound to the signal.
func Input(kind string, value *hx.SignalT[string]) hx.INode {
	return hx.Input().Type(kind).Class("input").BindValue(value).BindOnInput(value)
}

// Textarea is a multi line text.
func Textarea(value *hx.SignalT[string]) hx.INode {
	return hx.TextArea().Class("textarea").BindValue(value).BindOnInput(value)
}

// Select is the native select. It has an empty first option with the text of
// blank (none if blank is empty). For a select with a popup see SelectMenu.
func Select(value *hx.SignalT[string], blank string, options []Option) hx.INode {
	var nodes []hx.INode
	if blank != "" {
		nodes = append(nodes, hx.Option().Value("").Text(blank))
	}
	for _, o := range options {
		nodes = append(nodes, hx.Option().Value(o.Value).Text(o.Label))
	}
	return hx.Select().Class("select").BodyList(nodes).BindValue(value).BindOnChange(value)
}

// Checkbox is a checkbox with its label, bound to the signal.
func Checkbox(label string, value *hx.SignalT[bool]) hx.INode {
	return hx.Label().Class("label").Body(
		hx.Input().Type("checkbox").Class("input").BindChecked(value).BindOnChecked(value),
		hx.Span().Text(label))
}

// Switch is a toggle with its label, bound to the signal.
func Switch(label string, value *hx.SignalT[bool]) hx.INode {
	return hx.Label().Class("label").Body(
		hx.Input().Type("checkbox").Class("input").Attribute("role", "switch").BindChecked(value).BindOnChecked(value),
		hx.Span().Text(label))
}

// Radio is one radio button: it is checked while the signal has its value, and
// it sets the signal when chosen.
func Radio(name, value, label string, selected *hx.SignalT[string]) hx.INode {
	input := hx.Input().Type("radio").Class("input").Attribute("name", name).Value(value)
	input.Effect(func() { input.Checked(selected.Get() == value) })
	input.OnChange(func(hx.EventContext) { selected.Set(value) })
	return hx.Label().Class("label").Body(input, hx.Span().Text(label))
}

// RadioGroup is a Radio for each option, in a group that screen readers announce.
func RadioGroup(name string, selected *hx.SignalT[string], options []Option) hx.INode {
	nodes := make([]hx.INode, len(options))
	for i, o := range options {
		nodes[i] = Radio(name, o.Value, o.Label, selected)
	}
	return hx.Div().Attribute("role", "radiogroup").Class("grid", "gap-3").BodyList(nodes)
}

// Slider is a range input. The signal holds the number as text.
func Slider(minimum, maximum, step int, value *hx.SignalT[string]) hx.INode {
	return hx.Input().Type("range").Class("input").
		Attribute("min", strconv.Itoa(minimum)).Attribute("max", strconv.Itoa(maximum)).Attribute("step", strconv.Itoa(step)).
		BindValue(value).BindOnInput(value)
}

// InputGroup joins an input with text, icons and buttons around it. Put the
// extras in Addon, with Start or End, and use With(.., Vertical) for a column.
func InputGroup(children ...hx.INode) hx.INode {
	return hx.Div().Class("input-group").Body(children...)
}

// Addon is text or an icon that goes at the Start or the End of an InputGroup.
func Addon(children ...hx.INode) hx.INode {
	return hx.Span().Attribute("data-align", "start").Body(children...)
}

// FileInput is a file picker. Read the files with hxhttp.FilesOf(ctx.Target)
// in the browser.
func FileInput(onChange func(hx.EventContext)) hx.INode {
	f := hx.Input().Type("file").Class("input")
	if onChange != nil {
		f.OnChange(onChange)
	}
	return f
}
