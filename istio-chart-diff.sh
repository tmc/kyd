#!/bin/bash
set -euo pipefail
# This script helps generate diffs representing changes to helm chart configuration.
# Useful for applying subsets of Istio chart options.
# Note: expects to be run from within an Istio checkout.

# Inspiration via https://www.youtube.com/watch?v=FbYBO7Pi2d8
# wget istio-1.1.7.tar.gz

COMPONENT="${1:-${component}}"

for flag in false true; do
    set_args=()
    for arg in $*; do
        set_args+=("--set ${arg}=${flag}")
    done

    # If an extra-values yaml exists, include it.
    if test -f "${COMPONENT}-extra-values.yaml"; then
        set_args+=("--values ${COMPONENT}-extra-values.yaml")
    fi

    helm template \
    ${set_args[*]} \
    --namespace istio-system \
    install/kubernetes/helm/istio \
    > "${COMPONENT}-${flag}.yaml"
done

#diff --changed-group-format='%>' --unchanged-group-format='' "${COMPONENT}-false.yaml" "${COMPONENT}-true.yaml" > "${COMPONENT}-config.yaml"
#diff --line-format=%L "${COMPONENT}-true.yaml" "${COMPONENT}-false.yaml" > "${COMPONENT}-config-lf.yaml"
command -v kyd >/dev/null || go get github.com/tmc/kyd
kyd "${COMPONENT}-false.yaml" "${COMPONENT}-true.yaml" > "${COMPONENT}-config.yaml"
