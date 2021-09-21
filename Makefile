all:

build:
	docker build -t nicholasmnovak/timekeeper-morty:latest .
	docker push nicholasmnovak/timekeeper-morty:latest
