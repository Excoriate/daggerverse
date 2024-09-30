mod args;
mod command_utils;
mod configuration;
mod dagger_commands;
mod dagger_json;
mod dagger_utils;
mod git;
mod github_actions;
mod naming;
mod readme_and_docs;
mod templating;
mod utils;

use args::Args;
use clap::Parser;
use command_utils::{run_command_with_output, run_go_fmt};
use configuration::{get_module_configurations, NewDaggerModule};
use dagger_commands::run_dagger_develop;
use dagger_json::{update_dagger_json, update_examples_dagger_json, update_tests_dagger_json};
use git::find_git_root;
use github_actions::generate_github_actions_workflow;
use naming::to_pascal_case;
use readme_and_docs::copy_readme_and_license;
use std::env;
use std::fs;
use std::io::{Error, ErrorKind};
use std::path::Path;
use templating::{copy_and_process_templates, replace_module_name};
use utils::copy_dir_all;
use dagger_utils::{find_dagger_modules, dagger_module_exists};

fn main() -> Result<(), Error> {
    let args: Args = Args::parse();

    match args.task.as_str() {
        "create" => create_module_task(&args),
        "develop" => develop_modules(),
        _ => {
            eprintln!("Unknown task: {}", args.task);
            Err(Error::new(ErrorKind::InvalidInput, "Unknown task"))
        }
    }
}

fn create_module_task(args: &Args) -> Result<(), Error> {
    match &args.module {
        Some(module) => {
            let module_type = args.module_type.as_deref().unwrap_or("full");
            create_module(module, module_type)
        }
        None => {
            eprintln!("Module name is required for 'create' task");
            Err(Error::new(ErrorKind::InvalidInput, "Module name is required"))
        }
    }
}

fn create_module(module: &str, module_type: &str) -> Result<(), Error> {
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

fn print_module_info(new_module: &NewDaggerModule) {
    println!("Module path: {}", new_module.path);
    println!("Module src path: {}", new_module.module_src_path);
    println!("Module test src path: {}", new_module.module_test_src_path);
    println!(
        "GitHub Actions workflow path: {}",
        new_module.github_actions_workflow_path
    );
}

fn format_code(new_module: &NewDaggerModule) -> Result<(), Error> {
    println!("Running go fmt and ensuring the code is formatted correctly 🧹");
    run_go_fmt(&new_module.path)?;
    run_go_fmt(&format!("{}/examples/go", new_module.path))?;
    run_go_fmt(&new_module.module_test_src_path)
}

fn print_success_message(new_module: &NewDaggerModule) {
    println!("Module \"{}\" initialized successfully 🎉", new_module.name);
    println!("Don't forget to add it to GitHub Actions workflow 'release.yml' when your module is ready for release.");
    println!("It's recommended to run just cilocal <newmodule> to test the module locally before releasing it.");
}


fn develop_modules() -> Result<(), Error> {
    // Ensure we're in a Git repository
    if !Path::new(".git").exists() {
        return Err(Error::new(
            ErrorKind::NotFound,
            "Error: This script must be run from the root of a Git repository.",
        ));
    }

    println!("Git repository detected. Proceeding...");

    // Find all directories containing a 'dagger.json' file
    let modules = find_dagger_modules()?;

    if modules.is_empty() {
        println!("No modules found.");
        return Ok(());
    }

    // Initialize counters
    let total_modules = modules.len();
    let mut successful_modules = 0;
    let mut failed_modules = 0;

    println!("Identifying modules with dagger.json files...");
    for dir in &modules {
        println!("Module identified: {}", dir);
    }

    println!("\nRunning dagger develop in identified modules...\n");

    for dir in &modules {
        print!("Developing module: {}... ", dir);

        if Path::new(&format!("{}/dagger.json", dir)).exists() {
            println!("Entering directory: {}", dir);
            match run_dagger_develop(dir) {
                Ok(_) => {
                    println!("✅ Successfully developed module: {}", dir);
                    successful_modules += 1;
                }
                Err(e) => {
                    println!("❌ Failed to develop module: {}", dir);
                    eprintln!("Error: {}", e);
                    failed_modules += 1;
                }
            }
        } else {
            println!("Skipped 🚫 No dagger.json found in: {}", dir);
        }
    }

    println!("\n");

    if successful_modules == total_modules {
        println!(
            "Dagger develop completed for all {} modules successfully! 🎉",
            total_modules
        );
    } else if failed_modules > 0 {
        println!(
            "Dagger develop completed with {} successes ✅ and {} failures ❌.",
            successful_modules, failed_modules
        );
        return Err(Error::new(
            ErrorKind::Other,
            "Some modules failed to develop",
        ));
    } else {
        println!(
            "Dagger develop completed with {} successes ✅. Please check the output above.",
            successful_modules
        );
    }

    Ok(())
}

fn initialize_module(module_cfg: &NewDaggerModule) -> Result<(), Error> {
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
fn initialize_examples(module_cfg: &NewDaggerModule) -> Result<(), Error> {
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
fn initialize_tests(module_cfg: &NewDaggerModule) -> Result<(), Error> {
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


