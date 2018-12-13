package winput

import (
	"unsafe"

	"github.com/janmir/go-win32"
)

const (
	//HotKeyCopy Shortcut for copy
	HotKeyCopy = iota + 1
	//HotKeyPaste Shortcut for paste
	HotKeyPaste
	//HotKeyCut Shortcut for cut
	HotKeyCut
	//HotKeySelectAll Shortcut for select all
	HotKeySelectAll
	//HotKeySave Shortcut for save
	HotKeySave
	//HotKeyRedo Shortcut for red
	HotKeyRedo
	//HotKeyUndo Shortcut for undo
	HotKeyUndo
	//HotKeyStart Shortcut for start of text
	HotKeyStart
	//HotKeyEnd Shortcut for end
	HotKeyEnd
	//HotKeyAlt Shortcut for alt key press
	HotKeyAlt
	//HotKeyBackspace Shortcut for backspace
	HotKeyBackspace
	//HotKeySpace Shortcut for space key
	HotKeySpace
	//HotKeyTab Shortcut for tab
	HotKeyTab
	//HotKeyEnter Shortcut for Return key
	HotKeyEnter
	//HotKeyCapsLock Shortcut to toggle capslock
	HotKeyCapsLock

	//ModifierShift shift key
	ModifierShift = 1
	//ModifierCtrl control key
	ModifierCtrl = 2
	//ModifierAlt alternate key
	ModifierAlt = 4
	//ModifierHankaku Hankaku key
	ModifierHankaku = 6
	//ModifierReserved1 Reserved key
	ModifierReserved1 = 16
	//ModifierReserved2 Reserved key
	ModifierReserved2 = 32

	//VKShift shift key
	VKShift = 0x10
	//VKControl ctrl key
	VKControl = 0x11
	//VKAlt alt key
	VKAlt = 0x12
	//VKLeft left arrow key
	VKLeft = 0x25
	//VKRight right arrow key
	VKRight = 0x27
	//VKBackspace backspace key
	VKBackspace = 0x08
	//VKEnter Enter key
	VKEnter = 0x0D
	//VKSpace Space key
	VKSpace = 0x20
	//VKTab Tab key
	VKTab = 0x09
	//VKCapsLock Capslock key key
	VKCapsLock = 0x14

	_maxASCII = 255
)

var (
	inputKSize uintptr
)

//Input Handles windows input
type Input struct {
	win win32.Win32
	kbd win32.HKL
}

func init() {
	inputKSize = unsafe.Sizeof(win32.KEYBD_INPUT{})
}

//New creates a new Input structure
func New() Input {
	x := Input{}
	w := win32.New()

	//retrievs the active keyboard for the locale
	x.kbd = w.GetKeyboardLayout(0) //0 -> for current thread

	x.win = w
	return x
}

//HotKey keyboard shortcut keys
func (i Input) HotKey(hot int) bool {
	var (
		inputs   = make([]win32.KEYBD_INPUT, 2)
		key      rune
		modifier = win32.VK_CONTROL
		vcode    int
	)

	switch hot {
	case HotKeyCopy:
		key = 'c' //c for copy
	case HotKeyPaste:
		key = 'v' //v for paste
	case HotKeyCut:
		key = 'x' //x for cut
	case HotKeySelectAll:
		key = 'a' //a for select all
	case HotKeySave:
		key = 's' //s for save
	case HotKeyRedo:
		key = 'y' //y for redo
	case HotKeyUndo:
		key = 'z' //z for undo
	case HotKeyStart:
		vcode = VKLeft //left key
	case HotKeyEnd:
		vcode = VKRight //right key
	case HotKeyAlt:
		vcode = VKAlt //Alt key
		modifier = 0
	case HotKeyBackspace:
		vcode = VKBackspace //Backspace key
		modifier = 0
	case HotKeySpace:
		vcode = VKSpace //Space bar key
		modifier = 0
	case HotKeyEnter:
		vcode = VKEnter //Enter key
		modifier = 0
	case HotKeyTab:
		vcode = VKTab //Tab key
		modifier = 0
	case HotKeyCapsLock:
		vcode = VKCapsLock //Tab key
		modifier = 0
	}

	if vcode == 0 {
		_, vcode = i.getVKey(key)
	}
	inputs = append(inputs, i.KeyDown(vcode, modifier)...)
	inputs = append(inputs, i.KeyUp(vcode, modifier)...)

	return i.press(inputs)
}

