# -*- mode: Python -*-
load('ext://min_tilt_version', 'min_tilt_version')
min_tilt_version('0.33.1')

load('ext://dotenv', 'dotenv')
dotenv()

# Manage Contexts
context = os.environ.get('TILT_K8S_CONTEXT', 'docker-desktop')
allow_k8s_contexts(context)
current_context = k8s_context()
if current_context != context:
  warn('current k8s context is "{}" needs "{}". switching...'.format(current_context, context))
  local('kubectl config use-context {}'.format(context))


# Manage Registries
docker_registry = os.environ.get('TILT_DOCKER_REGISTRY', None)
if docker_registry:
  default_registry(docker_registry)

# Load resources
include('resources/ingress-nginx/Tiltfile')
include('resources/linkerd/Tiltfile')

# Load apps
include('apps/foo/Tiltfile')
include('apps/bar/Tiltfile')
include('apps/baz/Tiltfile')
