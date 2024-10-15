use crate::args::Args;
use crate::configuration::{get_module_configurations, NewDaggerModule};
use std::fs;
use std::io::{Error, ErrorKind, Write};
use std::path::Path;
use similar::{ChangeTag, TextDiff};

#[derive(Debug)]
struct FileChange {
    status: ChangeStatus,
    path: String,
    diff: Option<String>,
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
    let detailed = args.detailed.unwrap_or(false);

    println!("Syncing modules...");
    println!("Inspect type: {}", inspect_type);
    println!("Dry run: {}", dry_run);
    println!("Detailed: {}", detailed);

    sync_changes(inspect_type, dry_run, detailed)
}

pub fn inspect_modules_task(args: &Args) -> Result<(), Error> {
    let inspect_type = &args.inspect_type;
    let dry_run = args.dry_run.unwrap_or(false);
    let detailed = args.detailed.unwrap_or(false);

    println!("Inspecting modules...");
    println!("Inspect type: {}", inspect_type);
    println!("Dry run: {}", dry_run);
    println!("Detailed: {}", detailed);

    inspect_changes(inspect_type, dry_run, detailed)
}

fn sync_changes(inspect_type: &str, dry_run: bool, detailed: bool) -> Result<(), Error> {
    let modules_to_sync = get_modules_to_process(inspect_type)?;

    for module_type in modules_to_sync {
        println!(
            "Syncing changes for {} module type (dry run: {})",
            module_type, dry_run
        );
        let config = get_module_configurations(&format!("module-template-{}", module_type), module_type)?;
        let changes = detect_changes(&config, detailed)?;

        if !changes.is_empty() {
            println!("The following changes will be synced:");
            for change in &changes {
                println!("  {}: {}", change.status, change.path);
            }

            if !dry_run {
                if confirm_sync()? {
                    update_template_files(changes, &config)?;
                    println!("Changes synced successfully for {} module type", module_type);
                } else {
                    println!("Sync cancelled for {} module type", module_type);
                }
            } else {
                println!("Dry run: Changes would be synced for {} module type", module_type);
            }
        } else {
            println!("No changes detected for {} module type", module_type);
        }
    }

    Ok(())
}

fn confirm_sync() -> Result<bool, Error> {
    print!("Do you want to proceed with the sync? (y/N): ");
    io::stdout().flush()?;

    let mut input = String::new();
    io::stdin().read_line(&mut input)?;

    Ok(input.trim().to_lowercase() == "y")
}

