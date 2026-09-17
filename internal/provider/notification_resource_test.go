package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestNotificationEmailRecipientIsOptionalForNonEmailTargets(t *testing.T) {
	var response resource.SchemaResponse
	(&NotificationResource{}).Schema(context.Background(), resource.SchemaRequest{}, &response)

	targetConfigs, ok := response.Schema.Attributes["target_configs"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("target_configs schema type = %T, want schema.ListNestedAttribute", response.Schema.Attributes["target_configs"])
	}
	conf, ok := targetConfigs.NestedObject.Attributes["conf"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("target_configs.conf schema type = %T, want schema.SingleNestedAttribute", targetConfigs.NestedObject.Attributes["conf"])
	}
	emailRecipient, ok := conf.Attributes["email_recipient"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("target_configs.conf.email_recipient schema type = %T, want schema.SingleNestedAttribute", conf.Attributes["email_recipient"])
	}
	if emailRecipient.IsRequired() || !emailRecipient.IsOptional() || !emailRecipient.IsComputed() {
		t.Fatalf("email_recipient flags = required:%v optional:%v computed:%v, want false/true/true", emailRecipient.IsRequired(), emailRecipient.IsOptional(), emailRecipient.IsComputed())
	}
}
