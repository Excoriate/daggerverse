use std::fs;
use std::io::Error;
use std::path::Path;
use crate::configuration::NewDaggerModule;
use crate::templating::replace_module_name;

pub fn copy_readme_and_license(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    let readme_dest_path = format!("{}/README.md", module_cfg.path);
    let license_dest_path = format!("{}/LICENSE", module_cfg.path);
    println!(
        "Copying README.md and LICENSE files 📄: {}",
        module_cfg.name
    );

    // Ensure the destination directory exists
    fs::create_dir_all(&module_cfg.path)?;

    // Copy the README.md and LICENSE files from the template directory to the module path
    fs::copy(".daggerx/templates/README.md", &readme_dest_path)?;
    fs::copy(".daggerx/templates/LICENSE", &license_dest_path)?;

    // Replace placeholders in README.md if any
    let readme_content = fs::read_to_string(&readme_dest_path)?;
    let new_readme_content = readme_content.replace("[@MODULE_NAME]", &module_cfg.name);
    fs::write(readme_dest_path, new_readme_content)?;

    Ok(())
}

pub fn update_readme_content(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    let readme_path = format!("{}/README.md", module_cfg.path);
    println!("Updating README.md content 📄: {}", module_cfg.name);

    if !Path::new(&readme_path).exists() {
        return Err(Error::new(
            std::io::ErrorKind::NotFound,
            format!("README.md file not found in {}", module_cfg.path),
        ));
    }

    let readme_content = fs::read_to_string(&readme_path)?;
    let new_content = replace_module_name(&readme_content, &module_cfg.name);
    fs::write(&readme_path, new_content)?;

    Ok(())
}
