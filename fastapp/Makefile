build-dev:
	# make update-requirements
	docker build -f docker/fast-app.dev.Dockerfile -t mlsdk-fast-app .

run-dev:
	docker run -e FASTKEY=abcd123 -p 8001:8001 -v $$(pwd)/examples/my_study.hyperpack:/hyperpackage -v $$(pwd)/fastapp:/app/fastapp -v $$(pwd)/upload:/app/upload -it mlsdk-fast-app

update-requirements:
	conda env export --from-history > environment.yml
