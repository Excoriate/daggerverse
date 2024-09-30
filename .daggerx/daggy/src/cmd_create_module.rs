use std::env;
use std::fs;
use std::io::{Error, ErrorKind};

use crate::command_utils::run_command_with_output;
use crate::configuration::{get_module_configurations, NewDaggerModule};
use crate::dagger_json::{
    update_dagger_json, update_examples_dagger_json, update_tests_dagger_json,
};
use crate::dagger_utils::dagger_module_exists;
use crate::git::find_git_root;
use crate::github_actions::generate_github_actions_workflow;
use crate::readme_and_docs::copy_readme_and_license;
use crate::templating::copy_and_process_templates;
use crate::utils::copy_dir_all;

// Remove this line
// use std::path::PathBuf;

pub fn create_module(module: &str, module_type: &str) -> Result<(), Error> {
    println!("Creating module 🚀: {}", module);
    dagger_module_exists(module)?;

    let git_root = find_git_root()?;
    env::set_current_dir(git_root)?;

    let new_module = get_module_configurations(module, module_type)?;
    print_module_info(&new_module);

    initialize_module(&new_module)?;
    initialize_tests(&new_module)?;
    initialize_examples(&new_module)?;
    copy_readme_and_license(&new_module)?;
    generate_github_actions_workflow(&new_module)?;

    format_code(&new_module)?;

    print_success_message(&new_module);

    Ok(())
}

pub fn initialize_module(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    // Create the module directory
    fs::create_dir_all(&module_cfg.path)?;
    println!("Creating parent module 📦: {}", module_cfg.name);

    // Change to the module directory
    env::set_current_dir(&module_cfg.path)?;

    // Run dagger init
    run_command_with_output(
        &format!("dagger init --sdk go --name {} --source .", module_cfg.name),
        ".",
    )?;

    // Copy and process templates
    copy_and_process_templates(module_cfg, "../.daggerx/templates/module", ".")?;

    // Update dagger.json
    update_dagger_json(module_cfg)?;

    // Edit go.mod to set the correct module path
    let go_mod_edit_command = format!(
        "go mod edit -module github.com/Excoriate/daggerverse/{}",
        module_cfg.name
    );
    run_command_with_output(&go_mod_edit_command, ".")?;

    // Handle different module types
    match module_cfg.module_type.as_str() {
        "full" => {
            // Full module initialization logic
            run_command_with_output("dagger develop", ".")?;
        }
        "light" => {
            // Light module initialization logic
            println!("Initializing light module type...");
            // TODO: Add any specific logic for light module type here
        }
        _ => return Err(Error::new(ErrorKind::InvalidInput, "Invalid module type")),
    }

    // Change back to the root directory
    env::set_current_dir("..")?;

    Ok(())
}

// New function
// Modified function
pub fn initialize_examples(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    let examples_path = format!("{}/examples/go", module_cfg.path);
    println!(
        "Creating examples module (recipes)  📄: {}",
        module_cfg.name
    );

    // Create the examples directory
    fs::create_dir_all(&examples_path)?;

    // Change to the examples directory
    env::set_current_dir(&examples_path)?;

    // Run dagger init
    run_command_with_output("dagger init --sdk go --name go --source .", ".")?;

    // Copy and process templates
    copy_and_process_templates(module_cfg, "../../../.daggerx/templates/examples/go", ".")?;

    // Update dagger.json
    update_examples_dagger_json(module_cfg)?;

    // Edit go.mod
    let go_mod_edit_command = format!(
        "go mod edit -module github.com/Excoriate/daggerverse/{}/examples/go",
        module_cfg.name
    );
    run_command_with_output(&go_mod_edit_command, ".")?;

    // Copy testdata/common directory
    let src_testdata = "../../../.daggerx/templates/examples/go/testdata/common";
    let dest_testdata = "testdata/common";
    copy_dir_all(src_testdata, dest_testdata)?;

    // Run dagger install and develop
    run_command_with_output("dagger install ../../", ".")?;
    // run_command_with_output("dagger develop -m go", ".")?;
    run_command_with_output("dagger develop", ".")?;

    // Change back to the root directory
    env::set_current_dir("../../..")?;

    Ok(())
}

// Modified function
pub fn initialize_tests(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    let tests_path = format!("{}/tests", module_cfg.path);
    println!("Creating tests module (tests) 🧪: {}", module_cfg.name);

    // Create the tests directory
    fs::create_dir_all(&tests_path)?;

    // Change to the tests directory
    env::set_current_dir(&tests_path)?;

    // Run dagger init
    run_command_with_output("dagger init --sdk go --name tests --source .", ".")?;

    // Copy and process templates
    copy_and_process_templates(module_cfg, "../../.daggerx/templates/tests", ".")?;

    // Update dagger.json
    update_tests_dagger_json(module_cfg)?;

    // Edit go.mod
    let go_mod_edit_command = format!(
        "go mod edit -module github.com/Excoriate/daggerverse/{}/tests",
        module_cfg.name
    );
    run_command_with_output(&go_mod_edit_command, ".")?;

    // Copy testdata/common directory
    let src_testdata = "../../.daggerx/templates/tests/testdata/common";
    let dest_testdata = "testdata/common";
    copy_dir_all(src_testdata, dest_testdata)?;

    // Run dagger install and develop
    run_command_with_output("dagger install ../", ".")?;
    // run_command_with_output("dagger develop -m tests", ".")?;
    run_command_with_output("dagger develop", ".")?;

    // Handle different module types
    match module_cfg.module_type.as_str() {
        "full" => {
            // Full module initialization logic
            run_command_with_output("dagger develop", ".")?;
        }
        "light" => {
            // Light module initialization logic
            println!("Initializing light module type...");
            // Add any specific logic for light module type here
        }
        _ => return Err(Error::new(ErrorKind::InvalidInput, "Invalid module type")),
    }

    // Change back to the root directory
    env::set_current_dir("../..")?;

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
