module "awsnetworks" {
  source      = "git::ssh://git@bitbucket.org/cribl/criblcoffee-terraform.git//modules/aws/awsnetwork?ref=main"
  say_my_name = local.say_my_name
  usecase     = "terraform-provider-criblio-test"
  type        = "Development"
  cidr_block  = "0.0.0.0/0"  # Adjust as needed for your security requirements
}

# Store the user_data script in a local file to prevent instance recreation
# This file will only be created/updated when you explicitly want to update the instances
resource "local_file" "user_data_script" {
  content  = module.cribl_worker.user_data_script
  filename = "${path.module}/user_data_script.sh"
  
  lifecycle {
    ignore_changes = [content]
  }
}

module "linux_instances" {
  source = "git::ssh://git@bitbucket.org/cribl/criblcoffee-terraform.git//modules/aws/linux?ref=main"
  depends_on = [module.awsnetworks,criblio_group.my_group_defaulthybrid,module.bootstrap_token]
  # Required variables
  vpc_id    = module.awsnetworks.vpcid
  subnet_id = module.awsnetworks.subnet_id
  key_name  = var.aws_key_name  # Use variable instead of hardcoded value
  usecase   = "terraform-provider-criblio-test"
  type      = "Development"
  owner     = local.say_my_name
  
  # Optional variables
  instance_count = var.instance_count
  instance_type  = var.instance_type
  name_prefix    = "CriblioTest"
  assign_public_ip = true
  
  # If the module supports ami_id parameter, uncomment this:
  # ami_id = var.ami_id != "" ? var.ami_id : null
  
  # Security group ingress rules
  ingress = {
    22   = ["52.20.159.210/32","3.133.183.154/32","13.57.82.243/32","35.84.240.155/32","13.57.82.243/32","141.155.142.137/32"],
    4200 = ["10.0.0.0/8"]  # Cribl port from VPC
  }
  
  # Optional: Create IAM role for instances
  create_instance_role = true
  aws_iam_role_policy_attachment = [
    "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore",
    "arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"
  ]
  # Use the local file content instead of the module output directly
  # This prevents instance recreation when the module output changes
  template_file = local_file.user_data_script.content
}
