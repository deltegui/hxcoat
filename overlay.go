package coat

import (
	"strings"

	"github.com/deltegui/hx"
)

// A closed overlay is not in the tree: the content of a Dialog, an AlertDialog
// and a Drawer is built each time it opens. Basecoat draws the dim background
// with the ::backdrop of a dialog opened by JavaScript; here the dialog is
// opened with the open attribute, which has no backdrop, so the package puts
// its own behind it.

func backdrop(onClick func()) hx.INode {
	b := hx.Div().Attribute("aria-hidden", "true").
		Style("position", "fixed").Style("inset", "0").Style("z-index", "40").
		Style("background", "color-mix(in srgb, #000 50%, transparent)")
	if onClick != nil {
		b.OnClick(func(hx.EventContext) { onClick() })
	}
	return b
}

func build(f func() hx.INode) hx.INode {
	if f == nil {
		return nil
	}
	return f()
}

func closeButton(onClose func()) hx.INode {
	x := hx.Svg().Attribute("viewBox", "0 0 24 24").Attribute("fill", "none").Attribute("stroke", "currentColor").
		Attribute("stroke-width", "2").Attribute("stroke-linecap", "round").Attribute("stroke-linejoin", "round").
		Style("width", "1rem").Style("height", "1rem").
		Body(hx.Path().Attribute("d", "M18 6 6 18"), hx.Path().Attribute("d", "m6 6 12 12"))
	b := hx.Button().Type("button").Attribute("aria-label", "Close").
		Style("position", "absolute").Style("top", "1rem").Style("right", "1rem").Body(x)
	if onClose != nil {
		b.OnClick(func(hx.EventContext) { onClose() })
	}
	return b
}

func panel(sectionClass, title, description string, body, footer func() hx.INode, extra ...hx.INode) hx.INode {
	head := []hx.INode{hx.H2().Text(title)}
	if description != "" {
		head = append(head, hx.P().Text(description))
	}
	parts := []hx.INode{hx.El("header").BodyList(head)}
	if b := build(body); b != nil {
		section := hx.El("section").Body(b)
		if sectionClass != "" {
			section.Class(sectionClass)
		}
		parts = append(parts, section)
	}
	if f := build(footer); f != nil {
		parts = append(parts, hx.Footer().Body(f))
	}
	return hx.Div().BodyList(append(parts, extra...))
}

func overlay(open hx.Gettable[bool], create func() (hx.INode, hx.INode)) hx.INode {
	return hx.Dynamic(func() hx.INode {
		if !open.Get() {
			return nil
		}
		var shade, dialog hx.INode
		hx.Untrack(func() { shade, dialog = create() })
		return hx.Noop().Body(shade, dialog)
	})
}

// closer is the function that closes something: onClose, or, when it is nil and
// open is a signal, a function that sets it to false.
func closer(open hx.Gettable[bool], onClose func()) func() {
	if onClose != nil {
		return onClose
	}
	if s, ok := open.(hx.Settable[bool]); ok {
		return func() { s.Set(false) }
	}
	return nil
}

// DialogOptions is what a Dialog needs. Only Open is required.
type DialogOptions struct {
	// Open says if the dialog is in the tree. If it is a signal, the dialog closes
	// itself with it (see OnClose).
	Open               hx.Gettable[bool]
	Title, Description string
	// Body and Footer are built each time the dialog opens. Either can be nil.
	Body, Footer func() hx.INode
	// OnClose runs when the user clicks the background or the corner button. By
	// default it sets Open to false, when Open is a signal.
	OnClose func()
}

// Dialog is a window over the page:
//
//	coat.Dialog(coat.DialogOptions{Open: open, Title: "Edit profile", Body: form, Footer: buttons})
func Dialog(o DialogOptions) hx.INode {
	onClose := closer(o.Open, o.OnClose)
	return overlay(o.Open, func() (hx.INode, hx.INode) {
		return backdrop(onClose), hx.El("dialog").Class("dialog").Attribute("open", "").Attribute("aria-label", o.Title).
			Body(panel("", o.Title, o.Description, o.Body, o.Footer, closeButton(onClose)))
	})
}

