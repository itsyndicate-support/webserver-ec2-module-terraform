terraform {
    source = "../../../infrastructure"
}

generate "outputs_for_tests" {
    path = "${get_terragrunt_dir()}/../../../infrastructure/110-test-outputs.tf"
    if_exists = "overwrite"
    contents = <<EOF
output "http_instances_state" {
    value = alltrue([
        for instance in aws_instance.http :
            instance.instance_state == "running" ? true : false
    ])
}

output "db_instances_state" {
    value = alltrue([
        for instance in aws_instance.db :
            instance.instance_state == "running" ? true : false
    ])
}

# Must be false
output "db_public_access" {
    value = anytrue([
        for db in aws_instance.db :
            db.associate_public_ip_address
    ])
}
EOF
}


inputs = {
    vpc_cidr = "200.100.0.0/16"
    network_http = {
        subnet_name = "test_subnet_http"
        cidr        = "200.100.1.0/24"
    }
    network_db = {
        subnet_name = "test_subnet_db"
        cidr        = "200.100.2.0/24"
    }
    http_instance_names = ["test-instance-http-1", "test-instance-http-2"]
    db_instance_names = ["test-instance-db-1", "test-instance-db-2", "test-instance-db-3"]
    public_key_name = "ci_test_key"
    publick_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC4VBjwn+I08o8R5CluS+rby7wX99AfHh+WYJz0eQgAK8IA2A8PFEXBHekszigoGVYismEDy75ZnbnvbbSx00W0IqaUOgpV4tcM+Odwx1SaMjv7CISFWncWV64xKNxM0ieyuLzrtwInBPdSIMpsQCwhx4WhMFHu8PFtAlCCpwAdxb7k0T2StVIFgf2R6A/yhDuZ6hK2QL4L7vM5MfHyxDdB538oswBZ07NZnfVMg9ZfxxRsJbO6z7Xpjemn5+Hyna9KmV7y6X9bRCAfZNDJWYq79F7Os0HcMuHq+/PEyzbqYu/wyKpHxmyCPvpyShyWRD1j3Caio7io6spHUMXonjXMpQgsBeGTpfE7QVgGCHhprTMZAbHkREPXGivKpxE1oVuuOzElCunXPqALogvXwHun0W0SjpoNycBIldDy0xR1hFXf2rjabPovwu5Is2YfvrqWvGUtKWyap0MOlOYv2Gpk/6TA6k8e7IeAfJ1WI47r2c82Huq/j9/dF+TdNXTZUe8= kinkin@KinKinoV"
}