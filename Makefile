build-auth-image:
	docker build -t fadygamil/auth-service -f auth-service/Dockerfile .

run-auth-container;
	docker run -it -p 3030:8080 --name auth-service auth-service