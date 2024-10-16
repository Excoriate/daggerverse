use crate::args::Args;
use crate::configuration::{get_module_configurations, NewDaggerModule};
use std::fs;
use std::io::{Error, ErrorKind};
use std::path::{Path, PathBuf};

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
        let config = get_module_configurations(
            &format!(
                "module-template{}",
                if module_type == "light" { "-light" } else { "" }
            ),
            module_type,
        )?;

        // Debug logging
        println!("Debug: module_src_path = {}", config.module_src_path);
        println!(
            "Debug: template_path_by_type = {}",
            config.template_path_by_type
        );

        if !dry_run {
            // Sync module files
            sync_directory(
                &config.module_src_path,
                &config.template_path_by_type,
                "module",
                &config,
            )?;

            // Sync test files
            let test_src_path = Path::new(&config.module_src_path).join("tests");
            sync_directory(
                &test_src_path.to_string_lossy(),
                &config.template_path_by_type,
                "tests",
                &config,
            )?;

            // Sync example files
            let example_src_path = Path::new(&config.module_src_path)
                .join("examples")
                .join("go");
            sync_directory(
                &example_src_path.to_string_lossy(),
                &config.template_path_by_type,
                "examples/go",
                &config,
            )?;

            println!(
                "Changes synced successfully for {} module type",
                module_type
            );
        } else {
            println!(
                "Dry run: Changes would be synced for {} module type",
                module_type
            );
        }
    }

    Ok(())
}

