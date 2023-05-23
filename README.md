# files
A tool to count files in subdirectories of a specified path and sort them

## Install
1. `$ git clone https://github.com/dmnyu/files`
2. `$ cd files`
3. `$ make tidy` // runs go mod tidy
4. `$ make test` // runs the test fixtures
5. `$ make build` //runs go build -o files main/main.go
6. `$ sudo make install` // this will install the bin `files` to /usr/local/bin

## Run
$files --path path-to-directory [options]<br>
Options:<br>
&nbsp;&nbsp;--help&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;print this help message<br>
&nbsp;&nbsp;--report&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;output a tsv file listing<br>
&nbsp;&nbsp;--output-file&nbsp;&nbsp;&nbsp;&nbsp;path/to/report/file, default: `./filecounts.tsv`<br>
&nbsp;&nbsp;--verbose&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;verbose output to stdout<br>
&nbsp;&nbsp;--workers&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;number of threads to run simultaneously, default: 8<br>


## Sample output
<pre>
$ files --path /usr/local
total number of files in /usr/local: 49255

num files       path
---------       ----
48754           /usr/local/go
469             /usr/local/include
27              /usr/local/bin
5               /usr/local/lib
0               /usr/local/share
0               /usr/local/src
0               /usr/local/lib64
0               /usr/local/libexec
0               /usr/local/man
0               /usr/local/sbin
0               /usr/local
</pre>
