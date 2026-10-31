package coat

import (
	"strconv"

	"github.com/deltegui/hx"
)

// Button is a button that does not submit: Button("Save", Destructive, Sm).
func Button(label string, mods ...Mod) hx.INode {
	return With(hx.Button().Type("button").Class("btn").Text(label), mods...)
}

// SubmitButton submits the form it is in.
func SubmitButton(label string, mods ...Mod) hx.INode {
	return With(hx.Button().Type("submit").Class("btn").Text(label), mods...)
}

// ButtonLink is an anchor that looks like a button.
func ButtonLink(href, label string, mods ...Mod) hx.INode {
	return With(hx.A().Href(href).Class("btn").Text(label), mods...)
}

// IconButton is a button with only an icon (give it a label for screen readers).
func IconButton(label string, icon hx.INode, mods ...Mod) hx.INode {
	return With(hx.Button().Type("button").Class("btn").Attribute("aria-label", label).Body(icon), append([]Mod{Icon, Outline}, mods...)...)
}

// ButtonGroup sticks buttons, inputs and menus together. Use With(.., Vertical).
func ButtonGroup(children ...hx.INode) hx.INode {
	return hx.Div().Class("button-group").Attribute("role", "group").Body(children...)
}

// Badge is a small label: Badge("New", Secondary).
func Badge(text string, mods ...Mod) hx.INode {
	return With(hx.Span().Class("badge").Text(text), mods...)
}

// Alert is a message box with a title and a description (either can be empty).
// Use the Destructive mod for an error. An icon goes first when you pass one.
func Alert(title, description string, mods ...Mod) hx.INode {
	a := hx.Div().Class("alert").Attribute("role", "alert")
	var parts []hx.INode
	if title != "" {
		parts = append(parts, hx.H2().Text(title))
	}
	if description != "" {
		parts = append(parts, hx.El("section").Text(description))
	}
	return With(a.BodyList(parts), mods...)
}

// Card is a box with parts: put CardHeader, CardSection and CardFooter in it.
// Use With(.., Small) for a compact one.
func Card(children ...hx.INode) hx.INode { return hx.Div().Class("card").Body(children...) }

// CardHeader is the title, the description (empty for none) and an action at
// the right (nil for none).
func CardHeader(title, description string, action hx.INode) hx.INode {
	parts := []hx.INode{hx.H2().Text(title)}
	if description != "" {
		parts = append(parts, hx.P().Text(description))
	}
	if action != nil {
		parts = append(parts, hx.Div().Attribute("data-slot", "card-action").Body(action))
	}
	return hx.El("header").BodyList(parts)
}

func CardSection(children ...hx.INode) hx.INode { return hx.El("section").Body(children...) }
func CardFooter(children ...hx.INode) hx.INode  { return hx.Footer().Body(children...) }

// Avatar is a round picture. With no src it shows fallback (the initials).
func Avatar(src, alt, fallback string, mods ...Mod) hx.INode {
	a := hx.Span().Class("avatar")
	if src != "" {
		a.Body(hx.Img().Src(src).Attribute("alt", alt))
	} else {
		a.Body(hx.Span().Text(fallback))
	}
	return With(a, mods...)
}

// AvatarGroup overlaps avatars. A last Span with "data-count" shows "+3".
func AvatarGroup(children ...hx.INode) hx.INode {
	return hx.Div().Class("avatar-group").Body(children...)
}

// Kbd is a key of the keyboard.
func Kbd(keys string) hx.INode { return hx.El("kbd").Class("kbd").Text(keys) }

// Empty is what a list shows when it has nothing: a title, a description and the actions.
func Empty(title, description string, actions ...hx.INode) hx.INode {
	e := hx.Div().Class("empty")
	header := hx.El("header").Body(hx.H2().Text(title), hx.P().Text(description))
	if len(actions) == 0 {
		return e.Body(header)
	}
	return e.Body(header, hx.El("section").Body(actions...))
}

// Item is a row with an optional media on the left, a title and description,
// and actions on the right. Use Outline or Muted for the look, and Sm or XS.
func Item(title, description string, media, actions hx.INode, mods ...Mod) hx.INode {
	var parts []hx.INode
	if media != nil {
		parts = append(parts, hx.El("figure").Body(media))
	}
	text := []hx.INode{hx.H3().Text(title)}
	if description != "" {
		text = append(text, hx.P().Text(description))
	}
	parts = append(parts, hx.El("section").BodyList(text))
	if actions != nil {
		parts = append(parts, hx.Aside().Body(actions))
	}
	return With(hx.Div().Class("item").BodyList(parts), mods...)
}

// ItemGroup groups Items; put separators (Hr) between them if you want.
func ItemGroup(children ...hx.INode) hx.INode { return hx.Div().Class("item-group").Body(children...) }

// Progress is a bar that is value percent full (0 to 100).
func Progress(value int) hx.INode {
	bar := hx.Span()
	setBar(bar, value)
	return progressBox().Attribute("aria-valuenow", strconv.Itoa(value)).Body(bar)
}

// ProgressOf is a Progress that follows a signal.
func ProgressOf(value hx.Gettable[int]) hx.INode {
	bar := hx.Span()
	box := progressBox().Body(bar)
	return box.Effect(func() {
		v := value.Get()
		box.Attribute("aria-valuenow", strconv.Itoa(v))
		setBar(bar, v)
	})
}

func progressBox() hx.INode {
	return hx.Div().Class("progress").Attribute("role", "progressbar").
		Attribute("aria-valuemin", "0").Attribute("aria-valuemax", "100").
		Style("height", "0.5rem").Style("border-radius", "9999px").
		Style("background", "color-mix(in oklab, var(--primary) 20%, transparent)")
}

