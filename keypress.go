package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

/** syscalls_windows ****************************************************/
const (
	foreground_blue          = 0x1
	foreground_green         = 0x2
	foreground_red           = 0x4
	foreground_intensity     = 0x8
	background_blue          = 0x10
	background_green         = 0x20
	background_red           = 0x40
	background_intensity     = 0x80
	std_input_handle         = -0xa
	std_output_handle        = -0xb
	key_event                = 0x1
	mouse_event              = 0x2
	window_buffer_size_event = 0x4
	enable_window_input      = 0x8
	enable_mouse_input       = 0x10
	enable_extended_flags    = 0x80

	vk_f1          = 0x70
	vk_f2          = 0x71
	vk_f3          = 0x72
	vk_f4          = 0x73
	vk_f5          = 0x74
	vk_f6          = 0x75
	vk_f7          = 0x76
	vk_f8          = 0x77
	vk_f9          = 0x78
	vk_f10         = 0x79
	vk_f11         = 0x7a
	vk_f12         = 0x7b
	vk_insert      = 0x2d
	vk_delete      = 0x2e
	vk_home        = 0x24
	vk_end         = 0x23
	vk_pgup        = 0x21
	vk_pgdn        = 0x22
	vk_arrow_up    = 0x26
	vk_arrow_down  = 0x28
	vk_arrow_left  = 0x25
	vk_arrow_right = 0x27
	vk_backspace   = 0x8
	vk_tab         = 0x9
	vk_enter       = 0xd
	vk_esc         = 0x1b
	vk_space       = 0x20

	left_alt_pressed   = 0x2
	left_ctrl_pressed  = 0x8
	right_alt_pressed  = 0x1
	right_ctrl_pressed = 0x4
	shift_pressed      = 0x10

	generic_read            = 0x80000000
	generic_write           = 0x40000000
	console_textmode_buffer = 0x1
)

/** api_common **********************************************************/
type (
	InputMode int
	EventType uint8
	Modifier  uint8
	Key       uint16
	Attribute uint16
)

// This type represents a termbox event. The 'Mod', 'Key' and 'Ch' fields are
// valid if 'Type' is EventKey. The 'Width' and 'Height' fields are valid if
// 'Type' is EventResize. The 'Err' field is valid if 'Type' is EventError.
type Event struct {
	Type   EventType // one of Event* constants
	Mod    Modifier  // one of Mod* constants or 0
	Key    Key       // one of Key* constants, invalid if 'Ch' is not 0
	Ch     rune      // a unicode character
	Width  int       // width of the screen
	Height int       // height of the screen
	Err    error     // error in case if input failed
	MouseX int       // x coord of mouse
	MouseY int       // y coord of mouse
}

// Key constants, see Event.Key field.
const (
	KeyF1 Key = 0xFFFF - iota
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyInsert
	KeyDelete
	KeyHome
	KeyEnd
	KeyPgup
	KeyPgdn
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	key_min // see terminfo
	MouseLeft
	MouseMiddle
	MouseRight
)

const (
	KeyCtrlTilde      Key = 0x00
	KeyCtrl2          Key = 0x00
	KeyCtrlSpace      Key = 0x00
	KeyCtrlA          Key = 0x01
	KeyCtrlB          Key = 0x02
	KeyCtrlC          Key = 0x03
	KeyCtrlD          Key = 0x04
	KeyCtrlE          Key = 0x05
	KeyCtrlF          Key = 0x06
	KeyCtrlG          Key = 0x07
	KeyBackspace      Key = 0x08
	KeyCtrlH          Key = 0x08
	KeyTab            Key = 0x09
	KeyCtrlI          Key = 0x09
	KeyCtrlJ          Key = 0x0A
	KeyCtrlK          Key = 0x0B
	KeyCtrlL          Key = 0x0C
	KeyEnter          Key = 0x0D
	KeyCtrlM          Key = 0x0D
	KeyCtrlN          Key = 0x0E
	KeyCtrlO          Key = 0x0F
	KeyCtrlP          Key = 0x10
	KeyCtrlQ          Key = 0x11
	KeyCtrlR          Key = 0x12
	KeyCtrlS          Key = 0x13
	KeyCtrlT          Key = 0x14
	KeyCtrlU          Key = 0x15
	KeyCtrlV          Key = 0x16
	KeyCtrlW          Key = 0x17
	KeyCtrlX          Key = 0x18
	KeyCtrlY          Key = 0x19
	KeyCtrlZ          Key = 0x1A
	KeyEsc            Key = 0x1B
	KeyCtrlLsqBracket Key = 0x1B
	KeyCtrl3          Key = 0x1B
	KeyCtrl4          Key = 0x1C
	KeyCtrlBackslash  Key = 0x1C
	KeyCtrl5          Key = 0x1D
	KeyCtrlRsqBracket Key = 0x1D
	KeyCtrl6          Key = 0x1E
	KeyCtrl7          Key = 0x1F
	KeyCtrlSlash      Key = 0x1F
	KeyCtrlUnderscore Key = 0x1F
	KeySpace          Key = 0x20
	KeyBackspace2     Key = 0x7F
	KeyCtrl8          Key = 0x7F
)

