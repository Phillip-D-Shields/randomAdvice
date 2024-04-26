APP_NAME=randomAdvice.app
IMAGE_NAME=random-advice
CONTAINER_NAME=randomAdvice

.PHONY: build run stop clean

build:
	docker build -t $(IMAGE_NAME) .

run:
	docker run -d --name $(CONTAINER_NAME) -p 8080:8080 -v $(PWD)/advice.db:/app/advice.db $(IMAGE_NAME)

stop:
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)

clean:
	docker rmi $(IMAGE_NAME)