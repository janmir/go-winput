package winput

import (
	"testing"
	"time"

	util "github.com/janmir/go-util"
	"go.uber.org/goleak"
)

func TestAll(t *testing.T) {
	defer goleak.VerifyNoLeaks(t)
	in := New()
	time.Sleep(time.Second * 2)

	word := "Hello World!!@#$%^&*()_++_)(*&^%$#こんにちは世界"
	ok := in.Type(word)
	if !ok {
		t.Fail()
	}

	util.Logger("Select All")
	time.Sleep(time.Second * 2)
	ok = in.HotKey(HotKeySelectAll)
	if !ok {
		t.Fail()
	}
	util.Logger("Copy")
	time.Sleep(time.Second * 2)

	ok = in.HotKey(HotKeyCopy)
	if !ok {
		t.Fail()
	}
	util.Logger("Move to right")
	time.Sleep(time.Second * 2)
	ok = in.HotKey(HotKeyEnd)
	if !ok {
		t.Fail()
	}
	util.Logger("Paste")
	time.Sleep(time.Second * 2)
	ok = in.HotKey(HotKeyPaste)
	if !ok {
		t.Fail()
	}

	util.Logger("Select All")
	time.Sleep(time.Second * 2)
	ok = in.HotKey(HotKeySelectAll)
	if !ok {
		t.Fail()
	}

	util.Logger("Cut")
	time.Sleep(time.Second * 2)
	ok = in.HotKey(HotKeyCut)
	if !ok {
		t.Fail()
	}
	util.Logger("Paste")
	time.Sleep(time.Second * 2)
	ok = in.HotKey(HotKeyPaste)
	if !ok {
		t.Fail()
	}
}

func TestKeyType(t *testing.T) {
	defer goleak.VerifyNoLeaks(t)

	in := New()
	word := "Hello World!...????@#$%^&*()"
	ok := in.Type(word)
	if !ok {
		t.Fail()
	}
}

func TestKeyUnicode(t *testing.T) {
	defer goleak.VerifyNoLeaks(t)

	in := New()
	word := "Hello World!こんにちは世界"
	time.Sleep(2 * time.Second)
	ok := in.Type(word)
	if !ok {
		t.Fail()
	}
}

func TestKeyCaps(t *testing.T) {
	defer goleak.VerifyNoLeaks(t)

	in := New()
	word := "hello world!"

	time.Sleep(2 * time.Second)
	ok := in.Type(word)
	if !ok {
		t.Fail()
	}

	//toggle capslock
	ok = in.HotKey(HotKeyCapsLock)
	if !ok {
		t.Fail()
	}

	time.Sleep(2 * time.Second)
	ok = in.Type(word)
	if !ok {
		t.Fail()
	}
}