// Alt modifier constant, see Event.Mod field and SetInputMode function.
const (
	ModAlt Modifier = 0x01
)

// Input mode. See SetInputMode function.
const (
	InputEsc InputMode = 1 << iota
	InputAlt
	InputMouse
	InputCurrent InputMode = 0
)

// Event type. See Event.Type field.
const (
	EventKey EventType = iota
	EventResize
	EventMouse
	EventError
)

/** api_windows *********************************************************/

// Initializes termbox library. This function should be called before any other functions.
// After successful initialization, the library must be finalized using 'Close' function.
func Init() error {
	var err error

	in, err = syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	if err != nil {
		return err
	}
	out, err = syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil {
		return err
	}

	consolewin = true

	go input_event_producer()

	return nil
}

// Finalizes termbox library, should be called after successful initialization
// when termbox's functionality isn't required anymore.
func Close() {
	// we ignore errors here, because we can't really do anything about them
	syscall.Close(out)
}

// Wait for an event and return it. This is a blocking function call.
func PollEvent() Event {
	return <-input_comm
}

// Sets termbox input mode. Termbox has two input modes:
//
// 1. Esc input mode. When ESC sequence is in the buffer and it doesn't match
// any known sequence. ESC means KeyEsc. This is the default input mode.
//
// 2. Alt input mode. When ESC sequence is in the buffer and it doesn't match
// any known sequence. ESC enables ModAlt modifier for the next keyboard event.
//
// Both input modes can be OR'ed with Mouse mode. Setting Mouse mode bit up will
// enable mouse button click events.
//
// If 'mode' is InputCurrent, returns the current input mode. See also Input*
// constants.
func SetInputMode(mode InputMode) InputMode {
	if mode == InputCurrent {
		return input_mode
	}
	if consolewin {
		if mode&InputMouse != 0 {
			err := set_console_mode(in, enable_window_input|enable_mouse_input|enable_extended_flags)
			if err != nil {
				panic(err)
			}
		} else {
			err := set_console_mode(in, enable_window_input)
			if err != nil {
				panic(err)
			}
		}
	}

	input_mode = mode
	return input_mode
}

/** termbox *************************************************************/
type input_event struct {
	data []byte
	err  error
}

/** termbox_windows *****************************************************/
type (
	wchar uint16
	short int16
	dword uint32
	word  uint16
	coord struct {
		x short
		y short
	}
	small_rect struct {
		left   short
		top    short
		right  short
		bottom short
	}
	console_screen_buffer_info struct {
		size                coord
		cursor_position     coord
		attributes          word
		window              small_rect
		maximum_window_size coord
	}
	console_cursor_info struct {
		size    dword
		visible int32
	}
	input_record struct {
		event_type word
		_          [2]byte
		event      [16]byte
	}
	key_event_record struct {
		key_down          int32
		repeat_count      word
		virtual_key_code  word
		virtual_scan_code word
		unicode_char      wchar
		control_key_state dword
	}
	window_buffer_size_record struct {
		size coord
	}
	mouse_event_record struct {
		mouse_pos         coord
		button_state      dword
		control_key_state dword
		event_flags       dword
	}
)

