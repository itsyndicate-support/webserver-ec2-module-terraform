# Self-hosted runner evaluation

## Setup and workload

I configured a self-hosted GitHub Actions runner and tested it with the existing Terraform pipeline.

The runner used a separate EC2 `c7i-flex.large` instance in `eu-north-1`: 2 vCPU, 4 GiB RAM, Ubuntu 24.04, 20 GB gp3 storage and 4 GB swap. It ran as a systemd service under the `ubuntu` user.

I added a workflow input to choose between `ubuntu-latest` and `terraform-runner`. Self-hosted deployment runs remained manual and restricted to `main`. Pull-request checks stayed on GitHub-hosted runners. AWS authentication continued through OIDC.

For capacity and cost planning, I assumed 10 developers, 22 working days and four jobs per developer daily. This gives 880 jobs per month. At an assumed eight minutes per job, the workload requires 7040 minutes, or 117.33 hours.

One runner processes one job at a time. Peak concurrency and queue times were not tested.

## Performance comparison

| Runner | Operation | Job duration | Result |
|---|---|---|---|
| GitHub-hosted | apply | 1m 7s | Passed |
| Self-hosted | apply | 1m 59s | Passed |
| Self-hosted | destroy | 1m 12s | Passed |

Both apply runs used commit `e55cf22` and the `test` environment. The self-hosted run took 52 seconds longer.

These results are preliminary: I measured one successful apply per runner, and cache state, infrastructure state and machine configurations were not fully controlled.

Run links: [GitHub-hosted apply](https://github.com/wai-AI/webserver-ec2-module-terraform/actions/runs/43766239393), [self-hosted apply](https://github.com/wai-AI/webserver-ec2-module-terraform/actions/runs/43765396907), [self-hosted destroy](https://github.com/wai-AI/webserver-ec2-module-terraform/actions/runs/43767052883).

## Cost comparison

The planning scenario assumes a private repository on GitHub Team. Using the standard Linux x64 2-core rate of $0.006 per minute, 7040 minutes cost $42.24. With the full $18 monthly compute allowance available, the additional compute charge would be $24.24.

The eight-minute assumption belongs to this planning model, rather than the measured lab results.

For continuous self-hosted operation, I estimated 730 hours per month:

| Resource | Monthly cost |
|---|---:|
| EC2 `c7i-flex.large` | $66.26 |
| 20 GB gp3 storage | $1.67 |
| One public IPv4 address | $3.65 |
| Total | $71.58 |

The estimates exclude taxes, chargeable traffic and administration time. AWS credits, the GitHub Team subscription and GitHub storage charges are excluded. Application infrastructure costs are separate.

The actual lab repository is public, where standard GitHub-hosted execution is free. The experiment therefore did not establish cost savings.

## Challenges and solutions

The initial `t3.small` instance had 2 GiB RAM. During Terratest dependency compilation, the OOM killer terminated the compiler, and the runner service stopped.

I expanded the disk from 8 GB to 20 GB and added 4 GB swap. Compilation remained slow because of heavy swapping. After moving to an instance with 4 GiB RAM, the pipeline completed successfully.

This showed that choosing the smallest instance can create additional troubleshooting work and longer execution times.

## Conclusion

Self-hosted runners provide control over installed tools, persistent caches and private network access. They also require maintenance, security updates and capacity planning, while idle servers still incur costs.

For this project, I would keep GitHub-hosted runners. The tested self-hosted configuration did not improve execution time or demonstrate savings. Scheduled or on-demand runners could reduce idle costs, but I did not implement those options.

After testing, I destroyed the test infrastructure through the workflow and manually terminated the runner EC2 instance.

Sources: [GitHub Actions billing](https://docs.github.com/en/billing/concepts/product-billing/github-actions), [runner pricing](https://docs.github.com/en/billing/reference/actions-runner-pricing), [EC2 C7i-flex](https://aws.amazon.com/ec2/instance-types/c7i/), [IPv4 pricing](https://aws.amazon.com/vpc/pricing/).
