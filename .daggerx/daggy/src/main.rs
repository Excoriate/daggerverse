use std::env;
use std::fs::{self, File};
use std::io::{Write, Error, ErrorKind};
use std::path::Path;
use std::process::{Command, Output, Stdio};
use clap::Parser;
use serde::Deserialize;
use regex::Regex;

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

    // Update README content
    update_readme_content(&new_module)?;

    // Generate GitHub Actions workflow
    generate_github_actions_workflow(&new_module)?;

    println!("Module \"{}\" initialized successfully 🎉", new_module.name);
    println!("Don't forget to add it to GitHub Actions workflow 'release.yml' when your module is ready for release.");
    println!("It's recommended to run just cilocal <newmodule> to test the module locally before releasing it.");

    Ok(())
}

fn update_dagger_json(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    let dagger_json_path = format!("{}/dagger.json", module_cfg.path);
    let dagger_json_content = fs::read_to_string(&dagger_json_path)?;
    let new_dagger_json_content = dagger_json_content.replace(
        "}",
        r#",
        "exclude": ["../.direnv", "../.devenv", "../go.work", "../go.work.sum", "tests"]
    }"#,
    );

    fs::write(dagger_json_path, new_dagger_json_content)?;

    Ok(())
}

fn initialize_module(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    // Creates the destination directory for the new module.
    fs::create_dir_all(&module_cfg.path)?;
    run_command_with_output(&format!("dagger init --sdk go --name {}", module_cfg.name), &module_cfg.path)?;

    // Update dagger.json to exclude some files
    update_dagger_json(module_cfg)?;

    // Running dagger develop to initialize the module
    run_command_with_output(&format!("dagger develop -m {}", module_cfg.name), &module_cfg.path)?;

    let template_dir = ".daggerx/templates/module";
    let destination_dir = format!("{}/dagger", module_cfg.path);

    fs::create_dir_all(&destination_dir)?;
    copy_and_replace_templates(template_dir, &destination_dir, &module_cfg.name)?;

    // replace_readme_content(module_cfg)?;

    Ok(())
}

