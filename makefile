build:
	sudo docker build . --tag `git config user.name | tr A-Z a-z`/simdep-admin:`git rev-parse HEAD`
	sudo docker tag `git config user.name | tr A-Z a-z `/simdep-admin:`git rev-parse HEAD` `git config user.name | tr A-Z a-z `/simdep-admin:latest
push:
	sudo docker push `git config user.name | tr A-Z a-z `/simdep-admin:`git rev-parse HEAD`
	sudo docker push `git config user.name | tr A-Z a-z `/simdep-admin
debug:
	sudo docker run -it --rm -v `pwd`:/workspace -w /workspace golang:1.20 /bin/bash