use std::io::Error;
use std::process::Command;

pub fn find_git_root() -> Result<String, Error> {
    let output = Command::new("git")
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
