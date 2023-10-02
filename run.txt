This project was written as a recruitment assignment.
It's aim is to download and process data from the api.nbp.pl website.

To run the project, simply run `make`.
Alternatively, if the `make` tool is unavailable, one can execute 
```
go build -o nbp main.go
./nbp
```
to run the project.

Additionaly, flags can be provided to modify the program's behavior. Supported flags are
```
-numberOfChecks - this flag controls the number of checks that will be run concurrently with a cadence set by the interval flag;
 default value is -numberOfChecks=10. 
-interval - this flag sets the cadence in seconds in which checks should be done;
 default value is -interval=5. 
```

Usage example:
```
./nbp -numberOfChecks=1 interval=3
```
This will make the program run 1 check every 3 seconds.