//Type loads all characters of the string one-by-one
func (i Input) Type(word string) bool {
	inputs := make([]win32.KEYBD_INPUT, 0)

	//load characters as inputs
	for _, v := range word {
		mod, vcode := i.getVKey(v)

		down := i.KeyDown(vcode, mod)
		inputs = append(inputs, down...)
		if vcode < _maxASCII {
			up := i.KeyUp(vcode, mod)
			inputs = append(inputs, up...)
		}
	}

	//press all the loaded keys
	return i.press(inputs)
}

//press press key up and then down
func (i Input) press(inputs []win32.KEYBD_INPUT) bool {
	inputLen := len(inputs)
	count := -1
	if inputLen > 0 {
		c := i.win.SendInput(inputLen, unsafe.Pointer(&inputs[0]), inputKSize)
		count = int(c)
	}
	return count == inputLen
}

//KeyDown press key up and then down
func (i Input) KeyDown(vcode, mod int) []win32.KEYBD_INPUT {
	return i.input(vcode, mod, win32.KEYEVENTF_KEYDOWN)
}

//KeyUp press key up and then down
func (i Input) KeyUp(vcode, mod int) []win32.KEYBD_INPUT {
	return i.input(vcode, mod, win32.KEYEVENTF_KEYUP)
}

//input input keys
func (i Input) input(key, mod, event int) []win32.KEYBD_INPUT {
	var press win32.KEYBDINPUT
	inputs := make([]win32.KEYBD_INPUT, 0)

	//For modifier key down
	if event == win32.KEYEVENTF_KEYDOWN && mod > 0 {
		press = win32.KEYBDINPUT{
			WVk:     uint16(mod), //keycode
			DwFlags: uint32(win32.KEYEVENTF_KEYDOWN),
		}

		inputs = append(inputs, win32.KEYBD_INPUT{
			Type: win32.INPUT_KEYBOARD,
			Ki:   press,
		})
	}

	//The key itself
	var vkey = uint16(key)
	var scan = uint16(0)
	var eevent = uint32(event)
	if vkey >= _maxASCII {
		vkey = scan
		scan = uint16(key)
		eevent = uint32(win32.KEYEVENTF_UNICODE)
	}
	press = win32.KEYBDINPUT{
		WVk:     vkey, //keycode
		WScan:   scan,
		DwFlags: eevent,
	}

	inputs = append(inputs, win32.KEYBD_INPUT{
		Type: win32.INPUT_KEYBOARD,
		Ki:   press,
	})

	//For modifier key up
	if event == win32.KEYEVENTF_KEYUP && mod > 0 {
		press = win32.KEYBDINPUT{
			WVk:     uint16(mod), //keycode
			DwFlags: uint32(win32.KEYEVENTF_KEYUP),
		}

		inputs = append(inputs, win32.KEYBD_INPUT{
			Type: win32.INPUT_KEYBOARD,
			Ki:   press,
		})
	}

	return inputs[:]
}

//getVKey translate a character rune to its win virtual keycode
//and shift state
func (i Input) getVKey(key rune) (int, int) {
	mod, vcode := i.win.VkKeyScanEx(key, i.kbd)

	switch mod {
	case ModifierShift: //shift
		mod = VKShift
	case ModifierCtrl: //ctrl
		mod = VKControl
	case ModifierAlt: //alt
		mod = VKAlt
	case ModifierHankaku: //Hankaku unsupported
		fallthrough
	case ModifierReserved1: //Reserved
		fallthrough
	case ModifierReserved2: //Reserved
		fallthrough
	default:
		mod = 0
	}
	if vcode >= _maxASCII {
		vcode = int(key)
	}

	return mod, vcode
}
