package keyboard

import (
	"github.com/eiannone/keyboard"

	"github.com/thirdmartini/gogui/pkg/ux"
)

type Keyboard struct {
}

func (k *Keyboard) Listen(OnEvent func(ev *ux.Event)) error {
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		switch key {
		case keyboard.KeyCtrlC:
			OnEvent(ux.NewSystemEvent(ux.EventKindQuit))

		case keyboard.KeyArrowRight:
			OnEvent(ux.NewKeyPressEvent(ux.KeyPressRight))
		case keyboard.KeyArrowLeft:
			OnEvent(ux.NewKeyPressEvent(ux.KeyPressLeft))
		case keyboard.KeyArrowUp:
			OnEvent(ux.NewKeyPressEvent(ux.KeyPressUp))
		case keyboard.KeyArrowDown:
			OnEvent(ux.NewKeyPressEvent(ux.KeyPressDown))
		}
	}
}

func (k *Keyboard) Close() error {
	return keyboard.Close()
}

func NewKeyboard() (*Keyboard, error) {
	if err := keyboard.Open(); err != nil {
		return nil, err
	}

	k := &Keyboard{}

	return k, nil
}
