use std::fs::{self, File};
use std::io::{Write, Error, ErrorKind};
use std::path::Path;
use std::process::{Command, Output, Stdio};
use clap::Parser;
use serde::Deserialize;

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
struct Args {
    /// Task is the name of the task to run
    #[arg(short = 't', long = "task")]
    task: String,

    /// Module is the name of the dagger module to generate.
    #[arg(short = 'm', long = "module")]
    module: String,
}

#[derive(Deserialize)]
struct NewDaggerModule {
    path: String,
    name: String,
    github_actions_workflow_path: String,
    github_actions_workflow: String,
}

fn main() -> Result<(), Error> {
    let args: Args = Args::parse();

    match args.task.as_str() {
        "create" => {
            create_module(&args.module)?;
        },
        _ => {
            eprintln!("Unknown task: {}", args.task);
            std::process::exit(1);
        }
    }

    Ok(())
}

// Create a new module in the root of the current directory.
fn create_module(module: &str) -> Result<(), Error> {
    println!("Creating module 🚀: {}", module);
    let new_module = dagger_module_exists(module)?;

    // Initialize the new module
    initialize_module(&new_module)?;

    // Initialize tests for the module
    initialize_tests(&new_module)?;

    // Copy README and LICENSE files
    copy_readme_and_license(&new_module)?;

    // Generate GitHub Actions workflow
    generate_github_actions_workflow(&new_module)?;

    println!("Module \"{}\" initialized successfully 🎉", new_module.name);
    println!("Don't forget to add it to GitHub Actions when you are ready! 🚀");

    Ok(())
}

fn initialize_module(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    fs::create_dir_all(&module_cfg.path)?;
    run_command_with_output(&format!("dagger init --sdk go --name {}", module_cfg.name), &module_cfg.path)?;

    let dagger_json_path = format!("{}/dagger.json", module_cfg.path);
    let dagger_json_content = fs::read_to_string(&dagger_json_path)?;
    let new_dagger_json_content = dagger_json_content.replace(
        "}",
        r#",
        "exclude": ["../.direnv", "../.devenv", "../go.work", "../go.work.sum", "tests"]
    }"#,
    );
    fs::write(dagger_json_path, new_dagger_json_content)?;

    run_command_with_output(&format!("dagger develop -m {}", module_cfg.name), &module_cfg.path)?;

    Ok(())
}

fn initialize_tests(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    let tests_path = format!("{}/tests", module_cfg.path);
    fs::create_dir_all(&tests_path)?;
    run_command_with_output("dagger init --sdk go --name tests", &tests_path)?;

    let tests_dagger_json_path = format!("{}/tests/dagger.json", module_cfg.path);
    let tests_dagger_json_content = fs::read_to_string(&tests_dagger_json_path)?;
    let new_tests_dagger_json_content = tests_dagger_json_content.replace(
        "}",
        r#",
        "exclude": ["../../.direnv", "../../.devenv", "../../.go.work", "../../.go.work.sum"]
    }"#,
    );
    fs::write(tests_dagger_json_path, new_tests_dagger_json_content)?;

    fs::copy(".daggerx/templates/tests/main.go.tmpl", format!("{}/tests/dagger/main.go", module_cfg.path))?;
    run_command_with_output("dagger install ../", &tests_path)?;
    run_command_with_output("dagger develop -m tests", &tests_path)?;

    Ok(())
}

fn copy_readme_and_license(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    fs::copy(".daggerx/templates/README.md", format!("{}/README.md", module_cfg.path))?;
    fs::copy(".daggerx/templates/LICENSE", format!("{}/LICENSE", module_cfg.path))?;

    let readme_path = format!("{}/README.md", module_cfg.path);
    let readme_content = fs::read_to_string(&readme_path)?;
    let new_readme_content = readme_content.replace("[@MODULE_NAME]", &module_cfg.name);
    fs::write(readme_path, new_readme_content)?;

    let testdata_path = format!("{}/tests/testdata", module_cfg.path);
    fs::create_dir_all(testdata_path)?;

    Ok(())
}

fn generate_github_actions_workflow(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    fs::create_dir_all(&module_cfg.github_actions_workflow_path)?;
    let template_path = ".daggerx/templates/github/workflows/mod-template-ci.yaml.tmpl";
    let output_path = &module_cfg.github_actions_workflow;

    let template_content = fs::read_to_string(template_path)?;
    let new_content = template_content.replace("{{ mod }}", &module_cfg.name);
    let mut output_file = File::create(output_path)?;
    output_file.write_all(new_content.as_bytes())?;

    Ok(())
}

fn dagger_module_exists(module: &str) -> Result<NewDaggerModule, Error> {
    if module.is_empty() {
        return Err(Error::new(ErrorKind::InvalidInput, "Module name cannot be empty"));
    }

    // Check if the module already exists in the root of this directory.
    if Path::new(module).exists() {
        return Err(Error::new(ErrorKind::AlreadyExists, "Module already exists"));
    }

    let module_path_full = std::env::current_dir()?.join(module);
    let current_root_dir = std::env::current_dir()?;

    Ok(NewDaggerModule {
        path: module_path_full.to_string_lossy().to_string(),
        name: module.to_string(),
        github_actions_workflow_path: current_root_dir.join(".github/workflows").to_string_lossy().to_string(),
        github_actions_workflow: current_root_dir.join(".github/workflows").join(format!("mod-{}-ci.yaml", module)).to_string_lossy().to_string(),
    })
}

fn run_command_with_output(command: &str, current_dir: &str) -> Result<Output, Error> {
    println!("Running command: {}", command);

    let output = Command::new("sh")
        .arg("-c")
        .arg(command)
        .current_dir(current_dir)
        .stdout(Stdio::inherit())
        .stderr(Stdio::inherit())
        .output()?;

    if !output.status.success() {
        return Err(Error::new(ErrorKind::Other, format!("Command failed with exit code: {}", output.status)));
    }

    Ok(output)
}
