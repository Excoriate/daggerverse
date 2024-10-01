use std::io::{Error, ErrorKind};
use std::path::Path;

use crate::dagger_commands::run_dagger_develop;
use crate::dagger_utils::find_dagger_modules;

pub fn develop_modules() -> Result<(), Error> {
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
