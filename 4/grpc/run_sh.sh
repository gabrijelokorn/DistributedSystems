#!/bin/bash
#SBATCH --nodes=1
#SBATCH --array=0-9
#SBATCH --reservation=psistemi
#SBATCH --output=../Database/Database-%a.db

srun grpc -s localhost -p 8100 -id $SLURM_ARRAY_TASK_ID 