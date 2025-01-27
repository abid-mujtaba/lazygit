package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MenuContext struct {
	*MenuViewModel
	*ListContextTrait
}

var _ types.IListContext = (*MenuContext)(nil)

func NewMenuContext(
	view *gocui.View,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.HelperCommon,
	getOptionsMap func() map[string]string,
) *MenuContext {
	viewModel := NewMenuViewModel()

	return &MenuContext{
		MenuViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				ViewName:        "menu",
				Key:             "menu",
				Kind:            types.PERSISTENT_POPUP,
				OnGetOptionsMap: getOptionsMap,
				Focusable:       true,
			}), ContextCallbackOpts{
				OnFocus:        onFocus,
				OnFocusLost:    onFocusLost,
				OnRenderToMain: onRenderToMain,
			}),
			getDisplayStrings: viewModel.GetDisplayStrings,
			list:              viewModel,
			viewTrait:         NewViewTrait(view),
			c:                 c,
		},
	}
}

// TODO: remove this thing.
func (self *MenuContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.DisplayString
}

type MenuViewModel struct {
	menuItems []*types.MenuItem
	*BasicViewModel[*types.MenuItem]
}

func NewMenuViewModel() *MenuViewModel {
	self := &MenuViewModel{
		menuItems: nil,
	}

	self.BasicViewModel = NewBasicViewModel(func() []*types.MenuItem { return self.menuItems })

	return self
}

func (self *MenuViewModel) SetMenuItems(items []*types.MenuItem) {
	self.menuItems = items
}

// TODO: move into presentation package
func (self *MenuViewModel) GetDisplayStrings(startIdx int, length int) [][]string {
	stringArrays := make([][]string, len(self.menuItems))
	for i, item := range self.menuItems {
		if item.DisplayStrings == nil {
			styledStr := item.DisplayString
			if item.OpensMenu {
				styledStr = presentation.OpensMenuStyle(styledStr)
			}
			stringArrays[i] = []string{styledStr}
		} else {
			stringArrays[i] = item.DisplayStrings
		}
	}

	return stringArrays
}
