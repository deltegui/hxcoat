package coat_test

import (
	"strings"
	"testing"

	"github.com/deltegui/hx"
	"github.com/deltegui/hx/ssr"
	"github.com/deltegui/hxcoat"
)

func TestComponents(t *testing.T) {
	text := hx.Signal("hi")
	flag := hx.Signal(true)
	off := hx.Signal(false)
	closed := hx.Signal(false)
	num := hx.Signal(3)
	tab := hx.Signal("b")
	choice := hx.Signal("x")
	pick := hx.Signal("b")
	toasts := hx.Signal([]coat.Toast{{ID: 1, Title: "Saved", Description: "All good"}})
	nop := func() {}
	node := func() hx.INode { return hx.P().Text("x") }

	cases := []struct {
		name  string
		build func() hx.INode
		want  []string
	}{
		{"button", func() hx.INode { return coat.Button("Go", coat.Destructive, coat.Sm) }, []string{`class="btn"`, `data-size="sm"`, `data-variant="destructive"`, `type="button"`, `>Go</button>`}},
		{"submit link icon", func() hx.INode {
			return coat.ButtonGroup(coat.SubmitButton("a"), coat.ButtonLink("/x", "b", coat.Link), coat.IconButton("More", hx.Span().Text("+")))
		}, []string{`class="button-group"`, `type="submit"`, `href="/x"`, `data-size="icon"`, `aria-label="More"`}},
		{"badge", func() hx.INode { return coat.Badge("New", coat.Secondary) }, []string{`class="badge"`, `data-variant="secondary"`}},
		{"alert", func() hx.INode { return coat.Alert("Oops", "It broke", coat.Destructive) }, []string{`class="alert"`, `role="alert"`, `<h2>Oops</h2>`, `<section>It broke</section>`, `data-variant="destructive"`}},
		{"card", func() hx.INode {
			return coat.With(coat.Card(coat.CardHeader("T", "D", coat.Button("A")), coat.CardSection(node()), coat.CardFooter(node())), coat.Small)
		}, []string{`class="card"`, `data-size="sm"`, `<header><h2>T</h2><p>D</p>`, `data-slot="card-action"`, `<section><p>x</p></section>`, `<footer>`}},
		{"avatar", func() hx.INode {
			return coat.AvatarGroup(coat.Avatar("/a.png", "A", ""), coat.Avatar("", "", "AB", coat.Lg))
		}, []string{`class="avatar-group"`, `<img`, `alt="A"`, `<span>AB</span>`, `data-size="lg"`}},
		{"kbd empty", func() hx.INode {
			return hx.Div().Body(coat.Kbd("⌘K"), coat.Empty("None", "Nothing", coat.Button("Add")))
		}, []string{`class="kbd"`, `class="empty"`, `<h2>None</h2>`, `<section>`}},
		{"item", func() hx.INode {
			return coat.ItemGroup(coat.Item("Title", "Desc", coat.Avatar("", "", "A"), coat.Button("Go"), coat.Outline))
		}, []string{`class="item-group"`, `class="item"`, `data-variant="outline"`, `<figure>`, `<h3>Title</h3>`, `<aside>`}},
		{"progress", func() hx.INode { return coat.ProgressOf(num) }, []string{`class="progress"`, `role="progressbar"`, `aria-valuenow="3"`, `width:3%`}},
		{"skeleton", func() hx.INode { return coat.Skeleton("100%", "2rem") }, []string{`class="bg-accent rounded-md skeleton"`, `height:2rem`, `width:100%`}},
		{"skeleton circle", func() hx.INode { return coat.Skeleton("3rem", "3rem", coat.Circle) }, []string{`bg-accent rounded-full skeleton`}},
		{"table", func() hx.INode { return coat.TableContainer(coat.SimpleTable([]string{"A", "B"}, hx.Tr())) }, []string{`class="table-container"`, `class="table"`, `<th scope="col">A</th>`}},
		{"tooltip", func() hx.INode { return coat.Tooltip(coat.Button("x"), "Hello", coat.Bottom) }, []string{`data-tooltip="Hello"`, `data-side="bottom"`}},
		{"accordion", func() hx.INode {
			return coat.Accordion("faq", coat.AccordionItem{Title: "One", Body: node(), Open: true}, coat.AccordionItem{Title: "Two", Body: node()})
		}, []string{`class="accordion"`, `<details name="faq" open="">`, `<summary>One</summary>`, `<details name="faq"><summary>Two`}},
		{"breadcrumb", func() hx.INode { return coat.Breadcrumb(coat.Crumb{"Home", "/"}, coat.Crumb{"Now", "/n"}) }, []string{`class="breadcrumb"`, `<a href="/">Home</a>`, `aria-hidden="true"`, `aria-current="page"`}},
		{"pagination", func() hx.INode { return coat.Pagination(hx.Signal(5), hx.Signal(10), func(int) {}) }, []string{`aria-label="pagination"`, `aria-current="page"`, `…`, `Previous`, `Next`}},
		{"pagination first", func() hx.INode { return coat.Pagination(hx.Signal(1), hx.Signal(1), func(int) {}) }, []string{`disabled=""`}},
		{"field", func() hx.INode {
			return coat.FormField("Name", "name", coat.Input("text", text), hx.Signal([]string{"bad", "worse"}))
		}, []string{`class="field"`, `class="label"`, `for="name"`, `id="name"`, `class="input"`, `role="alert"`, `bad · worse`}},
		{"controls", func() hx.INode {
			return hx.Div().Body(coat.Textarea(text), coat.Select(text, "—", coat.Plain("a", "b")), coat.Checkbox("ok", off), coat.Switch("on", off),
				coat.RadioGroup("r", choice, coat.Plain("x", "y")), coat.Slider(0, 10, 1, text), coat.FileInput(nil), coat.Description("hint"))
		}, []string{`class="textarea"`, `class="select"`, `type="checkbox"`, `role="switch"`, `role="radiogroup"`, `type="radio"`, `type="range"`, `max="10"`, `type="file"`, `<p>hint</p>`}},
		{"fieldset input group", func() hx.INode {
			return coat.Fieldset("Account", coat.FieldSeparator("or"), coat.With(coat.InputGroup(coat.Addon(hx.Span().Text("$")), coat.Input("text", text)), coat.Vertical))
		}, []string{`class="fieldset"`, `<legend>Account</legend>`, `class="field-separator"`, `class="input-group"`, `data-orientation="vertical"`, `data-align="start"`}},
		{"dialog open", func() hx.INode { return coat.Dialog(flag, nop, "Edit", "Change it", node, node) }, []string{`class="dialog"`, `open=""`, `<h2>Edit</h2>`, `<p>Change it</p>`, `<section><p>x</p></section>`, `<footer>`, `aria-label="Close"`, `position:fixed`}},
		{"dialog closed", func() hx.INode { return hx.Div().Body(coat.Dialog(closed, nop, "Edit", "", node, node)) }, []string{`<div></div>`}},
		{"alert dialog", func() hx.INode { return coat.AlertDialog(flag, "Sure?", "No undo", node) }, []string{`class="alert-dialog"`, `role="alertdialog"`, `<h2>Sure?</h2>`}},
		{"drawer", func() hx.INode { return coat.Drawer(flag, nop, "Menu", "", node, nil, coat.Left) }, []string{`class="drawer"`, `data-side="left"`}},
		{"dropdown open", func() hx.INode {
			return coat.DropdownMenu(flag, coat.Button("Menu"), []hx.INode{coat.MenuHeading("Account"), coat.MenuItem("Profile", nil), coat.MenuSeparator(), coat.MenuItem("Delete", nop, coat.Destructive), coat.MenuCheckbox("Show", flag)}, coat.End)
		}, []string{`class="dropdown-menu"`, `aria-expanded="true"`, `aria-haspopup="menu"`, `data-popover=""`, `aria-hidden="false"`, `data-align="end"`, `role="menu"`, `role="heading"`, `role="menuitem"`, `role="separator"`, `aria-checked="true"`, `position:fixed`}},
		{"dropdown closed", func() hx.INode { return coat.DropdownMenu(off, coat.Button("Menu"), nil) }, []string{`aria-expanded="false"`, `aria-hidden="true"`}},
		{"popover", func() hx.INode { return coat.Popover(flag, coat.Button("Info"), hx.P().Text("Hi")) }, []string{`class="popover"`, `aria-haspopup="dialog"`}},
		{"select menu", func() hx.INode { return coat.SelectMenu(pick, "Pick", []coat.Option{{"a", "Alpha"}, {"b", "Beta"}}) }, []string{`class="select"`, `aria-haspopup="listbox"`, `>Beta<`, `role="listbox"`, `role="option"`, `aria-selected="true"`, `data-value="a"`}},
		{"combobox", func() hx.INode { return coat.Combobox(pick, "Pick", "Search", coat.Plain("a", "b")) }, []string{`class="combobox"`, `role="combobox"`, `placeholder="Search"`, `role="option"`}},
		{"command", func() hx.INode {
			return coat.Command(hx.Signal("cal"), "Type", coat.CommandGroup{"Suggestions", []coat.CommandItem{{Label: "Calendar", Shortcut: "⌘C"}, {Label: "Search"}}})
		}, []string{`class="command"`, `role="group"`, `role="heading"`, `>Calendar<`, `data-shortcut=""`}},
		{"toaster", func() hx.INode { return coat.Toaster(toasts, func(int) {}) }, []string{`class="toaster"`, `class="toast"`, `class="toast-content"`, `<h2>Saved</h2>`, `<p>All good</p>`, `data-toast-cancel`}},
		{"tabs", func() hx.INode {
			return coat.Tabs(tab, coat.Tab{ID: "a", Label: "A", Panel: func() hx.INode { return hx.P().Text("PANEL-A") }}, coat.Tab{ID: "b", Label: "B", Panel: func() hx.INode { return hx.P().Text("PANEL-B") }})
		}, []string{`class="tabs"`, `role="tablist"`, `role="tab"`, `aria-selected="true"`, `role="tabpanel"`, `PANEL-B`}},
		{"line tabs", func() hx.INode { return coat.LineTabs(tab, coat.Tab{ID: "a", Label: "A"}) }, []string{`data-variant="line"`}},
		{"sidebar", func() hx.INode {
			return coat.Sidebar(flag, hx.Span().Text("Brand"), nil, coat.SidebarGroup("Main", coat.SidebarLink("Home", "/", true), coat.SidebarLink("About", "/a", false)))
		}, []string{`class="sidebar"`, `aria-hidden="false"`, `data-sidebar-initialized`, `<h3>Main</h3>`, `aria-current="page"`}},
		{"layout", func() hx.INode { return coat.Stack(4, coat.Row(2, coat.Grid(3, 2))) }, []string{`flex flex-col gap-4`, `flex flex-wrap gap-2 items-center`, `gap-2 grid grid-cols-3`}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			html := ssr.Render(c.build)
			for _, w := range c.want {
				if !strings.Contains(html, w) {
					t.Errorf("missing %q in\n%s", w, html)
				}
			}
		})
	}
	if html := ssr.Render(func() hx.INode {
		return coat.Tabs(tab, coat.Tab{ID: "a", Label: "A", Panel: func() hx.INode { return hx.P().Text("PANEL-A") }})
	}); strings.Contains(html, "PANEL-A") {
		t.Errorf("Tabs shows an inactive panel: %s", html)
	}
	if html := ssr.Render(func() hx.INode {
		return coat.Command(hx.Signal("zzz"), "Type", coat.CommandGroup{"G", []coat.CommandItem{{Label: "Calendar"}}})
	}); strings.Contains(html, "Calendar") {
		t.Errorf("Command does not filter: %s", html)
	}
}
