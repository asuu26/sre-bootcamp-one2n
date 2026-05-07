.PHONY: deploy teardown

include .env

deploy:
	echo $(GHCR_TOKEN) | docker login ghcr.io -u $(GHCR_USER) --password-stdin
	IMAGE_TAG=$(IMAGE_TAG) docker compose up -d

teardown:
	docker compose down -v
