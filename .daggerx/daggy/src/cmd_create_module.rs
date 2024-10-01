use std::env;
use std::fs;
use std::io::Error;

use crate::command_utils::run_command_with_output;
use crate::configuration::{get_module_configurations, NewDaggerModule};
use crate::dagger_json::{
    update_dagger_json, update_examples_dagger_json, update_tests_dagger_json,
};
use crate::dagger_utils::dagger_module_exists;
use crate::git::find_git_root;
use crate::github_actions::generate_github_actions_workflow;
use crate::readme_and_docs::copy_readme_and_license;
use crate::templating::{copy_and_process_templates, copy_and_process_templates_with_exclusions};
use crate::utils::{calculate_relative_path, copy_dir_all};

// Remove this line
// use std::path::PathBuf;

pub fn create_module(module_name: &str, module_type: &str) -> Result<(), Error> {
    println!("Creating module 🚀: {}", module_name);
    println!("Checking if module already exists 🔍");
    dagger_module_exists(module_name)?;

    println!("Resolving git root 🔍");
    let git_root = find_git_root()?;
    env::set_current_dir(git_root.clone())?;

    println!("Git root resolved 🎉: {}", git_root.to_string());

    println!("Getting module configurations 🔍");
    let new_module = get_module_configurations(module_name, module_type)?;
    print_module_info(&new_module);

    println!("Module Configuration 📋:");
    println!("*************************************************");
    println!("Name 🏷️: {}", new_module.name);
    println!("Type 🛠️: {}", new_module.module_type);
    println!("Path 📂: {}", new_module.path);
    println!("Src Path 📂: {}", new_module.module_src_path);
    println!("Test Src Path 🧪: {}", new_module.module_test_src_path);
    println!("Template Path 📁: {}", new_module.template_path_by_type);
    println!("*************************************************");

    // Initializers
    println!("Starting Initializers 🚀");
    initialize_module(&new_module)?;
    initialize_tests(&new_module)?;
    initialize_examples(&new_module)?;
    println!("Initializers completed 🎉");
    // Post-initializers
    copy_readme_and_license(&new_module)?;
    generate_github_actions_workflow(&new_module)?;

    format_code(&new_module)?;
    print_success_message(&new_module);

    Ok(())
}

pub fn initialize_module(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    println!(
        "Creating & initializing parent module 📦: {}",
        module_cfg.name
    );

    // Create the module directory
    fs::create_dir_all(&module_cfg.path)?;

    // Change to the module directory
    println!(
        "Changing to module directory to run Dagger Init🔗: {}",
        module_cfg.path
    );
    env::set_current_dir(&module_cfg.path)?;

    // Run dagger init
    run_command_with_output(
        &format!("dagger init --sdk go --name {} --source .", module_cfg.name),
        ".",
    )?;

    println!(
        "Dagger Initialization completed 🎉 for module: {}",
        module_cfg.name
    );

    println!(
        "Copying & processing templates 📁: {}",
        module_cfg.template_path_by_type
    );

    // Resolving source path where the templates resides
    let source_template_path = format!("{}/module", module_cfg.template_path_by_type);
    println!("Source template resolved path 📂: {}", source_template_path);

    // Copying & processing templates
    println!(
        "Copying & processing templates 📁 from {} to . (current dir)",
        source_template_path
    );
    copy_and_process_templates(module_cfg, &source_template_path, ".")?;
    println!("Templates copied & processed 🎉");

    // Update dagger.json
    println!("Updating dagger.json 📝");
    update_dagger_json(module_cfg)?;
    println!("dagger.json updated 🎉");

    // Edit go.mod to set the correct module path
    println!("Editing go.mod to set the correct module path 🔗");
    let go_mod_edit_command = format!(
        "go mod edit -module github.com/Excoriate/daggerverse/{}",
        module_cfg.name
    );
    run_command_with_output(&go_mod_edit_command, ".")?;
    println!("go.mod edited 🎉");

    // Full module initialization logic
    println!("Running dagger develop 🔗 in current dir");
    run_command_with_output("dagger develop", ".")?;
    println!("dagger develop completed 🎉");

    // Change back to the root directory
    println!("Changing back to the root directory 🔙");
    env::set_current_dir("..")?;

    println!("Module (Parent) initialized and created Successfully 🎉");
    Ok(())
}

