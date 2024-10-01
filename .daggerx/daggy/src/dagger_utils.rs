use std::io::{Error, ErrorKind};
use std::path::Path;
use std::process::Command;

pub fn find_dagger_modules() -> Result<Vec<String>, Error> {
    let output = Command::new("find")
        .args(&[".", "-type", "f", "-name", "dagger.json"])
        .output()?;

    if !output.status.success() {
        return Err(Error::new(
            ErrorKind::Other,
            "Failed to execute find command",
        ));
    }

    let modules = String::from_utf8_lossy(&output.stdout)
        .lines()
        .map(|line| {
            Path::new(line)
                .parent()
                .unwrap()
                .to_string_lossy()
                .into_owned()
        })
        .collect::<Vec<String>>();

    Ok(modules)
}

pub fn dagger_module_exists(module: &str) -> Result<(), Error> {
    if module.is_empty() {
        return Err(Error::new(
            ErrorKind::InvalidInput,
            "Module name cannot be empty",
        ));
    }

    // Check if the module already exists in the root of this directory.
    if Path::new(module).exists() {
        return Err(Error::new(
            ErrorKind::AlreadyExists,
            "Module already exists",
        ));
    }

    Ok(())
}
