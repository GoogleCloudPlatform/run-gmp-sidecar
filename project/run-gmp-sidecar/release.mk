.PHONY: local-container-goreleaser
local-container-goreleaser:
	docker buildx build \
		--progress=plain \
		-t rungmpcol-build \
		-f Dockerfile.goreleaser_releaser \
		..
	CONTAINER_ID=$$(docker create rungmpcol-build) && \
		docker cp $$CONTAINER_ID:/run-gmp-sidecar/dist . &&\
		docker rm --force $$CONTAINER_ID