fn inspect_changes(inspect_type: &str, dry_run: bool, detailed: bool) -> Result<(), Error> {
    let modules_to_inspect = get_modules_to_process(inspect_type)?;

    for module_type in modules_to_inspect {
        println!(
            "Inspecting changes for {} module type (dry run: {})",
            module_type, dry_run
        );
        let config = get_module_configurations(&format!("module-template-{}", module_type), module_type)?;
        let changes = detect_changes(&config, detailed)?;

        if !changes.is_empty() {
            let summary = generate_summary(&changes, detailed);
            println!("Changes detected for {} module type:", module_type);
            println!("{}", summary);
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
        _ => Err(Error::new(
            ErrorKind::InvalidInput,
            "Invalid inspect type. Must be 'full', 'light', or 'all'.",
        )),
    }
}

fn detect_changes(config: &NewDaggerModule, detailed: bool) -> Result<Vec<FileChange>, Error> {
    let mut changes = Vec::new();

    // Check main module files
    check_directory_changes(
        &config.module_src_path,
        &config.template_path_by_type,
        &mut changes,
    )?;

    // Check test files
    let test_src_path = Path::new(&config.path).join("tests");
    let test_template_path = Path::new(&config.template_path_by_type).join("tests");
    check_directory_changes(
        &test_src_path.to_string_lossy(),
        &test_template_path.to_string_lossy(),
        &mut changes,
        detailed,
    )?;

    // Check example files
    let example_src_path = Path::new(&config.path).join("examples").join("go");
    let example_template_path = Path::new(&config.template_path_by_type)
        .join("examples")
        .join("go");
    check_directory_changes(
        &example_src_path.to_string_lossy(),
        &example_template_path.to_string_lossy(),
        &mut changes,
        detailed,
    )?;

    Ok(changes)
}

fn check_directory_changes(
    src_path: &str,
    template_path: &str,
    changes: &mut Vec<FileChange>,
    detailed: bool,
) -> Result<(), Error> {
    let src_dir = Path::new(src_path);
    let template_dir = Path::new(template_path);

    if !src_dir.exists() || !template_dir.exists() {
        return Ok(());
    }

    for entry in fs::read_dir(src_dir)? {
        let entry = entry?;
        let path = entry.path();
        if path.is_file() && is_relevant_go_file(&path) {
            let relative_path = path.strip_prefix(src_dir).unwrap();
            let template_file = template_dir.join(relative_path.with_extension("go.tmpl"));

            if !template_file.exists() {
                changes.push(FileChange {
                    status: ChangeStatus::Added,
                    path: relative_path.to_string_lossy().to_string(),
                    diff: if detailed { Some(fs::read_to_string(&path)?) } else { None },
                });
            } else if files_differ_ignoring_templates(&path, &template_file)? {
                let diff = if detailed {
                    Some(generate_diff_ignoring_templates(&path, &template_file)?)
                } else {
                    None
                };
                changes.push(FileChange {
                    status: ChangeStatus::Modified,
                    path: relative_path.to_string_lossy().to_string(),
                    diff,
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
                    diff: if detailed { Some(fs::read_to_string(&path)?) } else { None },
                });
            }
        }
    }

    Ok(())
}

fn is_relevant_go_file(path: &Path) -> bool {
    if let Some(extension) = path.extension() {
        if extension != "go" {
            return false;
        }
    } else {
        return false;
    }

    let file_name = path.file_name().unwrap().to_str().unwrap();
    if file_name == "dagger.gen.go" {
        return false;
    }

    !path.to_str().unwrap().contains("/internal/")
}

fn files_differ_ignoring_templates(file1: &Path, file2: &Path) -> Result<bool, Error> {
    let content1 = fs::read_to_string(file1)?;
    let content2 = fs::read_to_string(file2)?;
    
    let normalized1 = normalize_content(&content1);
    let normalized2 = normalize_content(&content2);

    Ok(normalized1 != normalized2)
}

fn normalize_content(content: &str) -> String {
    content
        .replace("ModuleTemplate", "{{.module_name}}")
        .replace("module-template", "{{.module_name_pkg}}")
}

fn generate_diff_ignoring_templates(file1: &Path, file2: &Path) -> Result<String, Error> {
    let content1 = fs::read_to_string(file1)?;
    let content2 = fs::read_to_string(file2)?;
    
    let normalized1 = normalize_content(&content1);
    let normalized2 = normalize_content(&content2);

    let diff = TextDiff::from_lines(&normalized1, &normalized2);

    let mut diff_output = String::new();
    for change in diff.iter_all_changes() {
        let sign = match change.tag() {
            ChangeTag::Delete => "-",
            ChangeTag::Insert => "+",
            ChangeTag::Equal => " ",
        };
        diff_output.push_str(&format!("{}{}\n", sign, change));
    }

    Ok(diff_output)
}

fn update_template_files(changes: Vec<FileChange>, config: &NewDaggerModule) -> Result<(), Error> {
    for change in changes {
        let src_file = Path::new(&config.module_src_path).join(&change.path);
        let template_file = Path::new(&config.template_path_by_type).join(&change.path).with_extension("go.tmpl");

        match change.status {
            ChangeStatus::Added | ChangeStatus::Modified => {
                let content = fs::read_to_string(&src_file)?;
                let updated_content = replace_template_variables(&content, config.module_type);
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

fn generate_summary(changes: &[FileChange], detailed: bool) -> String {
    let mut summary = String::new();
    for change in changes {
        let status = match change.status {
            ChangeStatus::Added => "Added",
            ChangeStatus::Modified => "Modified",
            ChangeStatus::Deleted => "Deleted",
        };
        summary.push_str(&format!("{}: {}\n", status, change.path));
        if detailed {
            if let Some(diff) = &change.diff {
                summary.push_str("Diff:\n");
                summary.push_str(diff);
                summary.push_str("\n");
            }
        }
    }
    summary
}