pub fn initialize_examples(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    let examples_path = format!("{}/examples/go", module_cfg.path);
    println!("Creating examples module (recipes) 📄: {}", module_cfg.name);

    // Create the examples directory
    fs::create_dir_all(&examples_path)?;
    println!("Examples directory created 📂: {}", examples_path);

    // Change to the examples directory
    println!("Changing to examples directory 🔄: {}", examples_path);
    env::set_current_dir(&examples_path)?;

    // Run dagger init
    println!("Running dagger init for examples module 🔗");
    run_command_with_output("dagger init --sdk go --name go --source .", ".")?;
    println!("Dagger init completed for examples module 🎉");

    // Copy and process templates
    let source_template_path = format!("{}/examples/go", module_cfg.template_path_by_type);
    println!(
        "Copying & processing templates 📁 from {} to . (current dir)",
        source_template_path
    );
    copy_and_process_templates(module_cfg, &source_template_path, ".")?;
    println!("Templates copied & processed for examples module 🎉");

    // Update dagger.json
    println!("Updating dagger.json for examples module 📝");
    update_examples_dagger_json(module_cfg)?;
    println!("dagger.json updated for examples module 🎉");

    // Edit go.mod
    println!("Editing go.mod to set the correct module path 🔗");
    let go_mod_edit_command = format!(
        "go mod edit -module github.com/Excoriate/daggerverse/{}/examples/go",
        module_cfg.name
    );
    run_command_with_output(&go_mod_edit_command, ".")?;
    println!("go.mod edited for examples module 🎉");

    // Copy testdata/common directory
    let src_testdata = format!(
        "{}/examples/go/testdata/common",
        module_cfg.template_path_by_type
    );
    let dest_testdata = "testdata/common";
    println!(
        "Copying testdata/common directory from {} to {}",
        src_testdata, dest_testdata
    );
    copy_dir_all(src_testdata, dest_testdata)?;
    println!("testdata/common directory copied 🎉");

    // Run dagger install and develop
    println!("Running dagger install for examples module 🔗");
    run_command_with_output("dagger install ../../", ".")?;
    println!("dagger install completed for examples module 🎉");

    println!("Running dagger develop for examples module 🔗");
    run_command_with_output("dagger develop", ".")?;
    println!("dagger develop completed for examples module 🎉");

    // Change back to the root directory
    println!("Changing back to the root directory 🔙");
    env::set_current_dir("../../..")?;

    println!("Examples module initialized and created successfully 🎉");
    Ok(())
}

pub fn initialize_tests(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    let tests_path = format!("{}/tests", module_cfg.path);
    println!("Creating tests module (tests) 🧪: {}", module_cfg.name);

    // Create the tests directory
    fs::create_dir_all(&tests_path)?;
    println!("Tests directory created 📂: {}", tests_path);

    // Change to the tests directory
    println!(
        "Changing current directory to tests directory 🔄: {}",
        tests_path
    );
    env::set_current_dir(&tests_path)?;

    // Run dagger init
    println!(
        "Running dagger init for tests module 🔗 in current dir {}",
        tests_path
    );
    run_command_with_output("dagger init --sdk go --name tests --source .", ".")?;
    println!(
        "Dagger init completed for tests module 🎉 in current dir {}",
        tests_path
    );

    // Copy and process templates
    let source_template_path = format!("{}/tests", module_cfg.template_path_by_type);
    println!(
        "Copying & processing templates 📁 from {} to {}",
        source_template_path, tests_path
    );

    if let Err(e) = copy_and_process_templates_with_exclusions(
        module_cfg,
        &source_template_path,
        ".",
        Some(vec![]),
    ) {
        println!(
            "❌ Error copying and processing templates from {} to {}: {}. Current dir {}",
            source_template_path,
            tests_path,
            e,
            env::current_dir()?.to_string_lossy()
        );
        return Err(e);
    }

    println!(
        "Templates copied & processed from {} to {} 🎉",
        source_template_path, tests_path
    );
    // Update dagger.json
    println!("Updating dagger.json for tests module 📝");
    update_tests_dagger_json(module_cfg)?;
    println!("dagger.json updated for tests module 🎉");

    // Edit go.mod
    println!("Editing go.mod to set the correct module path 🔗");
    let go_mod_edit_command = format!(
        "go mod edit -module github.com/Excoriate/daggerverse/{}/tests",
        module_cfg.name
    );
    run_command_with_output(&go_mod_edit_command, ".")?;
    println!("go.mod edited for tests module 🎉");

    // Copy testdata/common directory
    let src_testdata = format!("{}/tests/testdata/common", module_cfg.template_path_by_type);
    let dest_testdata = "testdata/common";
    println!(
        "Copying testdata/common directory from {} to {}",
        src_testdata, dest_testdata
    );
    copy_dir_all(src_testdata, dest_testdata)?;
    println!("testdata/common directory copied 🎉");

    // Run dagger install and develop
    println!("Running dagger install for tests module 🔗");
    run_command_with_output("dagger install ../", ".")?;
    println!("dagger install completed for tests module 🎉");

    println!("Running dagger develop for tests module 🔗");
    run_command_with_output("dagger develop", ".")?;
    println!("dagger develop completed for tests module 🎉");

    // Change back to the root directory
    println!("Changing back to the root directory 🔙");
    env::set_current_dir("../..")?;

    println!("Tests module initialized and created successfully 🎉");
    Ok(())
}

// Add these functions from main.rs
pub fn print_module_info(new_module: &NewDaggerModule) {
    println!("Module path: {}", new_module.path);
    println!("Module src path: {}", new_module.module_src_path);
    println!("Module test src path: {}", new_module.module_test_src_path);
    println!(
        "GitHub Actions workflow path: {}",
        new_module.github_actions_workflow_path
    );
}

pub fn format_code(new_module: &NewDaggerModule) -> Result<(), Error> {
    println!("Running go fmt and ensuring the code is formatted correctly 🧹");
    crate::command_utils::run_go_fmt(&new_module.path)?;
    crate::command_utils::run_go_fmt(&format!("{}/examples/go", new_module.path))?;
    crate::command_utils::run_go_fmt(&new_module.module_test_src_path)
}

pub fn print_success_message(new_module: &NewDaggerModule) {
    println!("Module \"{}\" initialized successfully 🎉", new_module.name);
    println!("Don't forget to add it to GitHub Actions workflow 'release.yml' when your module is ready for release.");
    println!("It's recommended to run just cilocal <newmodule> to test the module locally before releasing it.");
}
