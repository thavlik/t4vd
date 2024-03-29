# t4vd (Tool for Video Data)
[<img src="https://img.shields.io/badge/maintenance%20status-as%20is-yellow">](https://github.com/thavlik/t4vd)
[<img src="https://img.shields.io/badge/Language-golang-blue.svg">](https://go.dev/)
[<img src="https://img.shields.io/badge/License-Apache_2.0-orange.svg">](./LICENSE)
[<img src="https://img.shields.io/badge/License-MIT-lightblue.svg">](./LICENSE-MIT)

This project is an open source platform for simplifying the creation of video datasets. The core functionality is currently centered around labeling individual frames as "keep" or "discard", facilitating the training of models that filter "junk" frames. There are plans to implement a variety of other labeling strategies.

## Building
### Docker Images
All of the modules use the [kindest](https://github.com/midcontinentcontrols/kindest/) toolchain to streamline building and deploying. Individual microservices can also be manually built with [Docker](https://www.docker.com/):

```bash
# in this example, we're going to build the microservice
# named 'seer', which is responsible for all communication
# with YouTube and caching the results so we don't get
# IP banned :)
$ cd seer

# build seer using the kindest toolchain (easiest)
$ kindest build

# `kindest build` is shorthand for the following command:
$ docker build -t thavlik/t4vd-seer:latest -f Dockerfile ..
```

### Flutter App
The [app/ directory](app/) contains a [Flutter](https://github.com/flutter/flutter) front end. Building it for any platform should be straightforward as the plugin dependencies are minimal. Please open an issue if you encounter any problems.

Note that the front end can also be built and deployed as a web app. This is the default behavior when deploying with `kindest`.

## Deploying
The kubernetes deployment is packaged as a [Helm](https://helm.sh/) chart and is located in the [chart/ directory](chart/). Like building, deploying can done with either `kindest` or `helm`:

```bash
# the values.yaml file is embedded in the root
# kindest.yaml for the repository, and is ideal
# for development as it will skip deploying when
# the chart hasn't changed:
$ kindest deploy

# you can also install the chart manually:
$ export RELEASE_NAME=t4vd
$ export NAMESPACE=t4vd
$ helm install \
    -f myvalues.yaml \
    -n $NAMESPACE \
    $RELEASE_NAME \
    ./chart
```

While not officially supported, the backend images *should* be compatible with [docker-compose](https://docs.docker.com/compose/). Porting would entail converting the chart to a `docker-compose.yaml` file. In any case, I strongly urge you to use kubernetes.

## FAQ
### Why does the project not include ingress?
It is the opinion of many kubernetes experts, including myself, that the k8s-native ingress lacks the features required by most enterprise-scale projects. Ingress configuration is usually ad-hoc, so adding `Ingress` resources to the chart would only complicate deployment without bringing you any closer to actually getting it working. This has been the case for years, but as k8s-native ingress improves, the situation may change.

Note that you **may** get ingress working with the k8s-native `Ingress` resource. All that is required is HTTP support, which the default k8s nginx ingress has. However, I am currently using [traefik-helm-chart](https://github.com/traefik/traefik-helm-chart) for my ingress, though I wouldn't necessarily recommend it for a simple setup. Whatever you choose, I highly recommend [cert-manager](https://cert-manager.io/) for TLS certificate generation. Do it right and it'll be effortless.

### Who are you?
I am a professional software engineer dedicated to the ethical application of machine learning. When I discovered [PyTorch](https://pytorch.org/) and [TensorFlow](https://www.tensorflow.org/), I realized my career will inevitably be distinguished by this powerful technology, [even as a physician](https://github.com/thavlik/machine-learning-portfolio).

### Why ethical machine learning?
I am no moral authority, but this question directly invokes my conviction to live forthrightly.

When celebrated computer vision researcher [Joseph Redmon resigned from research citing military and privacy concerns](https://twitter.com/pjreddie/status/1230524770350817280?lang=en), I found it impossible to not admire. Joseph was confronted with the Cain & Abel paradigm, and with no way around it, he chose in accordance with his better nature. 

It is estimated that the story of Cain & Abel may be mankind's oldest, pre-dating all known religions and qualifying it as among the most relevant to the human condition. While life isn't black and white, *our intentions can be*. Deep learning's vanguard must refuse to be blinded by ambition, lest we bring about a metaphorical hell-on-earth with hyper-advanced military & surveillance states.

And can you even fathom the paradise we'd create if a majority of deep learning researchers struggled with this?

## License
All code in this repository is released under [MIT](LICENSE-MIT) / [Apache 2.0](LICENSE-Apache) dual license, which is extremely permissive. Please open an issue if somehow these terms are insufficient.

All videos are property of their respective creators.
