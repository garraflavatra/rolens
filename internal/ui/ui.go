package ui

import (
	"context"
	"runtime"

	"github.com/gen2brain/beeep"
	"github.com/ncruces/zenity"
)

type UI struct {
	ctx          context.Context
	progressBars map[uint]zenity.ProgressDialog
}

func New() *UI {
	return &UI{
		progressBars: make(map[uint]zenity.ProgressDialog),
	}
}

func (u *UI) Startup(ctx context.Context) {
	u.ctx = ctx
}

func (u *UI) Beep() {
	if runtime.GOOS == "windows" {
		return
	}
	beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
}
