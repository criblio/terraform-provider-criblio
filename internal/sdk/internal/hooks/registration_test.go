package hooks

import (
	"testing"
)

func TestInitHooks(t *testing.T) {
	var newHook Hooks
	var preBRH, postBRH, preAEH, postAEH, preASH, postASH int

	for range newHook.beforeRequestHook {
		preBRH++
	}
	for range newHook.afterErrorHook {
		preAEH++
	}
	for range newHook.afterSuccessHook {
		preASH++
	}

	initHooks(&newHook)

	for range newHook.beforeRequestHook {
		postBRH++
	}
	for range newHook.afterErrorHook {
		postAEH++
	}
	for range newHook.afterSuccessHook {
		postASH++
	}

	if preBRH != 0 {
		t.Errorf("pre-init Hook.beforeRequestHook struct returned %d, expected %d", preBRH, 0)
	}
	if preAEH != 0 {
		t.Errorf("pre-init Hook.afterErrorHook struct returned %d, expected %d", preBRH, 0)
	}
	if preASH != 0 {
		t.Errorf("pre-init Hook.afterSuccessHook struct returned %d, expected %d", preBRH, 0)
	}
	if postBRH != 1 {
		t.Errorf("post-init Hook.beforeRequestHook struct returned %d, expected %d", preBRH, 1)
	}
	if postAEH != 1 {
		t.Errorf("post-init Hook.afterErrorHook struct returned %d, expected %d", preBRH, 1)
	}
	if postASH != 0 {
		t.Errorf("post-init Hook.afterSuccessHook struct returned %d, expected %d", preBRH, 0)
	}

}
