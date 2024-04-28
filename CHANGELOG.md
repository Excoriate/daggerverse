# Changelog

## [1.11.0](https://github.com/Excoriate/daggerverse/compare/v1.10.0...v1.11.0) (2024-04-28)


### Features

* **tflint:** Update TFLint instructions in README.md ([#44](https://github.com/Excoriate/daggerverse/issues/44)) ([8d07963](https://github.com/Excoriate/daggerverse/commit/8d0796388a017b51171bc4c021452ab679daa73d))

## [1.10.0](https://github.com/Excoriate/daggerverse/compare/v1.9.0...v1.10.0) (2024-04-28)


### Features

* Improve Tflint Docker container initialization and command execution ([#43](https://github.com/Excoriate/daggerverse/issues/43)) ([052c2de](https://github.com/Excoriate/daggerverse/commit/052c2de5b63debee148283eabf0548c0f0d65223))


### Refactoring

* Simplify parameter declaration in addCMDsToContainer function ([#39](https://github.com/Excoriate/daggerverse/issues/39)) ([c004989](https://github.com/Excoriate/daggerverse/commit/c004989cc2d5e099f7e899ff574940a246dbb6af))


### Other

* fix hooks issues ([fcd55b3](https://github.com/Excoriate/daggerverse/commit/fcd55b30113fc13cbc9255e83eceb7db7db7e6e5))

## [1.9.0](https://github.com/Excoriate/daggerverse/compare/v1.8.0...v1.9.0) (2024-04-27)


### Features

* **dagger:** Add Release method to run 'goreleaser release' ([#36](https://github.com/Excoriate/daggerverse/issues/36)) ([5b4ff83](https://github.com/Excoriate/daggerverse/commit/5b4ff83d7dba76c3fb1666fa65e09dd53fbd057a))
* Enhance buildArgs function and add new methods in Goreleaser struct ([#32](https://github.com/Excoriate/daggerverse/issues/32)) ([b724076](https://github.com/Excoriate/daggerverse/commit/b72407686ac75d81e0fc8f43de20f4804955a2ac))
* Update .gitignore to exclude go.work ([#34](https://github.com/Excoriate/daggerverse/issues/34)) ([ecf519d](https://github.com/Excoriate/daggerverse/commit/ecf519ddaef295b62a71a60a882b8d97ea1d7ad9))


### Other

* add example project for goreleaesr module ([4be5287](https://github.com/Excoriate/daggerverse/commit/4be528785e2c7cca30c1dad6b37ebdb5eff24555))

## [1.8.0](https://github.com/Excoriate/daggerverse/compare/v1.7.0...v1.8.0) (2024-04-25)


### Features

* Add dagger generated files and query builder to gitignore ([#23](https://github.com/Excoriate/daggerverse/issues/23)) ([8640987](https://github.com/Excoriate/daggerverse/commit/864098737169133cdb42fd9d5c89a67c703913ff))
* Update labeler configuration with new rules ([#26](https://github.com/Excoriate/daggerverse/issues/26)) ([b5cb8ec](https://github.com/Excoriate/daggerverse/commit/b5cb8ec27531ab07a91a8af46b3d6fa96a5725ca))
* Update module dependencies ([#30](https://github.com/Excoriate/daggerverse/issues/30)) ([f2ca164](https://github.com/Excoriate/daggerverse/commit/f2ca16456bc74dd3c15c59d0b2277f24ece49fc5))
* Upgrade Go to version 1.21.7, update tflint engineVersion to v0.11.1, ([#31](https://github.com/Excoriate/daggerverse/issues/31)) ([01c88d7](https://github.com/Excoriate/daggerverse/commit/01c88d723425d18fface16016fae2ec33d96f2f0))


### Other

* Update gqlgen and go-logr versions for dagger internal telemetry folder. ([#25](https://github.com/Excoriate/daggerverse/issues/25)) ([59fb15b](https://github.com/Excoriate/daggerverse/commit/59fb15b6a961982e0bc6c339cc01cd72b0497635))


### Refactoring

* **cli:** improve handling of secret variables in get command ([#29](https://github.com/Excoriate/daggerverse/issues/29)) ([82e10c9](https://github.com/Excoriate/daggerverse/commit/82e10c90c0a4454bd9997e96cc16790fc0ec8fa4))
* Fix indentation in workflows and labeler files ([#27](https://github.com/Excoriate/daggerverse/issues/27)) ([e2c4a88](https://github.com/Excoriate/daggerverse/commit/e2c4a88d5614fbb76ce2ea839800e728e3a9969a))
* Refactor main.go to use secret token for GitLab client authentication. ([#28](https://github.com/Excoriate/daggerverse/issues/28)) ([25dbf73](https://github.com/Excoriate/daggerverse/commit/25dbf73e8cbe4ef443b3d90a1abdf77422367cbd))

## [1.7.0](https://github.com/Excoriate/daggerverse/compare/v1.6.0...v1.7.0) (2024-03-05)


### Features

* add terratest mvp dagger module ([640334b](https://github.com/Excoriate/daggerverse/commit/640334b0958b77ce908db8415044ff0a9b4f4c18))
* draft terratest module ([2f8f0f6](https://github.com/Excoriate/daggerverse/commit/2f8f0f6c4985256191a7972f61e15d249adcc415))


### Docs

* add better doc on gitlab-cicd-vars module ([ae32152](https://github.com/Excoriate/daggerverse/commit/ae32152112c411f90e78dab35b10ddef645d6a83))

## [1.6.0](https://github.com/Excoriate/daggerverse/compare/v1.5.0...v1.6.0) (2024-03-03)


### Features

* add dagger-module for gitlab cicdvars ([0fc7d73](https://github.com/Excoriate/daggerverse/commit/0fc7d73b38920bf82781e72d276bfab097b3628a))

## [1.5.0](https://github.com/Excoriate/daggerverse/compare/v1.4.1...v1.5.0) (2024-03-03)


### Features

* Add common terragrunt functions ([5f95ba7](https://github.com/Excoriate/daggerverse/commit/5f95ba763efda36c807e67b41de14c911c78010f))
* Add common terragrunt functions ([220908d](https://github.com/Excoriate/daggerverse/commit/220908d3da331fec2a181b2aeb3a3f0f5196a5c7))
* add enablePrivateGit flag in IacTerragrunt struct ([e86e583](https://github.com/Excoriate/daggerverse/commit/e86e5837c8d83f38990460f63fd7c80e49e8fead))
* Add new apis supported for secrets, git ssh, and cache ([8ebb47b](https://github.com/Excoriate/daggerverse/commit/8ebb47b1831841861d42130f638a0333396b4970))
* Add repo setup ([6b7c984](https://github.com/Excoriate/daggerverse/commit/6b7c984550df5e9725027ee469ef051a374a6bc5))
* Add support for environment variables ([8b4e37c](https://github.com/Excoriate/daggerverse/commit/8b4e37c800af932eadedea6d7cfab7d5ab9f4c22))
* Add sync utility ([6f93818](https://github.com/Excoriate/daggerverse/commit/6f93818f4ab76e19d92faabd2cb03a0da444f864))
* add terraform dagger-module ([260c3d3](https://github.com/Excoriate/daggerverse/commit/260c3d3c112f443493e4f11069e609263538506c))
* add terraform dagger-module ([9755d3f](https://github.com/Excoriate/daggerverse/commit/9755d3fdb883f5f7a6ea6e0841f8426473da41a8))
* Add tg init, add reusable functions ([d6d87c8](https://github.com/Excoriate/daggerverse/commit/d6d87c886782997d8b0c820d6bc67a6afe740f19))
* Add working version ([1aa9e19](https://github.com/Excoriate/daggerverse/commit/1aa9e19e38acb329e71553e55bfa9ce65ca0aecf))
* first commit ([242b312](https://github.com/Excoriate/daggerverse/commit/242b312b5bfdc84f69828660844fce6636dd28f6))
* Fix order of git-ssh api ([f56789f](https://github.com/Excoriate/daggerverse/commit/f56789f387cff10313519e3709b10e4a23ffd68b))


### Bug Fixes

* dead-links ([274b13e](https://github.com/Excoriate/daggerverse/commit/274b13e2a9d6232b60e40b41137172c0edb3316b))
* dead-links ([99eb6df](https://github.com/Excoriate/daggerverse/commit/99eb6dfffca266664370737b85bdc98f2896d3fe))
* explicit dependency for daggerx ([471e1b2](https://github.com/Excoriate/daggerverse/commit/471e1b2e289ab79fa016985226a0ab261f4febb0))
* Fix commands ([75a5962](https://github.com/Excoriate/daggerverse/commit/75a5962227c95f4a30d34d2b2147ddab29eea0ea))
* fix workflows ([a4d54a5](https://github.com/Excoriate/daggerverse/commit/a4d54a557716ca942f7e0391752f4aa3d531997d))
* markdown link checker ([4da7915](https://github.com/Excoriate/daggerverse/commit/4da79154595b0c62fbbe8e159b587417bc8fbe4b))
* update sync task to use 'dagger develop' instead of 'dagger mod sync -m' command ([#17](https://github.com/Excoriate/daggerverse/issues/17)) ([c26c720](https://github.com/Excoriate/daggerverse/commit/c26c72065e9179aabeb4ab73646f66ea8a5b7a02))
* workflows and incorrect interpolation ([03c5ab2](https://github.com/Excoriate/daggerverse/commit/03c5ab29cd9a04fb4f8311fc902514507a0bc322))
* workflows and incorrect interpolation ([2816973](https://github.com/Excoriate/daggerverse/commit/28169731e1fc715a15d1a5c3a5ab26fc2740f6b5))
* workflows and incorrect interpolation ([f4e83a4](https://github.com/Excoriate/daggerverse/commit/f4e83a401833acb79c22aa01828ef8f6b0b92a5b))


### Docs

* Add better docs, and usage ([b78a703](https://github.com/Excoriate/daggerverse/commit/b78a7035535f6c2dcd7b435b106ea1027a1849ef))
* Add better docs, and usage ([f75e9c3](https://github.com/Excoriate/daggerverse/commit/f75e9c314e26427ad26b714149902857a87213aa))
* Add final docs ([f7dd4ed](https://github.com/Excoriate/daggerverse/commit/f7dd4ed3469c49c2bd419d3d1f5200146529ac66))
* Add readme, add logo ([5f074f4](https://github.com/Excoriate/daggerverse/commit/5f074f411dd4db5a936fa22df9c65be98653a5fd))
* add README.md, add reusable taskfile for dagger common operations ([a09f56b](https://github.com/Excoriate/daggerverse/commit/a09f56b63c91be5c7ff1b856fe46de1bede38d89))
* add README.md, add reusable taskfile for dagger common operations ([8a21d4f](https://github.com/Excoriate/daggerverse/commit/8a21d4fbe76ef0a8ba861aebf1af3c53340acf5d))


### Other

* Add Dagger tasks, and other repo settings ([8a0cd49](https://github.com/Excoriate/daggerverse/commit/8a0cd494e29ccd0c66a2801277045822e4cc4bbe))
* add example with an external terraform module ([284f3a7](https://github.com/Excoriate/daggerverse/commit/284f3a7167f24a27bc81e9896a91d1e2fb935373))
* Add missing .gitattributes ([184b1df](https://github.com/Excoriate/daggerverse/commit/184b1df8d378a9697405c80358097f673d432633))
* Add required ignored files from jetbrains ([34369db](https://github.com/Excoriate/daggerverse/commit/34369db2b2bfc0b146611153c708a5abdd0fef59))
* Fix dagger CLI link ([5537e21](https://github.com/Excoriate/daggerverse/commit/5537e210744c146a07d89b6fcbdc73895184eaa1))
* Fix link that points to the actual iac-terragrunt module ([2ab9790](https://github.com/Excoriate/daggerverse/commit/2ab979079a46eee13ac32b5a43a4fd2f285f9507))
* **main:** release 1.0.0 ([#1](https://github.com/Excoriate/daggerverse/issues/1)) ([fc35a7f](https://github.com/Excoriate/daggerverse/commit/fc35a7f889669b8ddebb1f5c3f79c65a5aa71b59))
* **main:** release 1.0.1 ([#2](https://github.com/Excoriate/daggerverse/issues/2)) ([c24676b](https://github.com/Excoriate/daggerverse/commit/c24676b18252c174b756c572109f5b0df2d5b045))
* **main:** release 1.1.0 ([#3](https://github.com/Excoriate/daggerverse/issues/3)) ([152950e](https://github.com/Excoriate/daggerverse/commit/152950ef9326e227203e5f1f5429149faf289b70))
* **main:** release 1.1.0 ([#4](https://github.com/Excoriate/daggerverse/issues/4)) ([88d3729](https://github.com/Excoriate/daggerverse/commit/88d3729b25a95ea8146d2b73f9b6957717c7be7b))
* **main:** release 1.1.1 ([#5](https://github.com/Excoriate/daggerverse/issues/5)) ([8435d55](https://github.com/Excoriate/daggerverse/commit/8435d55b85b7fcce729f929da50b2fde87ec9baf))
* **main:** release 1.1.2 ([#6](https://github.com/Excoriate/daggerverse/issues/6)) ([840542b](https://github.com/Excoriate/daggerverse/commit/840542be4e47b54555d76a7ddc8d3b89bdf0699c))
* **main:** release 1.1.3 ([#7](https://github.com/Excoriate/daggerverse/issues/7)) ([456e1e1](https://github.com/Excoriate/daggerverse/commit/456e1e1ccfbda3be747fc403416ac2c51e1bbcc1))
* **main:** release 1.1.4 ([#8](https://github.com/Excoriate/daggerverse/issues/8)) ([0c2495c](https://github.com/Excoriate/daggerverse/commit/0c2495c0ca1fecce366734e0eb6a97c8ceaf3adb))
* **main:** release 1.2.0 ([#9](https://github.com/Excoriate/daggerverse/issues/9)) ([2a52044](https://github.com/Excoriate/daggerverse/commit/2a52044c1d34332474cb725cf1a366eac5a227ec))
* **main:** release 1.3.0 ([#14](https://github.com/Excoriate/daggerverse/issues/14)) ([6afe1de](https://github.com/Excoriate/daggerverse/commit/6afe1de075be4196d68ef697a9d8f29c97a0f906))
* **main:** release 1.3.1 ([#15](https://github.com/Excoriate/daggerverse/issues/15)) ([c230fac](https://github.com/Excoriate/daggerverse/commit/c230fac4365c4f31d9bcf978ca15bd79fe897192))
* **main:** release 1.4.0 ([#16](https://github.com/Excoriate/daggerverse/issues/16)) ([55fe99d](https://github.com/Excoriate/daggerverse/commit/55fe99d9744f60ec4fa69df7f73d8ff1f51b64ee))
* Remove go-build from allowed hooks ([559072f](https://github.com/Excoriate/daggerverse/commit/559072f330ac87199e8528d40f3c57b8c3dd9010))
* update dependencies, test publish ([92627e9](https://github.com/Excoriate/daggerverse/commit/92627e9c3dddab31bbcc5febb5c2866c98a9444d))
* update go tasks to run in specific module directory ([#19](https://github.com/Excoriate/daggerverse/issues/19)) ([924dd4d](https://github.com/Excoriate/daggerverse/commit/924dd4dfe24f9b6638a399f17e313516065d1e36))

## [1.4.1](https://github.com/Excoriate/daggerverse/compare/v1.4.0...v1.4.1) (2024-03-03)


### Bug Fixes

* fix workflows ([a4d54a5](https://github.com/Excoriate/daggerverse/commit/a4d54a557716ca942f7e0391752f4aa3d531997d))
* update sync task to use 'dagger develop' instead of 'dagger mod sync -m' command ([#17](https://github.com/Excoriate/daggerverse/issues/17)) ([c26c720](https://github.com/Excoriate/daggerverse/commit/c26c72065e9179aabeb4ab73646f66ea8a5b7a02))


### Other

* update go tasks to run in specific module directory ([#19](https://github.com/Excoriate/daggerverse/issues/19)) ([924dd4d](https://github.com/Excoriate/daggerverse/commit/924dd4dfe24f9b6638a399f17e313516065d1e36))

## [1.4.0](https://github.com/Excoriate/daggerverse/compare/v1.3.1...v1.4.0) (2024-03-02)


### Features

* add terraform dagger-module ([260c3d3](https://github.com/Excoriate/daggerverse/commit/260c3d3c112f443493e4f11069e609263538506c))
* add terraform dagger-module ([9755d3f](https://github.com/Excoriate/daggerverse/commit/9755d3fdb883f5f7a6ea6e0841f8426473da41a8))

## [1.3.1](https://github.com/Excoriate/daggerverse/compare/v1.3.0...v1.3.1) (2024-02-05)


### Other

* add example with an external terraform module ([284f3a7](https://github.com/Excoriate/daggerverse/commit/284f3a7167f24a27bc81e9896a91d1e2fb935373))

## [1.3.0](https://github.com/Excoriate/daggerverse/compare/v1.2.0...v1.3.0) (2024-02-05)


### Features

* add enablePrivateGit flag in IacTerragrunt struct ([e86e583](https://github.com/Excoriate/daggerverse/commit/e86e5837c8d83f38990460f63fd7c80e49e8fead))
* Add new apis supported for secrets, git ssh, and cache ([8ebb47b](https://github.com/Excoriate/daggerverse/commit/8ebb47b1831841861d42130f638a0333396b4970))
* Fix order of git-ssh api ([f56789f](https://github.com/Excoriate/daggerverse/commit/f56789f387cff10313519e3709b10e4a23ffd68b))


### Bug Fixes

* Fix commands ([75a5962](https://github.com/Excoriate/daggerverse/commit/75a5962227c95f4a30d34d2b2147ddab29eea0ea))


### Other

* Add missing .gitattributes ([184b1df](https://github.com/Excoriate/daggerverse/commit/184b1df8d378a9697405c80358097f673d432633))

## [1.2.0](https://github.com/Excoriate/daggerverse/compare/v1.1.4...v1.2.0) (2024-02-01)


### Features

* Add support for environment variables ([8b4e37c](https://github.com/Excoriate/daggerverse/commit/8b4e37c800af932eadedea6d7cfab7d5ab9f4c22))

## [1.1.4](https://github.com/Excoriate/daggerverse/compare/v1.1.3...v1.1.4) (2024-01-30)


### Docs

* Add better docs, and usage ([b78a703](https://github.com/Excoriate/daggerverse/commit/b78a7035535f6c2dcd7b435b106ea1027a1849ef))
* Add better docs, and usage ([f75e9c3](https://github.com/Excoriate/daggerverse/commit/f75e9c314e26427ad26b714149902857a87213aa))

## [1.1.3](https://github.com/Excoriate/daggerverse/compare/v1.1.2...v1.1.3) (2024-01-30)


### Docs

* Add final docs ([f7dd4ed](https://github.com/Excoriate/daggerverse/commit/f7dd4ed3469c49c2bd419d3d1f5200146529ac66))


### Other

* Fix dagger CLI link ([5537e21](https://github.com/Excoriate/daggerverse/commit/5537e210744c146a07d89b6fcbdc73895184eaa1))
* Fix link that points to the actual iac-terragrunt module ([2ab9790](https://github.com/Excoriate/daggerverse/commit/2ab979079a46eee13ac32b5a43a4fd2f285f9507))

## [1.1.2](https://github.com/Excoriate/daggerverse/compare/v1.1.1...v1.1.2) (2024-01-30)


### Other

* update dependencies, test publish ([92627e9](https://github.com/Excoriate/daggerverse/commit/92627e9c3dddab31bbcc5febb5c2866c98a9444d))

## [1.1.1](https://github.com/Excoriate/daggerverse/compare/v1.1.0...v1.1.1) (2024-01-30)


### Bug Fixes

* explicit dependency for daggerx ([471e1b2](https://github.com/Excoriate/daggerverse/commit/471e1b2e289ab79fa016985226a0ab261f4febb0))


### Other

* **main:** release 1.1.0 ([#4](https://github.com/Excoriate/daggerverse/issues/4)) ([88d3729](https://github.com/Excoriate/daggerverse/commit/88d3729b25a95ea8146d2b73f9b6957717c7be7b))

## [1.1.0](https://github.com/Excoriate/daggerverse/compare/v1.0.1...v1.1.0) (2024-01-30)


### Features

* Add common terragrunt functions ([5f95ba7](https://github.com/Excoriate/daggerverse/commit/5f95ba763efda36c807e67b41de14c911c78010f))
* Add common terragrunt functions ([220908d](https://github.com/Excoriate/daggerverse/commit/220908d3da331fec2a181b2aeb3a3f0f5196a5c7))
* Add sync utility ([6f93818](https://github.com/Excoriate/daggerverse/commit/6f93818f4ab76e19d92faabd2cb03a0da444f864))
* Add tg init, add reusable functions ([d6d87c8](https://github.com/Excoriate/daggerverse/commit/d6d87c886782997d8b0c820d6bc67a6afe740f19))


### Other

* Add required ignored files from jetbrains ([34369db](https://github.com/Excoriate/daggerverse/commit/34369db2b2bfc0b146611153c708a5abdd0fef59))

## [1.0.1](https://github.com/Excoriate/daggerverse/compare/v1.0.0...v1.0.1) (2024-01-22)


### Docs

* add README.md, add reusable taskfile for dagger common operations ([a09f56b](https://github.com/Excoriate/daggerverse/commit/a09f56b63c91be5c7ff1b856fe46de1bede38d89))
* add README.md, add reusable taskfile for dagger common operations ([8a21d4f](https://github.com/Excoriate/daggerverse/commit/8a21d4fbe76ef0a8ba861aebf1af3c53340acf5d))

## 1.0.0 (2024-01-22)


### Features

* Add repo setup ([6b7c984](https://github.com/Excoriate/daggerverse/commit/6b7c984550df5e9725027ee469ef051a374a6bc5))
* Add working version ([1aa9e19](https://github.com/Excoriate/daggerverse/commit/1aa9e19e38acb329e71553e55bfa9ce65ca0aecf))
* first commit ([242b312](https://github.com/Excoriate/daggerverse/commit/242b312b5bfdc84f69828660844fce6636dd28f6))


### Docs

* Add readme, add logo ([5f074f4](https://github.com/Excoriate/daggerverse/commit/5f074f411dd4db5a936fa22df9c65be98653a5fd))


### Other

* Add Dagger tasks, and other repo settings ([8a0cd49](https://github.com/Excoriate/daggerverse/commit/8a0cd494e29ccd0c66a2801277045822e4cc4bbe))
* Remove go-build from allowed hooks ([559072f](https://github.com/Excoriate/daggerverse/commit/559072f330ac87199e8528d40f3c57b8c3dd9010))
