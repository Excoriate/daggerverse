use crate::configuration::NewDaggerModule;
use crate::naming::{to_camel_case, to_pascal_case};
use pathdiff::diff_paths;
use std::fs;
use std::io::Error;
use std::path::Path; // Ensure pathdiff is imported
use walkdir::WalkDir;

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
    dest_dir: &str,
    module_name: &str,
) -> Result<(), Error> {
    let mut dirs_to_process = vec![(template_dir.to_string(), dest_dir.to_string())];

    while let Some((current_src_dir, current_dest_dir)) = dirs_to_process.pop() {
        println!("Processing directory: {}", current_src_dir);

        for entry in fs::read_dir(&current_src_dir)? {
            let entry = entry?;
            let path = entry.path();
            let dest_path = Path::new(&current_dest_dir).join(entry.file_name());

            if path.is_dir() {
                dirs_to_process.push((
                    path.to_str().unwrap().to_string(),
                    dest_path.to_str().unwrap().to_string(),
                ));
            } else {
                let content = fs::read_to_string(&path)?;
                let replaced_content = replace_module_name(&content, module_name);
                fs::write(&dest_path, replaced_content)?;
            }
        }
    }

    println!(
        "Templates copied and processed successfully from {} to {}",
        template_dir, dest_dir
    );
    Ok(())
}

pub fn replace_module_name(content: &str, module_name: &str) -> String {
    content.replace("{{module_name}}", module_name)
}

pub fn copy_and_process_templates_with_exclusions(
    module_cfg: &NewDaggerModule,
    template_dir: &str,
    dest_dir: &str,
    exclusions: Option<Vec<String>>,
) -> Result<(), Error> {
    let exclusions = exclusions.unwrap_or_default();
    println!(
        "Starting to copy and process templates from {} to {}",
        template_dir, dest_dir
    );

    let mut dirs_to_process = vec![(template_dir.to_string(), dest_dir.to_string())];

    while let Some((current_src_dir, current_dest_dir)) = dirs_to_process.pop() {
        println!("Processing directory: {}", current_src_dir);

        for entry in fs::read_dir(&current_src_dir)? {
            let entry = entry?;
            let path = entry.path();
            let file_name = entry.file_name().to_string_lossy().to_string();

            if exclusions.contains(&file_name) {
                println!("Skipping excluded file: {}", file_name);
                continue;
            }

            if path.is_dir() {
                let new_dest_dir = format!("{}/{}", current_dest_dir, file_name);
                println!("Creating directory: {}", new_dest_dir);
                fs::create_dir_all(&new_dest_dir)?;
                dirs_to_process.push((path.to_string_lossy().to_string(), new_dest_dir));
            } else {
                let content = fs::read_to_string(&path)?;
                let new_content = process_template_content(&content, module_cfg);

                let dest_file_name = file_name.replace(".tmpl", "");
                let dest_path = format!("{}/{}", current_dest_dir, dest_file_name);
                println!("Writing file: {}", dest_path);
                fs::write(dest_path, new_content)?;
            }
        }
    }

    println!(
        "Templates copied and processed successfully from {} to {}",
        template_dir, dest_dir
    );
    Ok(())
}

pub fn copy_template_files(template_dir: &str, target_dir: &str) -> Result<(), Error> {
    // Adjust the copying logic to handle both full and light module structures
    for entry in WalkDir::new(template_dir)
        .into_iter()
        .filter_map(|e| e.ok())
    {
        let path = entry.path();
        let relative_path = path.strip_prefix(template_dir).unwrap();
        let target_path = Path::new(target_dir).join(relative_path);

        if path.is_dir() {
            fs::create_dir_all(&target_path)?;
        } else {
            fs::copy(path, &target_path)?;
        }
    }

    Ok(())
}
