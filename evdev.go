package talkiepi

import (
  //"evdev"
  "github.com/gvalkov/golang-evdev"
  "fmt"
  "time"
)

func EVOpen() {
  device, _ := evdev.Open("/dev/input/event5")
  fmt.Println(device)
}

// Listing accessible input devices.
func EVListInputDevices() {
	devices, _ := evdev.ListInputDevices()

	for _, dev := range devices {
		fmt.Printf("%s %s %s", dev.Fn, dev.Name, dev.Phys)
	}
}

func format_event(ev *evdev.InputEvent) string {
	var res, f, code_name string

	code := int(ev.Code)
	etype := int(ev.Type)

	switch ev.Type {
	case evdev.EV_SYN:
		if ev.Code == evdev.SYN_MT_REPORT {
			f = "time %d.%-8d +++++++++ %s ++++++++"
		} else {
			f = "time %d.%-8d --------- %s --------"
		}
		return fmt.Sprintf(f, ev.Time.Sec, ev.Time.Usec, evdev.SYN[code])
	case evdev.EV_KEY:
		val, haskey := evdev.KEY[code]
		if haskey {
			code_name = val
		} else {
			val, haskey := evdev.BTN[code]
			if haskey {
				code_name = val
			} else {
				code_name = "?"
			}
		}
	default:
		m, haskey := evdev.ByEventType[etype]
		if haskey {
			code_name = m[code]
		} else {
			code_name = "?"
		}
	}

	evfmt := "time %d.%-8d type %d (%s), code %-3d (%s), value %d"
	res = fmt.Sprintf(evfmt, ev.Time.Sec, ev.Time.Usec, etype,
		evdev.EV[int(ev.Type)], ev.Code, code_name, ev.Value)

	return res
}

func (b *Talkiepi) EV() {
  var events []evdev.InputEvent
  var err error
  var evdevice string = "/dev/input/event5"//todo make this a config option
  var key string = "KEY_B"

  if evdev.IsInputDevice(evdevice) {
	   device, _ := evdev.Open(evdevice)

    go func() {
      for {
    		events, err = device.Read()
    		for i := range events {
          code := int(events[i].Code)
          switch events[i].Type {
          case evdev.EV_KEY: //just in case that we want other events in the future
            val, haskey := evdev.KEY[code]
            if haskey {
              if val == key { //todo: make this an config option
                switch events[i].Value {
                  case 0: //keyUP
                    b.TransmitStop()
                  case 1: //keyDOWN
                    b.TransmitStart()
                  case 2: //keyHOLD
                    break
                }
              }
            } else {
              str := format_event(&events[i])
              fmt.Println(str)
            }
          }
      	}
        time.Sleep(500 * time.Millisecond)
  	}
}()
} else {
    fmt.Println("input device is not a propper carracter event device")
}
}
