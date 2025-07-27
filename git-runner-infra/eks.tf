resource "aws_eks_cluster" "eks" {
  name = "${var.project_name}"
  role_arn = aws_iam_role.eks_cluster_role.arn

  vpc_config {
    subnet_ids = [
        aws_subnet.public_subnets[0].id,
        aws_subnet.public_subnets[1].id,
    ]
  }

  depends_on = [ 
    aws_iam_role_policy_attachment.eks_cluster_attach
   ]
}

resource "aws_eks_node_group" "node_group" {
  cluster_name = aws_eks_cluster.eks.name
  node_group_name = "${var.project_name}-nodes"
  node_role_arn = aws_iam_role.eks_node_role.arn
  subnet_ids = [
    aws_subnet.public_subnets[0].id,
    aws_subnet.public_subnets[1].id,
  ]
  scaling_config {
    desired_size = 2
    max_size = 4
    min_size = 1
  }

  instance_types = ["t3.small"]

  depends_on = [ 
    aws_iam_role_policy_attachment.eks_node_attach1,
    aws_iam_role_policy_attachment.eks_node_attach2,
    aws_iam_role_policy_attachment.eks_node_attach3
   ]
}