fn copy_and_replace_templates(template_dir: &str, destination_dir: &str, module_name: &str) -> Result<(), Error> {
    for entry in fs::read_dir(template_dir)? {
        let entry = entry?;
        let path = entry.path();

        if path.is_dir() {
            let new_dir = format!("{}/{}", destination_dir, entry.file_name().to_string_lossy());
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

fn initialize_tests(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    let tests_path = format!("{}/tests", module_cfg.path);
    fs::create_dir_all(&tests_path)?;
    run_command_with_output("dagger init --sdk go --name tests", &tests_path)?;

    // remove the main.go that's automatically generated if exists.
    if Path::new(&format!("{}/main.go", tests_path)).exists() {
        println!("Removing main.go from tests module");
        fs::remove_file(format!("{}/main.go", tests_path))?;
    }

    let tests_dagger_json_path = format!("{}/tests/dagger.json", module_cfg.path);
    let tests_dagger_json_content = fs::read_to_string(&tests_dagger_json_path)?;
    let new_tests_dagger_json_content = tests_dagger_json_content.replace(
        "}",
        r#",
        "exclude": ["../../.direnv", "../../.devenv", "../../.go.work", "../../.go.work.sum"]
    }"#,
    );
    fs::write(tests_dagger_json_path, new_tests_dagger_json_content)?;

    let test_template_dir = ".daggerx/templates/module/tests";
    let test_destination_dir = format!("{}/tests/dagger", module_cfg.path);  // Corrected destination path

    copy_and_replace_templates(&test_template_dir, &test_destination_dir, &module_cfg.name)?;

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
    let new_content = replace_module_name_lowercase(&template_content, &module_cfg.name);
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

    let module_path_full = env::current_dir()?.join(module);
    let current_root_dir = env::current_dir()?;

    // Get the current dir, if the env var CURRENT_DIR_OVERRIDE is set, use that instead
    Ok(NewDaggerModule {
        path: module_path_full.to_string_lossy().to_string(),
        name: module.to_string(),
        github_actions_workflow_path: current_root_dir.join(".github/workflows").to_string_lossy().to_string(),
        github_actions_workflow: current_root_dir.join(".github/workflows").join(format!("mod-{}-ci.yaml", module)).to_string_lossy().to_string(),
    })
}

fn run_command_with_output(command: &str, target_dir: &str) -> Result<Output, Error> {
    println!("Running command: {}", command);
    let output = Command::new("sh")
        .arg("-c")
        .arg(command)
        .current_dir(target_dir)
        .stdout(Stdio::inherit())
        .stderr(Stdio::inherit())
        .output()?;

    if !output.status.success() {
        return Err(Error::new(ErrorKind::Other, format!("Command failed with exit code: {} and with error: {}", output.status, String::from_utf8_lossy(&output.stderr))));
    }

    Ok(output)
}

fn replace_module_name(content: &str, module_name: &str) -> String {
    let pascal_case_name = to_pascal_case(module_name);
    let camel_case_name = to_camel_case(module_name);

    let re_pascal = Regex::new(r"\{\{\s*\.\s*module_name\s*\}\}").unwrap();
    let re_camel = Regex::new(r"\{\{\s*\.\s*module_name_camel\s*\}\}").unwrap();
    let re_lowercase = Regex::new(r"\{\{\s*\.\s*module_name_lowercase\s*\}\}").unwrap();

    let content = re_pascal.replace_all(content, &pascal_case_name);
    let content = re_camel.replace_all(&content, &camel_case_name);
    re_lowercase.replace_all(&content, module_name).to_string()
}

fn capitalize_module_name(module_name: &str) -> String {
    let mut chars = module_name.chars();
    match chars.next() {
        None => String::new(),
        Some(first) => first.to_uppercase().collect::<String>() + chars.as_str(),
    }
}

fn replace_module_name_lowercase(content: &str, module_name: &str) -> String {
    let re = Regex::new(r"\{\{\s*\.\s*module_name\s*\}\}").unwrap();
    re.replace_all(content, module_name).to_string()
}

fn update_readme_content(module_cfg: &NewDaggerModule) -> Result<(), Error> {
    let readme_path = format!("{}/README.md", module_cfg.path);

    if !Path::new(&readme_path).exists() {
        return Err(Error::new(ErrorKind::NotFound, format!("README.md file not found in {}", module_cfg.path)));
    }

    let readme_content = fs::read_to_string(&readme_path)?;
    let new_content = replace_module_name_smart(&readme_content, &module_cfg.name);
    fs::write(&readme_path, new_content)?;

    Ok(())
}

fn replace_module_name_smart(content: &str, module_name: &str) -> String {
    let pascal_case_name = to_pascal_case(module_name);
    let lowercase_name = module_name.to_lowercase();

    let re = Regex::new(r"```[\s\S]*?```|`[^`\n]+`|\{\{\s*\.\s*module_name\s*\}\}").unwrap();

    re.replace_all(content, |caps: &regex::Captures| {
        let matched = caps.get(0).unwrap().as_str();
        if matched.starts_with("```") || matched.starts_with("`") {
            // Inside code blocks, use lowercase with hyphens
            matched.replace("{{.module_name}}", &lowercase_name)
        } else {
            // Outside code blocks, use PascalCase without hyphens
            matched.replace("{{.module_name}}", &pascal_case_name)
        }
    }).to_string()
}

fn to_camel_case(s: &str) -> String {
    s.split('-')
        .enumerate()
        .map(|(i, part)| {
            if i == 0 {
                part.to_lowercase()
            } else {
                capitalize_module_name(part)
            }
        })
        .collect()
}

fn to_pascal_case(s: &str) -> String {
    s.split('-')
        .map(capitalize_module_name)
        .collect()
}

fn capitalize_first_letter(s: &str) -> String {
    let mut c = s.chars();
    match c.next() {
        None => String::new(),
        Some(f) => f.to_uppercase().collect::<String>() + c.as_str(),
    }
}
