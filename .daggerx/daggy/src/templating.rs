use crate::configuration::NewDaggerModule;
use crate::naming::{to_camel_case, to_pascal_case};

pub fn process_template_content(content: &str, module_cfg: &NewDaggerModule) -> String {
    let pkg_name = module_cfg
        .name
        .to_string()
        .to_lowercase()
        .trim()
        .replace(" ", "-");
    let pascal_case_name = to_pascal_case(&module_cfg.name);
    let camel_case_name = to_camel_case(&module_cfg.name);
    let lowercase_name = module_cfg.name.to_lowercase();

    let content = content.replace("{{.module_name_pkg}}", &pkg_name);
    let content = content.replace("{{.module_name}}", &pascal_case_name);
    let content = content.replace("{{.module_name_camel}}", &camel_case_name);
    content.replace("{{.module_name_lowercase}}", &lowercase_name)
}
