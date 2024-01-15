#!/bin/bash
#SBATCH --nodes=1
#SBATCH --array=0-9
#SBATCH --reservation=psistemi
#SBATCH --output=/d/hpc/home/go7745/4/procesi/Database/%a.db

srun /d/hpc/home/go7745/4/procesi/grpc/grpc -s localhost -p 8100 -id $SLURM_ARRAY_TASK_ID -n $SLURM_ARRAY_TASK_COUNT