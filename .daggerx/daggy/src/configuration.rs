use std::env;
use std::io::Error;

#[derive(Debug)]
pub struct NewDaggerModule {
    pub path: String,
    pub name: String,
    pub module_src_path: String,
    pub module_test_src_path: String,
    pub github_actions_workflow_path: String,
    pub github_actions_workflow: String,
    pub module_type: String,
}

pub fn get_module_configurations(
    module: &str,
    module_type: &str,
) -> Result<NewDaggerModule, Error> {
    let module_path_full = env::current_dir()?.join(module);
    let current_root_dir = env::current_dir()?;

    Ok(NewDaggerModule {
        path: module_path_full.to_string_lossy().to_string(),
        module_src_path: module_path_full.to_string_lossy().to_string(),
        module_test_src_path: module_path_full.join("tests").to_string_lossy().to_string(),
        name: module.to_string(),
        github_actions_workflow_path: current_root_dir
            .join(".github/workflows")
            .to_string_lossy()
            .to_string(),
        github_actions_workflow: current_root_dir
            .join(".github/workflows")
            .join(format!("ci-mod-{}.yaml", module))
            .to_string_lossy()
            .to_string(),
        module_type: module_type.to_string(),
    })
}
