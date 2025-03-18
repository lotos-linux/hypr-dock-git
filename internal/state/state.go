package state

import (
	"hypr-dock/internal/item"
	"sync"

	"github.com/gotk3/gotk3/gtk"
)

type State struct {
	Window      *gtk.Window
	DetectArea  *gtk.Window
	ItemsBox    *gtk.Box
	Orientation gtk.Orientation
	PreventHide bool
	AddedApps   AddedApps
	Special     bool
	mu          sync.Mutex
}

type AddedApps struct {
	List map[string]*item.Item
	mu   sync.Mutex
}

func New() *State {
	return &State{
		AddedApps: NewAddedApps(),
	}
}

func NewAddedApps() AddedApps {
	return AddedApps{
		List: make(map[string]*item.Item),
	}
}

func (s *State) GetAddedApps() *AddedApps {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &s.AddedApps
}

func (aa *AddedApps) Add(className string, item *item.Item) {
	aa.mu.Lock()
	defer aa.mu.Unlock()
	aa.List[className] = item
}

func (aa *AddedApps) Remove(className string) {
	aa.mu.Lock()
	defer aa.mu.Unlock()
	delete(aa.List, className)
}

func (s *State) Update(fn func(*State)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(s)
}

func (s *State) SetWindow(window *gtk.Window) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Window = window
}

func (s *State) SetDetectArea(window *gtk.Window) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DetectArea = window
}

func (s *State) SetItemsBox(box *gtk.Box) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ItemsBox = box
}

func (s *State) SetOrientation(orientation gtk.Orientation) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Orientation = orientation
}

func (s *State) SetPreventHide(is bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.PreventHide = is
}

func (s *State) SetSpecial(is bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Special = is
}

func (s *State) GetWindow() *gtk.Window {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Window
}

func (s *State) GetDetectArea() *gtk.Window {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.DetectArea
}

func (s *State) GetItemsBox() *gtk.Box {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ItemsBox
}

func (s *State) GetOrientation() gtk.Orientation {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Orientation
}

func (s *State) GetPreventHide() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.PreventHide
}

func (s *State) GetSpecial() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Special
}
