package coat

import "github.com/deltegui/hx"

// Tab is an entry of Tabs and the panel it shows. The panel is built when the
// tab is chosen and is not in the tree otherwise.
type Tab struct {
	ID, Label string
	Panel     func() hx.INode
}

// Tabs is a row of tabs and the panel of the chosen one. active holds the ID
// and clicking a tab sets it.
func Tabs(active *hx.SignalT[string], tabs ...Tab) hx.INode { return tabs_(active, "", tabs) }

// LineTabs is Tabs with the line look (an underline instead of a pill).
func LineTabs(active *hx.SignalT[string], tabs ...Tab) hx.INode { return tabs_(active, "line", tabs) }

func tabs_(active *hx.SignalT[string], look string, tabs []Tab) hx.INode {
	buttons := make([]hx.INode, len(tabs))
	for i, t := range tabs {
		b := hx.Button().Type("button").Attribute("role", "tab").Text(t.Label)
		b.OnClick(func(hx.EventContext) { active.Set(t.ID) })
		b.Effect(func() { b.Attribute("aria-selected", boolAttr(active.Get() == t.ID)) })
		buttons[i] = b
	}
	list := hx.Div().Attribute("role", "tablist").Attribute("aria-orientation", "horizontal").BodyList(buttons)
	if look != "" {
		list.Attribute("data-variant", look)
	}
	panels := hx.Dynamic(func() hx.INode {
		id := active.Get()
		var shown hx.INode
		hx.Untrack(func() {
			for _, t := range tabs {
				if t.ID == id && t.Panel != nil {
					shown = hx.Div().Attribute("role", "tabpanel").Body(t.Panel())
				}
			}
		})
		return shown
	})
	return hx.Div().Class("tabs").Body(list, panels)
}

// Sidebar is the menu at the side of an app. header and footer can be nil; put
// SidebarGroups in groups. While open is false it is hidden on a phone.
func Sidebar(open hx.Gettable[bool], header, footer hx.INode, groups ...hx.INode) hx.INode {
	var parts []hx.INode
	if header != nil {
		parts = append(parts, hx.Header().Body(header))
	}
	parts = append(parts, hx.El("section").Class("scrollbar").Body(groups...))
	if footer != nil {
		parts = append(parts, hx.Footer().Body(footer))
	}
	side := hx.Aside().Class("sidebar").Attribute("data-sidebar-initialized", "true").Body(hx.Nav().BodyList(parts))
	return side.Effect(func() { side.Attribute("aria-hidden", boolAttr(!open.Get())) })
}

// SidebarGroup is a titled list of SidebarLinks.
func SidebarGroup(title string, links ...hx.INode) hx.INode {
	return hx.Div().Attribute("role", "group").Body(hx.H3().Text(title), hx.Ul().Body(links...))
}

// SidebarLink is a link of a SidebarGroup. current marks the page you are on.
func SidebarLink(label, href string, current bool) hx.INode {
	a := hx.A().Href(href).Text(label)
	if current {
		a.Attribute("aria-current", "page")
	}
	return hx.Li().Body(a)
}