const (
	mouse_lmb = 0x1
	mouse_rmb = 0x2
	mouse_mmb = 0x4 | 0x8 | 0x10
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	proc_set_console_active_screen_buffer = kernel32.NewProc("SetConsoleActiveScreenBuffer")
	proc_set_console_screen_buffer_size   = kernel32.NewProc("SetConsoleScreenBufferSize")
	proc_create_console_screen_buffer     = kernel32.NewProc("CreateConsoleScreenBuffer")
	proc_get_console_screen_buffer_info   = kernel32.NewProc("GetConsoleScreenBufferInfo")
	proc_write_console_output_character   = kernel32.NewProc("WriteConsoleOutputCharacterW")
	proc_write_console_output_attribute   = kernel32.NewProc("WriteConsoleOutputAttribute")
	proc_set_console_cursor_info          = kernel32.NewProc("SetConsoleCursorInfo")
	proc_set_console_cursor_position      = kernel32.NewProc("SetConsoleCursorPosition")
	proc_read_console_input               = kernel32.NewProc("ReadConsoleInputW")
	proc_get_console_mode                 = kernel32.NewProc("GetConsoleMode")
	proc_set_console_mode                 = kernel32.NewProc("SetConsoleMode")
	proc_fill_console_output_character    = kernel32.NewProc("FillConsoleOutputCharacterW")
	proc_fill_console_output_attribute    = kernel32.NewProc("FillConsoleOutputAttribute")
)

var (
	orig_mode   dword
	orig_screen syscall.Handle
	//back_buffer  cellbuf
	//front_buffer cellbuf
	termw      int
	termh      int
	input_mode = InputEsc
	//cursor_x     = cursor_hidden
	//cursor_y     = cursor_hidden
	//foreground   = ColorDefault
	//background   = ColorDefault
	in       syscall.Handle
	out      syscall.Handle
	attrsbuf []word
	charsbuf []wchar
	//diffbuf      []diff_msg
	beg_x        = -1
	beg_y        = -1
	beg_i        = -1
	input_comm   = make(chan Event)
	alt_mode_esc = false
	consolewin   = false

	// these ones just to prevent heap allocs at all costs
	tmp_info  console_screen_buffer_info
	tmp_arg   dword
	tmp_coord = coord{0, 0}
)

