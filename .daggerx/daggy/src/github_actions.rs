use std::fs;
use std::io::Error;
use crate::configuration::NewDaggerModule;
use crate::templating::process_template_content;

pub fn generate_github_actions_workflow(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    println!("Generating GitHub Actions workflow 🚀: {}", module_cfg.name);
    fs::create_dir_all(&module_cfg.github_actions_workflow_path)?;
    let template_path = ".daggerx/templates/github/workflows/mod-template-ci.yaml.tmpl";
    let output_path = &module_cfg.github_actions_workflow;

    let template_content = fs::read_to_string(template_path)?;
    let new_content = process_template_content(&template_content, module_cfg);
    fs::write(output_path, new_content)?;

    Ok(())
}