// AlertDialogOptions is what an AlertDialog needs.
type AlertDialogOptions struct {
	Open               hx.Gettable[bool]
	Title, Description string
	// Footer has the buttons, which are what closes it.
	Footer func() hx.INode
}

// AlertDialog is a Dialog that asks for an answer: it has no corner button and
// the background does not close it.
func AlertDialog(o AlertDialogOptions) hx.INode {
	return overlay(o.Open, func() (hx.INode, hx.INode) {
		return backdrop(nil), hx.El("dialog").Class("alert-dialog").Attribute("open", "").Attribute("role", "alertdialog").Attribute("aria-label", o.Title).
			Body(panel("", o.Title, o.Description, nil, o.Footer))
	})
}

// DrawerOptions is what a Drawer needs. Only Open is required.
type DrawerOptions struct {
	// Open says if the drawer is in the tree. If it is a signal, the drawer closes
	// itself with it (see OnClose).
	Open               hx.Gettable[bool]
	Title, Description string
	// Body and Footer are built each time the drawer opens. Either can be nil.
	Body, Footer func() hx.INode
	// OnClose runs when the user clicks the background or the corner button. By
	// default it sets Open to false, when Open is a signal.
	OnClose func()
	// Mods: the side it comes from (Top, Bottom, the default, Left or Right) and anything else.
	Mods []Mod
}

// Drawer is a panel that comes from a side of the page:
//
//	coat.Drawer(coat.DrawerOptions{Open: open, Title: "Filters", Mods: []coat.Mod{coat.Right}, Body: filters})
func Drawer(o DrawerOptions) hx.INode {
	onClose := closer(o.Open, o.OnClose)
	return overlay(o.Open, func() (hx.INode, hx.INode) {
		d := hx.El("dialog").Class("drawer").Attribute("open", "").Attribute("aria-label", o.Title).Style("z-index", "50").
			Body(panel("px-4", o.Title, o.Description, o.Body, o.Footer, closeButton(onClose)))
		return backdrop(onClose), With(d, o.Mods...)
	})
}

// popover is a trigger that opens a floating box. It reads and writes open, and a
// transparent layer over the page closes it when the user clicks outside.
func popover(class string, open *hx.SignalT[bool], trigger, content hx.INode, mods []Mod, haspopup string) hx.INode {
	trigger.Attribute("aria-haspopup", haspopup).OnClick(func(hx.EventContext) { open.Set(!open.Get()) })
	trigger.Effect(func() { trigger.Attribute("aria-expanded", boolAttr(open.Get())) })
	box := With(hx.Div().Attribute("data-popover", "").Body(content), mods...)
	box.Effect(func() { box.Attribute("aria-hidden", boolAttr(!open.Get())) })
	outside := hx.Dynamic(func() hx.INode {
		if !open.Get() {
			return nil
		}
		return hx.Div().Attribute("aria-hidden", "true").
			Style("position", "fixed").Style("inset", "0").Style("z-index", "40").
			OnClick(func(hx.EventContext) { open.Set(false) })
	})
	return hx.Div().Class(class).Body(trigger, box, outside)
}

// Popover shows content next to the trigger while open is true. The content
// can have a header with a title and a description. Pass Top, Bottom, Left,
// Right and Start, Center, End.
func Popover(open *hx.SignalT[bool], trigger, content hx.INode, mods ...Mod) hx.INode {
	return popover("popover", open, trigger, content, mods, "dialog")
}

// DropdownMenu is a trigger (usually a Button) that opens a menu of MenuItem,
// MenuHeading and MenuSeparator. Choosing an item closes it.
func DropdownMenu(open *hx.SignalT[bool], trigger hx.INode, items []hx.INode, mods ...Mod) hx.INode {
	menu := hx.Div().Attribute("role", "menu").BodyList(items).
		OnClick(func(hx.EventContext) { open.Set(false) })
	return popover("dropdown-menu", open, trigger, menu, mods, "menu")
}

