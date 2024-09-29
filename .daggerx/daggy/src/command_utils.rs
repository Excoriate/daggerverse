use crate::git::find_git_root;
use std::io::Error;
use std::process::{Command, Output, Stdio};

pub fn run_command_with_output(command: &str, target_dir: &str) -> Result<Output, Error> {
    println!("Running command: {}", command);
    let target_directory = if target_dir.is_empty() {
        find_git_root()?
    } else {
        target_dir.to_string()
    };

    println!("Running command in directory: {}", target_directory);
    let output = Command::new("sh")
        .arg("-c")
        .arg(command)
        .current_dir(target_directory)
        .stdout(Stdio::inherit())
        .stderr(Stdio::inherit())
        .output()?;

    if !output.status.success() {
        return Err(Error::new(
            std::io::ErrorKind::Other,
            format!(
                "Command failed with exit code: {} and with error: {}",
                output.status,
                String::from_utf8_lossy(&output.stderr)
            ),
        ));
    }

    Ok(output)
}

pub fn run_go_fmt(module_path: &str) -> Result<(), Error> {
    run_command_with_output("go fmt ./...", module_path)?;
    Ok(())
}
