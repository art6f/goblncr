load('ext://dotenv', 'dotenv')

DIR_ROOT = os.getcwd()
DIR_DEPLOYMENTS = os.path.join(DIR_ROOT, 'deployments', '/')

allow_k8s_contexts('default')

# configs
k8s_yaml([
    DIR_DEPLOYMENTS + "/k8s/namespace.yaml",
    
    DIR_DEPLOYMENTS + "/k8s/cluster-role.yaml",
    DIR_DEPLOYMENTS + "/k8s/cluster-role-binding.yaml",

    DIR_DEPLOYMENTS + "/k8s/deployment.yaml",
    DIR_DEPLOYMENTS + "/k8s/service.yaml",
])

# generate config from envs
dotenv('.env')
configmap_subs = local("envsubst < "  + DIR_DEPLOYMENTS + "/k8s/configmap.template.yaml")
k8s_yaml(configmap_subs)

local_resource(
    'goblncr-build',
    cmd='CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./build/goblncr ./cmd/goblncr/goblncr.go',
    deps=['.'],
    ignore=['./build/', './deployments/'],
    labels=['Balancer'],
)

docker_build(
    "goblncr:dev",
    DIR_ROOT,
    dockerfile=DIR_DEPLOYMENTS + "/docker/Dockerfile.dev",
    container_args={
        'K8S_PODS_SELECTOR': "app=goblncr",
    },
)

k8s_resource(
    'goblncr',
    objects=[
        'goblncr:namespace',
        'goblncr-role:clusterrole',
        'goblncr-global:clusterrolebinding',
        'goblncr-config:configmap',
    ],
    labels='Resources'
)

# Link image to deployment
k8s_resource(
  "goblncr",
  port_forwards=8080,
  labels=['Balancer'],
)

watch_file("deployments/docker/Dockerfile")
watch_file("cmd")
watch_file("internal")
watch_file("pkg")