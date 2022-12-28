# t4vd (Tool for Video Data)
This project is a collection of tools for simplifying the creation of video datasets.

## Building
### Docker Images
All of the modules use the [kindest](https://github.com/midcontinentcontrols/kindest/) toolchain to streamline building and deploying. Individual microservices can also be manually built with [Docker](https://www.docker.com/):

```bash
$ cd seer

# build using the kindest toolchain
$ kindest build

# `kindest build` is shorthand for the following command:
$ docker build -t thavlik/t4vd-seer:latest -f Dockerfile ..
```

### Flutter App
The [app/ directory](app/) contains a standard [Flutter](https://github.com/flutter/flutter) front end. Building it for any platform should be straightforward as the plugin dependencies are minimal. Please open an issue if you encounter any problems.

Note that the front end can also be built and deployed as a web app. This is the default behavior with `kindest`.

## Deploying
The kubernetes deployment is packaged as a [Helm](https://helm.sh/) chart and is located in the [chart/ directory](chart/). Like building, deploying can done with either `kindest` or `helm`:

```bash
# the values.yaml file is embedded in the root
# kindest.yaml for the repository, and is ideal
# for development
$ kindest deploy

# we can also install the chart manually:
$ export RELEASE_NAME=t4vd
$ export NAMESPACE=t4vd
$ helm install \
    -f myvalues.yaml \
    -n $NAMESPACE \
    $RELEASE_NAME \
    ./chart
```

While not officially supported, the backend images should be compatible with [docker-compose](https://docs.docker.com/compose/). Porting would entail converting the chart to a `docker-compose.yaml` file. In any case, I strongly urge you to use kubernetes.

## License
All videos are property of their respective creators.

This project's code is released under [MIT](LICENSE-MIT) / [Apache 2.0](LICENSE-Apache) dual license, which is extremely permissive.

You are encouraged to fork this project and use it for your own purposes, even commercial.

