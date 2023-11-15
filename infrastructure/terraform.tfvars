vpc_cidr = "192.168.0.0/16"

network_http = {
  subnet_name = "subnet_http"
  cidr        = "192.168.1.0/24"
}

http_instance_names = ["instance-http-1", "instance-http-2"]

network_db = {
  subnet_name = "subnet_db"
  cidr        = "192.168.2.0/24"
}

db_instance_names = ["instance-db-1", "instance-db-2", "instance-db-3"]

#pub_key_filename = "virt.pub"

#pub_key_path = "/home/vanadium/.ssh/"