// MenuItem is an entry of a DropdownMenu. Pass Destructive for a dangerous one.
func MenuItem(label string, onClick func(), mods ...Mod) hx.INode {
	item := hx.Div().Attribute("role", "menuitem").Attribute("tabindex", "0").Text(label)
	if onClick != nil {
		item.OnClick(func(hx.EventContext) { onClick() })
	}
	return With(item, mods...)
}

// MenuCheckbox is an entry that is checked or not, bound to the signal.
func MenuCheckbox(label string, value *hx.SignalT[bool]) hx.INode {
	item := hx.Div().Attribute("role", "menuitemcheckbox").Attribute("tabindex", "0").Text(label)
	item.Effect(func() { item.Attribute("aria-checked", boolAttr(value.Get())) })
	return item.OnClick(func(hx.EventContext) { value.Set(!value.Get()) })
}

// MenuHeading is a title between groups of items.
func MenuHeading(text string) hx.INode {
	return hx.Div().Attribute("role", "heading").Text(text)
}

func MenuSeparator() hx.INode { return hx.Hr().Attribute("role", "separator") }

// SelectMenu is a select with a popup list, bound to the signal. For the native
// select see Select.
func SelectMenu(value *hx.SignalT[string], placeholder string, options []Option, mods ...Mod) hx.INode {
	open := hx.Signal(false)
	label := hx.Span().Class("truncate")
	label.Effect(func() { label.Text(labelOf(options, value.Get(), placeholder)) })
	trigger := With(hx.Button().Type("button").Body(label, chevronDown()), Outline).Class("btn")
	list := hx.Div().Attribute("role", "listbox").Body(optionNodes(options, value, func(v string) {
		value.Set(v)
		open.Set(false)
	})...)
	return popover("select", open, trigger, list, mods, "listbox")
}

func labelOf(options []Option, value, placeholder string) string {
	for _, o := range options {
		if o.Value == value {
			return o.Label
		}
	}
	return placeholder
}

func optionNodes(options []Option, value hx.Gettable[string], choose func(string)) []hx.INode {
	nodes := make([]hx.INode, len(options))
	for i, o := range options {
		n := hx.Div().Attribute("role", "option").Attribute("data-value", o.Value).Text(o.Label)
		n.Effect(func() { n.Attribute("aria-selected", boolAttr(value.Get() == o.Value)) })
		nodes[i] = n.OnClick(func(hx.EventContext) { choose(o.Value) })
	}
	return nodes
}

func chevronDown() hx.INode {
	return hx.Svg().Attribute("viewBox", "0 0 24 24").Attribute("fill", "none").Attribute("stroke", "currentColor").
		Attribute("stroke-width", "2").Attribute("stroke-linecap", "round").Attribute("stroke-linejoin", "round").
		Body(hx.Path().Attribute("d", "m6 9 6 6 6-6"))
}

func matches(label, query string) bool {
	return strings.Contains(strings.ToLower(label), strings.ToLower(strings.TrimSpace(query)))
}

// ComboboxOptions is what a Combobox needs.
type ComboboxOptions struct {
	// Value is the value of the chosen option ("" for none).
	Value *hx.SignalT[string]
	// Placeholder shows while nothing is chosen; SearchPlaceholder is in the search box.
	Placeholder, SearchPlaceholder string
	Options                        []Option
	Mods                           []Mod
}

