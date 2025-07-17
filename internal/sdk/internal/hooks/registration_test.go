package hooks

import (
	"testing"
)

func TestInitHooks(t *testing.T) {
	var newHook Hooks
	var preBRH, postBRH, preAEH, postAEH, preASH, postASH int

	for _, _ = range newHook.beforeRequestHook {
		preBRH++
	}
	for _, _ = range newHook.afterErrorHook {
		preAEH++
	}
	for _, _ = range newHook.afterSuccessHook {
		preASH++
	}

	initHooks(&newHook)

	for _, _ = range newHook.beforeRequestHook {
		postBRH++
	}
	for _, _ = range newHook.afterErrorHook {
		postAEH++
	}
	for _, _ = range newHook.afterSuccessHook {
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
