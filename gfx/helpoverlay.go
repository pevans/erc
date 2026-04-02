package gfx

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	helpPadding  = 16
	helpFontSize = 14
	helpTitle    = "Keyboard Shortcuts"
)

type helpEntry struct {
	key  string
	desc string
}

var helpEntries = []helpEntry{
	{"Ctrl-A ESC", "Pause / Resume"},
	{"Ctrl-A +/-", "Speed up / down"},
	{"Ctrl-A [/]", "Volume down / up"},
	{"Ctrl-A V", "Mute / unmute"},
	{"Ctrl-A W", "Write protect"},
	{"Ctrl-A B", "Debugger"},
	{"Ctrl-A C", "Caps lock"},
	{"Ctrl-A 0-9", "Select state slot"},
	{"Ctrl-A S", "Save state"},
	{"Ctrl-A L", "Load state"},
	{"Ctrl-A N/P", "Next / previous disk"},
	{"Ctrl-A Q", "Quit"},
	{"Ctrl-A ?/H", "This help screen"},
}

// HelpOverlay displays a modal with shortcut help text.
type HelpOverlay struct {
	active       bool
	page         int
	totalPages   int
	perPage      int
	screenWidth  int
	screenHeight int
	faceSource   *text.GoTextFaceSource
}

// HelpModal is the global help overlay instance.
var HelpModal *HelpOverlay

func init() {
	HelpModal = &HelpOverlay{}
}

func (h *HelpOverlay) ensureFont() {
	if h.faceSource == nil {
		fs, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
		if err != nil {
			return
		}
		h.faceSource = fs
	}
}

// Show activates the help modal.
func (h *HelpOverlay) Show(width, height int) {
	h.active = true
	h.page = 0
	h.screenWidth = width
	h.screenHeight = height
	h.ensureFont()
	h.calculatePagination()
}

// Hide deactivates the help modal.
func (h *HelpOverlay) Hide() {
	h.active = false
}

// IsActive returns whether the help modal is currently shown.
func (h *HelpOverlay) IsActive() bool {
	return h.active
}

func (h *HelpOverlay) calculatePagination() {
	face := &text.GoTextFace{
		Source: h.faceSource,
		Size:   helpFontSize,
	}

	_, lineHeight := text.Measure("X", face, 0)

	// Title + gap + entries + possible footer
	titleHeight := lineHeight + helpPadding
	footerHeight := lineHeight + helpPadding
	availableHeight := float64(h.screenHeight)*0.8 - 2*helpPadding - titleHeight - footerHeight

	h.perPage = int(availableHeight / lineHeight)
	if h.perPage < 1 {
		h.perPage = 1
	}
	if h.perPage >= len(helpEntries) {
		h.perPage = len(helpEntries)
		h.totalPages = 1
	} else {
		h.totalPages = (len(helpEntries) + h.perPage - 1) / h.perPage
	}
}

// HandleKey processes a key event while the modal is active. Returns true if
// ESC was pressed (meaning dismiss).
func (h *HelpOverlay) HandleKey(key rune) bool {
	const (
		leftKey  = 0x08
		rightKey = 0x15
	)

	switch key {
	case 0x1B: // ESC
		return true
	case leftKey:
		if h.page > 0 {
			h.page--
		}
	case rightKey:
		if h.page < h.totalPages-1 {
			h.page++
		}
	}
	return false
}

// Draw renders the help modal.
func (h *HelpOverlay) Draw(screen *ebiten.Image) {
	if !h.active || h.faceSource == nil {
		return
	}

	face := &text.GoTextFace{
		Source: h.faceSource,
		Size:   helpFontSize,
	}

	_, lineHeight := text.Measure("X", face, 0)

	// Calculate modal dimensions
	modalWidth := float64(h.screenWidth) * 0.8
	modalHeight := float64(h.screenHeight) * 0.8
	modalX := (float64(h.screenWidth) - modalWidth) / 2
	modalY := (float64(h.screenHeight) - modalHeight) / 2

	// Draw black background
	vector.DrawFilledRect(screen, float32(modalX), float32(modalY),
		float32(modalWidth), float32(modalHeight),
		color.Black, false)

	// Draw white border (1px)
	white := color.White
	vector.StrokeRect(screen, float32(modalX), float32(modalY),
		float32(modalWidth), float32(modalHeight),
		1, white, false)

	// Draw title centered
	titleWidth, _ := text.Measure(helpTitle, face, 0)
	titleX := modalX + (modalWidth-titleWidth)/2
	titleY := modalY + helpPadding

	titleOpts := &text.DrawOptions{}
	titleOpts.GeoM.Translate(titleX, titleY)
	text.Draw(screen, helpTitle, face, titleOpts)

	// Draw entries for current page
	entryY := titleY + lineHeight + helpPadding
	startIdx := h.page * h.perPage
	endIdx := startIdx + h.perPage
	if endIdx > len(helpEntries) {
		endIdx = len(helpEntries)
	}

	keyColumnX := modalX + helpPadding
	descColumnX := modalX + helpPadding + 140

	for i := startIdx; i < endIdx; i++ {
		entry := helpEntries[i]

		keyOpts := &text.DrawOptions{}
		keyOpts.GeoM.Translate(keyColumnX, entryY)
		text.Draw(screen, entry.key, face, keyOpts)

		descOpts := &text.DrawOptions{}
		descOpts.GeoM.Translate(descColumnX, entryY)
		text.Draw(screen, entry.desc, face, descOpts)

		entryY += lineHeight
	}

	// Draw footer if paginated
	if h.totalPages > 1 {
		footerY := modalY + modalHeight - helpPadding - lineHeight
		footer := "LEFT/RIGHT to page, ESC to close"

		pageInfo := "Page " + itoa(h.page+1) + " of " + itoa(h.totalPages)
		pageInfoWidth, _ := text.Measure(pageInfo, face, 0)
		pageInfoX := modalX + (modalWidth-pageInfoWidth)/2

		pageOpts := &text.DrawOptions{}
		pageOpts.GeoM.Translate(pageInfoX, footerY)
		text.Draw(screen, pageInfo, face, pageOpts)

		footerWidth, _ := text.Measure(footer, face, 0)
		footerX := modalX + (modalWidth-footerWidth)/2
		footerOpts := &text.DrawOptions{}
		footerOpts.GeoM.Translate(footerX, footerY+lineHeight)
		text.Draw(screen, footer, face, footerOpts)
	}
}

func itoa(n int) string {
	if n < 0 {
		return "-" + itoa(-n)
	}
	if n < 10 {
		return string(rune('0' + n))
	}
	return itoa(n/10) + string(rune('0'+n%10))
}
