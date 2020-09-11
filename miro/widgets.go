package miro

const (
	widgetsPath = "widgets"
)

// WidgetsService handles communication to Miro Widgets API.
//
// API doc: https://developers.miro.com/reference#board-object
type WidgetsService service

type Widget interface {
}

// Sticker object represents Miro Sticker.
//
// API doc: https://developers.miro.com/reference#sticker
//go:generate gomodifytags -file $GOFILE -struct Sticker -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Sticker -add-tags json -w -transform camelcase
type Sticker struct {
}

// Shape object represents Miro Shape.
//
// API doc: https://developers.miro.com/reference#shape
//go:generate gomodifytags -file $GOFILE -struct Shape -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Shape -add-tags json -w -transform camelcase
type Shape struct {
}

// Text object represents Miro Text.
//
// API doc: https://developers.miro.com/reference#text
//go:generate gomodifytags -file $GOFILE -struct Text -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Text -add-tags json -w -transform camelcase
type Text struct {
}

// Line object represents Miro Line.
//
// API doc: https://developers.miro.com/reference#line
//go:generate gomodifytags -file $GOFILE -struct Line -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Line -add-tags json -w -transform camelcase
type Line struct {
}

// Card object represents Miro Card.
//
// API doc: https://developers.miro.com/reference#card
//go:generate gomodifytags -file $GOFILE -struct Card -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Card -add-tags json -w -transform camelcase
type Card struct {
}
