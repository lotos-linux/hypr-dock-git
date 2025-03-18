package ipc

type Workspace struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Monitor         string `json:"monitor"`
	Windows         int    `json:"windows"`
	Hasfullscreen   bool   `json:"hasfullscreen"`
	Lastwindow      string `json:"lastwindow"`
	Lastwindowtitle string `json:"lastwindowtitle"`
}

type Monitor struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Make        string  `json:"make"`
	Model       string  `json:"model"`
	Serial      string  `json:"serial"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	RefreshRate float64 `json:"refreshRate"`
	X           int     `json:"x"`
	Y           int     `json:"y"`

	ActiveWorkspace struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"activeWorkspace"`

	Reserved   []int   `json:"reserved"`
	Scale      float64 `json:"scale"`
	Transform  int     `json:"transform"`
	Focused    bool    `json:"focused"`
	DpmsStatus bool    `json:"dpmsStatus"`
	Vrr        bool    `json:"vrr"`
}

type Client struct {
	Address string `json:"address"`
	Mapped  bool   `json:"mapped"`
	Hidden  bool   `json:"hidden"`
	At      []int  `json:"at"`
	Size    []int  `json:"size"`

	Workspace struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"workspace"`

	Floating         bool          `json:"floating"`
	Pseudo           bool          `json:"pseudo"`
	Monitor          int           `json:"monitor"`
	Class            string        `json:"class"`
	Title            string        `json:"title"`
	InitialClass     string        `json:"initialClass"`
	InitialTitle     string        `json:"initialTitle"`
	Pid              int           `json:"pid"`
	Xwayland         bool          `json:"xwayland"`
	Pinned           bool          `json:"pinned"`
	Fullscreen       int           `json:"fullscreen"`
	FullscreenClient int           `json:"fullscreenClient"`
	Grouped          []interface{} `json:"grouped"`
	Tags             []interface{} `json:"tags"`
	Swallowing       string        `json:"swallowing"`
	FocusHistoryID   int           `json:"focusHistoryID"`
	InhibitingIdle   bool          `json:"inhibitingIdle"`
}
