# Terraform 🏗️

![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.0-0f0f19.svg?style=flat-square)


A simple [Dagger](https://dagger.io) module to manage **terraform**.

## Configuration 🛠️

Through the [Dagger CLI](https://docs.dagger.io/cli/465058/install), or by using it directly reusing it within your module, you can configure the following options:

* ⚙️ `version`: The version of [Terraform](https://www.terraform.io/) to use.  Default is `latest`.
* 📁 `src`: The path to the source code of the module.
* 🐳 `image`: The container image to use. If not provided, the default one will be used, which is `hashicorp/terraform`.
* 🚢 `ctr`: The dagger container that can be passed to the module. It's an optional parameter. If not provided, the default one will be used.


---

## Features 🎨

| Command or functionality                                         | Command      | Example                                                                       | Status |
|------------------------------------------------------------------|--------------|-------------------------------------------------------------------------------|--------|
| Initialize a Terraform module (support arguments)                | **init**     | `dagger call --src="." init --tfmod="mydir/module" --args="-backend=false"`   | ✅      |
| Create an execution plan (support arguments)                     | **plan**     | `dagger call --src="." plan --tfmod="mydir/module" --args="-var='foo=bar'"`   | ✅      |
| Apply changes (support arguments)                                | **apply**    | `dagger call --src="." apply --tfmod="mydir/module" --args="-auto-approve"`   | ✅      |
| Destroy Terraform-managed infrastructure (support arguments)     | **destroy**  | `dagger call --src="." destroy --tfmod="mydir/module" --args="-auto-approve"` | ✅      |
| Validate the Terraform files (support arguments)                 | **validate** | `dagger call --src="." validate --tfmod="mydir/module" --args=""`             | ✅      |
| Format Terraform files to a canonical format (support arguments) | **fmt**      | `dagger call --src="." fmt --tfmod="mydir/module" --args=""`                  | ✅      |

>NOTE: The commands `plan`, `apply`, and `destroy` supports an extra argument called `--init-args` to pass additional arguments to the command **init** command.

---

## Usage 🚀

Using the published module

  ```bash
  # Assuming you're in a directory with a terraform module in the path terraform/test/tf-module-1
dagger -m github.com/Excoriate/daggerverse/terraform@v1.4.0 call --src="." plan \
--tfmod="terraform/test/tf-module-1" \
--args="-var is_enabled=false, -refresh=false" \
--init-args="-backend=false" stdout  ```
