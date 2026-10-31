# hxcoat

Every component of [Basecoat](https://basecoatui.com) (the shadcn/ui design system
in plain HTML and CSS) as a Hix node. Only Basecoat's CSS is used. What Basecoat
does with JavaScript (dialogs, menus, tabs, toasts) is a `Signal` here, so it works
with `hx/ssr`, with `hx/wasm` and when the browser hydrates what the server wrote.
It imports only `hx`.

A page needs the stylesheet and Tailwind (Basecoat's CSS has the components, not
the utility classes like `flex` or `gap-4` that go around them):

```html
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/basecoat-css@1.1.0/dist/basecoat.cdn.min.css">
<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4.1.13/dist/index.global.js"></script>
```

(`coat.CSS` and `coat.Tailwind` have those addresses. Dark mode is the class `dark` on `<html>`.)

```go
import "github.com/deltegui/hxcoat"

open := hx.Signal(false)
coat.Stack(4,
	coat.Button("Delete", coat.Destructive).OnClick(func(hx.EventContext) { open.Set(true) }),
	coat.AlertDialog(open, "Are you sure?", "This cannot be undone.", func() hx.INode {
		return hx.Noop().Body(
			coat.Button("Cancel", coat.Outline).OnClick(func(hx.EventContext) { open.Set(false) }),
			coat.Button("Continue").OnClick(func(hx.EventContext) { open.Set(false) }))
	}),
)
```

Basecoat chooses the look with attributes (`data-variant`, `data-size`), so a modifier
is a `Mod`. A component whose arguments are plain text takes them last
(`coat.Button("Save", coat.Secondary, coat.Sm)`); one that takes children is
given them with `coat.With(node, mods...)`. `coat.Classes(...)` and `coat.Attr(k, v)` are
Mods too.

| File | Components |
| --- | --- |
| `basic.go` | Button, SubmitButton, ButtonLink, IconButton, ButtonGroup, Badge, Alert, Card, Avatar, Kbd, Empty, Item, Progress, Skeleton, Table, Tooltip, Accordion, Breadcrumb, Pagination |
| `form.go` | Field, Label, FormField, Fieldset, FieldSeparator, Input, Textarea, Select, Checkbox, Switch, Radio(Group), Slider, InputGroup, FileInput |
| `overlay.go` | Dialog, AlertDialog, Drawer, Popover, DropdownMenu, SelectMenu, Combobox, Command, Toaster |
| `nav.go` | Tabs, LineTabs, Sidebar |
| `coat.go` | the Mods (variants, sizes, sides), `With`, `BindAttr`, `BindClass`, Stack/Row/Grid |

An overlay is in the tree only while its signal says so, and a click outside a
menu closes it. Basecoat draws the dim background of a dialog with the `::backdrop`
of a dialog opened by JavaScript; here the dialog uses the `open` attribute, so
the package draws that background itself. Not included: the chart and the theme
switcher, which are JavaScript.

`OWL_DEMO/feature/coat` shows all of them in `/coat`; `make check-coat` there drives it in Chromium.
