build-dev:
	make update-requirements
	docker build -f docker/fast-app.dev.Dockerfile -t mlsdk-fast-app .

run-dev:
	docker run -p 8001:8001 -v $$(pwd)/test-data/my_study.hyperpack:/hyperpackage -v $$(pwd)/fastapp:/app/fastapp -it mlsdk-fast-app 

update-requirements:
	conda env export --from-history > environment.yml
