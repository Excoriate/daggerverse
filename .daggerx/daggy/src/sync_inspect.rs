use crate::args::Args;
use crate::configuration::{get_module_configurations, NewDaggerModule};
use std::fs;
use std::io::{Error, ErrorKind};
use std::path::Path;

#[derive(Debug)]
struct FileChange {
    status: ChangeStatus,
    path: String,
}

#[derive(Debug, PartialEq)]
enum ChangeStatus {
    Modified,
    Added,
    Deleted,
}

pub fn sync_modules_task(args: &Args) -> Result<(), Error> {
    let inspect_type = &args.inspect_type;
    let dry_run = args.dry_run.unwrap_or(false);

    println!("Syncing modules...");
    println!("Inspect type: {}", inspect_type);
    println!("Dry run: {}", dry_run);

    sync_changes(inspect_type, dry_run)
}

pub fn inspect_modules_task(args: &Args) -> Result<(), Error> {
    let inspect_type = &args.inspect_type;
    let dry_run = args.dry_run.unwrap_or(false);

    println!("Inspecting modules...");
    println!("Inspect type: {}", inspect_type);
    println!("Dry run: {}", dry_run);

    inspect_changes(inspect_type, dry_run)
}

fn sync_changes(inspect_type: &str, dry_run: bool) -> Result<(), Error> {
    let modules_to_sync = get_modules_to_process(inspect_type)?;

    for module_type in modules_to_sync {
        println!("Syncing changes for {} module type (dry run: {})", module_type, dry_run);
        let config = get_module_configurations(&format!("module-template-{}", module_type), module_type)?;
        let changes = detect_changes(&config)?;

        if !changes.is_empty() {
            println!("The following changes will be synced:");
            for change in &changes {
                println!("{:?}: {}", change.status, change.path);
            }

            if !dry_run {
                update_template_files(changes, &config)?;
                println!("Changes synced successfully for {} module type", module_type);
            } else {
                println!("Dry run: Changes would be synced for {} module type", module_type);
            }
        } else {
            println!("No changes detected for {} module type", module_type);
        }
    }

    Ok(())
}

fn inspect_changes(inspect_type: &str, dry_run: bool) -> Result<(), Error> {
    let modules_to_inspect = get_modules_to_process(inspect_type)?;

    for module_type in modules_to_inspect {
        println!("Inspecting changes for {} module type (dry run: {})", module_type, dry_run);
        let config = get_module_configurations(&format!("module-template-{}", module_type), module_type)?;
        let changes = detect_changes(&config)?;

        if !changes.is_empty() {
            println!("Changes detected for {} module type:", module_type);
            for change in &changes {
                println!("{:?}: {}", change.status, change.path);
            }
        } else {
            println!("No changes detected for {} module type", module_type);
        }
    }

    Ok(())
}

fn get_modules_to_process(inspect_type: &str) -> Result<Vec<&str>, Error> {
    match inspect_type {
        "full" => Ok(vec!["full"]),
        "light" => Ok(vec!["light"]),
        "all" => Ok(vec!["full", "light"]),
        _ => Err(Error::new(ErrorKind::InvalidInput, "Invalid inspect type. Must be 'full', 'light', or 'all'.")),
    }
}

fn detect_changes(config: &NewDaggerModule) -> Result<Vec<FileChange>, Error> {
    let mut changes = Vec::new();

    // Check main module files
    check_directory_changes(&config.module_src_path, &config.template_path_by_type, &mut changes)?;

    // Check test files
    let test_src_path = Path::new(&config.path).join("tests");
    let test_template_path = Path::new(&config.template_path_by_type).join("tests");
    check_directory_changes(&test_src_path.to_string_lossy(), &test_template_path.to_string_lossy(), &mut changes)?;

    // Check example files
    let example_src_path = Path::new(&config.path).join("examples").join("go");
    let example_template_path = Path::new(&config.template_path_by_type).join("examples").join("go");
    check_directory_changes(&example_src_path.to_string_lossy(), &example_template_path.to_string_lossy(), &mut changes)?;

    Ok(changes)
}

fn check_directory_changes(src_path: &str, template_path: &str, changes: &mut Vec<FileChange>) -> Result<(), Error> {
    let src_dir = Path::new(src_path);
    let template_dir = Path::new(template_path);

    if !src_dir.exists() || !template_dir.exists() {
        return Ok(());
    }

    for entry in fs::read_dir(src_dir)? {
        let entry = entry?;
        let path = entry.path();
        if path.is_file() && path.extension().map_or(false, |ext| ext == "go") {
            let relative_path = path.strip_prefix(src_dir).unwrap();
            let template_file = template_dir.join(relative_path.with_extension("go.tmpl"));

            if !template_file.exists() {
                changes.push(FileChange {
                    status: ChangeStatus::Added,
                    path: relative_path.to_string_lossy().to_string(),
                });
            } else if files_differ(&path, &template_file)? {
                changes.push(FileChange {
                    status: ChangeStatus::Modified,
                    path: relative_path.to_string_lossy().to_string(),
                });
            }
        }
    }

    // Check for deleted files
    for entry in fs::read_dir(template_dir)? {
        let entry = entry?;
        let path = entry.path();
        if path.is_file() && path.extension().map_or(false, |ext| ext == "tmpl") {
            let relative_path = path.strip_prefix(template_dir).unwrap().with_extension("");
            let src_file = src_dir.join(&relative_path);

            if !src_file.exists() {
                changes.push(FileChange {
                    status: ChangeStatus::Deleted,
                    path: relative_path.to_string_lossy().to_string(),
                });
            }
        }
    }

    Ok(())
}

fn files_differ(file1: &Path, file2: &Path) -> Result<bool, Error> {
    let content1 = fs::read_to_string(file1)?;
    let content2 = fs::read_to_string(file2)?;
    Ok(content1 != content2)
}

fn update_template_files(changes: Vec<FileChange>, config: &NewDaggerModule) -> Result<(), Error> {
    for change in changes {
        let src_file = Path::new(&config.module_src_path).join(&change.path);
        let template_file = Path::new(&config.template_path_by_type)
            .join(&change.path)
            .with_extension("go.tmpl");

        match change.status {
            ChangeStatus::Added | ChangeStatus::Modified => {
                let content = fs::read_to_string(&src_file)?;
                let updated_content = replace_template_variables(&content, &config.module_type);
                fs::write(&template_file, updated_content)?;
                println!("Updated: {}", template_file.display());
            }
            ChangeStatus::Deleted => {
                fs::remove_file(&template_file)?;
                println!("Deleted: {}", template_file.display());
            }
        }
    }
    Ok(())
}

fn replace_template_variables(content: &str, module_type: &str) -> String {
    let (module_name, module_name_pkg) = match module_type {
        "full" => ("ModuleTemplate", "module-template"),
        "light" => ("ModuleTemplateLight", "module-template-light"),
        _ => unreachable!(),
    };

    content
        .replace(module_name, "{{.module_name}}")
        .replace(module_name_pkg, "{{.module_name_pkg}}")
}