func setBar(bar hx.INode, value int) {
	bar.Style("width", strconv.Itoa(max(0, min(100, value)))+"%").Style("background", "var(--primary)")
}

// Skeleton is a gray block that pulses, to stand in for content that loads.
// width and height are CSS lengths ("100%", "2rem"); the look is Tailwind's,
// so the Circle mod makes a circle.
func Skeleton(width, height string, mods ...Mod) hx.INode {
	return With(hx.Div().Class("skeleton", "bg-accent", "rounded-md").Style("width", width).Style("height", height), mods...)
}

// Circle is a Mod for a Skeleton that stands in for an avatar.
func Circle(n hx.INode) { n.RemoveClass("rounded-md").Class("rounded-full") }

// Table is a table. Put THead, TBody and TFoot in it, or use SimpleTable.
func Table(children ...hx.INode) hx.INode { return hx.Table().Class("table").Body(children...) }

// SimpleTable is a table with a header row of headers and the rows in a TBody.
func SimpleTable(headers []string, rows ...hx.INode) hx.INode {
	cells := make([]hx.INode, len(headers))
	for i, h := range headers {
		cells[i] = hx.Th().Attribute("scope", "col").Text(h)
	}
	return Table(hx.THead().Body(hx.Tr().BodyList(cells)), hx.TBody().Body(rows...))
}

// TableContainer lets a wide table scroll sideways.
func TableContainer(table hx.INode) hx.INode { return hx.Div().Class("table-container").Body(table) }

// Tooltip shows text when the pointer is over n (or n has the focus). Pass
// Top, Bottom, Left or Right to say where.
func Tooltip(n hx.INode, text string, mods ...Mod) hx.INode {
	return With(n.Attribute("data-tooltip", text), mods...)
}

// AccordionItem is a section of an Accordion.
type AccordionItem struct {
	Title string
	Body  hx.INode
	Open  bool
}

// Accordion is a list of sections that fold. It is the native <details>, so it
// needs no signal. With a name only one is open at a time.
func Accordion(name string, items ...AccordionItem) hx.INode {
	nodes := make([]hx.INode, len(items))
	for i, item := range items {
		d := hx.El("details").Body(hx.El("summary").Text(item.Title), item.Body)
		if name != "" {
			d.Attribute("name", name)
		}
		if item.Open {
			d.Attribute("open", "")
		}
		nodes[i] = d
	}
	return hx.El("section").Class("accordion").BodyList(nodes)
}

// Crumb is an item of a Breadcrumb. The last one is the current page.
type Crumb struct{ Label, Href string }

// Breadcrumb is the trail "Home / Products / Edit". The last Crumb is not a link.
func Breadcrumb(crumbs ...Crumb) hx.INode {
	var items []hx.INode
	for i, c := range crumbs {
		if i > 0 {
			items = append(items, hx.Li().Attribute("role", "presentation").Attribute("aria-hidden", "true").Body(chevron()))
		}
		if i == len(crumbs)-1 {
			items = append(items, hx.Li().Body(hx.Span().Attribute("role", "link").Attribute("aria-current", "page").Attribute("aria-disabled", "true").Text(c.Label)))
		} else {
			items = append(items, hx.Li().Body(hx.A().Href(c.Href).Text(c.Label)))
		}
	}
	return hx.Nav().Class("breadcrumb").Attribute("aria-label", "breadcrumb").Body(hx.Ol().BodyList(items))
}

func chevron() hx.INode {
	return hx.Svg().Attribute("viewBox", "0 0 24 24").Attribute("fill", "none").Attribute("stroke", "currentColor").
		Attribute("stroke-width", "2").Attribute("stroke-linecap", "round").Attribute("stroke-linejoin", "round").
		Body(hx.Path().Attribute("d", "m9 18 6-6-6-6"))
}

// Pagination shows previous, a window of pages around the current one with the
// first and the last, and next. page and pages count from 1.
func Pagination(page, pages hx.Gettable[int], goTo func(int)) hx.INode {
	item := func(label string, target int, enabled bool, mods ...Mod) hx.INode {
		b := With(hx.Button().Type("button").Class("btn").Text(label), mods...)
		if enabled {
			b.OnClick(func(hx.EventContext) { goTo(target) })
		} else {
			b.Attribute("disabled", "")
		}
		return hx.Li().Body(b)
	}
	return hx.Nav().Attribute("role", "navigation").Attribute("aria-label", "pagination").Body(hx.Dynamic(func() hx.INode {
		current, last := page.Get(), pages.Get()
		numbered := func(p int) hx.INode {
			if p == current {
				return hx.Li().Body(With(hx.Button().Type("button").Class("btn").Attribute("aria-current", "page").Text(strconv.Itoa(p)), Outline, Sm))
			}
			return item(strconv.Itoa(p), p, true, Ghost, Sm)
		}
		dots := func() hx.INode {
			return hx.Li().Body(hx.Span().Style("padding", "0 0.5rem").Attribute("aria-hidden", "true").Text("…"))
		}
		items := []hx.INode{item("‹ Previous", current-1, current > 1, Ghost, Sm)}
		from, to := max(1, current-2), min(last, current+2)
		if from > 1 {
			items = append(items, numbered(1))
			if from > 2 {
				items = append(items, dots())
			}
		}
		for p := from; p <= to; p++ {
			items = append(items, numbered(p))
		}
		if to < last {
			if to < last-1 {
				items = append(items, dots())
			}
			items = append(items, numbered(last))
		}
		items = append(items, item("Next ›", current+1, current < last, Ghost, Sm))
		return hx.Ul().Class("flex", "flex-row", "items-center", "gap-1").BodyList(items)
	}))
}
