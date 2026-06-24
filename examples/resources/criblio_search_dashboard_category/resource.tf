resource "criblio_search_dashboard_category" "my_searchdashboardcategory" {
  description = "My dashboard category description"
  group_id    = "default_search"
  id          = "myDashboardCategoryId"
  is_pack     = false
  name        = "MyDashboardCategoryName"
}
