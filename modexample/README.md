# Module modexample for Dagger

![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.0-0f0f19.svg?style=flat-square)

A simple [Dagger](https://dagger.io) _place the description of the module here_

## Configuration 🛠️

Through the [Dagger CLI](https://docs.dagger.io/cli/465058/install), or by using it directly within your module, you can configure the following options:

* ⚙️ `ctr`: The container to use as a base container. If not specified, a new container is created.
* ⚙️ `version`: The version of the Go image to use. Defaults to `latest`.
* ⚙️ `image`: The Go image to use. Defaults to `golang:alpine`.

---

## Features 🎨

| Command or functionality | Command | Example        | Status |
|--------------------------|---------|----------------|--------|
| Run Go Tests             | **run** | `dagger call ` | ✅      |


## Using the modexample Module 🚀



_Place the description of the module here_

---

### Usage through the Dagger CLI 🚀



List all the functions available in the module:

  ```bash
  # enter into the module's directory
  cd <module-path>



  # list all the functions available in the module

  dagger develop && dagger functions


```

Call a function:



  ```bash
  # call a function
  dagger call <function-name> [arguments]


```

---

## Testing 🧪


This module includes a [testing](tests) module that aims to test the functionality of the [@MODULE] module. The tests are written in Go and can be run using the following command:

```bash
## Run the tests using the just command
just test gotest
```

>NOTE: The `just` command entails the use of the [**Justfile**](https://just.systems) for task automation.
