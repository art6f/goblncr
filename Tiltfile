DIR_ROOT = os.getcwd()
DIR_DEPLOYMENTS = os.path.join(DIR_ROOT, 'deployments', '/')

allow_k8s_contexts('default')

# configs
k8s_yaml([
    DIR_DEPLOYMENTS + "/k8s/namespace.yaml",
    DIR_DEPLOYMENTS + "/k8s/configmap.yaml",
    DIR_DEPLOYMENTS + "/k8s/deployment.yaml",
    DIR_DEPLOYMENTS + "/k8s/service.yaml",
])

# image build (no push)
default_registry('localhost:5000')
docker_build(
    "goblncr:dev",
    DIR_ROOT,
    dockerfile=DIR_DEPLOYMENTS + "/docker/Dockerfile",
)


# Link image to deployment
k8s_resource(
  "goblncr",
  port_forwards=8080,
)

watch_file("deployments/docker/Dockerfile")
watch_file("cmd")
watch_file("internal")
watch_file("pkg")