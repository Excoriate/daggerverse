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
    pub template_cfg: TemplateCfg,
}

#[derive(Debug)]
pub struct TemplateCfg {
    pub github_actions_workflow_ci: String,
    pub github_actions_workflow_ci_template: String,
    pub module_test_template: String,
    pub module_template: String,
    pub module_light_template_dir: String,
    pub module_full_template_dir: String,
}

const GITHUB_ACTIONS_WORKFLOW_DIR: &str = ".github/workflows";
const MODULE_TESTS_DIR: &str = "tests";
const TEMPLATE_DIR: &str = ".daggerx/templates";
const GITHUB_ACTIONS_WORKFLOW_CI_TEMPLATE: &str = "mod-template-ci.yaml.tmpl";
const MODULE_LIGHT_TEMPLATE_DIR: &str = "mod-light";
const MODULE_FULL_TEMPLATE_DIR: &str = "mod-full";

pub fn get_module_configurations(
    module: &str,
    module_type: &str,
) -> Result<NewDaggerModule, Error> {
    let module_path_full = env::current_dir()?.join(module);
    let current_root_dir = env::current_dir()?;

    Ok(NewDaggerModule {
        path: module_path_full.to_string_lossy().to_string(),
        module_src_path: module_path_full.to_string_lossy().to_string(),
        module_test_src_path: module_path_full
            .join(MODULE_TESTS_DIR)
            .to_string_lossy()
            .to_string(),
        name: module.to_string(),
        github_actions_workflow_path: current_root_dir
            .join(GITHUB_ACTIONS_WORKFLOW_DIR)
            .to_string_lossy()
            .to_string(),
        github_actions_workflow: current_root_dir
            .join(GITHUB_ACTIONS_WORKFLOW_DIR)
            .join(format!("ci-mod-{}.yaml", module))
            .to_string_lossy()
            .to_string(),
        module_type: module_type.to_string(),
        template_cfg: TemplateCfg {
            github_actions_workflow_ci: current_root_dir
                .join(GITHUB_ACTIONS_WORKFLOW_DIR)
                .join(GITHUB_ACTIONS_WORKFLOW_CI_TEMPLATE)
                .to_string_lossy()
                .to_string(),
            github_actions_workflow_ci_template: current_root_dir
                .join(TEMPLATE_DIR)
                .join(GITHUB_ACTIONS_WORKFLOW_DIR)
                .join(GITHUB_ACTIONS_WORKFLOW_CI_TEMPLATE)
                .to_string_lossy()
                .to_string(),
            module_test_template: current_root_dir
                .join(TEMPLATE_DIR)
                .join(GITHUB_ACTIONS_WORKFLOW_CI_TEMPLATE)
                .to_string_lossy()
                .to_string(),
            module_template: current_root_dir
                .join(TEMPLATE_DIR)
                .join(GITHUB_ACTIONS_WORKFLOW_CI_TEMPLATE)
                .to_string_lossy()
                .to_string(),
            module_light_template_dir: current_root_dir
                .join(TEMPLATE_DIR)
                .join(MODULE_LIGHT_TEMPLATE_DIR)
                .to_string_lossy()
                .to_string(),
            module_full_template_dir: current_root_dir
                .join(TEMPLATE_DIR)
                .join(MODULE_FULL_TEMPLATE_DIR)
                .to_string_lossy()
                .to_string(),
        },
    })
}
