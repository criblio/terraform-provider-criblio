package hooks

/*
 * This file is only ever generated once on the first generation and then is free to be modified.
 * Any hooks you wish to add should be registered in the initHooks function. Feel free to define
 * your hooks in this file or in separate files in the hooks package.
 *
 * Hooks are registered per SDK instance, and are valid for the lifetime of the SDK instance.
 */

func initHooks(h *Hooks) {
	// Register combined Cribl Terraform hook for both authentication and URL routing
	criblHook := NewCriblTerraformHook()

	// Clear existing hooks to ensure our hook is first
	h.sdkInitHooks = []sdkInitHook{criblHook}
	h.beforeRequestHook = []beforeRequestHook{criblHook}
	h.afterErrorHook = []afterErrorHook{criblHook}
}
