docker-image:
	docker build -t kouda --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} .
