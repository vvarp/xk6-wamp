build:
	xk6 build --with github.com/vvarp/xk6-wamp=.
	./k6 run example.js