// Combobox is a SelectMenu with a search box that filters the options as the
// user types.
func Combobox(o ComboboxOptions) hx.INode {
	value, placeholder, searchPlaceholder, options, mods := o.Value, o.Placeholder, o.SearchPlaceholder, o.Options, o.Mods
	open := hx.Signal(false)
	query := hx.Signal("")
	label := hx.Span().Class("truncate")
	label.Effect(func() { label.Text(labelOf(options, value.Get(), placeholder)) })
	trigger := With(hx.Button().Type("button").Body(label, chevronDown()), Outline).Class("btn")
	search := hx.Input().Type("text").Attribute("role", "combobox").Placeholder(searchPlaceholder).BindValue(query).BindOnInput(query)
	list := hx.Div().Attribute("role", "listbox").Body(hx.Dynamic(func() hx.INode {
		var shown []Option
		for _, o := range options {
			if matches(o.Label, query.Get()) {
				shown = append(shown, o)
			}
		}
		var nodes []hx.INode
		hx.Untrack(func() {
			nodes = optionNodes(shown, value, func(v string) {
				value.Set(v)
				query.Set("")
				open.Set(false)
			})
		})
		return hx.Noop().BodyList(nodes)
	}))
	return popover("combobox", open, trigger, hx.Noop().Body(search, list), mods, "listbox")
}

// CommandItem is an entry of a Command. Shortcut is optional text on the right.
type CommandItem struct {
	Label, Shortcut string
	OnSelect        func()
}

// CommandGroup is a titled set of CommandItems.
type CommandGroup struct {
	Heading string
	Items   []CommandItem
}

// Command is a search box with a list that filters as the user types (a
// command palette). query holds the text; put it in a Dialog to make it modal.
func Command(query *hx.SignalT[string], placeholder string, groups ...CommandGroup) hx.INode {
	search := hx.Input().Type("text").Attribute("role", "combobox").Placeholder(placeholder).
		BindValue(query).BindOnInput(query)
	menu := hx.Div().Attribute("role", "menu").Body(hx.Dynamic(func() hx.INode {
		q := query.Get()
		var shown []hx.INode
		hx.Untrack(func() {
			for _, g := range groups {
				var items []hx.INode
				for _, it := range g.Items {
					if !matches(it.Label, q) {
						continue
					}
					content := []hx.INode{hx.Span().Text(it.Label)}
					if it.Shortcut != "" {
						content = append(content, hx.Span().Attribute("data-shortcut", "").Text(it.Shortcut))
					}
					n := hx.Div().Attribute("role", "menuitem").Attribute("tabindex", "0").BodyList(content)
					if it.OnSelect != nil {
						sel := it.OnSelect
						n.OnClick(func(hx.EventContext) { sel() })
					}
					items = append(items, n)
				}
				if len(items) == 0 {
					continue
				}
				group := hx.Div().Attribute("role", "group")
				if g.Heading != "" {
					items = append([]hx.INode{hx.Div().Attribute("role", "heading").Text(g.Heading)}, items...)
				}
				shown = append(shown, group.BodyList(items))
			}
		})
		return hx.Noop().BodyList(shown)
	}))
	return hx.Div().Class("command").Body(hx.El("header").Body(search), menu)
}

// Toast is a short message that shows in the corner.
type Toast struct {
	ID                 int
	Title, Description string
}

// Toaster is the corner where toasts show. toasts is the list on screen (add and
// remove with a signal; a toast is not removed by itself) and dismiss runs when
// the user closes one.
func Toaster(toasts hx.Gettable[[]Toast], dismiss func(id int)) hx.INode {
	return hx.Div().Class("toaster").Body(hx.Dynamic(func() hx.INode {
		list := toasts.Get()
		var nodes []hx.INode
		hx.Untrack(func() {
			for _, t := range list {
				text := []hx.INode{hx.H2().Text(t.Title)}
				if t.Description != "" {
					text = append(text, hx.P().Text(t.Description))
				}
				nodes = append(nodes, hx.Div().Class("toast").Attribute("role", "status").Attribute("aria-atomic", "true").Attribute("aria-hidden", "false").Body(
					hx.Div().Class("toast-content").Body(
						hx.El("section").BodyList(text),
						hx.Footer().Body(hx.Button().Type("button").Attribute("data-toast-cancel", "").Text("Dismiss").
							OnClick(func(hx.EventContext) { dismiss(t.ID) })))))
			}
		})
		return hx.Noop().BodyList(nodes)
	}))
}
