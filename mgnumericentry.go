// Copyright (C) 2025-2026 Murilo Gomes Julio
// SPDX-License-Identifier: MIT

// Site: https://mugomes.github.io

package mgnumericentry

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// ---------------- MGSmallButton (botão pequeno com auto-repeat) ----------------

type MGSmallButton struct {
	widget.BaseWidget
	Label    string
	OnTapped func()
	hover    bool
	repeat   bool
	stop     chan struct{}
}

func newMGSmallButton(label string, onTap func()) *MGSmallButton {
	b := &MGSmallButton{
		Label:    label,
		OnTapped: onTap,
		stop:     make(chan struct{}, 1),
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *MGSmallButton) CreateRenderer() fyne.WidgetRenderer {
	txt := canvas.NewText(b.Label, color.Black)
	txt.Alignment = fyne.TextAlignCenter
	bg := canvas.NewRectangle(color.NRGBA{201, 201, 201, 255})

	objects := []fyne.CanvasObject{bg, txt}
	return &mgSmallButtonRenderer{
		button:  b,
		bg:      bg,
		txt:     txt,
		objects: objects,
	}
}

func (b *MGSmallButton) MinSize() fyne.Size {
	return fyne.NewSize(18, 14)
}

func (b *MGSmallButton) Tapped(*fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *MGSmallButton) MouseIn(*desktop.MouseEvent) {
	b.hover = true
	b.Refresh()
}

func (b *MGSmallButton) MouseOut() {
	b.hover = false
	b.Refresh()
}

func (b *MGSmallButton) MouseDown(ev *desktop.MouseEvent) {
	if ev.Button != desktop.MouseButtonPrimary {
		return
	}

	// start repeat after short delay
	b.repeat = true
	go func() {
		// initial delay
		select {
		case <-time.After(400 * time.Millisecond):
		case <-b.stop:
			if !b.repeat {
				return
			}
		}
		for {
			select {
			case <-b.stop:
				return
			default:
				if b.OnTapped != nil {
					b.OnTapped()
				}
				time.Sleep(70 * time.Millisecond)
			}
		}
	}()
}

func (b *MGSmallButton) MouseUp(ev *desktop.MouseEvent) {
	if !b.repeat {
		return
	}
	b.repeat = false
	// stop goroutine
	select {
	case b.stop <- struct{}{}:
	default:
	}
}

type mgSmallButtonRenderer struct {
	button  *MGSmallButton
	bg      *canvas.Rectangle
	txt     *canvas.Text
	objects []fyne.CanvasObject
}

func (r *mgSmallButtonRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.txt.Resize(size)
	r.txt.Move(fyne.NewPos(0, 0))
}

func (r *mgSmallButtonRenderer) MinSize() fyne.Size {
	return r.button.MinSize()
}

func (r *mgSmallButtonRenderer) Refresh() {
	if r.button.hover {
		r.bg.FillColor = color.NRGBA{200, 200, 200, 255}
	} else {
		r.bg.FillColor = color.NRGBA{230, 230, 230, 255}
	}
	r.bg.Refresh()
	r.txt.Refresh()
}

func (r *mgSmallButtonRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *mgSmallButtonRenderer) Destroy()                     {}

// ---------------- MGNumericEntry (Entry numérico) ----------------

type MGNumericEntry struct {
	widget.Entry
	min       int
	max       int
	step      int
	value     int
	OnChanged func(int)
}

func NewMGNumericEntry(min, max, initial int) *MGNumericEntry {
	e := &MGNumericEntry{
		min:   min,
		max:   max,
		step:  1,
		value: initial,
	}
	e.ExtendBaseWidget(e)

	// set initial text (clamped)
	if initial < min {
		e.value = min
	} else if initial > max {
		e.value = max
	} else {
		e.value = initial
	}
	e.Entry.SetText(fmt.Sprintf("%d", e.value))

	return e
}

// Permite apenas digitação de dígitos e sinal '-' (se min < 0)
func (e *MGNumericEntry) TypedRune(r rune) {
	// permitir digitar '-' só no início e quando min < 0
	if r == '-' {
		if e.min >= 0 {
			return
		}
		// allow '-' only if it's first char
		if len(e.Text) != 0 {
			return
		}
		e.Entry.TypedRune(r)
		return
	}
	if r >= '0' && r <= '9' {
		e.Entry.TypedRune(r)
	} else {
		// ignora outros runes
		return
	}
	e.updateValueFromText()
}

func (e *MGNumericEntry) TypedKey(k *fyne.KeyEvent) {
	// permitir backspace/delete
	if k.Name == fyne.KeyBackspace || k.Name == fyne.KeyDelete {
		e.Entry.TypedKey(k)
		e.updateValueFromText()
		return
	}
	// permitir ← → etc repassando
	e.Entry.TypedKey(k)
	// no geral, atualiza
	e.updateValueFromText()
}

func (e *MGNumericEntry) updateValueFromText() {
	txt := strings.TrimSpace(e.Text)
	if txt == "" || txt == "-" {
		// treat empty or lone '-' as min (or 0?) -> não atualizar ainda
		return
	}
	v, err := strconv.Atoi(txt)
	if err != nil {
		return
	}
	if v < e.min {
		v = e.min
	}
	if v > e.max {
		v = e.max
	}
	if v != e.value {
		e.value = v
		if e.OnChanged != nil {
			e.OnChanged(e.value)
		}
	}
	// keep displayed text clamped (se necessário)
	if fmt.Sprintf("%d", e.value) != txt {
		e.Entry.SetText(fmt.Sprintf("%d", e.value))
	}
}

func (e *MGNumericEntry) GetValue() int {
	return e.value
}

func (e *MGNumericEntry) SetValue(v int) {
	if v < e.min {
		v = e.min
	}
	if v > e.max {
		v = e.max
	}
	e.value = v
	fyne.Do(func() {
		e.Entry.SetText(fmt.Sprintf("%d", e.value))
	})
	if e.OnChanged != nil {
		e.OnChanged(e.value)
	}
}

func (e *MGNumericEntry) increment() {
	next := e.value + e.step
	if next > e.max {
		next = e.max
	}
	e.SetValue(next)
}

func (e *MGNumericEntry) decrement() {
	next := e.value - e.step
	if next < e.min {
		next = e.min
	}
	e.SetValue(next)
}

// ---------------- Combinação: Entry + SpinButtons ----------------

func NewMGNumericEntryWithButtons(min, max, initial int) (fyne.CanvasObject, *MGNumericEntry) {
	entry := NewMGNumericEntry(min, max, initial)

	up := newMGSmallButton("▲", func() { entry.increment() })
	down := newMGSmallButton("▼", func() { entry.decrement() })

	// container vertical com sem espaçamento extra
	spin := container.NewVBox(up, down)
	// estreitar margens internas: usar NewHBox com entry + spin
	return container.NewHBox(entry, spin), entry
}
