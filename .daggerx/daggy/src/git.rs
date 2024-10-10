use std::io::Error;
use std::path::{Path, PathBuf};
use std::process::Command;

pub fn find_git_root() -> Result<String, Error> {
    find_git_root_from_path(&std::env::current_dir()?)
}

pub fn find_git_root_from_path(start_path: &Path) -> Result<String, Error> {
    let output = Command::new("git")
        .current_dir(start_path)
        .args(&["rev-parse", "--show-toplevel"])
        .output()?;

    if output.status.success() {
        Ok(String::from_utf8_lossy(&output.stdout).trim().to_string())
    } else {
        Err(Error::new(
            std::io::ErrorKind::Other,
            "Not in a git repository",
        ))
    }
}
