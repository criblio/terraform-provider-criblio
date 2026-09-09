variable "existing_lakehouse_id" {
  description = "ID of an existing Cribl Lake Lakehouse. New Lakehouses can no longer be created."
  type        = string
}

data "criblio_cribl_lake_house" "existing" {
  id = var.existing_lakehouse_id
}