fn inspect_changes(inspect_type: &str, dry_run: bool, detailed: bool) -> Result<(), Error> {
    let modules_to_inspect = get_modules_to_process(inspect_type)?;

    for module_type in modules_to_inspect {
        println!(
            "Inspecting changes for {} module type (dry run: {})",
            module_type, dry_run
        );
        let config =
            get_module_configurations(&format!("module-template-{}", module_type), module_type)?;
        let changes = detect_changes(&config, detailed)?;

        if !changes.is_empty() {
            println!("Changes detected for {} module type:", module_type);
            for change in &changes {
                println!("{:?}: {}", change.status, change.path);
                if detailed && change.diff.is_some() {
                    println!("Diff:\n{}", change.diff.as_ref().unwrap());
                }
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
        _ => Err(Error::new(
            ErrorKind::InvalidInput,
            "Invalid inspect type. Must be 'full', 'light', or 'all'.",
        )),
    }
}

fn detect_changes(config: &NewDaggerModule, detailed: bool) -> Result<Vec<FileChange>, Error> {
    let mut changes = Vec::new();

    // Debug logging
    println!("Debug: Detecting changes in {}", config.module_src_path);

    // Check main module files
    check_directory_changes(
        &config.module_src_path,
        &config.template_path_by_type,
        &mut changes,
        detailed,
        "",
    )?;

    // Check test files
    let test_src_path = Path::new(&config.module_src_path).join("tests");
    let test_template_path = Path::new(&config.template_path_by_type).join("tests");
    check_directory_changes(
        &test_src_path.to_string_lossy(),
        &test_template_path.to_string_lossy(),
        &mut changes,
        detailed,
        "tests/",
    )?;

    // Check example files
    let example_src_path = Path::new(&config.module_src_path)
        .join("examples")
        .join("go");
    let example_template_path = Path::new(&config.template_path_by_type)
        .join("examples")
        .join("go");
    check_directory_changes(
        &example_src_path.to_string_lossy(),
        &example_template_path.to_string_lossy(),
        &mut changes,
        detailed,
        "examples/go/",
    )?;

    Ok(changes)
}

fn check_directory_changes(
    src_path: &str,
    template_path: &str,
    changes: &mut Vec<FileChange>,
    detailed: bool,
    prefix: &str,
) -> Result<(), Error> {
    let src_dir = Path::new(src_path);
    let template_dir = if prefix.is_empty() {
        Path::new(template_path).join("module")
    } else {
        Path::new(template_path).join(prefix.trim_end_matches('/'))
    };

    // Debug logging
    println!("Debug: Checking directory changes");
    println!("Debug: src_dir = {}", src_dir.display());
    println!("Debug: template_dir = {}", template_dir.display());

    if !src_dir.exists() || !template_dir.exists() {
        println!("Debug: One of the directories does not exist");
        return Ok(());
    }

    for entry in fs::read_dir(src_dir)? {
        let entry = entry?;
        let path = entry.path();
        if path.is_file() && path.extension().map_or(false, |ext| ext == "go") {
            let relative_path = path.strip_prefix(src_dir).unwrap();
            let template_file = template_dir.join(relative_path).with_extension("go.tmpl");

            // Debug logging
            println!("Debug: Checking file: {}", path.display());
            println!("Debug: Template file: {}", template_file.display());

            // Ignore dagger.gen.go files and files in internal/ directories
            if relative_path.file_name().unwrap() == "dagger.gen.go"
                || relative_path
                    .components()
                    .any(|c| c.as_os_str() == "internal")
            {
                continue;
            }

            if !template_file.exists() {
                changes.push(FileChange {
                    status: ChangeStatus::Added,
                    path: format!("{}{}", prefix, relative_path.to_string_lossy()),
                    diff: None,
                });
            } else if files_differ(&path, &template_file)? {
                let diff = if detailed {
                    Some(generate_diff(&path, &template_file)?)
                } else {
                    None
                };
                changes.push(FileChange {
                    status: ChangeStatus::Modified,
                    path: format!("{}{}", prefix, relative_path.to_string_lossy()),
                    diff,
                });
            }
        }
    }

    // Check for deleted files
    for entry in fs::read_dir(template_dir.clone())? {
        let entry = entry?;
        let path = entry.path();
        if path.is_file() && path.extension().map_or(false, |ext| ext == "tmpl") {
            let relative_path = path.strip_prefix(&template_dir).unwrap().with_extension("");
            let src_file = src_dir.join(&relative_path);

            if !src_file.exists() {
                changes.push(FileChange {
                    status: ChangeStatus::Deleted,
                    path: format!("{}{}", prefix, relative_path.to_string_lossy()),
                    diff: None,
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
        let template_file = if change.path.starts_with("tests/") {
            Path::new(&config.template_path_by_type)
                .join(&change.path)
                .with_extension("go.tmpl")
        } else if change.path.starts_with("examples/go/") {
            Path::new(&config.template_path_by_type)
                .join(&change.path)
                .with_extension("go.tmpl")
        } else {
            Path::new(&config.template_path_by_type)
                .join("module")
                .join(&change.path)
                .with_extension("go.tmpl")
        };

        match change.status {
            ChangeStatus::Added | ChangeStatus::Modified => {
                let content = fs::read_to_string(&src_file)?;
                let updated_content = replace_template_variables(&content, &config.module_type);
                fs::create_dir_all(template_file.parent().unwrap())?;
                fs::write(&template_file, updated_content)?;
                println!("Updated: {}", template_file.display());
            }
            ChangeStatus::Deleted => {
                if template_file.exists() {
                    fs::remove_file(&template_file)?;
                    println!("Deleted: {}", template_file.display());
                }
            }
        }
    }

    // Copy testdata directories
    copy_testdata_directories(config)?;

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

fn copy_testdata_directories(config: &NewDaggerModule) -> Result<(), Error> {
    let src_testdata = Path::new(&config.module_src_path)
        .join("tests")
        .join("testdata");
    let dest_testdata = Path::new(&config.template_path_by_type)
        .join("tests")
        .join("testdata");

    if src_testdata.exists() {
        fs::create_dir_all(&dest_testdata)?;
        copy_dir_contents(&src_testdata, &dest_testdata)?;
    }

    let src_example_testdata = Path::new(&config.module_src_path)
        .join("examples")
        .join("go")
        .join("testdata");
    let dest_example_testdata = Path::new(&config.template_path_by_type)
        .join("examples")
        .join("go")
        .join("testdata");

    if src_example_testdata.exists() {
        fs::create_dir_all(&dest_example_testdata)?;
        copy_dir_contents(&src_example_testdata, &dest_example_testdata)?;
    }

    Ok(())
}

fn copy_dir_contents(src: &Path, dest: &Path) -> Result<(), Error> {
    for entry in fs::read_dir(src)? {
        let entry = entry?;
        let path = entry.path();
        if path.is_dir() {
            let new_dest = dest.join(path.file_name().unwrap());
            fs::create_dir_all(&new_dest)?;
            copy_dir_contents(&path, &new_dest)?;
        } else {
            let new_dest = dest.join(path.file_name().unwrap());
            fs::copy(&path, &new_dest)?;
        }
    }
    Ok(())
}

fn generate_diff(_file1: &Path, _file2: &Path) -> Result<String, Error> {
    // Implement diff generation logic here
    // You can use a crate like `diff` or implement a simple line-by-line comparison
    unimplemented!("Diff generation not implemented yet")
}

fn sync_directory(
    src_path: &str,
    template_base_path: &str,
    dir_type: &str,
    config: &NewDaggerModule,
) -> Result<(), Error> {
    let src_dir = Path::new(src_path);
    let template_dir = Path::new(template_base_path).join(dir_type);

    println!("Debug: Syncing directory: {}", src_dir.display());
    println!("Debug: To template directory: {}", template_dir.display());

    if !src_dir.exists() {
        println!(
            "Debug: Source directory does not exist: {}",
            src_dir.display()
        );
        return Ok(());
    }

    fs::create_dir_all(&template_dir)?;

    for entry in fs::read_dir(src_dir)? {
        let entry = entry?;
        let path = entry.path();
        if path.is_file() && path.extension().map_or(false, |ext| ext == "go") {
            let file_name = path.file_name().unwrap().to_str().unwrap();
            if file_name == "dagger.gen.go"
                || path.components().any(|c| c.as_os_str() == "internal")
            {
                continue;
            }

            let template_file = template_dir.join(file_name).with_extension("go.tmpl");
            let content = fs::read_to_string(&path)?;
            let updated_content = replace_template_variables(&content, &config.module_type);
            fs::write(&template_file, updated_content)?;
            println!("Updated: {}", template_file.display());
        }
    }

    Ok(())
}
