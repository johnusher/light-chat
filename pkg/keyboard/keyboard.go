package keyboard

import (
	"github.com/nsf/termbox-go"
	log "github.com/sirupsen/logrus"
)

// Based on: https://github.com/nsf/termbox-go/blob/master/_demos/keyboard.go

type Keyboard struct {
	keys chan<- rune
}

func Init(keys chan<- rune) (*Keyboard, error) {
	return &Keyboard{keys}, nil
}

func (k *Keyboard) Run() error {
	err := termbox.Init()
	if err != nil {
		log.Errorf("termbox init failed; %s", err)
		return err
	}

	defer func() {
		close(k.keys)
		termbox.Close()
	}()

	termbox.Flush()

	for {

		switch ev := termbox.PollEvent(); ev.Type {

		case termbox.EventKey:
			// close termbox when "q" is pressed
			if string(ev.Ch) == "q" {
				termbox.Close()
				return nil
			}

			// if (string(ev.Ch) == "0") || (string(ev.Ch) == "1") {
			// 	k.keys <- ev.Ch
			// 	termbox.Flush()
			// }

			k.keys <- ev.Ch
			termbox.Flush()

		case termbox.EventError:
			log.Errorf("termbox error: %s", ev.Err)
			return ev.Err
		}
	}
}
