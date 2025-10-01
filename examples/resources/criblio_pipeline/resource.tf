resource "criblio_pipeline" "my_pipeline" {
  conf = {
    async_func_timeout = 300
    description        = "myPipelineDescription"
    functions = [
      {
        conf = {
          key = jsonencode("value")
        }
        description = "My pipeline fuction configuration description"
        disabled    = true
        filter      = "truthy"
        final       = true
        group_id    = "myUniqueGroupId"
        id          = "myPipelineFunctionConf"
      }
    ]
    groups = {
      key = {
        description = "My short description for this pipeline group"
        disabled    = true
        name        = "myGroupName"
      }
    }
    output = "myOutputDestination"
    streamtags = [
      "my",
      "tags",
    ]
  }
  group_id = "myExistingGroupId"
  id       = "myPipelineId"
}