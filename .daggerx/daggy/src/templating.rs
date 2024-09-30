use crate::configuration::NewDaggerModule;
use crate::naming::{to_camel_case, to_pascal_case};
use std::fs;
use std::io::Error;

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

pub fn copy_and_process_templates(
    module_cfg: &NewDaggerModule,
    template_dir: &str,
    dest_dir: &str,
) -> Result<(), Error> {
    for entry in fs::read_dir(template_dir)? {
        let entry = entry?;
        let path = entry.path();

        if path.is_dir() {
            let new_dir = format!("{}/{}", dest_dir, entry.file_name().to_string_lossy());
            fs::create_dir_all(&new_dir)?;
            copy_and_process_templates(module_cfg, &path.to_string_lossy(), &new_dir)?;
        } else {
            let content = fs::read_to_string(&path)?;
            let new_content = process_template_content(&content, module_cfg);

            let dest_file_name = entry.file_name().to_string_lossy().replace(".tmpl", "");
            let dest_path = format!("{}/{}", dest_dir, dest_file_name);
            fs::write(dest_path, new_content)?;
        }
    }

    Ok(())
}

pub fn copy_and_replace_templates(
    template_dir: &str,
    destination_dir: &str,
    module_name: &str,
) -> Result<(), Error> {
    for entry in fs::read_dir(template_dir)? {
        let entry = entry?;
        let path = entry.path();

        if path.is_dir() {
            let new_dir = format!(
                "{}/{}",
                destination_dir,
                entry.file_name().to_string_lossy()
            );
            fs::create_dir_all(&new_dir)?;
            copy_and_replace_templates(&path.to_string_lossy(), &new_dir, module_name)?;
        } else {
            let content = fs::read_to_string(&path)?;
            let new_content = if path.extension().map_or(false, |ext| ext == "go") {
                replace_module_name(&content, &to_pascal_case(&module_name))
            } else {
                replace_module_name(&content, module_name)
            };

            let dest_file_name = entry.file_name().to_string_lossy().replace(".tmpl", "");
            let dest_path = format!("{}/{}", destination_dir, dest_file_name);
            fs::write(dest_path, new_content)?;
        }
    }

    Ok(())
}

pub fn replace_module_name(content: &str, module_name: &str) -> String {
    let pascal_case_name = to_pascal_case(module_name);
    let camel_case_name = to_camel_case(module_name);

    let re_pascal = regex::Regex::new(r"\{\{\s*\.\s*module_name\s*\}\}").unwrap();
    let re_camel = regex::Regex::new(r"\{\{\s*\.\s*module_name_camel\s*\}\}").unwrap();
    let re_lowercase = regex::Regex::new(r"\{\{\s*\.\s*module_name_lowercase\s*\}\}").unwrap();

    let content = re_pascal.replace_all(content, &pascal_case_name);
    let content = re_camel.replace_all(&content, &camel_case_name);
    re_lowercase.replace_all(&content, module_name).to_string()
}
