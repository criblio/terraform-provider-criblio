resource "criblio_lookup_file" "my_lookups" {
  content     = <<-CSV
column1,column2,column3,column4
value1,value2,value3,value4
alpha,beta,gamma,delta
CSV
  description = "my_description"
  group_id    = "default"
  id          = "my_id"
  mode        = "memory"
  tags        = "my_tags"
}
