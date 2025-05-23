// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package sdk

type V5Billing struct {
	Consumption *Consumption
	Invoices    *Invoices

	sdkConfiguration sdkConfiguration
}

func newV5Billing(sdkConfig sdkConfiguration) *V5Billing {
	return &V5Billing{
		sdkConfiguration: sdkConfig,
		Consumption:      newConsumption(sdkConfig),
		Invoices:         newInvoices(sdkConfig),
	}
}
