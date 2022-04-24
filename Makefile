botGo:
		go build -o botGo && ./botGo
rrun:
		brew services start redis
rstop:
		brew services stop redis
