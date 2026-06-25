#!/bin/bash

cd /tmp
curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash

# get the kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl.sha256"

echo "$(cat kubectl.sha256)  kubectl" | sha256sum --check
if [[ "$?" -ne 0 ]]; then
    echo "Checksum failed!"
    exit 1
fi

sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# helm?
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-4 | bash