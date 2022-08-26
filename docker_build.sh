#!/usr/bin/env bash
# Docker multiplatform build

[ -z "$DOCKER_REPO" ] && { echo "You must set DOCKER_REPO"; exit 1; }

image="go-yolink"
platforms="arm64 amd64"
repoimage="${DOCKER_REPO}/${image}"
version="v1"

imglist=""
for platform in $platforms; do
    echo "Building for $platform:"
    docker buildx build . -t ${image}:${platform} -t ${repoimage}:${platform} --platform linux/${platform}
    docker push ${repoimage}:${platform}
    imglist+="${repoimage}:${platform} "
done

docker manifest rm ${repoimage}:${version}
docker manifest create ${repoimage}:${version} $imglist
docker manifest push ${repoimage}:${version}
