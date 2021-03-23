# directory contents
Each directory is its own module and contains two Go source files with the same name as the directory, a corresponding file with tests as well as a `go.mod` and `go.sum`.
There are also shell scripts included to make it easier to configure and run the entire project as well as the individual components.
Some components include extra files that which will be discussed in the section about that specific component.
The directory structure looks like this:
```
.
├── ARCHITECTURE.md
├── Dockerfile
├── LICENSE
├── Makefile
├── README.md
├── TODO.md
├── frontend
│   ├── Makefile
│   ├── frontend.go
│   ├── frontend_test.go
│   ├── go.mod
│   ├── go.sum
│   ├── runfront.sh
│   ├── static
│   └── templates
├── gifcreator
│   ├── gifcreator.go
│   ├── gifcreator_test.go
│   ├── go.mod
│   ├── go.sum
│   ├── rungifcreator.sh
│   ├── scene
│   ├── server.sh
│   └── worker.sh
├── proto
│   ├── gifcreator.pb.go
│   ├── gifcreator.proto
│   ├── render.pb.go
│   └── render.proto
├── render
│   ├── go.mod
│   ├── go.sum
│   ├── render.go
│   └── render_test.go
├── rundeps.sh
├── rundocker.sh
├── runminio.sh
├── runredis.sh
└── upstream-README.md
```

# components
## `render`
Renders each frame of the gif and uploads them to minio..
## `gifcreator`
Gifcreator is generally responsible for coordinating the process of sending the correct assets to the renderer as well as stitching the rendered frames together to form the final gif.
At least two instances of the code in this package need to be running at once -- one in server mode and the other in worker mode.
### server mode
Server mode duties:
### worker mode
Worker mode duties:
## `frontend`
Web frontend for submitting user preferences upon which the gif will be created.
