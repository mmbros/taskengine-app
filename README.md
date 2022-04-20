# taskengine-app

Demo app to show mmbros/taskengine package usage.

Available commands:

- `demo` Execute a demo scenario, show progress and output results in json format
- `server` Start an http server to show json files containing the results of execution
- `version` Print version information

## `demo` command

Performs a demo scenario and save results in json format.

Options:

- `workers` number of workers
- `instances` instances of each worker
- `tasks` number of tasks
- `progress` show progress of execution
- `seed` random seed generator
- `spread` perc of how many workers executes each tasks (0..100)
- `output` pathname of the output file
- `force` overwrite already existing output file

Random Result options:

- `mean` mean value
- `stddev` standard deviation
- `errperc` perc of task error (0..100)

## `server` command

Start an http server to view a page with graphs based upon
the taskengine json files created with the demo command.

Options:

- `address` server address and port
- `folder` folder containing the json files
- `recursive` search recursively all the json files of the sub-folders

## `version` command

Prints version informations.

Options:

- `build-options` print verion with build options

Example:

``` shell
$ taskengine-app version --build-options 
taskengine-app version dev-20220420T141155
taskengine package version v0.3.0
go version go1.18 linux/amd64
build date: 2022-04-20 14:11:55 +0200
git commit: 4f6b962
os/arch: Linux x86_64
```
