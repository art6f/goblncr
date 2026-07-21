load('ext://dotenv', 'dotenv')
dotenv('.env')

DIR_ROOT = os.getcwd()
DIR_DEPLOYMENTS = os.path.join(DIR_ROOT, 'deployments', '/')

allow_k8s_contexts('default')

# generate config from envs
configmap_subs = local("envsubst < "  + DIR_DEPLOYMENTS + "/k8s/configmap.template.yaml")
k8s_yaml(configmap_subs)

GO_OS = os.getenv('GOOS')
GO_ARCH = os.getenv('GOARCH')
BUILD_CMD='CGO_ENABLED=0 GOOS={} GOARCH={} go build -gcflags="all=-N -l" -o ./build/goblncr-dbg ./cmd/goblncr/goblncr.go' \
    .format(GO_OS, GO_ARCH)

local_resource(
    'goblncr-build',
    cmd=BUILD_CMD,
    deps=['.'],
    ignore=['./build/', './deployments/'],
    labels=['Balancer'],
)

# configs
k8s_yaml([
    DIR_DEPLOYMENTS + "/k8s/namespace.yaml",
    
    DIR_DEPLOYMENTS + "/k8s/cluster-role.yaml",
    DIR_DEPLOYMENTS + "/k8s/cluster-role-binding.yaml",

    DIR_DEPLOYMENTS + "/k8s/deployment.yaml",
    DIR_DEPLOYMENTS + "/k8s/service.yaml",
])

docker_build(
    "goblncr:dev",
    DIR_ROOT,
    dockerfile=DIR_DEPLOYMENTS + "/docker/Dockerfile.dev",
)

# Link image to deployment
k8s_resource(
    "goblncr",
    port_forwards=["2345:2345"],
    labels=['Balancer'],

    objects=[
        'goblncr:namespace',
        'goblncr-role:clusterrole',
        'goblncr-global:clusterrolebinding',
        'goblncr-config:configmap',
    ],
)

watch_file("deployments/docker/Dockerfile")
watch_file("cmd")
watch_file("internal")
watch_file("pkg")