func set_console_mode(h syscall.Handle, mode dword) (err error) {
	r0, _, e1 := syscall.Syscall(proc_set_console_mode.Addr(), 2, uintptr(h), uintptr(mode), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func read_console_input(h syscall.Handle, record *input_record) (err error) {
	r0, _, e1 := syscall.Syscall6(proc_read_console_input.Addr(), 4, uintptr(h), uintptr(unsafe.Pointer(record)), 1, uintptr(unsafe.Pointer(&tmp_arg)), 0, 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func key_event_record_to_event(r *key_event_record) (Event, bool) {
	if r.key_down == 0 {
		return Event{}, false
	}

	e := Event{Type: EventKey}
	if input_mode&InputAlt != 0 {
		if alt_mode_esc {
			e.Mod = ModAlt
			alt_mode_esc = false
		}
		if r.control_key_state&(left_alt_pressed|right_alt_pressed) != 0 {
			e.Mod = ModAlt
		}
	}

	ctrlpressed := r.control_key_state&(left_ctrl_pressed|right_ctrl_pressed) != 0

	if r.virtual_key_code >= vk_f1 && r.virtual_key_code <= vk_f12 {
		switch r.virtual_key_code {
		case vk_f1:
			e.Key = KeyF1
		case vk_f2:
			e.Key = KeyF2
		case vk_f3:
			e.Key = KeyF3
		case vk_f4:
			e.Key = KeyF4
		case vk_f5:
			e.Key = KeyF5
		case vk_f6:
			e.Key = KeyF6
		case vk_f7:
			e.Key = KeyF7
		case vk_f8:
			e.Key = KeyF8
		case vk_f9:
			e.Key = KeyF9
		case vk_f10:
			e.Key = KeyF10
		case vk_f11:
			e.Key = KeyF11
		case vk_f12:
			e.Key = KeyF12
		default:
			panic("unreachable")
		}

		return e, true
	}

	if r.virtual_key_code <= vk_delete {
		switch r.virtual_key_code {
		case vk_insert:
			e.Key = KeyInsert
		case vk_delete:
			e.Key = KeyDelete
		case vk_home:
			e.Key = KeyHome
		case vk_end:
			e.Key = KeyEnd
		case vk_pgup:
			e.Key = KeyPgup
		case vk_pgdn:
			e.Key = KeyPgdn
		case vk_arrow_up:
			e.Key = KeyArrowUp
		case vk_arrow_down:
			e.Key = KeyArrowDown
		case vk_arrow_left:
			e.Key = KeyArrowLeft
		case vk_arrow_right:
			e.Key = KeyArrowRight
		case vk_backspace:
			if ctrlpressed {
				e.Key = KeyBackspace2
			} else {
				e.Key = KeyBackspace
			}
		case vk_tab:
			e.Key = KeyTab
		case vk_enter:
			e.Key = KeyEnter
		case vk_esc:
			switch {
			case input_mode&InputEsc != 0:
				e.Key = KeyEsc
			case input_mode&InputAlt != 0:
				alt_mode_esc = true
				return Event{}, false
			}
		case vk_space:
			if ctrlpressed {
				// manual return here, because KeyCtrlSpace is zero
				e.Key = KeyCtrlSpace
				return e, true
			} else {
				e.Key = KeySpace
			}
		}

		if e.Key != 0 {
			return e, true
		}
	}

	if ctrlpressed {
		if Key(r.unicode_char) >= KeyCtrlA && Key(r.unicode_char) <= KeyCtrlRsqBracket {
			e.Key = Key(r.unicode_char)
			if input_mode&InputAlt != 0 && e.Key == KeyEsc {
				alt_mode_esc = true
				return Event{}, false
			}
			return e, true
		}
		switch r.virtual_key_code {
		case 192, 50:
			// manual return here, because KeyCtrl2 is zero
			e.Key = KeyCtrl2
			return e, true
		case 51:
			if input_mode&InputAlt != 0 {
				alt_mode_esc = true
				return Event{}, false
			}
			e.Key = KeyCtrl3
		case 52:
			e.Key = KeyCtrl4
		case 53:
			e.Key = KeyCtrl5
		case 54:
			e.Key = KeyCtrl6
		case 189, 191, 55:
			e.Key = KeyCtrl7
		case 8, 56:
			e.Key = KeyCtrl8
		}

		if e.Key != 0 {
			return e, true
		}
	}

	if r.unicode_char != 0 {
		e.Ch = rune(r.unicode_char)
		return e, true
	}

	return Event{}, false
}

func input_event_producer() {
	var r input_record
	var err error
	var last_button Key
	var last_state = dword(0)
	for {
		err = read_console_input(in, &r)
		if err != nil {
			input_comm <- Event{Type: EventError, Err: err}
		}

		switch r.event_type {
		case key_event:
			kr := (*key_event_record)(unsafe.Pointer(&r.event))
			ev, ok := key_event_record_to_event(kr)
			if ok {
				for i := 0; i < int(kr.repeat_count); i++ {
					input_comm <- ev
				}
			}
		case window_buffer_size_event:
			sr := *(*window_buffer_size_record)(unsafe.Pointer(&r.event))
			input_comm <- Event{
				Type:   EventResize,
				Width:  int(sr.size.x),
				Height: int(sr.size.y),
			}
		case mouse_event:
			mr := *(*mouse_event_record)(unsafe.Pointer(&r.event))

			// single or double click
			switch mr.event_flags {
			case 0:
				cur_state := mr.button_state
				switch {
				case last_state&mouse_lmb == 0 && cur_state&mouse_lmb != 0:
					last_button = MouseLeft
				case last_state&mouse_rmb == 0 && cur_state&mouse_rmb != 0:
					last_button = MouseRight
				case last_state&mouse_mmb == 0 && cur_state&mouse_mmb != 0:
					last_button = MouseMiddle
				default:
					last_state = cur_state
					continue
				}
				last_state = cur_state
				fallthrough
			case 2:
				input_comm <- Event{
					Type:   EventMouse,
					Key:    last_button,
					MouseX: int(mr.mouse_pos.x),
					MouseY: int(mr.mouse_pos.y),
				}
			}
		}
	}
}

func main() {
	err := Init()
	if err != nil {
		panic(err)
	}
	defer Close()

	SetInputMode(InputEsc | InputMouse)

mainloop:
	for {
		switch ev := PollEvent(); ev.Type {
		case EventKey:
			switch ev.Key {
			case KeyEsc:
				break mainloop
			default:
				fmt.Printf("Char string: %s\n", string(ev.Ch))
			}
		case EventMouse:
			fmt.Printf("EventMouse: %s\n", string(ev.Ch))
		case EventError:
			panic(ev.Err)
		}
	}
}
