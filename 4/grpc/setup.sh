module load Go
module load protobuf
module load binutils/2.39-GCCcore-12.2.0

export PATH="$PATH:$(go env GOPATH)/bin"