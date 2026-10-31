// Package coat wraps Basecoat (https://basecoatui.com, the shadcn/ui design
// system in plain HTML and CSS) in Hix nodes: every component of Basecoat is a
// function that returns an hx.INode with the right tags, roles and attributes.
//
// It uses only Basecoat's CSS. Basecoat ships a little JavaScript for what needs
// behavior (dialogs, menus, tabs, toasts); here that behavior is a Signal, so
// the same code runs on the server (hx/ssr), in the browser (hx/wasm) and when
// the browser hydrates what the server wrote. It imports only hx.
//
// A page needs the stylesheet, and Tailwind for the utility classes that
// Basecoat expects around its components (flex, grid, gap-4...):
//
//	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/basecoat-css@1.1.0/dist/basecoat.cdn.min.css">
//	<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4.1.13/dist/index.global.js"></script>
//
// (the constants CSS and Tailwind have those addresses). Dark mode is the class
// "dark" on <html>, as in shadcn/ui.
//
// Basecoat picks the look of a component with attributes (data-variant,
// data-size), so a modifier is a Mod, and the components whose first arguments
// are plain text take them at the end:
//
//	coat.Button("Delete", coat.Destructive, coat.Sm)
//	coat.Badge("New", coat.Secondary)
//
// A component that takes children can be given a Mod with With:
//
//	coat.With(coat.Card(...), coat.Small)
package coat

import (
	"strconv"

	"github.com/deltegui/hx"
)

// CSS is the stylesheet of the version of Basecoat that this package is written for.
const CSS = "https://cdn.jsdelivr.net/npm/basecoat-css@1.1.0/dist/basecoat.cdn.min.css"

// Tailwind is the script that makes the utility classes (flex, gap-4, p-6...)
// that the page uses. Basecoat's stylesheet has the components, not the utilities.
const Tailwind = "https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4.1.13/dist/index.global.js"

// Mod changes a node: the look of a component (a variant, a size) or anything
// else you want to set on it.
type Mod func(hx.INode)

// With applies the mods to a node and returns it.
func With(n hx.INode, mods ...Mod) hx.INode {
	for _, m := range mods {
		m(n)
	}
	return n
}

// Attr is a Mod that sets an attribute.
func Attr(key, value string) Mod { return func(n hx.INode) { n.Attribute(key, value) } }

// Classes is a Mod that adds classes (Tailwind utilities, for example).
func Classes(classes ...string) Mod { return func(n hx.INode) { n.Class(classes...) } }

func variant(v string) Mod { return Attr("data-variant", v) }
func size(s string) Mod    { return Attr("data-size", s) }

// Variants of buttons, badges, items and alerts. The first look of each
// component is its default, which needs no Mod.
var (
	Primary     = variant("primary")
	Secondary   = variant("secondary")
	Outline     = variant("outline")
	Ghost       = variant("ghost")
	Link        = variant("link")
	Destructive = variant("destructive")
	Muted       = variant("muted")
)

// Sizes. The Icon ones are for a button that has only an icon.
var (
	XS     = size("xs")
	Sm     = size("sm")
	Lg     = size("lg")
	Icon   = size("icon")
	IconXS = size("icon-xs")
	IconSm = size("icon-sm")
	IconLg = size("icon-lg")
	Small  = size("sm") // for a Card or an Item
)

// Where a popover, a tooltip or a drawer shows up.
var (
	Top    = Attr("data-side", "top")
	Bottom = Attr("data-side", "bottom")
	Left   = Attr("data-side", "left")
	Right  = Attr("data-side", "right")
	Start  = Attr("data-align", "start")
	Center = Attr("data-align", "center")
	End    = Attr("data-align", "end")
)

// Vertical makes a ButtonGroup or an InputGroup go top to bottom.
var Vertical = Attr("data-orientation", "vertical")

// Option is an item of a Select or a group of Radios: what it sends and what
// the user reads.
type Option struct{ Value, Label string }

// Plain makes options whose label is their value.
func Plain(values ...string) []Option {
	options := make([]Option, len(values))
	for i, v := range values {
		options[i] = Option{v, v}
	}
	return options
}

// BindAttr sets the attribute of the node to what the function says, and
// follows the signals the function reads: it is how a component reacts to a
// Signal without being rebuilt.
func BindAttr(n hx.INode, key string, value func() string) hx.INode {
	return n.Effect(func() { n.Attribute(key, value()) })
}

// BindClass adds the class to the node while the signal is true.
func BindClass(n hx.INode, class string, on hx.Gettable[bool]) hx.INode {
	return n.Effect(func() { n.ClassIf(class, on.Get()) })
}

func boolAttr(on bool) string {
	if on {
		return "true"
	}
	return "false"
}

// Layout helpers with Tailwind classes (they need the Tailwind script).

// Stack puts the children in a column with a gap of n (4 is 1rem).
func Stack(gap int, children ...hx.INode) hx.INode {
	return hx.Div().Class("flex", "flex-col", "gap-"+strconv.Itoa(gap)).Body(children...)
}

// Row puts the children in a line with a gap of n, centered vertically.
func Row(gap int, children ...hx.INode) hx.INode {
	return hx.Div().Class("flex", "flex-wrap", "items-center", "gap-"+strconv.Itoa(gap)).Body(children...)
}

// Grid puts the children in a grid of cols columns with a gap of n.
func Grid(cols, gap int, children ...hx.INode) hx.INode {
	return hx.Div().Class("grid", "grid-cols-"+strconv.Itoa(cols), "gap-"+strconv.Itoa(gap)).Body(children...)
}
